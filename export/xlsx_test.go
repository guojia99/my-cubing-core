package exports

import (
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	coreModel "github.com/guojia99/my-cubing-core"
)

func _testDb() *gorm.DB {
	db, _ := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN: "root:my123456@tcp(127.0.0.1:3306)/mycube2?charset=utf8&parseTime=True&loc=Local",
			},
		), &gorm.Config{
			Logger: logger.Discard,
		},
	)
	return db
}
func TestExportContestScoreXlsx(t *testing.T) {
	db := _testDb()

	core := coreModel.NewCore(db, false, time.Second)

	err := ExportContestScoreXlsx(core, 37, "test.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}
