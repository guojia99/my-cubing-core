package core

import (
	"testing"
	"time"
)

func TestClient_resetRecords(t *testing.T) {

	c := NewCore(_testDb(), false, time.Minute)
	cc := c.(*Client)

	err := cc.resetRecords()
	if err != nil {
		t.Fatal(err)
	}
}
