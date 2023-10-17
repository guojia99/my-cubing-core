package model

const (
	FinishDetailRecord  = "record"  // 记录
	FinishDetailNeglect = "neglect" // 忽略
	FinishDetailDelete  = "delete"  // 删除
)

type PreScore struct {
	Score

	ContestName string `json:"ContestName" gorm:"column:contest_name"` // 比赛名称
	RoundName   string `json:"RoundName" gorm:"column:round_name"`     // 轮次名

	Recorder     string `json:"Recorder" gorm:"column:recorder"`          // 记录人
	Processor    string `json:"Processor" gorm:"column:processor"`        // 处理人
	Finish       bool   `json:"Finish" gorm:"column:finish"`              // 是否处理
	FinishDetail string `json:"FinishDetail" gorm:"column:finish_detail"` // 处理结果
	Source       string `json:"Source" gorm:"column:source"`              // 来源 QQ web ...
}
