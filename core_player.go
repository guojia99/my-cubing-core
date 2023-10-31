package core

import (
	"errors"
	"regexp"
	"sort"

	"github.com/guojia99/my-cubing-core/model"
)

func checkName(name string) bool {
	pattern := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]{1,15}$`)
	return pattern.MatchString(name)
}

// 添加玩家
func (c *Client) addPlayer(player model.Player) error {
	if !checkName(player.Name) {
		return errors.New("invalid name")
	}

	if err := c.db.Create(&player).Error; err != nil {
		return err
	}
	return nil
}

// 更新玩家信息
func (c *Client) updatePlayer(playerId uint, player model.Player) error {
	var p model.Player
	if err := c.db.Where("id = ?", playerId).First(&p).Error; err != nil {
		return err
	}

	if !checkName(player.Name) {
		return errors.New("invalid name")
	}

	p.Name = player.Name
	p.ActualName = player.ActualName
	p.WcaID = player.WcaID
	p.SetTitles(player.TitlesVal)

	if err := c.db.Save(&p).Error; err != nil {
		return err
	}
	// 更新玩家所有成绩的name
	if err := c.db.Model(&model.Score{}).
		Where("player_id = ?", player.ID).
		UpdateColumn("player_name", p.Name).Error; err != nil {
		return err
	}
	return nil
}

// 删除玩家
func (c *Client) removePlayer(playerId uint) error {
	var player model.Player
	if err := c.db.Where("id = ?", playerId).First(&player).Error; err != nil {
		return err
	}

	var count int64
	if err := c.db.Model(&model.Score{}).Where("player_id = ?", playerId).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("can't not delete has score player")
	}

	err := c.db.Delete(&player).Error
	return err
}

// 获取玩家详细信息
func (c *Client) getPlayer(playerId uint) (PlayerDetail, error) {
	var player model.Player
	if err := c.db.First(&player, "id = ?", playerId).Error; err != nil {
		return PlayerDetail{}, err
	}

	var contestIDs []uint64
	c.db.
		Model(&model.Score{}).
		Distinct("contest_id").
		Where("player_id = ?", playerId).
		Pluck("contest_id", &contestIDs)

	out := PlayerDetail{
		Player:        player,
		ContestNumber: len(contestIDs),
	}

	var score []model.Score
	c.db.Model(&model.Score{}).Find(&score, "player_id = ?", playerId)
	for _, s := range score {
		if s.Project.RouteType() == model.RouteTypeRepeatedly {
			out.RecoveryNumber += 1
			if s.DBest() {
				out.ValidRecoveryNumber += 1
			}
			continue
		}

		rs := s.GetResult()
		out.RecoveryNumber += len(rs)
		for _, val := range rs {
			if val <= model.DNF {
				out.ValidRecoveryNumber += 1
			}
		}
	}
	return out, nil
}

// 获取玩家列表 (分页)
func (c *Client) getPlayers(page, size int) (int64, []model.Player, error) {
	if page == 0 {
		page = 1
	}
	if size == 0 || size > 100 {
		size = 100
	}

	offset := (page - 1) * size
	limit := size
	var out []model.Player
	if err := c.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&out).Error; err != nil {
		return 0, nil, err
	}

	var count int64
	if err := c.db.Model(&model.Player{}).Count(&count).Error; err != nil {
		return 0, nil, err
	}

	return count, out, nil
}

// 获取玩家全项目最佳成绩
func (c *Client) getPlayerBestScore(playerId uint) (bestSingle, bestAvg map[model.Project]RankScore) {
	bestSingle, bestAvg = make(map[model.Project]RankScore), make(map[model.Project]RankScore)
	allBest, allAvg := c.GetBestScore()

	// todo 双层map
	for _, val := range allBest {
		for i := 0; i < len(val); i++ {
			if val[i].PlayerID == playerId {
				bestSingle[val[i].Project] = RankScore{
					Rank:  i + 1,
					Score: val[i],
				}
				break
			}
		}
	}
	for _, val := range allAvg {
		for i := 0; i < len(val); i++ {
			if val[i].PlayerID == playerId {
				bestAvg[val[i].Project] = RankScore{
					Rank:  i + 1,
					Score: val[i],
				}
				break
			}
		}
	}
	return bestSingle, bestAvg
}

// 获取玩家领奖台
func (c *Client) getPlayerPodiums(playerID uint) Podiums {
	var player model.Player
	if err := c.db.Where("id = ?", playerID).First(&player).Error; err != nil {
		return Podiums{}
	}

	var out = Podiums{Player: player}

	// 查选手参加过的所有比赛且结束的
	var cacheContestId []uint
	c.db.
		Model(&model.Score{}).
		Distinct("contest_id").
		Where("player_id = ?", playerID).
		Pluck("player_id", &cacheContestId)
	if len(cacheContestId) == 0 {
		return out
	}
	var contests []model.Contest
	c.db.Where("is_end = ?", 1).Find(&contests)

	// 查选手所有比赛的成绩
	for _, contest := range contests {
		topThree := c.getContestTop(contest.ID, 3)
		for _, pj := range model.AllProjectRoute() {
			score, ok := topThree[pj]
			if !ok {
				continue
			}
			// todo 名次相等
			for idx, val := range score {
				if val.PlayerID == playerID {
					switch val.Rank {
					case 1:
						out.Gold += 1
					case 2:
						out.Silver += 1
					case 3:
						out.Bronze += 1
					}
					out.PodiumsResults = append(out.PodiumsResults, PodiumsResult{
						Contest: contest,
						Score:   score[idx],
					})
				}
			}
		}
	}
	return out
}

// 获取玩家记录
func (c *Client) getPlayerRecord(playerID uint) []RecordMessage {
	var out []RecordMessage

	var player model.Player
	if err := c.db.Find(&player, "id = ?", playerID).Error; err != nil {
		return out
	}

	var records []model.Record
	if err := c.db.Where("player_id = ?", playerID).Find(&records).Error; err != nil {
		return out
	}

	for _, record := range records {
		var contest model.Contest
		var score model.Score
		_ = c.db.First(&contest, "id = ?", record.ContestID).Error
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

// 获取玩家所有成绩
func (c *Client) getPlayerScore(playerID uint) (bestSingle, bestAvg []model.Score, scoresByContest []ScoresByContest) {
	var scores []model.Score
	c.db.Where("player_id = ?", playerID).Find(&scores)
	if len(scores) == 0 {
		return
	}

	var (
		cache     = make(map[uint][]model.Score)
		avgCache  = make(map[model.Project]model.Score)
		bestCache = make(map[model.Project]model.Score)
	)
	for _, score := range scores {
		if _, ok := cache[score.ContestID]; !ok {
			cache[score.ContestID] = make([]model.Score, 0)
		}
		cache[score.ContestID] = append(cache[score.ContestID], score)

		if got, ok := bestCache[score.Project]; !ok || score.IsBestScore(got) {
			bestCache[score.Project] = score
		}
		if got, ok := avgCache[score.Project]; !ok || score.IsBestAvgScore(got) {
			avgCache[score.Project] = score
		}
	}

	for key, val := range cache {
		var contest model.Contest
		if err := c.db.
			Where("id = ?", key).
			Where("is_end = ?", 1).
			First(&contest).Error; err != nil {
			continue
		}
		var rounds []model.Round
		c.db.Find(&rounds, "contest_id = ?", contest.ID)

		sort.Slice(val, func(i, j int) bool { return val[i].ID > val[j].ID })
		scoresByContest = append(scoresByContest, ScoresByContest{
			Contest: contest,
			Rounds:  rounds,
			Scores:  val,
		})
	}

	for _, val := range avgCache {
		bestAvg = append(bestAvg, val)
	}
	for _, val := range bestCache {
		bestSingle = append(bestSingle, val)
	}

	// 给所有成绩排序
	sort.Slice(bestSingle, func(i, j int) bool { return bestSingle[i].ID > bestSingle[j].ID })
	sort.Slice(bestAvg, func(i, j int) bool { return bestAvg[i].ID > bestAvg[j].ID })
	sort.Slice(scoresByContest, func(i, j int) bool { return scoresByContest[i].Contest.ID > scoresByContest[j].Contest.ID })
	return
}

// 获取玩家sor信息
func (c *Client) getPlayerSor(playerID uint) (single, avg map[model.SorStatisticsKey]SorScore) {
	single, avg = make(map[model.SorStatisticsKey]SorScore, len(model.SorKeyMap())), make(map[model.SorStatisticsKey]SorScore, len(model.SorKeyMap()))
	singleCache, avgCache := c.getSor()

	for k, _ := range model.SorKeyMap() {
		if _, ok := singleCache[k]; ok {
			for idx, score := range singleCache[k] {
				if score.Player.ID == playerID {
					score.SingleRank = int64(idx + 1)
					single[k] = score
					break
				}
			}
		}
		if _, ok := avgCache[k]; ok {
			for idx, score := range avgCache[k] {
				if score.Player.ID == playerID {
					score.AvgRank = int64(idx + 1)
					avg[k] = score
					break
				}
			}
		}
	}
	return single, avg
}

// 获取宿敌信息
func (c *Client) getPlayerNemesis(playerId uint) NemesisDetails {
	var out = make(NemesisDetails, 0)

	var player model.Player
	if err := c.db.First(&player, "id = ?", playerId).Error; err != nil {
		return nil
	}
	var allPlayer []model.Player
	if err := c.db.Where("id != ?", playerId).Find(&allPlayer).Error; err != nil {
		return nil
	}

	allBest, allAvg := c.getBestScore()

	// todo 抽出来
	// 1. 缓存所有人的成绩
	var (
		bestCache = make(map[uint]map[model.Project]model.Score, len(allPlayer))
		avgCache  = make(map[uint]map[model.Project]model.Score, len(allPlayer))
	)
	for _, p := range allPlayer {
		bestCache[p.ID] = make(map[model.Project]model.Score)
		avgCache[p.ID] = make(map[model.Project]model.Score)
	}

	for _, pj := range model.AllProjectRoute() {
		for _, b := range allBest[pj] {
			if b.PlayerID != playerId {
				bestCache[b.PlayerID][pj] = b
			}
		}
		for _, a := range allAvg[pj] {
			if a.PlayerID != playerId {
				avgCache[a.PlayerID][pj] = a
			}
		}
	}

	// 2. 对比自己的成绩， 只要有一个成绩比其他人好，这个人就不是宿敌, 对比时， 如果自己没有这个项目， 则算输
	selfBest, selfAvg := c.GetPlayerBestScore(playerId)

loop:
	for _, p := range allPlayer {
		otherBest, otherAvg := bestCache[p.ID], avgCache[p.ID]

		for _, pj := range model.AllProjectRoute() {
			// 1. 如果都没有这个项目， 不需要对比
			selfB, ok1 := selfBest[pj]
			otherB, ok2 := otherBest[pj]
			if !ok1 && !ok2 {
				continue
			}
			// 2. 如果自己有这个项目， 他没有， 直接退出循环
			if ok1 && !ok2 {
				continue loop
			}
			// 3. 如果他有这个项目，自己没有
			if !ok1 && ok2 {
				continue
			}
			// 4. 对比自己的单次是否比他好
			if selfB.Rank <= otherB.Rank {
				continue loop
			}

			// 5. 如果单次比他差，查看是自己和他否有平均
			selfA, ok3 := selfAvg[pj]
			otherA, ok4 := otherAvg[pj]

			// 6. 都没有平均的时候
			if !ok3 && !ok4 {
				continue
			}
			// 7. 如果自己有平均， 他没有
			if ok3 && !ok4 {
				continue loop
			}
			// 8. 如果他有平均， 自己没有
			if !ok3 && ok4 {
				continue
			}
			// 9.对比平均排名
			if selfA.Rank <= otherA.Rank {
				continue loop
			}
		}

		out = append(out, NemesisDetail{
			Player: p,
			Single: otherBest,
			Avg:    otherAvg,
		})
	}
	return out
}

func (c *Client) getPlayerRelativeSor(playerID uint) map[model.SorStatisticsKey]RelativeSor {
	var out = make(map[model.SorStatisticsKey]RelativeSor)
	var player model.Player
	if err := c.db.Where("id = ?", playerID).First(&player).Error; err != nil {
		return out
	}
	bestSingle, bestAvg := c.getPlayerBestScore(playerID)

	var bs = make(map[model.Project][]model.Score)
	var ba = make(map[model.Project][]model.Score)
	for _, pj := range model.AllProjectRoute() {
		if s, ok := bestSingle[pj]; ok {
			bs[pj] = []model.Score{s.Score}
		}
		if a, ok := bestAvg[pj]; ok {
			ba[pj] = []model.Score{a.Score}
		}
	}

	sor := c.parserRelativeSor([]model.Player{player}, bs, ba)

	for k, _ := range model.SorKeyMap() {
		if _, ok := sor[k]; !ok || len(sor[k]) == 0 {
			continue
		}
		out[k] = sor[k][0]
	}
	return out
}

func (c *Client) getPlayerUser(player model.Player) model.PlayerUser {
	var out model.PlayerUser
	c.db.First(&out, "player_id = ?", player.ID)
	return out
}

func (c *Client) addPlayerUser(player model.Player, user model.PlayerUser) error {
	return c.UpdatePlayerUser(player, user)
}

func (c *Client) updatePlayerUser(player model.Player, user model.PlayerUser) error {
	if err := c.db.First(&player, "id = ?", player.ID).Error; err != nil {
		return errors.New("玩家不存在")
	}

	user.PlayerID = player.ID
	if !user.Valid() {
		return errors.New("校验错误")
	}

	var err error
	if playerUser := c.getPlayerUser(player); playerUser.ID == 0 { // 创建
		err = c.db.Create(&user).Error
	} else { //更新
		playerUser.QQ = user.QQ
		playerUser.WeChat = user.WeChat
		playerUser.Phone = user.Phone
		playerUser.LoginID = user.LoginID
		err = c.db.Save(&playerUser).Error
	}
	return err
}
