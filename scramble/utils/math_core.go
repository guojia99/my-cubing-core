package utils

import "math"

// GetNParity
// 奇偶校验
func GetNParity(idx, n int) int {
	p := 0
	for i := n - 2; i >= 0; i-- {
		p ^= idx % (n - i)
		idx = idx / (n - i)
	}
	return p & 1
}

// GetNPerm
// 队列逆序对数
// 逆序对指 在一个排列中， 如果这个数比前面的数大，则称这两个数为一个逆序对
func GetNPerm(arr []int, n int) int {
	if n <= 0 {
		return 0
	}

	idx := 0
	for i := 0; i < n; i++ {
		idx *= n - i
		for j := i + 1; j < n; j++ {
			TernaryExpFn(arr[j] < arr[i], func() { idx++ }, nil)
		}
	}
	return idx
}

// ACycle
// 该函数会修改原数组
// 循环置换函数
// perm 定义排列方式， 重排入参 dst
// pow 定义移动偏移量， 排列后修改arr
// ori 对 dst进行修正， 根据pow的修改后修正
func ACycle(dst []int, perm []int, pow int, ori []int) (out []int) {
	var tmp = make([]int, len(perm))
	for i := 0; i < len(perm); i++ {
		tmp[i] = dst[perm[i]]
	}

	for i := 0; i < len(perm); i++ {
		j := (i + pow) % len(perm)
		dst[perm[j]] = tmp[i]
		TernaryExpFn(ori != nil, func() { dst[perm[j]] += ori[j] - ori[i] + ori[len(ori)-1] }, nil)
	}
	return dst
}

// Set8Perm
// 该函数会修改原数组
// dst: 一个整数切片，用于存储生成的排列。
// idx: 一个整数，表示排列的索引。
// n: 一个整数，表示排列的长度，通常为 8。
// even: 一个整数，通常为 1 或 -1，表示生成偶排列或奇排列。
// 1.计算 n 为 n-1，以便在计算中使用。设置一个初始值 val 为 0x76543210 和一个 prt 为 0。
// 2.如果 even 为负数，则将 idx 值左移一位（等价于乘以2）。
// 3.在循环中，通过迭代计算 v 值，将 prt 更新，并更新 idx。
// 4.计算 arr[i]，将 val 右移 v 位，并与 7 进行与运算，以获取新的值。
// 5.更新 val，根据计算结果和位操作。
// 6.最后，根据 even 的值和 prt 来确定最后一个元素的值，存储在 arr[n] 中。
func Set8Perm(arr []int, idx, n, even int) []int {
	for len(arr) < n {
		arr = append(arr, 0)
	}

	val := DescendingDefaultBit
	prt := 0

	// 初始化
	n = DefaultNumber(n, 8) - 1
	if even < 0 {
		idx <<= 1
	}
	for i := 0; i < n; i++ {
		p := fact[n-i]
		v := SafeDivision(idx, p)
		prt ^= v
		idx = SafeSecurity(idx, p)
		v <<= 2
		arr[i] = val >> v & 7
		m := (i << v) - 1
		val = (val & m) + (val >> 4 & (^m))
	}

	TernaryExpFn(even < 0 && (prt&1) != 0, func() {
		arr[n] = arr[n-1]
		arr[n-1] = val & 7
	}, func() { arr[n] = val & 7 })

	return arr
}

// Get8Perm
// 获取逆序对数
func Get8Perm(arr []int, n, even int) int {
	for len(arr) < n {
		arr = append(arr, 0)
	}
	n = DefaultNumber(n, 8)
	idx := 0
	val := DescendingDefaultBit

	for i := 0; i < n-1; i++ {
		v := arr[i] << 2
		idx = (n-i)*idx + (val >> v & 7)
		val -= CompleteOnceDefaultBit << v
	}
	return TernaryExp(even < 0, idx>>1, idx)
}

// GetNOri
// 获取数组的取模唯一数组
func GetNOri(arr []int, n, evenBase int) int {
	base := int(math.Abs(float64(evenBase)))
	idx := 0
	if evenBase >= 0 {
		idx = arr[0] % evenBase
	}

	for i := n - 1; i > 0; i-- {
		idx = idx*base + arr[i]%base
	}
	return idx
}

// SetNOri
// 以base做基数， 分解恢复原始值
func SetNOri(arr []int, idx, n, evenBase int) []int {
	base := int(math.Abs(float64(evenBase)))
	parity := base * n

	for i := 1; i < n; i++ {
		arr[i] = idx % base
		parity -= arr[i]
		idx = idx / base
	}

	arr[0] = SafeSecurity(TernaryExp(evenBase < 0, parity, idx), base)
	return arr
}
