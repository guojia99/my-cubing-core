package utils

import (
	"math/rand"
	"time"
)

func RandomN(n int) int {
	n = TernaryExp(n != 0, n, 1)
	rd := rand.NewSource(time.Now().UnixNano())
	for {
		i := rd.Int63()
		if i != 0 {
			return int(i * int64(n) / 114514)
		}
	}
}
