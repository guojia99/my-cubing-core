package core

import (
	"errors"
	"fmt"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/guojia99/my-cubing-core/model"
)

// addScore 添加一条成绩
func (c *Client) addScore(playerID uint, contestID uint, project model.Project, roundID uint, result []float64, penalty model.ScorePenalty) (err error) {
	// 1. 确定比赛是否存在
	var contest model.Contest
	if err = c.db.Where("id = ?", contestID).First(&contest).Error; err != nil || contest.IsEnd {
		return fmt.Errorf("比赛不存在或已经结束")
	}

	// 2. 获取轮次信息
	var round model.Round
	if err = c.db.Where("id = ?", roundID).First(&round).Error; err != nil {
		return err
	}
	if !round.IsStart {
		return errors.New("该轮次未开始")
	}

	// 3. 玩家信息
	var player = model.Player{}
	if err = c.db.Where("id = ?", playerID).First(&player).Error; err != nil {
		return err
	}

	// 4. 尝试找到本场比赛成绩, 刷新后保存
	var score model.Score
	err = c.db.Model(&model.Score{}).
		Where("player_id = ?", player.ID).
		Where("contest_id = ?", contestID).
		Where("route_id = ?", round.ID).
		First(&score).Error

	if err != nil || score.ID == 0 {
		score = model.Score{
			PlayerID:   player.ID,
			PlayerName: player.Name,
			ContestID:  contestID,
			RouteID:    round.ID,
			Project:    project,
		}
	}
	score.SetResult(result, penalty)
	score.Penalty, _ = jsoniter.MarshalToString(penalty)
	score.IsBestSingle, score.IsBestAvg = false, false

	save := c.db.Save(&score)
	if err = save.Error; err != nil {
		return err
	}
	score.ID = uint(save.RowsAffected)

	// 5. 找到该玩家最佳成绩
	bestS, hasBest, BestA, hasAvg := c.getPlayerBestScoreByProject(playerID, project)

	// 最佳成绩对比
	// 0. 如果当前是D, 则直接不给.
	// 1. 最佳成绩是否是当前成绩ID, 如果是, 则直接给最佳成绩的字段标签
	// 2. 无最佳成绩, 且当前有成绩, 给最佳成绩标签
	// 3. 有最佳成绩, 且不是当前的ID, 对比当前成绩是否最佳.
	if !score.DBest() && ((!hasBest) ||
		(hasBest && bestS.Score.ID == score.ID) ||
		(hasBest && bestS.Score.ID != score.ID && score.IsBestScore(bestS.Score))) {
		score.IsBestSingle = true
	}
	if !score.DAvg() && ((!hasAvg) ||
		(hasAvg && BestA.Score.ID == score.ID) ||
		(hasAvg && BestA.Score.ID != score.ID && score.IsBestAvgScore(BestA.Score))) {
		score.IsBestAvg = true
	}

	c.db.Save(&score)

	return nil
}

// removeScore 删除一条成绩
func (c *Client) removeScore(scoreID uint) (err error) {
	var score model.Score
	if err = c.db.Model(&model.Score{}).Where("id = ?", scoreID).First(&score).Error; err != nil {
		return err
	}
	var contest model.Contest
	if err = c.db.First(&contest, "id = ?", score.ContestID).Error; err != nil {
		return err
	}

	if contest.IsEnd {
		return errors.New("contest is end")
	}

	return c.db.Delete(&score).Error
}

// endContestScore 结束一场比赛并获取记录
func (c *Client) endContestScore(contestID uint) (err error) {
	// 1. 确定比赛是否存在 且非结束的
	var contest model.Contest
	if err = c.db.Where("id = ?", contestID).First(&contest).Error; err != nil || contest.IsEnd {
		return fmt.Errorf("the contest id end or error %+v", err)
	}

	// 2. 获取本场比赛最佳
	thisContestBestSingle, thisContestBestAvg := c.getContestBestSingle(contestID, false), c.getContestBestAvg(contestID, false)
	oldContestBest, oldContestAvg := c.getContestBestSingle(contestID, true), c.getContestBestAvg(contestID, true)

	var records []model.Record
	for key, score := range thisContestBestSingle {
		if _, ok := oldContestBest[key]; ok && score.IsBestScore(oldContestBest[key]) {
			records = append(
				records, model.Record{
					RType:      model.RecordBySingle,
					ScoreId:    score.ID,
					PlayerID:   score.PlayerID,
					PlayerName: score.PlayerName,
					ContestID:  score.ContestID,
				},
			)
		}
	}

	for key, score := range thisContestBestAvg {
		switch score.Project.RouteType() {
		case model.RouteTypeRepeatedly, model.RouteType1rounds:
			continue
		}

		if _, ok := oldContestAvg[key]; ok && score.IsBestAvgScore(oldContestAvg[key]) {
			records = append(
				records, model.Record{
					RType:      model.RecordByAvg,
					ScoreId:    score.ID,
					PlayerID:   score.PlayerID,
					PlayerName: score.PlayerName,
					ContestID:  score.ContestID,
				},
			)
		}
	}
	_ = c.db.Save(&records)

	// 3. 统计排名
	var rounds []model.Round
	c.db.Where("id in ?", contest.GetRoundIds()).Find(&rounds)
	var roundCache = make(map[string][]model.Round)
	for i := 0; i < len(rounds); i++ {
		key := fmt.Sprintf("%s_%d", rounds[i].Project, rounds[i].Number)
		if _, ok := roundCache[key]; !ok {
			roundCache[key] = []model.Round{rounds[i]}
			continue
		}
		roundCache[key] = append(roundCache[key], rounds[i])
	}
	for _, val := range roundCache {
		var ids []uint
		for _, v := range val {
			ids = append(ids, v.ID)
		}
		var scores []model.Score
		c.db.Where("route_id in ?", ids).Find(&scores)
		model.SortScores(scores)
		c.db.Save(&scores)
	}

	// 4. 结束比赛
	contest.IsEnd = true
	contest.EndTime = time.Now()
	return c.db.Save(&contest).Error
}

// getScoreByPlayerContest 获取玩家的某场比赛的成绩列表
func (c *Client) getScoreByPlayerContest(playerId uint, contestId uint) ([]model.Score, error) {
	var player model.Player
	if err := c.db.Where("id = ?", playerId).First(&player).Error; err != nil {
		return nil, err
	}

	var contest model.Contest
	if err := c.db.Where("id = ?", contestId).First(&contest).Error; err != nil {
		return nil, err
	}

	var score []model.Score
	if err := c.db.Where("player_id = ?", playerId).Where("contest_id = ?", contestId).Find(&score).Error; err != nil {
		return nil, err
	}

	sort.Slice(
		score, func(i, j int) bool {
			return score[i].CreatedAt.Sub(score[j].CreatedAt) > 0
		},
	)

	for i, _ := range score {
		var round model.Round
		c.db.Where("id = ?", score[i].RouteID).First(&round)
		score[i].RouteValue = round
	}
	return score, nil
}
