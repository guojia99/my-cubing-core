package utils

import (
	"context"
	"runtime"
	"sync"
)

type MultiGoFn func(ctx context.Context, idx int)

func maxGo() int {
	p := runtime.GOMAXPROCS(0)
	if p != 0 {
		return p * 4
	}
	return 4
}

// MultiGo
// 并发处理封装，由函数自行控制上锁时机
func MultiGo(ctx context.Context, fns []MultiGoFn) {
	//maxG := maxGo()
	//n := len(fns)
	//for idx := 0; idx < n; idx += maxG {
	//	for i := idx; i < idx+maxG && i < n; i++ {
	//		var wg sync.WaitGroup
	//		wg.Add(1)
	//		go func(index int, fn MultiGoFn) {
	//			defer wg.Done()
	//			fn(ctx, index)
	//		}(i, fns[i])
	//		wg.Wait()
	//	}
	//}
	var wg sync.WaitGroup
	for i := 0; i < len(fns); i++ {
		wg.Add(1)
		go func(index int, fn MultiGoFn) {
			defer wg.Done()
			fn(ctx, index)
		}(i, fns[i])
		wg.Wait()
	}

	return
}
