package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(6, day6)
}

func day6() {
	age, err := parseCommaInts(strings.TrimSpace(aoc.MustString(6)))
	if err != nil {
		log.Fatal(err)
	}

	v := lanternfishsim(age, 80)
	fmt.Println("Day 6/1:", lanternfishsim2(age, 80))
	fmt.Println("Day 6/2:", lanternfishsim2(age, 256))
}

func sumints(v []int) int64 {
	var sum int64
	for _, x := range v {
		sum += int64(x)
	}
	return sum
}

func lanternfishsim(v0 []int, days int) []int {
	var v, w []int
	v = make([]int, len(v0))
	copy(v, v0)

	for day := 0; day < days; day++ {
		v, w = w[:0], v
		n := 0
		for _, x := range w {
			if x > 0 {
				v = append(v, x-1)
			} else {
				v = append(v, 6)
				n++
			}
		}
		v = append(v, intv(8, n)...)
	}

	return v
}

func lanternfishsim2(v0 []int, days int) int64 {
	va := make([]int, 9)
	wa := make([]int, 9)

	for _, x := range v0 {
		va[x]++
	}

	for day := 0; day < days; day++ {
		va, wa = wa, va
		copy(va[0:], wa[1:])
		va[6] += wa[0]
		va[8] = wa[0]
	}
	return sumints(va)
}

func intv(value, count int) []int {
	v := make([]int, count)
	for i := range v {
		v[i] = value
	}
	return v
}
