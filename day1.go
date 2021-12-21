package main

import (
	"fmt"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(1, day1)
}

func day1() {
	input := aoc.MustInts(1)

	n, lastv := 0, 0
	for i, v := range input {
		if i > 0 && v > lastv {
			n++
		}
		lastv = v
	}
	fmt.Println("Day 1/1:", n)

	lastsum := input[0] + input[1] + input[2]
	n = 0
	for i := 3; i < len(input); i++ {
		sum := lastsum + input[i] - input[i-3]
		if sum > lastsum {
			n++
		}
		lastsum = sum
	}
	fmt.Println("Day 1/2:", n)
}
