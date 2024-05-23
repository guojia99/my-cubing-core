package core

import (
	"fmt"
	"sort"

	"github.com/guojia99/my-cubing-core/model"
)

// todo 考虑多线程
func (c *Client) getContestPlayers(contestID uint, past bool) (playerIDs []uint64, players []model.Player) {
	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}

	// 查这场比赛所有选手
	if c.db.Model(&model.Score{}).Distinct("player_id").Where(conn, contestID).Pluck("player_id", &playerIDs); len(playerIDs) == 0 {
		return nil, nil
	}
	c.db.Where("id in ?", playerIDs).Find(&players)
	return
}

// playerByProjectMap 最佳成绩列表和平均成绩列表制作成基于项目的玩家映射表
func (c *Client) playerByProjectMap(bestSingle, bestAvg map[model.Project][]model.Score) (singlePlayerDict, avgPlayerDict map[model.Project]map[uint]model.Score) {
	singlePlayerDict = make(map[model.Project]map[uint]model.Score)
	avgPlayerDict = make(map[model.Project]map[uint]model.Score)

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
	return
}

func (c *Client) parserRelativeSor(players []model.Player, bestSingle, bestAvg map[model.Project][]model.Score) (allPlayerSor map[model.SorStatisticsKey][]RelativeSor) {
	allPlayerSor = make(map[model.SorStatisticsKey][]RelativeSor)

	bestSingle1, bestAvg1 := c.getAllProjectBestScores()                         // 全场最佳的成绩
	singlePlayerDict, avgPlayerDict := c.playerByProjectMap(bestSingle, bestAvg) // 个人成绩缓存 map[model.Project]map[uint]model.Score
	for sorKey, projects := range model.SorKeyMap() {
		var playerCache = make(map[uint]*RelativeSor)

		for _, player := range players {
			playerCache[player.ID] = &RelativeSor{Player: player}
			for _, pj := range projects {
				// 1. 如果该项目无最佳,则代表没人玩
				if _, ok := bestSingle1[pj]; !ok {
					continue
				}
				// 1, 计时项目:
				// 项目分: ((gr最佳 + 1 / 个人最佳 + 1)  + (gr平均 + 1  / 个人平均最佳 + 1)) * 10
				// 没有玩的项目直接给1分
				//2. 多盲项目:
				// 项目分:  (个人最佳分 + 1 / gr分 + 1 )  * gr个数

				// 2. 如果选手无最佳代表 平均都没有
				if _, ok := singlePlayerDict[pj][player.ID]; !ok || singlePlayerDict[pj][player.ID].Best <= model.DNF {
					switch pj.RouteType() {
					case model.RouteTypeRepeatedly:
						playerCache[player.ID].Sor += (1.0 / bestSingle1[pj].Best) * bestSingle1[pj].Result1
					case model.RouteType1rounds:
						playerCache[player.ID].Sor += 1
					default:
						playerCache[player.ID].Sor += 2
					}
					continue
				}

				// 3. 添加单次成绩
				switch pj.RouteType() {
				case model.RouteTypeRepeatedly:
					playerCache[player.ID].Sor += (singlePlayerDict[pj][player.ID].Best / bestSingle1[pj].Best) * bestSingle1[pj].Result1
					continue
				case model.RouteType1rounds:
					playerCache[player.ID].Sor += ((bestSingle1[pj].Best + 1) / (singlePlayerDict[pj][player.ID].Best + 1)) * 10
					continue
				default:
					playerCache[player.ID].Sor += ((bestSingle1[pj].Best + 1) / (singlePlayerDict[pj][player.ID].Best + 1)) * 10
				}

				// 4, 如果该项目无最佳平均或该项目无平均
				if _, ok := bestAvg1[pj]; !ok {
					continue
				}

				// 5. 如果选手有最佳但无平均
				if _, ok := avgPlayerDict[pj][player.ID]; !ok || avgPlayerDict[pj][player.ID].Avg <= model.DNF {
					playerCache[player.ID].Sor += 1
					continue
				}

				playerCache[player.ID].Sor += ((bestAvg1[pj].Avg + 1) / (singlePlayerDict[pj][player.ID].Avg + 1)) * 10
			}
		}

		var data []RelativeSor
		for _, val := range playerCache {
			data = append(
				data, RelativeSor{
					Player: val.Player,
					Sor:    val.Sor,
				},
			)
		}
		sort.Slice(data, func(i, j int) bool { return data[i].Sor >= data[j].Sor })
		allPlayerSor[sorKey] = data
	}
	return
}

func (c *Client) parserSorSort(players []model.Player, bestSingle, bestAvg map[model.Project][]model.Score) (single, avg map[model.SorStatisticsKey][]SorScore) {
	single, avg = make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap()))

	singlePlayerDict, avgPlayerDict := c.playerByProjectMap(bestSingle, bestAvg)

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

		sort.Slice(
			s, func(i, j int) bool {
				if s[i].SingleCount == s[j].SingleCount {
					return s[i].SingleProjects < s[j].SingleProjects
				}
				return s[i].SingleCount < s[j].SingleCount
			},
		)
		sort.Slice(
			a, func(i, j int) bool {
				if a[i].AvgCount == a[j].AvgCount {
					return a[i].AvgProjects < a[j].AvgProjects
				}
				return a[i].AvgCount < a[j].AvgCount
			},
		)

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
		for i := 0; i < len(lastScores.Scores) && i <= n; i++ {
			if i < n {
				s = append(s, lastScores.Scores[i])
			}
			if lastScores.Scores[i].DBest() {
				continue
			}
		}
		out[project] = s
	}
	return out
}

