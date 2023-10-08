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
	preScore := model.PreScore{
		Score: model.Score{
			PlayerID:  request.PlayerID,
			ContestID: request.ContestID,
			RouteID:   request.RoundId,
			Project:   request.Project,
			Result1:   request.Result[0],
			Result2:   request.Result[1],
			Result3:   request.Result[2],
			Result4:   request.Result[3],
			Result5:   request.Result[4],
		},
		Recorder: request.Recorder,
		Source:   request.Source,
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

func (c *Client) getPreScores(page, size int, useFinal, final bool) (int64, []model.PreScore, error) {
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

	if useFinal {
		err = c.db.Where("finish = ?", final).Offset(offset).Limit(limit).Find(&out).Error
		c.db.Where("finish = ?", final).Count(&count)
		return count, out, err
	}
	err = c.db.Offset(offset).Limit(limit).Find(&out).Error
	c.db.Count(&count)
	return count, out, err
}

func (c *Client) getPreScoresByContest(contestID uint) ([]model.PreScore, error) {
	var (
		err error
		out []model.PreScore
	)

	err = c.db.Where("contest_id = ?", contestID).Find(&out).Error
	return out, err
}
