/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/15 下午10:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"fmt"
	"runtime"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing-core/model"
)

var _ Core = &Client{}

type Client struct {
	debug     bool
	db        *gorm.DB
	cache     *cache.Cache
	cacheTime time.Duration
}

func (c *Client) reloadCache() {
	c.cache.Flush()
	runtime.GC()
}

// score Core

func (c *Client) AddScore(req AddScoreRequest) error {
	if err := c.addScore(req.PlayerID, req.ContestID, req.Project, req.RoundId, req.Result, req.Penalty); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) RemoveScore(scoreID uint) error {
	if err := c.removeScore(scoreID); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) EndContestScore(contestId uint) error {
	if err := c.endContestScore(contestId); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) GetScoreByPlayerContest(playerId uint, contestId uint) ([]model.Score, error) {
	key := fmt.Sprintf("GetScoreByPlayerContest_%v_%v", playerId, contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]model.Score), nil
	}

	out, err := c.getScoreByPlayerContest(playerId, contestId)
	if err != nil {
		return nil, err
	}
	_ = c.cache.Add(key, out, c.cacheTime)
	return out, nil
}

func (c *Client) AddPreScore(request AddPreScoreRequest) error {
	return c.addPreScore(request)
}

func (c *Client) ProcessPreScore(request ProcessPreScoreRequest) error {
	if err := c.processPreScore(request); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) GetPreScores(page, size int, final Bool) (int64, []model.PreScore, error) {
	return c.getPreScores(page, size, final)
}

func (c *Client) GetPreScoresByPlayer(playerID uint, page, size int, final Bool) (int64, []model.PreScore, error) {
	return c.getPreScoresByPlayer(playerID, page, size, final)
}

func (c *Client) GetPreScoresByContest(contestID uint, page, size int, final Bool) (int64, []model.PreScore, error) {
	return c.getPreScoresByContest(contestID, page, size, final)
}

// player core

func (c *Client) AddPlayer(player model.Player) error {
	if err := c.addPlayer(player); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) UpdatePlayer(playerId uint, player model.Player) error {
	if err := c.updatePlayer(playerId, player); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) RemovePlayer(playerId uint) error {
	if err := c.removePlayer(playerId); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) GetPlayer(playerId uint) (PlayerDetail, error) {
	key := fmt.Sprintf("GetPlayer_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(PlayerDetail), nil
	}

	out, err := c.getPlayer(playerId)
	if err != nil {
		return PlayerDetail{}, err
	}
	_ = c.cache.Add(key, out, c.cacheTime)
	return out, nil
}

func (c *Client) GetPlayers(page, size int) (int64, []model.Player, error) {
	key := fmt.Sprintf("GetPlayers_%v_%v", page, size)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(int64), result[1].([]model.Player), nil
	}

	count, out, err := c.getPlayers(page, size)
	if err != nil {
		return count, out, err
	}
	_ = c.cache.Add(key, [2]any{count, out}, c.cacheTime)
	return count, out, err
}

