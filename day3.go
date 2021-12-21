package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(3, day3)
}

func day3() {
	nlines := 0
	var ones []int

	scanner := bufio.NewScanner(aoc.Reader(3))
	for scanner.Scan() {
		if ones == nil {
			ones = make([]int, len(scanner.Text()))
		}

		for i, ch := range scanner.Text() {
			if ch == '1' {
				ones[i]++
			}
		}

		nlines++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("IO error: %w", err)
	}

	var γ, mask uint64
	for _, n := range ones {
		γ <<= 1
		if n*2 > nlines {
			γ |= 1
		}
		mask = (mask << 1) | 1
	}
	ε := ^γ & mask

	aoc.Logln(ones, nlines)
	aoc.Logf("γ = %#b\nε = %#b\n", γ, ε)
	fmt.Println("Day 3/1:", γ*ε)
}
