package exports

import (
	"fmt"

	excelize "github.com/xuri/excelize/v2"

	coreModel "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my-cubing-core/model"
	"github.com/guojia99/my-cubing-core/utils"
)

var sheetKeys = []string{
	"A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N",
}

var borders = []excelize.Border{
	{Type: "left", Color: "1F1F1F", Style: 1},
	{Type: "right", Color: "1F1F1F", Style: 1},
	{Type: "top", Color: "1F1F1F", Style: 1},
	{Type: "bottom", Color: "1F1F1F", Style: 1},
}

func setContestScoreToExcel(contestID uint, f *excelize.File, core coreModel.Core, sheet string) error {

	scores := core.GetContestScore(contestID)

	// style
	redStyle, _ := f.NewStyle(
		&excelize.Style{
			Fill:      excelize.Fill{Type: "gradient", Color: []string{"F1B8F1", "F1B8F1"}, Shading: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border:    borders,
		},
	)
	greenStyle, _ := f.NewStyle(
		&excelize.Style{
			Fill:      excelize.Fill{Type: "gradient", Color: []string{"DDFF95", "DDFF95"}, Shading: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border:    borders,
		},
	)
	centerStyle, _ := f.NewStyle(
		&excelize.Style{
			Border:    borders,
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		},
	)
	_ = f.SetColWidth(sheet, "A", "H", 6)
	_ = f.SetColWidth(sheet, "B", "B", 15)
	_ = f.SetColWidth(sheet, "E", "H", 15)
	_ = f.SetColWidth(sheet, "H", "H", 50)
	_ = f.SetColStyle(sheet, "A:H", centerStyle)

	// 标题头
	line := 1
	var header = []string{"序号", "项目", "轮次", "排名", "选手", "单次", "平均", "成绩"}
	for idx, head := range header {
		_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", sheetKeys[idx], line), head)
	}
	line += 1

	// 内容
	var recordCache = make(map[string]coreModel.RecordMessage)
	for _, val := range core.GetContestRecord(contestID) {
		recordCache[fmt.Sprintf("%d_%d", val.Score.ID, val.Record.RType)] = val
	}

	for _, pj := range model.AllProjectRoute() {
		score, ok := scores[pj]
		if !ok {
			continue
		}

		if len(score) == 0 || len(score[0].Scores) == 0 {
			continue
		}

		curLine := line
		for _, val := range score { // 分轮次
			curRound := line
			for _, s := range val.Scores { // 该轮次的成绩
				data := []interface{}{
					line - 1,                        // 序号 A
					s.Project.Cn(),                  // 项目 B
					val.Round[0].Number,             // 轮次 C
					s.Rank,                          // 排名 D
					s.PlayerName,                    // 选手 E
					utils.BestOrAvgParser(s, false), // best F
					utils.BestOrAvgParser(s, true),  // avg G
					utils.ScoresParser(s),           // 成绩 H
				}
				for idx, d := range data {
					_ = f.SetCellValue(sheet, fmt.Sprintf("%s%d", sheetKeys[idx], line), d)
				}

				// 设置单次 pb和gr
				if s.IsBestSingle {
					_ = f.SetCellStyle(sheet, fmt.Sprintf("F%d", line), fmt.Sprintf("F%d", line), greenStyle)
				}
				if _, ok = recordCache[fmt.Sprintf("%d_%d", s.ID, model.RecordBySingle)]; ok {
					_ = f.SetCellStyle(sheet, fmt.Sprintf("F%d", line), fmt.Sprintf("F%d", line), redStyle)
				}

				// 设置平均的 pb和gr
				if s.IsBestAvg {
					_ = f.SetCellStyle(sheet, fmt.Sprintf("G%d", line), fmt.Sprintf("G%d", line), greenStyle)
				}
				if _, ok = recordCache[fmt.Sprintf("%d_%d", s.ID, model.RecordByAvg)]; ok {
					_ = f.SetCellStyle(sheet, fmt.Sprintf("G%d", line), fmt.Sprintf("G%d", line), redStyle)
				}
				line += 1
			}

			_ = f.MergeCell(sheet, fmt.Sprintf("C%d", curRound), fmt.Sprintf("C%d", line-1))
		}

		_ = f.MergeCell(sheet, fmt.Sprintf("B%d", curLine), fmt.Sprintf("B%d", line-1))
	}
	return nil
}

func ExportContestScoreXlsx(core coreModel.Core, contestID uint, fileName string) error {
	f := excelize.NewFile()
	defer f.Close()

	// 查比赛
	contest, err := core.GetContest(contestID)
	if err != nil {
		return err
	}

	// 比赛成绩页面
	contestSheet, _ := f.NewSheet(contest.Name)
	f.SetActiveSheet(contestSheet)
	if err = setContestScoreToExcel(contestID, f, core, contest.Name); err != nil {
		return err
	}

	_ = f.DeleteSheet("Sheet1")
	err = f.SaveAs(fileName)
	return err
}
