/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/15 下午10:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing-core/model"
)

var _ Core = &Client{}

type Client struct {
	debug bool
	db    *gorm.DB

	statisticalCache *cache.Cache
	cacheTime        time.Duration
}

func (c *Client) BackToFile() error {

	return nil
}

func (c *Client) ReSetRecords() error { return c.resetRecords() }

// score Core
func (c *Client) AddScore(req AddScoreRequest) error {
	return c.addScore(req.PlayerID, req.ContestID, req.Project, req.RoundId, req.Result, req.Penalty)
}
func (c *Client) RemoveScore(scoreID uint) error       { return c.removeScore(scoreID) }
func (c *Client) EndContestScore(contestId uint) error { return c.endContestScore(contestId) }
func (c *Client) GetScoreByPlayerContest(playerId uint, contestId uint) ([]model.Score, error) {
	return c.getScoreByPlayerContest(playerId, contestId)
}

func (c *Client) AddPreScore(request AddPreScoreRequest) error {
	return c.addPreScore(request)
}

func (c *Client) ProcessPreScore(request ProcessPreScoreRequest) error {
	return c.processPreScore(request)
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
	return c.addPlayer(player)
}

func (c *Client) UpdatePlayer(playerId uint, player model.Player) error {
	return c.updatePlayer(playerId, player)
}

func (c *Client) RemovePlayer(playerId uint) error {
	return c.removePlayer(playerId)
}

func (c *Client) GetPlayer(playerId uint) (PlayerDetail, error) {
	return c.getPlayer(playerId)
}

func (c *Client) GetPlayers(page, size int) (int64, []model.Player, error) {
	return c.getPlayers(page, size)
}

func (c *Client) GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore) {
	return c.getPlayerBestScore(playerId)
}

func (c *Client) GetPlayerPodiums(playerId uint) Podiums {
	_, pds := c.getPodiums()
	return *pds[playerId]
}

func (c *Client) GetPlayerRecord(playerId uint) []RecordMessage {
	return c.getPlayerRecord(playerId)
}

func (c *Client) GetPlayerScore(playerId uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) {
	return c.getPlayerScore(playerId)
}

func (c *Client) GetPlayerSor(playerId uint) (single, avg map[model.SorStatisticsKey]SorScore) {
	return c.getPlayerSor(playerId)
}

func (c *Client) GetPlayerNemesis(playerID uint) NemesisDetails {
	key := fmt.Sprintf("GetPlayerNemesis_%v", playerID)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.(NemesisDetails)
	}

	out := c.getPlayerNemesis(playerID)
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetPlayerRelativeSor(playerID uint) map[model.SorStatisticsKey]RelativeSor {
	key := fmt.Sprintf("GetPlayerRelativeSor_%v", playerID)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey]RelativeSor)
	}

	out := c.getPlayerRelativeSor(playerID)
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

// contest core

func (c *Client) AddContest(request AddContestRequest) error { return c.addContest(request) }
func (c *Client) RemoveContest(contestId uint) error         { return c.removeContest(contestId) }
func (c *Client) GetContest(contestId uint) (contest model.Contest, err error) {
	return c.getContest(contestId)
}
func (c *Client) GetContests(page, size int, typ string) (int64, []model.Contest, error) {
	return c.getContests(page, size, typ)
}

func (c *Client) GetContestSor(contestId uint) (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := fmt.Sprintf("GetContestSor_%v", contestId)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.SorStatisticsKey][]SorScore), result[1].(map[model.SorStatisticsKey][]SorScore)
	}

	single, avg = c.getContestSor(contestId)
	_ = c.statisticalCache.Add(key, [2]any{single, avg}, c.cacheTime)
	return
}

