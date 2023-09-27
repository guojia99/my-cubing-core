package core

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/guojia99/my-cubing-core/model"
)

var defaultProjectRounds = func() []CreateContestRequestRound {
	var out []CreateContestRequestRound
	for _, p := range model.WCAProjectRoute() {
		out = append(out, CreateContestRequestRound{
			Project: p,
			Number:  1,
			Name:    fmt.Sprintf("%s单轮赛", p.Cn()),
			Part:    1,
			IsStart: true,
			Final:   true,
			Upsets:  []string{},
		})
	}
	return out
}()

// addContest 添加一场比赛
func (c *Client) addContest(req AddContestRequest) error {
	var contest model.Contest
	if err := c.db.Where("name = ?", req.Name).First(&contest).Error; err == nil {
		return errors.New("name error")
	}

	contest = model.Contest{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		StartTime: func() time.Time {
			if req.StartTime == 0 {
				return time.Now()
			}
			return time.Unix(req.StartTime, 0)
		}(),
		EndTime: func() time.Time {
			if req.EndTime == 0 {
				return time.Now().Add(time.Hour * 24 * 60) // 60 day
			}
			return time.Unix(req.EndTime, 0)
		}(),
	}

	if err := c.db.Save(&contest).Error; err != nil {
		return err
	}

	if len(req.Rounds) == 0 || req.Rounds == nil {
		req.Rounds = defaultProjectRounds
	}
	var rounds []model.Round
	for _, val := range req.Rounds {
		var round = model.Round{
			ContestID: contest.ID,
			Project:   val.Project,
			Number:    val.Number,
			Part:      val.Part,
			IsStart:   val.IsStart,
			Name:      val.Name,
			Final:     val.Final,
		}
		round.SetUpsets(val.Upsets)
		c.db.Create(&round)
		rounds = append(rounds, round)
	}
	contest.SetRoundIds(rounds)
	err := c.db.Save(&contest).Error
	return err
}

// removeContest 删除比赛
func (c *Client) removeContest(contestId uint) error {
	var contest model.Contest
	if err := c.db.Where("id = ?", contestId).First(&contest).Error; err != nil {
		return err
	}

	var count int64
	if err := c.db.Model(&model.Score{}).Where("contest_id = ?", contestId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("the contest has score, can't not delete")
	}

	err := c.db.Delete(&contest).Error
	return err
}

// getContest 获取比赛信息
func (c *Client) getContest(contestId uint) (contest model.Contest, err error) {
	if err = c.db.First(&contest, "id = ?", contestId).Error; err != nil {
		return
	}

	if err = c.db.Find(&contest.Rounds, "id in ?", contest.GetRoundIds()).Error; err != nil {
		return
	}
	return
}

// getContests 获取比赛列表
func (c *Client) getContests(page, size int, typ string) (int64, []model.Contest, error) {
	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 100
	}

	offset := (page - 1) * size
	limit := size

	var contests []model.Contest
	var err error
	var count int64

	if typ == "" {
		err = c.db.Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&contests).Error
		c.db.Model(&model.Contest{}).Count(&count)
	} else {
		err = c.db.Where("c_type = ?", typ).Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&contests).Error
		c.db.Model(&model.Contest{}).Where("c_type = ?", typ).Count(&count)
	}
	if err != nil {
		return 0, nil, err
	}

	for i := 0; i < len(contests); i++ {
		c.db.Where("id in ?", contests[i].GetRoundIds()).Find(&contests[i].Rounds)
	}

	return count, contests, nil
}

// getContestSor 获取比赛的Sor排名
func (c *Client) getContestSor(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore) {
	single, avg = make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey][]SorScore, len(model.SorKeyMap()))

	// 查这场比赛所有选手
	var playerIDs []uint64
	if c.db.Model(&model.Score{}).Distinct("player_id").Where("contest_id = ?", contestID).Pluck("player_id", &playerIDs); len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	bestSingleCache, bestAvgCache := c.getContestAllBestScores(contestID)
	fmt.Println(bestAvgCache[model.Cube333BF])
	single, avg = parserSorSort(players, bestSingleCache, bestAvgCache)
	return
}

