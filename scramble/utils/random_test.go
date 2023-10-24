package utils

import (
	"fmt"
	"testing"
)

func TestRandomN(t *testing.T) {

	for i := 0; i < 5; i++ {
		fmt.Println(RandomN(i))
	}
}
