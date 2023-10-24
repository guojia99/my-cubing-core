package core

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestClient_getSor(t *testing.T) {
	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN: "root:my123456@tcp(127.0.0.1:3306)/mycube2?charset=utf8&parseTime=True&loc=Local",
	}), &gorm.Config{
		Logger: logger.Discard,
	})

	c := NewCore(db, false, time.Minute)
	pd := c.GetPodiums()
	fmt.Println(pd)
}
