/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/15 下午10:09.
 *  * Author: guojia(https://github.com/guojia99)
 */

package core

import (
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my-cubing-core/model"
)

func NewCore(db *gorm.DB, debug bool, cacheTime time.Duration) *Client {
	if cacheTime == 0 {
		cacheTime = time.Minute * 15
	}
	return &Client{
		debug:     debug,
		db:        db,
		cache:     cache.New(cacheTime, cacheTime),
		cacheTime: cacheTime,
	}
}

type Core interface {
	ScoreCore
	PlayerCore
	ContestCore
	StatisticalCore
}

type ScoreCore interface {
	AddScore(AddScoreRequest) error                                               // 添加一条成绩
	RemoveScore(scoreID uint) error                                               // 移除一条成绩
	EndContestScore(contestId uint) error                                         // 结束比赛并统计比赛结果
	GetScoreByPlayerContest(playerId uint, contestId uint) ([]model.Score, error) // 获取玩家在某场比赛的成绩
}

type PlayerCore interface {
	AddPlayer(player model.Player) error                                                        // 添加玩家
	UpdatePlayer(playerId uint, player model.Player) error                                      // 更新玩家信息
	RemovePlayer(playerId uint) error                                                           // 删除玩家
	GetPlayer(playerId uint) (PlayerDetail, error)                                              // 获取玩家详细信息
	GetPlayers(page, size int) (int64, []model.Player, error)                                   // 获取玩家列表 (分页)
	GetPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore)         // 获取玩家
	GetPlayerPodiums(playerID uint) Podiums                                                     // 获取玩家领奖台
	GetPlayerRecord(playerID uint) []RecordMessage                                              // 获取玩家记录
	GetPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scores []ScoresByContest) // 获取玩家所有成绩
	GetPlayerSor(playerID uint) (single, avg map[model.SorStatisticsKey]SorScore)               // 获取玩家sor信息
	GetPlayerOldEnemy(playerID uint) OldEnemyDetails                                            // 获取玩家宿敌信息
}

type ContestCore interface {
	AddContest(AddContestRequest) error                                                 // 添加比赛
	RemoveContest(contestId uint) error                                                 // 删除比赛
	GetContest(contestId uint) (contest model.Contest, rounds []model.Round, err error) // 获取信息
	GetContests(page, size int, typ string) (int64, []model.Contest, error)             // 获取比赛列表 （分页）
	GetContestSor(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore)   // 获取比赛sor排名
	GetContestScore(contestID uint) map[model.Project][]RoutesScores                    // 获取比赛成绩列表
	GetContestPodiums(contestID uint) []Podiums                                         // 获取比赛领奖台
	GetContestRecord(contestID uint) []RecordMessage                                    // 获取比赛记录列表
}

type StatisticalCore interface {
	GetRecords(page, size int) (int64, []model.Record, error)                        // 获取所有记录（分页）
	GetBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)             // 获取所有项目每个人的最佳成绩汇总
	GetBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score) // 获取单项目每个人最佳成绩汇总成绩
	GetAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score)    // 获取所有项目最佳成绩
	GetSor() (single, avg map[model.SorStatisticsKey][]SorScore)                     // 获取sor排名汇总
	GetPodiums() []Podiums                                                           // 获取领奖台汇总
}
