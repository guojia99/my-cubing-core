package scramble

type DetailScrambleFn func(size int) string

type Scramble interface {
	Random() string // 打乱
	//Solve(seq string, n int) ([]string, error) // 求解， n为求解数, 最大10
}

const (
	U1 = iota
	U2
	U3
	R1
	R2
	R3
	F1
	F2
	F3
	D1
	D2
	D3
	L1
	L2
	L3
	B1
	B2
	B3
)
