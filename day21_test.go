package main

import "testing"

func TestDay21(t *testing.T) {
	p1, p2 := 4, 8
	const want = 739785
	got := sim21(detdice(100), p1, p2)
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
