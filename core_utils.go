package core

import (
	"sort"

	"github.com/guojia99/my-cubing-core/model"
)

// parserSorSort 解析sor
func parserSorSort(players []model.Player, bestSingle, bestAvg map[model.Project][]model.Score) (single, avg map[model.SorStatisticsKey][]SorScore) {
	single, avg = make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap()))

	// todo 抽出来
	var singlePlayerDict = make(map[model.Project]map[uint]model.Score)
	var avgPlayerDict = make(map[model.Project]map[uint]model.Score)

	// 1. 做map缓存
	for _, pj := range model.AllProjectRoute() {
		singlePlayerDict[pj] = make(map[uint]model.Score)
		avgPlayerDict[pj] = make(map[uint]model.Score)

		if _, ok := bestSingle[pj]; ok {
			for _, val := range bestSingle[pj] {
				singlePlayerDict[pj][val.PlayerID] = val
			}
		}
		if _, ok := bestAvg[pj]; ok {
			for _, val := range bestAvg[pj] {
				avgPlayerDict[pj][val.PlayerID] = val
			}
		}
	}

	// 2. 全项目排序
	for sorKey, projects := range model.SorKeyMap() {
		var s = make([]SorScore, 0)
		var a = make([]SorScore, 0)
		var playerCache = make(map[uint]*SorScore)

		for _, player := range players {
			playerCache[player.ID] = &SorScore{Player: player}

			for _, pj := range projects {
				if val, ok := singlePlayerDict[pj][player.ID]; ok {
					playerCache[val.PlayerID].SingleCount += int64(val.Rank)
					playerCache[val.PlayerID].SingleProjects += 1
				} else if len(bestSingle[pj]) > 0 {
					playerCache[player.ID].SingleCount += int64(bestSingle[pj][len(bestSingle[pj])-1].Rank + 1)
				}

				if val, ok := avgPlayerDict[pj][player.ID]; ok {
					playerCache[val.PlayerID].AvgCount += int64(val.Rank)
					playerCache[val.PlayerID].AvgProjects += 1
				} else if len(bestAvg[pj]) > 0 {
					playerCache[player.ID].AvgCount += int64(bestAvg[pj][len(bestAvg[pj])-1].Rank + 1)
				}
			}
		}

		for _, val := range playerCache {
			s = append(s, SorScore{Player: val.Player, SingleCount: val.SingleCount, SingleProjects: val.SingleProjects})
			a = append(a, SorScore{Player: val.Player, AvgCount: val.AvgCount, AvgProjects: val.AvgProjects})
		}

		sort.Slice(s, func(i, j int) bool {
			if s[i].SingleCount == s[j].SingleCount {
				return s[i].SingleProjects < s[j].SingleProjects
			}
			return s[i].SingleCount < s[j].SingleCount
		})
		sort.Slice(a, func(i, j int) bool {
			if a[i].AvgCount == a[j].AvgCount {
				return a[i].AvgProjects < a[j].AvgProjects
			}
			return a[i].AvgCount < a[j].AvgCount
		})

		single[sorKey] = s
		avg[sorKey] = a
	}
	return
}

// getPodiumsByPlayer 获取比赛前top N成绩, 会依据不同项目按最佳成绩或最佳平均来区分输出
func (c *Client) getContestTop(contestID uint, n int) map[model.Project][]model.Score {
	var contest model.Contest
	if err := c.db.Where("id = ? ", contestID).First(&contest).Error; err != nil || !contest.IsEnd {
		return nil
	}

	var out = make(map[model.Project][]model.Score)

	// todo 用全部差的方式会比较好

	allScores := c.getContestScore(contestID)
	for _, project := range model.AllProjectRoute() {
		scores, ok := allScores[project]
		if !ok {
			continue
		}
		if len(scores) == 0 {
			continue
		}

		// 只需要最后一轮的成绩
		// todo 考虑同名次
		lastScores := scores[0]

		var s []model.Score
		for i := 0; i < len(lastScores.Scores); i++ {
			if i < n {
				s = append(s, lastScores.Scores[i])
			}
			if i > n {
				break
			}
		}
		out[project] = s
	}
	return out
}

// getContestBestSingle 获取比赛每个项目的最佳成绩
func (c *Client) getContestBestSingle(contestID uint, past bool) map[model.Project]model.Score {
	var out = make(map[model.Project]model.Score)

	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}

	for _, project := range model.AllProjectRoute() {
		var score model.Score
		var err error

		switch project.RouteType() {
		case model.RouteTypeRepeatedly:
			err = c.db.
				Where(conn, contestID).
				Where("project = ?", project).
				Where("best > ?", model.DNF).
				Order("best DESC").
				Order("r2").
				Order("r3").
				Order("created_at").
				First(&score).Error
		default:
			err = c.db.
				Where(conn, contestID).
				Where("project = ?", project).
				Where("best > ?", model.DNF).
				Order("best").
				Order("created_at").
				First(&score).Error
		}
		if err != nil {
			continue
		}

		out[project] = score
	}
	return out
}

// getContestBestAvg 获取比赛每个项目的最佳平均成绩
func (c *Client) getContestBestAvg(contestID uint, past bool) map[model.Project]model.Score {
	var out = make(map[model.Project]model.Score)
	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}
	for _, project := range model.AllProjectRoute() {
		var score model.Score
		if err := c.db.
			Where(conn, contestID).
			Where("project = ?", project).
			Where("avg > ?", model.DNF).
			Order("avg").
			Order("created_at").
			First(&score).Error; err != nil {
			continue
		}
		out[project] = score
	}
	return out
}

// getContestAllBestScores 获取某比赛所有最佳成绩排名
func (c *Client) getContestAllBestScores(contestID uint) (single, avg map[model.Project][]model.Score) {
	single = make(map[model.Project][]model.Score)
	avg = make(map[model.Project][]model.Score)

	// 查这场比赛所有选手
	var playerIDs []uint64
	if c.db.Model(&model.Score{}).Distinct("player_id").Where("contest_id = ?", contestID).Pluck("player_id", &playerIDs); len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	for _, project := range model.AllProjectRoute() {
		single[project] = make([]model.Score, 0)
		avg[project] = make([]model.Score, 0)
		for _, player := range players {
			var b, a model.Score
			if project.RouteType() == model.RouteTypeRepeatedly {
				if err := c.db.
					Where("player_id = ?", player.ID).
					Where("project = ?", project).
					Where("best > ?", model.DNF).
					Where("contest_id = ?", contestID).
					Order("best").
					Order("r1 DESC").
					Order("r2").
					Order("r3").
					First(&b).Error; err == nil {
					single[project] = append(single[project], b)
				}
				continue
			}
			if err := c.db.
				Where("player_id = ?", player.ID).
				Where("project = ?", project).
				Where("best > ?", model.DNF).
				Where("contest_id = ?", contestID).
				Order("best").
				First(&b).Error; err == nil {
				single[project] = append(single[project], b)
			}
			if err := c.db.
				Where("player_id = ?", player.ID).
				Where("project = ?", project).
				Where("avg > ?", model.DNF).
				Where("contest_id = ?", contestID).
				Order("avg").
				First(&a).Error; err == nil {
				avg[project] = append(avg[project], a)
			}
		}
	}

	for _, project := range model.AllProjectRoute() {
		model.SortByBest(single[project])
		model.SortByAvg(avg[project])
	}
	return
}
