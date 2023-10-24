package utils

const (
	DescendingDefaultBit   = 0x76543210
	CompleteOnceDefaultBit = 0x11111110
)

func TernaryExp(expr bool, val1, val2 int) int {
	if expr {
		return val1
	}
	return val2
}

func TernaryExpFn(expr bool, f1, f2 func()) {
	var run func() = nil
	if expr {
		run = f1
	} else {
		run = f2
	}
	if run != nil {
		run()
	}
}

// fact
// 阶乘
// 这是int 类型中最多能容纳的数量 21个
var fact = func() []int {
	var out = make([]int, 21)
	out[0] = 1
	for i := 0; i < 20; i++ {
		out[i+1] = out[i] * (i + 1)
	}
	return out
}()

func SafeDivision(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func SafeSecurity(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

func DefaultNumber(orgNumber, defaultNum int) int {
	if orgNumber == 0 {
		return defaultNum
	}
	return orgNumber
}