func (c *Client) GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore) {
	key := fmt.Sprintf("GetPlayerBestScore_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project]RankScore), result[1].(map[model.Project]RankScore)
	}

	bestSingle, bestAvg = c.getPlayerBestScore(playerId)
	_ = c.cache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetPlayerPodiums(playerId uint) Podiums {
	key := fmt.Sprintf("GetPlayerPodiums_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(Podiums)
	}

	out := c.getPlayerPodiums(playerId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetPlayerRecord(playerId uint) []RecordMessage {
	key := fmt.Sprintf("GetPlayerRecord_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]RecordMessage)
	}

	out := c.getPlayerRecord(playerId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetPlayerScore(playerId uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	key := fmt.Sprintf("GetPlayerScore_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([3]any)
		return result[0].([]model.Score), result[1].([]model.Score), result[2].([]ScoresByContest)
	}

	bestSingle, bestAvg, scores = c.getPlayerScore(playerId)
	_ = c.cache.Add(key, [3]any{bestSingle, bestAvg, scores}, c.cacheTime)
	return
}

func (c *Client) GetPlayerSor(playerId uint) (single, avg map[model.SorStatisticsKey]SorScore) {
	key := fmt.Sprintf("GetPlayerSor_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.SorStatisticsKey]SorScore), result[1].(map[model.SorStatisticsKey]SorScore)
	}

	single, avg = c.getPlayerSor(playerId)
	_ = c.cache.Add(key, [2]any{single, avg}, c.cacheTime)
	return
}

func (c *Client) GetPlayerOldEnemy(playerId uint) OldEnemyDetails {
	key := fmt.Sprintf("GetPlayerOldEnemy_%v", playerId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(OldEnemyDetails)
	}

	out := c.getPlayerOldEnemy(playerId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetPlayerRelativeSor(playerID uint) map[model.SorStatisticsKey]RelativeSor {
	key := fmt.Sprintf("GetPlayerRelativeSor_%v", playerID)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey]RelativeSor)
	}

	out := c.getPlayerRelativeSor(playerID)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

// contest core

func (c *Client) AddContest(request AddContestRequest) error {
	if err := c.addContest(request); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) RemoveContest(contestId uint) error {
	if err := c.removeContest(contestId); err != nil {
		return err
	}
	c.reloadCache()
	return nil
}

func (c *Client) GetContest(contestId uint) (contest model.Contest, err error) {
	key := fmt.Sprintf("GetContest_%v", contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(model.Contest), nil
	}

	contest, err = c.getContest(contestId)
	if err != nil {
		return
	}
	_ = c.cache.Add(key, contest, c.cacheTime)
	return
}

func (c *Client) GetContests(page, size int, typ string) (int64, []model.Contest, error) {
	key := fmt.Sprintf("GetContests_%v_%v_%v", page, size, typ)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(int64), result[1].([]model.Contest), nil
	}

	count, contests, err := c.getContests(page, size, typ)
	if err != nil {
		return count, contests, err
	}
	_ = c.cache.Add(key, [2]any{count, contests}, c.cacheTime)
	return count, contests, err
}

func (c *Client) GetContestSor(contestId uint) (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := fmt.Sprintf("GetContestSor_%v", contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.SorStatisticsKey][]SorScore), result[1].(map[model.SorStatisticsKey][]SorScore)
	}

	single, avg = c.getContestSor(contestId)
	_ = c.cache.Add(key, [2]any{single, avg}, c.cacheTime)
	return
}

func (c *Client) GetContestScore(contestId uint) map[model.Project][]RoutesScores {
	key := fmt.Sprintf("GetContestScore_%v", contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.Project][]RoutesScores)
	}

	out := c.getContestScore(contestId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetContestPodiums(contestId uint) []Podiums {
	key := fmt.Sprintf("GetContestPodiums_%v", contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getContestPodiums(contestId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetContestRecord(contestId uint) []RecordMessage {
	key := fmt.Sprintf("GetContestRecord_%v", contestId)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]RecordMessage)
	}

	out := c.getContestRecord(contestId)
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

// st core

func (c *Client) GetRecords(page, size int) (int64, []model.Record, error) {
	key := fmt.Sprintf("GetRecords_%v_%v", page, size)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(int64), result[1].([]model.Record), nil
	}
	count, records, err := c.getRecords(page, size)
	if err != nil {
		return count, records, err
	}
	_ = c.cache.Add(key, [2]any{count, records}, c.cacheTime)
	return count, records, err
}

func (c *Client) GetBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := fmt.Sprintf("GetBestScore")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project][]model.Score), result[1].(map[model.Project][]model.Score)
	}
	bestSingle, bestAvg = c.getBestScore()
	_ = c.cache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	key := fmt.Sprintf("GetBestScoreByProject_%v", project)
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].([]model.Score), result[1].([]model.Score)
	}
	bestSingle, bestAvg = c.getBestScoreByProject(project)
	_ = c.cache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	key := fmt.Sprintf("GetAllProjectBestScores")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project]model.Score), result[1].(map[model.Project]model.Score)
	}
	bestSingle, bestAvg = c.getAllProjectBestScores()
	_ = c.cache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetSor() (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := fmt.Sprintf("GetSor")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.SorStatisticsKey][]SorScore), result[1].(map[model.SorStatisticsKey][]SorScore)
	}
	single, avg = c.getSor()
	_ = c.cache.Add(key, [2]any{single, avg}, c.cacheTime)
	return
}

func (c *Client) GetPodiums() []Podiums {
	key := fmt.Sprintf("GetPodiums")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getPodiums()
	_ = c.cache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetRelativeSor() (allPlayerSor map[model.SorStatisticsKey][]RelativeSor) {
	key := fmt.Sprintf("GetRelativeSor")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey][]RelativeSor)
	}
	allPlayerSor = c.getRelativeSor()
	_ = c.cache.Add(key, allPlayerSor, c.cacheTime)
	return
}

func (c *Client) GetAvgRelativeSor() map[model.SorStatisticsKey]RelativeSor {
	key := fmt.Sprintf("GetRelativeSor")
	if val, ok := c.cache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey]RelativeSor)
	}

	out := c.getAvgRelativeSor()
	_ = c.cache.Add(key, out, c.cacheTime)
	return out

}