// getContestBestSingle 获取比赛每个项目的最佳成绩, 如果使用 past则不是该成绩
func (c *Client) getContestBestSingle(contestID uint, past bool) map[model.Project]model.Score {
	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}

	var allScore []model.Score
	c.db.Where(conn, contestID).Find(&allScore)

	_, players := c.getContestPlayers(contestID, past)
	single, _ := c.getBestByScores(allScore, players)
	return single
}

// getContestBestAvg 获取比赛每个项目的最佳平均成绩
func (c *Client) getContestBestAvg(contestID uint, past bool) map[model.Project]model.Score {
	conn := "contest_id = ?"
	if past {
		conn = "contest_id != ?"
	}

	var allScore []model.Score
	c.db.Where(conn, contestID).Find(&allScore)

	_, players := c.getContestPlayers(contestID, past)
	_, avg := c.getBestByScores(allScore, players)
	return avg
}

func (c *Client) sortByScores(allScore []model.Score, players []model.Player) (bestSingle, bestAvg map[model.Project][]model.Score) {
	bestSingle, bestAvg = make(map[model.Project][]model.Score), make(map[model.Project][]model.Score)

	for _, project := range model.AllProjectRoute() {
		bestSingle[project] = make([]model.Score, 0)
		bestAvg[project] = make([]model.Score, 0)
	}

	// key 是 project + playerID
	var singleCache = make(map[string]model.Score)
	var avgCache = make(map[string]model.Score)

	// todo 分片
	for _, score := range allScore {
		key := fmt.Sprintf("%d_%s", score.PlayerID, score.Project)

		if _, ok := singleCache[key]; !ok || score.IsBestScore(singleCache[key]) {
			singleCache[key] = score
		}

		switch score.Project.RouteType() {
		case model.RouteType1rounds, model.RouteTypeRepeatedly:
			continue
		}

		if _, ok := avgCache[key]; !ok || score.IsBestAvgScore(avgCache[key]) {
			avgCache[key] = score
		}
	}

	for _, project := range model.AllProjectRoute() {
		for _, player := range players {
			key := fmt.Sprintf("%d_%s", player.ID, project)

			if single, ok := singleCache[key]; ok && !single.DBest() {
				bestSingle[project] = append(bestSingle[project], single)
			}
			if avg, ok := avgCache[key]; ok && !avg.DAvg() {
				bestAvg[project] = append(bestAvg[project], avg)
			}
		}

		model.SortByBest(bestSingle[project])
		model.SortByAvg(bestAvg[project])
	}
	return
}

func (c *Client) getBestByScores(allScore []model.Score, players []model.Player) (bestSingle, bestAvg map[model.Project]model.Score) {
	bestSingles, bestAvgs := c.sortByScores(allScore, players)

	bestSingle, bestAvg = make(map[model.Project]model.Score), make(map[model.Project]model.Score)

	for _, pj := range model.AllProjectRoute() {
		if val, ok := bestSingles[pj]; ok && len(val) > 0 {
			bestSingle[pj] = val[0]
		}

		if val, ok := bestAvgs[pj]; ok && len(val) > 0 {
			bestAvg[pj] = val[0]
		}
	}
	return
}

// getBestByScoresMuil 获取头部成绩， 且头部成绩有多个
func (c *Client) getBestByScoresMuil(allScore []model.Score, players []model.Player) (bestSingle, bestAvg map[model.Project][]model.Score) {
	bestSingles, bestAvgs := c.sortByScores(allScore, players)
	bestSingle, bestAvg = make(map[model.Project][]model.Score), make(map[model.Project][]model.Score)

	for _, pj := range model.AllProjectRoute() {
		if val, ok := bestSingles[pj]; ok && len(val) > 0 {
			var sco []model.Score
			sco = append(sco, val[0])
			for i := 1; i < len(sco); i++ {
				if val[i].IsBestScore(sco[0]) {
					sco = append(sco, val[i])
				}
			}
			bestSingles[pj] = sco
		}

		if val, ok := bestAvgs[pj]; ok && len(val) > 0 {
			var sco []model.Score
			sco = append(sco, val[0])
			for i := 1; i < len(sco); i++ {
				if val[i].IsBestAvgScore(sco[0]) {
					sco = append(sco, val[i])
				}
			}
			bestAvgs[pj] = sco
		}
	}
	return
}

// getContestAllBestScores 获取某比赛所有最佳成绩排名
func (c *Client) getContestAllBestScores(contestID uint) (bestSingle, bestAvg map[model.Project][]model.Score) {
	playerIDs, players := c.getContestPlayers(contestID, false)
	var allScore []model.Score
	c.db.Where("contest_id = ?").Where("player_id in ?", playerIDs).Find(&allScore)
	bestSingle, bestAvg = c.sortByScores(allScore, players)
	return
}

// getLastScoresMapByContest 获取每个比赛每个项目最后一轮的成绩
func (c *Client) getLastScoresMapByContest() (out map[uint]map[model.Project][]model.Score) {
	var roundIds []uint
	c.db.Model(&model.Round{}).Select("id").Where("final = ?", true).Find(&roundIds)

	var scores []model.Score
	c.db.Where("route_id in ?", roundIds).Find(&scores)

	var contests []model.Contest
	c.db.Find(&contests)

	out = make(map[uint]map[model.Project][]model.Score)
	for _, contest := range contests {
		out[contest.ID] = make(map[model.Project][]model.Score)
	}
	for _, score := range scores {
		if out[score.ContestID][score.Project] == nil {
			out[score.ContestID][score.Project] = make([]model.Score, 0)
		}
		out[score.ContestID][score.Project] = append(out[score.ContestID][score.Project], score)
	}

	for contestID, _ := range out {
		for pj, _ := range out[contestID] {
			model.SortScores(out[contestID][pj])
		}
	}

	return out
}
