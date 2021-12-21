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
		log.Fatalf("IO error: %v", err)
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

	o2 := filter3(vbits, mask, 1)
	co2 := filter3(vbits, mask, 0)
	fmt.Println("Day 3/3:", o2*co2)
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

func filter3(v0 []uint64, mask uint64, bit byte) uint64 {

	shift := uint(bits.OnesCount64(mask))

	var v, w []uint64
	v = append(v, v0...)
	for shift != 0 && len(v) > 1 {
		w, v = v, w[:0]
		shift--

		onecount := 0
		for _, x := range w {
			if byte((x>>shift)&1) == 1 {
				onecount++
			}
		}

		var wantbit byte
		if 2*onecount >= len(w) {
			wantbit = bit
		} else {
			wantbit = bit ^ 1
		}

		/*
			oned := 2*onecount - len(w)
			var wantbit byte
			if bit == 1 {
				if oned >= 0 {
					wantbit = 1
				} else {
					wantbit = 0
				}
			} else {
				if oned >= 0 {
					wantbit = 0
				} else {
					wantbit = 1
				}

				if oned < 0 {
					wantbit = 1
				} else {
					wantbit = 0
				}
			}
		*/

		for _, x := range w {
			if byte((x>>shift)&1) == wantbit {
				v = append(v, x)
			}
		}
		aoc.Logf("%d %d %d %d/%d\n", shift, wantbit, len(v), onecount, len(w))
	}

	if len(v) != 1 {
		log.Fatalln("Logic error 1 != ", len(v))
	}

	return v[0]
}
