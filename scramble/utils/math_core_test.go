package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNParity(t *testing.T) {
	type args struct {
		idx int
		n   int
	}
	tests := []struct {
		args args
		want int
	}{
		{
			args: args{
				idx: 4,
				n:   6,
			},
			want: 0,
		},
		{
			args: args{
				idx: 11,
				n:   3,
			},
			want: 1,
		},
		{
			args: args{
				idx: 11,
				n:   23,
			},
			want: 0,
		},
		{
			args: args{
				idx: 1,
				n:   2,
			},
			want: 1,
		},
		{
			args: args{
				idx: 2,
				n:   1,
			},
			want: 0,
		},
		{
			args: args{
				idx: 3,
				n:   2,
			},
			want: 1,
		}, {
			args: args{
				idx: 3,
				n:   3,
			},
			want: 0,
		},
		{
			args: args{
				idx: 5,
				n:   3,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%d", tt.args.idx, tt.args.n), func(t *testing.T) {
			if got := GetNParity(tt.args.idx, tt.args.n); got != tt.want {
				t.Errorf("GetNParity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNPerm(t *testing.T) {
	type args struct {
		arr []int
		n   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				arr: []int{2, 2, 8, 4, 5, 6},
				n:   5,
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				arr: []int{1, 2, 6, 2, 3, 4, 5, 6},
				n:   6,
			},
			want: 18,
		},
		{
			name: "3",
			args: args{
				arr: []int{1, 2, 6, 11, 2, 3, 1, 4, 5, 6},
				n:   6,
			},
			want: 16,
		},
		{
			name: "4",
			args: args{
				arr: []int{1, 2, 3, 2, 6, 11, 2, 3, 1, 4, 5, 6},
				n:   6,
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNPerm(tt.args.arr, tt.args.n); got != tt.want {
				t.Errorf("NPerm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestACycle(t *testing.T) {
	type args struct {
		dst  []int
		perm []int
		pow  int
		ori  []int
	}
	tests := []struct {
		name    string
		args    args
		wantOut []int
	}{
		{
			name: "1",
			args: args{
				dst:  []int{1, 2, 3, 4},
				perm: []int{2, 0, 3, 1},
				pow:  1,
				ori:  []int{10, 20, 30, 5, 0},
			},
			wantOut: []int{13, -21, 7, 11},
		},
		{
			name: "2",
			args: args{
				dst:  []int{2, 2, 3, 4},
				perm: []int{2, 0, 3, 1},
				pow:  1,
				ori:  []int{10, 20, 30, 5, 0},
			},
			wantOut: []int{13, -21, 7, 12},
		},
		{
			name: "3",
			args: args{
				dst:  []int{2, 4, 3, 4},
				perm: []int{2, 0, 3, 1},
				pow:  1,
				ori:  []int{10, 20, 30, 5, 0},
			},
			wantOut: []int{13, -21, 9, 12},
		},
		{
			name: "4",
			args: args{
				dst:  []int{2, 4, 3, 4},
				perm: []int{2, 1, 3, 2},
				pow:  1,
				ori:  []int{20, 40, 30, 5, 0},
			},
			wantOut: []int{2, 23, 18, -6},
		},
		{
			name: "4",
			args: args{
				dst:  []int{1, 2, 3, 4},
				perm: []int{2, 0, 3, 1},
				pow:  3,
				ori:  []int{10, 20, 30, 5, 0},
			},
			wantOut: []int{-6, -2, -9, 27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ACycle(tt.args.dst, tt.args.perm, tt.args.pow, tt.args.ori); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ACycle() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestSet8Perm(t *testing.T) {
	type args struct {
		dst  []int
		idx  int
		n    int
		even int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				dst:  []int{1, 2, 3, 4},
				idx:  1,
				n:    3,
				even: 2,
			},
			want: []int{0, 1, 0, 4},
		},
		{
			name: "2",
			args: args{
				dst:  []int{5, 6, 1, 3, 1, 2, 3, 4},
				idx:  2,
				n:    8,
				even: 1,
			},
			want: []int{0, 0, 1, 3, 3, 5, 7, 7},
		},
		{
			name: "3",
			args: args{
				dst:  []int{},
				idx:  1,
				n:    4,
				even: -1,
			},
			want: []int{0, 1, 2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Set8Perm(tt.args.dst, tt.args.idx, tt.args.n, tt.args.even); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set8Perm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet8Perm(t *testing.T) {
	type args struct {
		arr  []int
		n    int
		even int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				arr:  []int{4, 3, 2, 1},
				n:    4,
				even: 1,
			},
			want: 32,
		},
		{
			name: "2",
			args: args{
				arr:  []int{4, 3, 2, 1},
				n:    4,
				even: -1,
			},
			want: 16,
		},
		{
			name: "3",
			args: args{
				arr:  []int{4, 3, 2, 1},
				n:    3,
				even: -1,
			},
			want: 5,
		},
		{
			name: "4",
			args: args{
				arr:  []int{4, 3, 2, 1},
				n:    6,
				even: -1,
			},
			want: 283,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get8Perm(tt.args.arr, tt.args.n, tt.args.even); got != tt.want {
				t.Errorf("Get8Perm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNOri(t *testing.T) {
	type args struct {
		arr      []int
		n        int
		evenBase int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				n:        3,
				evenBase: 2,
			},
			want: 7,
		},
		{
			name: "2",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				n:        3,
				evenBase: 1,
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				n:        3,
				evenBase: -2,
			},
			want: 3,
		},
		{
			name: "4",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				n:        6,
				evenBase: -2,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNOri(tt.args.arr, tt.args.n, tt.args.evenBase); got != tt.want {
				t.Errorf("GetNOri() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNOri(t *testing.T) {
	type args struct {
		arr      []int
		idx      int
		n        int
		evenBase int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				idx:      1,
				n:        2,
				evenBase: 2,
			},
			want: []int{0, 1, 5, 2, 4, 6},
		},
		{
			name: "2",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				idx:      1,
				n:        2,
				evenBase: 1,
			},
			want: []int{0, 0, 5, 2, 4, 6},
		},
		{
			name: "3",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				idx:      1,
				n:        2,
				evenBase: -2,
			},
			want: []int{1, 1, 5, 2, 4, 6},
		},
		{
			name: "4",
			args: args{
				arr:      []int{1, 3, 5, 2, 4, 6},
				idx:      1,
				n:        6,
				evenBase: -2,
			},
			want: []int{1, 1, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetNOri(tt.args.arr, tt.args.idx, tt.args.n, tt.args.evenBase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetNOri() = %v, want %v", got, tt.want)
			}
		})
	}
}
