package main

import (
	"fmt"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(7, day7)
}

func day7() {
	vx := aoc.MustInts(7)

	minf := crabfuel(vx, 0)
	for x, to := imin(vx...), imax(vx...); x < to; x++ {
		f := crabfuel(vx, x)
		if f < minf {
			minf = f
		}
	}

	fmt.Println("Day 7/1:", minf)
}

func crabfuel(vx []int, target int) int {
	sum := 0
	for _, x := range vx {
		sum += iabs(x - target)
	}
	return sum
}

func imin(v ...int) int {
	if len(v) == 0 {
		panic("empty input")
	}
	z := v[0]
	for _, x := range v[1:] {
		if x < z {
			z = x
		}
	}
	return z
}

func imax(v ...int) int {
	if len(v) == 0 {
		panic("empty input")
	}
	z := v[0]
	for _, x := range v[1:] {
		if x > z {
			z = x
		}
	}
	return z
}

func iabs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}
