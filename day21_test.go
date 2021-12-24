package main

import (
	"testing"
)

func TestDay21(t *testing.T) {
	p1, p2 := 4, 8
	const want = 739785
	got := sim21(detdice(100), p1, p2)
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func TestDay21_2(t *testing.T) {
	tests := []struct {
		winscore int
		w1, w2   int64
	}{
		{1, 27, 0},
		{2, 183, 156},
		{3, 990, 207},
		{4, 2930, 971},
		{5, 7907, 2728},
		{6, 30498, 7203},
		{7, 127019, 152976},
		{8, 655661, 1048978},
		{9, 4008007, 4049420},
		{10, 18973591, 12657100},
	}

	logday21 = testing.Verbose()

	p1, p2 := 4, 8
	for _, tt := range tests {
		m := make(map[diracpl]diracwin)
		state := diracpl{
			p: [2]int{p1, p2},
		}
		r := diracwinsub(m, state, tt.winscore)
		w1, w2 := r[0], r[1]
		if w1 != tt.w1 || w2 != tt.w2 {
			t.Errorf("winscore p1: %d got %d/%d, want %d/%d",
				tt.winscore, w1, w2, tt.w1, tt.w2)
		}
	}
}
