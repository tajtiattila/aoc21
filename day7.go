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

	id := func(x int) int { return x }
	cf2 := func(x int) int { return (x*x + x) / 2 }

	fmt.Println("Day 7/1:", crabminfuel(vx, id))
	fmt.Println("Day 7/2:", crabminfuel(vx, cf2))
}

func crabminfuel(vx []int, fuelf func(dist int) int) int {
	minf := crabfuel(vx, 0, fuelf)
	for x, to := imin(vx...), imax(vx...); x < to; x++ {
		f := crabfuel(vx, x, fuelf)
		if f < minf {
			minf = f
		}
	}
	return minf
}

func crabfuel(vx []int, target int, fuelf func(dist int) int) int {
	sum := 0
	for _, x := range vx {
		sum += fuelf(iabs(x - target))
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