func (c *Client) GetContestScore(contestId uint) map[model.Project][]RoutesScores {
	key := fmt.Sprintf("GetContestScore_%v", contestId)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.(map[model.Project][]RoutesScores)
	}

	out := c.getContestScore(contestId)
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetContestPodiums(contestId uint) []Podiums {
	key := fmt.Sprintf("GetContestPodiums_%v", contestId)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out := c.getContestPodiums(contestId)
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetContestRecord(contestId uint) []RecordMessage {
	key := fmt.Sprintf("GetContestRecord_%v", contestId)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.([]RecordMessage)
	}

	out := c.getContestRecord(contestId)
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetAllContestStatics() (contests []ContestStatics) {
	key := fmt.Sprintf("GetAllContestStatics")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.([]ContestStatics)
	}

	out := c.getAllContestStatics()
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

// st core

func (c *Client) GetRecords(page, size int) (int64, []model.Record, error) {
	key := fmt.Sprintf("GetRecords_%v_%v", page, size)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(int64), result[1].([]model.Record), nil
	}
	count, records, err := c.getRecords(page, size)
	if err != nil {
		return count, records, err
	}
	_ = c.statisticalCache.Add(key, [2]any{count, records}, c.cacheTime)
	return count, records, err
}

func (c *Client) GetBestScore() (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := fmt.Sprintf("GetBestScore")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project][]model.Score), result[1].(map[model.Project][]model.Score)
	}
	bestSingle, bestAvg = c.getBestScore()
	_ = c.statisticalCache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}
func (c *Client) GetBestScoreByTimes(startTime, endTime time.Time) (bestSingle, bestAvg map[model.Project][]model.Score) {
	key := fmt.Sprintf("GetBestScoreByTimes_%v_%v", startTime, endTime)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project][]model.Score), result[1].(map[model.Project][]model.Score)
	}
	bestSingle, bestAvg = c.getBestScoreWithTime(startTime, endTime)
	_ = c.statisticalCache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) {
	key := fmt.Sprintf("GetBestScoreByProject_%v", project)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].([]model.Score), result[1].([]model.Score)
	}
	bestSingle, bestAvg = c.getBestScoreByProject(project)
	_ = c.statisticalCache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score) {
	key := fmt.Sprintf("GetAllProjectBestScores")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project]model.Score), result[1].(map[model.Project]model.Score)
	}
	bestSingle, bestAvg = c.getAllProjectBestScores()
	_ = c.statisticalCache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}
func (c *Client) GetAllProjectBestScoresByTimes(startTime, endTime time.Time) (bestSingle, bestAvg map[model.Project]model.Score) {
	key := fmt.Sprintf("GetAllProjectBestScoresByTimes_%v_%v", startTime, endTime)
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.Project]model.Score), result[1].(map[model.Project]model.Score)
	}
	bestSingle, bestAvg = c.getAllProjectBestScoresWithTime(startTime, endTime)
	_ = c.statisticalCache.Add(key, [2]any{bestSingle, bestAvg}, c.cacheTime)
	return
}

func (c *Client) GetSor() (single, avg map[model.SorStatisticsKey][]SorScore) {
	key := fmt.Sprintf("GetSor")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		result := val.([2]any)
		return result[0].(map[model.SorStatisticsKey][]SorScore), result[1].(map[model.SorStatisticsKey][]SorScore)
	}
	single, avg = c.getSor()
	_ = c.statisticalCache.Add(key, [2]any{single, avg}, c.cacheTime)
	return
}

func (c *Client) GetPodiums() []Podiums {
	key := fmt.Sprintf("GetPodiums")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.([]Podiums)
	}

	out, _ := c.getPodiums()
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out
}

func (c *Client) GetRelativeSor() (allPlayerSor map[model.SorStatisticsKey][]RelativeSor) {
	key := fmt.Sprintf("GetRelativeSor")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey][]RelativeSor)
	}
	allPlayerSor = c.getRelativeSor()
	_ = c.statisticalCache.Add(key, allPlayerSor, c.cacheTime)
	return
}

func (c *Client) GetAvgRelativeSor() map[model.SorStatisticsKey]RelativeSor {
	key := fmt.Sprintf("GetAvgRelativeSor")
	if val, ok := c.statisticalCache.Get(key); ok && !c.debug {
		return val.(map[model.SorStatisticsKey]RelativeSor)
	}

	out := c.getAvgRelativeSor()
	_ = c.statisticalCache.Add(key, out, c.cacheTime)
	return out

}

// player user

func (c *Client) GetPlayerUser(player model.Player) model.PlayerUser {
	return c.getPlayerUser(player)
}

func (c *Client) AddPlayerUser(player model.Player, user model.PlayerUser) error {
	return c.addPlayerUser(player, user)
}

func (c *Client) UpdatePlayerUser(player model.Player, user model.PlayerUser) error {
	return c.updatePlayerUser(player, user)
}
