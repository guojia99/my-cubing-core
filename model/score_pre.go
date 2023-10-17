package model

const (
	FinishDetailRecord  = "record"  // 记录
	FinishDetailNeglect = "neglect" // 忽略
	FinishDetailDelete  = "delete"  // 删除
)

type PreScore struct {
	Score

	ContestName  string `json:"contestName" gorm:"column:contest_name"`   // 比赛名称
	Recorder     string `json:"recorder" gorm:"column:recorder"`          // 记录人
	Processor    string `json:"processor" gorm:"column:processor"`        // 处理人
	Finish       bool   `json:"finish" gorm:"column:finish"`              // 是否处理
	FinishDetail string `json:"finishDetail" gorm:"column:finish_detail"` // 处理结果
	Source       string `json:"source" gorm:"column:source"`              // 来源 QQ web ...
}
