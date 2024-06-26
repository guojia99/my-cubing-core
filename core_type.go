package core

import (
	"sort"

	"github.com/guojia99/my-cubing-core/model"
)

type Bool uint

const (
	NotBool   Bool = 0
	FalseBool Bool = 1
	TrueBool  Bool = 2
)

type (
	AddScoreRequest struct {
		PlayerID  uint               `json:"PlayerID"`
		ContestID uint               `json:"ContestID"`
		Project   model.Project      `json:"Project"`
		RoundId   uint               `json:"RoundId"`
		Result    []float64          `json:"Result"`
		Penalty   model.ScorePenalty `json:"Penalty"`
	}

	CreateContestRequestRound struct {
		Project model.Project `json:"Project"`
		Number  int           `json:"Number"`
		Part    int           `json:"Part"`
		Name    string        `json:"Name"`
		IsStart bool          `json:"IsStart"`
		Final   bool          `json:"Final"`
		Upsets  []string      `json:"Upsets"`
	}

	AddContestRequest struct {
		Name        string                      `json:"Name"`
		Description string                      `json:"Description"`
		GroupID     string                      `json:"GroupID"`
		Rounds      []CreateContestRequestRound `json:"Rounds"`
		Type        string                      `json:"Type"`
		StartTime   int64                       `json:"StartTime"`
		EndTime     int64                       `json:"EndTime"`
	}

	AddPreScoreRequest struct {
		AddScoreRequest
		Source   string `json:"Source"`
		Recorder string `json:"Recorder"`
	}

	ProcessPreScoreRequest struct {
		Id           uint
		Processor    string
		FinishDetail string
	}
)

type RankScore struct {
	Rank  int         `json:"Rank,omitempty"` // 排名
	Score model.Score `json:"Score,omitempty"`
}

type RoutesScores struct {
	Final  bool          `json:"final,omitempty"`
	Round  []model.Round `json:"Round,omitempty"`
	Scores []model.Score `json:"Scores,omitempty"`
}

type ScoresByContest struct {
	Contest model.Contest `json:"Contest,omitempty"`
	Scores  []model.Score `json:"Scores,omitempty"`
	Rounds  []model.Round `json:"Rounds,omitempty"`
}

type Podiums struct {
	Player         model.Player    `json:"Player,omitempty"`
	Gold           int64           `json:"Gold,omitempty"`
	Silver         int64           `json:"Silver,omitempty"`
	Bronze         int64           `json:"Bronze,omitempty"`
	PodiumsResults []PodiumsResult `json:"PodiumsResults,omitempty"`
}

func (p *Podiums) Add(rank int) *Podiums {
	switch rank {
	case 1:
		p.Gold += 1
	case 2:
		p.Silver += 1
	case 3:
		p.Bronze += 1
	}
	return p
}

type PodiumsResult struct {
	Contest model.Contest `json:"Contest,omitempty"`
	Score   model.Score   `json:"Score,omitempty"`
}

type SorScore struct {
	Player         model.Player `json:"Player,omitempty"`
	SingleRank     int64        `json:"SingleRank,omitempty"`
	SingleCount    int64        `json:"SingleCount,omitempty"`
	SingleProjects int64        `json:"SingleProjects,omitempty"` // 参与项目数
	AvgRank        int64        `json:"AvgRank,omitempty"`
	AvgCount       int64        `json:"AvgCount,omitempty"`
	AvgProjects    int64        `json:"AvgProjects,omitempty"` // 参与项目数
}

func SortPodiums(in []Podiums) {
	sort.Slice(
		in, func(i, j int) bool {
			if in[i].Gold == in[j].Gold {
				if in[i].Silver == in[j].Silver {
					return in[i].Bronze > in[j].Bronze
				}
				return in[i].Silver > in[j].Silver
			}
			return in[i].Gold > in[j].Gold
		},
	)
}

type RecordMessage struct {
	Record  model.Record  `json:"Record,omitempty"`
	Player  model.Player  `json:"Player,omitempty"`
	Score   model.Score   `json:"Score,omitempty"`
	Contest model.Contest `json:"Contest,omitempty"`
}

type PlayerDetail struct {
	model.Player

	ContestNumber       int `json:"ContestNumber,omitempty"`
	ValidRecoveryNumber int `json:"ValidRecoveryNumber,omitempty"`
	RecoveryNumber      int `json:"RecoveryNumber,omitempty"`
}

type NemesisDetail struct {
	Player model.Player                  `json:"Player,omitempty"`
	Single map[model.Project]model.Score `json:"Single,omitempty"`
	Avg    map[model.Project]model.Score `json:"Avg,omitempty"`
}

type NemesisDetails []NemesisDetail

type RelativeSor struct {
	Player model.Player `json:"Player,omitempty"`
	// 个人成绩
	Sor float64 `json:"Sor,omitempty"`

	// 平均成绩会用到
	Avg  float64 `json:"Avg,omitempty"`  // 全部平均
	Top5 float64 `json:"Top5,omitempty"` // 前5平均, 大于100的取前5%
	Max  float64 `json:"Max,omitempty"`  // 最大值
}

type ContestStatics struct {
	model.Contest

	PlayerNum  int `json:"PlayerNum"`
	ProjectNum int `json:"ProjectNum"`
}
