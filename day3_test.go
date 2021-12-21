package main

import (
	"io"
	"strings"
	"testing"

	"github.com/tajtiattila/aoc"
)

func day3tr() io.Reader {
	return strings.NewReader(`00100
11110
10110
10111
10101
01111
00111
11100
10000
11001
00010
01010
`)
}

func TestDay3(t *testing.T) {
	vbits, mask, err := parsebinary(day3tr())
	if err != nil {
		t.Fatal(err)
	}

	if testing.Verbose() {
		aoc.Verbose = true
	}

	const wanto2 = 23
	o2 := filter3(vbits, mask, 1)
	if o2 != wanto2 {
		t.Fatalf("Got O₂ = %d, want %d", o2, wanto2)
	}

	const wantco2 = 10
	co2 := filter3(vbits, mask, 0)
	if co2 != wantco2 {
		t.Fatalf("Got CO₂ = %d, want %d", co2, wantco2)
	}
}
