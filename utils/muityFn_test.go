package utils

import (
	"context"
	"fmt"
	"testing"
)

func TestMultiGo(t *testing.T) {

	var data = make([]string, 100)

	var fn []MultiGoFn
	for i := 0; i < 100; i++ {
		fn = append(fn, func(ctx context.Context, idx int) {
			//lock.Lock()
			data[idx] = fmt.Sprintf("%d", idx)
			//lock.Unlock()
		})
	}

	MultiGo(context.Background(), fn)

	fmt.Println(len(data), data)
}
