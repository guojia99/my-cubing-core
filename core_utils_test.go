package core

import (
	"fmt"
	"testing"
	"time"
)

func TestClient_getLastScoresMapByContest(t *testing.T) {
	db := _testDb()

	c := NewCore(db, false, time.Second).(*Client)

	ts := time.Now()
	_ = c.getLastScoresMapByContest()
	fmt.Println(time.Now().Sub(ts))
}