// getContestScore 获取比赛所有成绩
func (c *Client) getContestScore(contestID uint) map[model.Project][]RoutesScores {
	var out = make(map[model.Project][]RoutesScores)

	var contest model.Contest
	if err := c.db.First(&contest, "id = ?", contestID).Error; err != nil {
		return nil
	}
	var rounds []model.Round
	if err := c.db.
		Model(&model.Round{}).
		Where("id in ?", contest.GetRoundIds()).
		Order("number DESC").
		Find(&rounds).Error; err != nil {
		return nil
	}

	// 按number分类
	var roundCache = make(map[string][]model.Round)
	for _, val := range rounds {
		key := fmt.Sprintf("%s_%d", val.Project, val.Number)
		if data, ok := roundCache[key]; ok {
			data = append(data, val)
			roundCache[key] = data
			continue
		}
		roundCache[key] = []model.Round{val}
	}

	// 查询所有成绩
	for _, rs := range roundCache {
		if len(rs) == 0 {
			continue
		}
		var pj = rs[0].Project
		var scores []model.Score
		var ids []uint
		for _, v := range rs {
			ids = append(ids, v.ID)
		}
		c.db.Where("route_id in ?", ids).Find(&scores)
		model.SortScores(scores)

		if _, ok := out[pj]; !ok {
			out[pj] = make([]RoutesScores, 0)
		}
		out[pj] = append(out[pj], RoutesScores{
			Final:  rs[0].Final,
			Round:  rs,
			Scores: scores,
		})
	}

	for _, pj := range model.AllProjectRoute() {
		if _, ok := out[pj]; !ok {
			continue
		}
		sort.Slice(out[pj], func(i, j int) bool { return out[pj][i].Final })
	}

	return out
}

// getContestPodiums 获取比赛奖牌榜
func (c *Client) getContestPodiums(contestID uint) (out []Podiums) {
	// 未结束的比赛无领奖台
	var contest model.Contest
	if err := c.db.Where("id = ? ", contestID).First(&contest).Error; err != nil || !contest.IsEnd {
		return
	}

	// 查这场比赛所有选手
	var playerIDs []uint64
	c.db.
		Model(&model.Score{}).
		Distinct("player_id").
		Where("contest_id = ?", contestID).
		Pluck("player_id", &playerIDs)
	if len(playerIDs) == 0 {
		return
	}
	var players []model.Player
	c.db.Where("id in ?", playerIDs).Find(&players)

	var cache = make(map[uint]*Podiums)
	for _, tt := range c.getContestTop(contestID, 3) {
		for _, val := range tt {
			if _, ok := cache[val.PlayerID]; !ok {
				cache[val.PlayerID] = &Podiums{}
			}

			switch val.Rank {
			case 1:
				cache[val.PlayerID].Gold += 1
			case 2:
				cache[val.PlayerID].Silver += 1
			case 3:
				cache[val.PlayerID].Bronze += 1
			}
		}
	}

	for _, player := range players {
		podiums := Podiums{
			Player: player,
		}
		if val, ok := cache[player.ID]; ok {
			podiums.Gold = val.Gold
			podiums.Silver = val.Silver
			podiums.Bronze = val.Bronze
		}
		out = append(out, podiums)
	}
	SortPodiums(out)
	return
}

// getContestRecord 获取比赛的记录
func (c *Client) getContestRecord(contestID uint) []RecordMessage {
	var out []RecordMessage

	var contest model.Contest
	if err := c.db.First(&contest, "id = ?", contestID).Error; err != nil {
		return out
	}

	var records []model.Record
	if err := c.db.Where("contest_id = ?", contestID).Find(&records).Error; err != nil {
		return out
	}

	for _, record := range records {
		var player model.Player
		var score model.Score
		_ = c.db.First(&player, "id = ?", record.PlayerID).Error
		_ = c.db.First(&score, "id = ?", record.ScoreId).Error

		out = append(out, RecordMessage{
			Record:  record,
			Player:  player,
			Score:   score,
			Contest: contest,
		})
	}
	return out
}
