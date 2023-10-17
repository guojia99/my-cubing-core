package core

import (
	"errors"

	jsoniter "github.com/json-iterator/go"

	"github.com/guojia99/my-cubing-core/model"
)

func (c *Client) addPreScore(request AddPreScoreRequest) error {
	for len(request.Result) < 5 {
		request.Result = append(request.Result, model.DNF)
	}

	// 基本数据校验
	var player model.Player
	if err := c.db.First(&player, "id = ?", request.PlayerID).Error; err != nil {
		return errors.New("找不到选手")
	}
	var contest model.Contest
	if err := c.db.First(&contest, "id = ?", request.ContestID).Error; err != nil {
		return errors.New("找不到比赛")
	}
	var round model.Round
	if err := c.db.First(&round, "id = ?", request.RoundId).Error; err != nil {
		return errors.New("找不到该轮次")
	}
	if round.ContestID != contest.ID {
		return errors.New("该轮次和比赛无法对应")
	}
	if round.Project != request.Project {
		return errors.New("该轮次和项目无法对应")
	}

	// 查询预录入数据
	var preScore model.PreScore
	err := c.db.Where("player_id = ?", request.PlayerID).
		Where("contest_id = ?", request.ContestID).
		Where("route_id = ?", request.RoundId).
		Where("project = ?", request.Project).
		Order("created_at DESC").
		First(&preScore).Error
	if err == nil && preScore.ID != 0 && !preScore.Finish {
		return errors.New("该预录入成绩已存在")
	}

	preScore = model.PreScore{
		Score: model.Score{
			PlayerID:   request.PlayerID,
			PlayerName: player.Name,
			ContestID:  request.ContestID,
			RouteID:    request.RoundId,
			Project:    request.Project,
			Result1:    request.Result[0],
			Result2:    request.Result[1],
			Result3:    request.Result[2],
			Result4:    request.Result[3],
			Result5:    request.Result[4],
		},
		ContestName: contest.Name,
		RoundName:   round.Name,
		Recorder:    request.Recorder,
		Source:      request.Source,
	}

	preScore.Penalty, _ = jsoniter.MarshalToString(request.Penalty)
	return c.db.Create(&preScore).Error
}

func (c *Client) processPreScore(request ProcessPreScoreRequest) error {
	var preScore model.PreScore
	if err := c.db.Where("id = ?", request.Id).First(&preScore).Error; err != nil {
		return errors.New("该记录未找到")
	}
	if preScore.Finish {
		return errors.New("该记录已经被处理")
	}

	switch request.FinishDetail {
	case model.FinishDetailRecord:
		var penalty model.ScorePenalty
		_ = jsoniter.UnmarshalFromString(preScore.Penalty, &penalty)
		err := c.addScore(preScore.PlayerID, preScore.ContestID, preScore.Project, preScore.RouteID, preScore.GetResult(), penalty)
		if err != nil {
			return err
		}
	case model.FinishDetailDelete:
	case model.FinishDetailNeglect:

	}
	preScore.Finish = true
	preScore.FinishDetail = request.FinishDetail
	preScore.Processor = request.Processor

	return c.db.Save(&preScore).Error
}

func (c *Client) getPreScores(page, size int, final Bool) (int64, []model.PreScore, error) {
	var (
		count int64
		err   error
		out   []model.PreScore
	)

	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 100
	}

	offset := (page - 1) * size
	limit := size

	if final > NotBool {
		err = c.db.Where("finish = ?", final == TrueBool).Offset(offset).Limit(limit).Find(&out).Error
		c.db.Model(&model.PreScore{}).Where("finish = ?", final == TrueBool).Count(&count)
	} else {
		err = c.db.Offset(offset).Limit(limit).Find(&out).Error
		c.db.Model(&model.PreScore{}).Count(&count)
	}

	return count, out, err
}

func (c *Client) getPreScoresByContest(contestID uint, page, size int, final Bool) (int64, []model.PreScore, error) {
	var (
		count int64
		err   error
		out   []model.PreScore
	)

	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 100
	}

	offset := (page - 1) * size
	limit := size

	if final > NotBool {
		err = c.db.Where("contest_id = ?", contestID).Where("finish = ?", final == TrueBool).Offset(offset).Limit(limit).Find(&out).Error
		c.db.Model(&model.PreScore{}).Where("contest_id = ?", contestID).Where("finish = ?", final == TrueBool).Count(&count)
		return count, out, err
	}
	err = c.db.Where("contest_id = ?", contestID).Offset(offset).Limit(limit).Find(&out).Error
	c.db.Model(&model.PreScore{}).Where("contest_id = ?", contestID).Count(&count)
	return count, out, err
}

func (c *Client) getPreScoresByPlayer(playerID uint, page, size int, final Bool) (int64, []model.PreScore, error) {
	var (
		count int64
		err   error
		out   []model.PreScore
	)

	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 100
	}

	offset := (page - 1) * size
	limit := size

	if final > NotBool {
		err = c.db.Where("player_id = ?", playerID).Where("finish = ?", final == TrueBool).Offset(offset).Limit(limit).Find(&out).Error
		c.db.Model(&model.PreScore{}).Where("player_id = ?", playerID).Where("finish = ?", final == TrueBool).Count(&count)
		return count, out, err
	}
	err = c.db.Offset(offset).Limit(limit).Find(&out).Error
	c.db.Model(&model.PreScore{}).Where("player_id = ?", playerID).Count(&count)
	return count, out, err
}
