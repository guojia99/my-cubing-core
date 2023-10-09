package core

import (
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
	bestSingle, bestAvg = make(map[model.Project][]model.Score), make(map[model.Project][]model.Score)

	var players []model.Player
	c.db.Find(&players)

	for _, project := range model.AllProjectRoute() {
		bestSingle[project] = make([]model.Score, 0)
		bestAvg[project] = make([]model.Score, 0)
	}

	for _, project := range model.AllProjectRoute() {
		for _, player := range players {
			var best, avg model.Score
			if project.RouteType() == model.RouteTypeRepeatedly {
				if err := c.db.
					Where("player_id = ?", player.ID).
					Where("project = ?", project).
					Where("best > ?", model.DNF).
					Order("best DESC").
					Order("r1 DESC").
					Order("r2").
					Order("r3").
					First(&best).Error; err == nil {
					var round model.Round
					c.db.Where("id = ?", best.RouteID).First(&round)
					best.RouteValue = round
					bestSingle[project] = append(bestSingle[project], best)
				}
				continue
			}
			if err := c.db.
				Where("player_id = ?", player.ID).
				Where("project = ?", project).
				Where("best > ?", model.DNF).
				Order("best").
				First(&best).Error; err == nil {
				var round model.Round
				c.db.Where("id = ?", best.RouteID).First(&round)
				best.RouteValue = round
				bestSingle[project] = append(bestSingle[project], best)
			}
			if err := c.db.
				Where("player_id = ?", player.ID).
				Where("project = ?", project).
				Where("avg > ?", model.DNF).
				Order("avg").
				First(&avg).Error; err == nil {
				var round model.Round
				c.db.Where("id = ?", avg.RouteID).First(&round)
				avg.RouteValue = round
				bestAvg[project] = append(bestAvg[project], avg)
			}
		}

		model.SortByBest(bestSingle[project])
		model.SortByAvg(bestAvg[project])
	}

	return
}

func (c *Client) getBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	b, a := c.getBestScore()
	return b[project], a[project]
}

func (c *Client) getAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	bestSingle, bestAvg = make(map[model.Project]model.Score), make(map[model.Project]model.Score)

	for _, project := range model.AllProjectRoute() {
		var best, avg model.Score
		if project.RouteType() == model.RouteTypeRepeatedly {
			if err := c.db.
				Where("project = ?", project).
				Where("best > ?", model.DNF).
				Order("best DESC").
				Order("r1 DESC").
				Order("r2").
				Order("r3").
				First(&best).Error; err == nil {
				bestSingle[project] = best
			}
			continue
		}
		if err := c.db.
			Where("best > ?", model.DNF).
			Where("project = ?", project).
			Order("best").
			First(&best).Error; err == nil {
			bestSingle[project] = best
		}
		if err := c.db.
			Where("avg > ?", model.DNF).
			Where("project = ?", project).
			Order("avg").
			First(&avg).Error; err == nil {
			bestAvg[project] = avg
		}
	}
	return
}

// getSor 获取所有玩家的Sor排名
func (c *Client) getSor() (single, avg map[model.SorStatisticsKey][]SorScore) {

	var players []model.Player
	if err := c.db.Find(&players).Error; err != nil {
		return
	}
	bestSingle, bestAvg := c.getBestScore()
	single, avg = parserSorSort(players, bestSingle, bestAvg)

	return
}

func (c *Client) getPodiums() []Podiums {
	var players []model.Player
	_ = c.db.Find(&players)
	var out []Podiums
	for _, player := range players {
		pd := c.getPlayerPodiums(player.ID)
		pd.PodiumsResults = nil
		out = append(out, pd)
	}
	SortPodiums(out)
	return out
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
