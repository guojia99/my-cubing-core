package core

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func _testDb() *gorm.DB {
	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN: "root:my123456@tcp(127.0.0.1:3306)/mycube2?charset=utf8&parseTime=True&loc=Local",
	}), &gorm.Config{
		Logger: logger.Discard,
	})
	return db
}

func TestClient_getSor(t *testing.T) {
	c := NewCore(_testDb(), false, time.Minute)
	pd := c.GetPodiums()
	fmt.Println(pd)
}

func BenchmarkClient_getBestScore(b *testing.B) {
	db := _testDb()
	c := NewCore(db, false, time.Minute).(*Client)

	t := time.Now()
	c.getBestScore()
	fmt.Println(time.Now().Sub(t))
}
