package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(3, day3)
}

func day3() {

	vbits, mask, err := parsebinary(aoc.Reader(3))
	if err != nil {
		log.Fatalf("IO error: %w", err)
	}

	var maxbit uint64
	if mask > 1 {
		maxbit = (mask >> 1) + 1
	} else {
		maxbit = 1
	}

	ones := make([]int, bits.OnesCount64(mask))
	for _, bits := range vbits {
		for i := range ones {
			if bits&(maxbit>>i) != 0 {
				ones[i]++
			}
		}
	}

	var γ uint64
	for _, n := range ones {
		γ <<= 1
		if n*2 > len(vbits) {
			γ |= 1
		}
	}
	ε := ^γ & mask

	aoc.Logln(ones, len(vbits))
	aoc.Logf("γ = %#b\nε = %#b\n", γ, ε)
	fmt.Println("Day 3/1:", γ*ε)
}

func parsebinary(r io.Reader) (v []uint64, mask uint64, err error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var w uint64
		for _, ch := range scanner.Text() {
			w <<= 1
			if ch == '1' {
				w |= 1
			}
		}
		v = append(v, w)
		if mask == 0 {
			for range scanner.Text() {
				mask = (mask << 1) | 1
			}
		}
	}
	return v, mask, scanner.Err()
}
