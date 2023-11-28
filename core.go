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

func NewCore(db *gorm.DB, debug bool, cacheTime time.Duration) Core {
	if cacheTime == 0 {
		cacheTime = time.Minute * 15
	}
	return &Client{
		debug:            debug,
		db:               db,
		statisticalCache: cache.New(cacheTime, cacheTime),
		cacheTime:        cacheTime,
	}
}

type Core interface {
	BackToFile() error

	ScoreCore
	PlayerCore
	PlayerUserCore
	ContestCore
	StatisticalCore
}

type ScoreCore interface {
	AddScore(AddScoreRequest) error                                               // 添加一条成绩
	RemoveScore(scoreID uint) error                                               // 移除一条成绩
	EndContestScore(contestId uint) error                                         // 结束比赛并统计比赛结果
	GetScoreByPlayerContest(playerId uint, contestId uint) ([]model.Score, error) // 获取玩家在某场比赛的成绩

	AddPreScore(AddPreScoreRequest) error                                                              // 添加一个预录入的成绩
	ProcessPreScore(ProcessPreScoreRequest) error                                                      // 处理一个预录入的成绩
	GetPreScores(page, size int, final Bool) (int64, []model.PreScore, error)                          // 获取预录入的成绩列表(分页), useFinal表示是否使用final筛选字段
	GetPreScoresByPlayer(playerID uint, page, size int, final Bool) (int64, []model.PreScore, error)   // 按玩家获取(分页)
	GetPreScoresByContest(contestID uint, page, size int, final Bool) (int64, []model.PreScore, error) // 按比赛获取(分页)
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
	GetPlayerNemesis(playerID uint) NemesisDetails                                              // 获取玩家宿敌信息
	GetPlayerRelativeSor(playerID uint) map[model.SorStatisticsKey]RelativeSor                  // 获取单个玩家的相对排位分
}

type PlayerUserCore interface {
	GetPlayerUser(player model.Player) model.PlayerUser                // 获取一个玩家的用户信息
	AddPlayerUser(player model.Player, user model.PlayerUser) error    // 添加一个用户信息
	UpdatePlayerUser(player model.Player, user model.PlayerUser) error // 更新一个用户信息
}

type ContestCore interface {
	AddContest(AddContestRequest) error                                               // 添加比赛
	RemoveContest(contestId uint) error                                               // 删除比赛
	GetAllContestStatics() (contests []ContestStatics)                                // 所有比赛的数据统计
	GetContest(contestId uint) (contest model.Contest, err error)                     // 获取信息
	GetContests(page, size int, typ string) (int64, []model.Contest, error)           // 获取比赛列表 （分页）
	GetContestSor(contestID uint) (single, avg map[model.SorStatisticsKey][]SorScore) // 获取比赛sor排名
	GetContestScore(contestID uint) map[model.Project][]RoutesScores                  // 获取比赛成绩列表
	GetContestPodiums(contestID uint) []Podiums                                       // 获取比赛领奖台
	GetContestRecord(contestID uint) []RecordMessage                                  // 获取比赛记录列表
}

type StatisticalCore interface {
	GetRecords(page, size int) (int64, []model.Record, error)                                                        // 获取所有记录（分页）
	GetBestScore() (bestSingle, bestAvg map[model.Project][]model.Score)                                             // 获取所有项目每个人的最佳成绩汇总
	GetBestScoreByTimes(startTime, endTime time.Time) (bestSingle, bestAvg map[model.Project][]model.Score)          // [Time]获取所有项目每个人的最佳成绩汇总
	GetBestScoreByProject(project model.Project) (bestSingle, bestAvg []model.Score)                                 // 获取单项目每个人最佳成绩汇总成绩
	GetAllProjectBestScores() (bestSingle, bestAvg map[model.Project]model.Score)                                    // 获取所有项目最佳成绩
	GetAllProjectBestScoresByTimes(startTime, endTime time.Time) (bestSingle, bestAvg map[model.Project]model.Score) // [Time]获取所有项目最佳成绩
	GetPodiums() []Podiums                                                                                           // 获取领奖台汇总
	GetSor() (single, avg map[model.SorStatisticsKey][]SorScore)                                                     // 获取sor排名汇总
	GetAvgRelativeSor() map[model.SorStatisticsKey]RelativeSor                                                       // 平均相对排位分
	GetRelativeSor() (allPlayerSor map[model.SorStatisticsKey][]RelativeSor)                                         // 相对排位分, 返回所有人的平均排位分,计算方式是用当前最佳成绩为标准, 计算与其差距
}
