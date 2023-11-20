package core

import (
	"time"

	"github.com/guojia99/my-cubing-core/model"
)

func (c *Client) getRecords(page, size int) (int64, []model.Record, error) {
	if size == 0 || size > 100 {
		size = 100
	}
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * size
	limit := size

	var (
		records []model.Record
		count   int64
	)

	if err := c.db.Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return 0, nil, err
	}
	if err := c.db.Model(&model.Record{}).Count(&count).Error; err != nil {
		return 0, nil, err
	}
	for i := 0; i < len(records); i++ {
		var score model.Score
		c.db.First(&score, "id = ?", records[i].ScoreId)
		records[i].ScoreValue = score

		var contest model.Contest
		c.db.First(&contest, "id = ?", records[i].ContestID)
		records[i].ContestValue = contest
	}
	return count, records, nil
}

// getBestScore 获取所有玩家的最佳成绩
func (c *Client) getBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	var players []model.Player
	c.db.Find(&players)

	var allScore []model.Score
	c.db.Find(&allScore)

	bestSingle, bestAvg = c.sortByScores(allScore, players)
	return
}

func (c *Client) getBestScoreWithTime(s, e time.Time) (bestSingle, bestAvg map[model.Project][]model.Score) {
	var allScore []model.Score
	c.db.Where("created_at > ?", s).Where("created_at < ?", e).Find(&allScore)

	var playerIds []uint
	c.db.Model(&model.Score{}).Distinct("player_id").Where("created_at > ?", s).Where("created_at < ?", e).Pluck("player_id", &playerIds)

	var players []model.Player
	c.db.Where("id in ?", playerIds).Find(&players)

	bestSingle, bestAvg = c.sortByScores(allScore, players)
	return
}

func (c *Client) getBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	b, a := c.getBestScore()
	return b[project], a[project]
}

// getAllProjectBestScores 获取所有项目各自的最佳成绩(最新的成绩为准)
func (c *Client) getAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	var players []model.Player
	c.db.Find(&players)

	var allScore []model.Score
	c.db.Find(&allScore)

	bestSingle, bestAvg = c.getBestByScores(allScore, players)
	return
}

func (c *Client) getAllProjectBestScoresWithTime(s, e time.Time) (bestSingle, bestAvg map[model.Project]model.Score) {
	var allScore []model.Score
	c.db.Where("created_at > ?", s).Where("created_at < ?", e).Find(&allScore)

	var playerIds []uint
	c.db.Model(&model.Score{}).Distinct("player_id").Where("created_at > ?", s).Where("created_at < ?", e).Pluck("player_id", &playerIds)

	var players []model.Player
	c.db.Where("id in ?", playerIds).Find(&players)

	bestSingle, bestAvg = c.getBestByScores(allScore, players)
	return
}

// getSor 获取所有玩家的Sor排名
func (c *Client) getSor() (single, avg map[model.SorStatisticsKey][]SorScore) {
	var players []model.Player
	if err := c.db.Find(&players).Error; err != nil {
		return
	}
	bestSingle, bestAvg := c.getBestScore()
	single, avg = c.parserSorSort(players, bestSingle, bestAvg)
	return
}

func (c *Client) getPodiumsSort(players []model.Player, scoreMaps map[uint]map[model.Project][]model.Score) ([]Podiums, map[uint]*Podiums) {
	var pds []Podiums

	// map key is player id
	var pdsMap = make(map[uint]*Podiums)
	for _, p := range players {
		pdsMap[p.ID] = &Podiums{Player: p}
	}

	// 这个map是  map[uint]map[model.Project][]model.Score
	// 第一层循环是
	for _, scoreMap := range scoreMaps {
		for _, scores := range scoreMap {
			if len(scores) < 3 {
				continue
			}

			maxRank := 3
			if len(scores) <= 5 {
				maxRank = 1
			}

			for _, score := range scores {
				if score.Rank > maxRank {
					break
				}
				podiums := pdsMap[score.PlayerID]
				podiums.Add(score.Rank)
			}
		}
	}

	for _, p := range pdsMap {
		pds = append(pds, *p)
	}
	SortPodiums(pds)

	return pds, pdsMap
}

func (c *Client) getPodiums() ([]Podiums, map[uint]*Podiums) {
	var players []model.Player
	_ = c.db.Find(&players)
	return c.getPodiumsSort(players, c.getLastScoresMapByContest())
}

func (c *Client) getRelativeSor() (allPlayerSor map[model.SorStatisticsKey][]RelativeSor) {
	allPlayerSor = make(map[model.SorStatisticsKey][]RelativeSor)
	var players []model.Player
	if err := c.db.Find(&players).Error; err != nil {
		return
	}
	allBest, allAvg := c.getBestScore()
	return c.parserRelativeSor(players, allBest, allAvg)
}

func (c *Client) getAvgRelativeSor() map[model.SorStatisticsKey]RelativeSor {
	all := c.getRelativeSor()

	var out = make(map[model.SorStatisticsKey]RelativeSor)
	for k, _ := range model.SorKeyMap() {
		if _, ok := all[k]; !ok {
			continue
		}
		if len(all[k]) == 0 {
			continue
		}
		data := RelativeSor{}

		// 前五
		top5 := 5
		if len(all[k]) > 100 {
			top5 = int(float64(len(all[k])) * 0.05)
		}
		if len(all[k]) <= 5 {
			top5 = len(all[k])
		}

		// 平均
		avgStart, avgEnd := 0, len(all[k])
		if len(all[k]) >= 20 {
			diff := int(float64(len(all[k])) * 0.1)
			avgStart = avgStart + diff
			avgEnd = avgEnd - diff
		}
		n := float64(avgEnd - avgStart)

		for idx, val := range all[k] {
			if idx == 0 {
				data.Max = val.Sor
			}
			if idx < top5 {
				data.Top5 += val.Sor / float64(top5)
			}
			if idx >= avgStart && idx < avgEnd {
				data.Avg += val.Sor / n
			}
		}
		out[k] = data
	}
	return out
}
