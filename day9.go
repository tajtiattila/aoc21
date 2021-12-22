package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/tajtiattila/aoc"
	"github.com/tajtiattila/aoc/ascmap"
)

func init() {
	aoc.Register(9, day9)
}

func day9() {
	m, err := ascmap.FromReader(aoc.Reader(9))
	if err != nil {
		log.Fatal(err)
	}

	sum := 0
	var basin []int
	for p := range mappoints(m) {
		if v := lowpointv(m, p.x, p.y); v != 0 {
			sum += v
			basin = append(basin, basinsize(m, p.x, p.y))
		}
	}

	fmt.Println("Day 9/1:", sum)

	sort.Ints(basin)
	v := basin[len(basin)-3:]

	fmt.Println("Day 9/2:", v[0]*v[1]*v[2])
}

func mappoints(m *ascmap.Ascmap) <-chan point {
	ch := make(chan point)
	go func() {
		for y := 0; y < m.Dy; y++ {
			for x := 0; x < m.Dx; x++ {
				ch <- point{x, y}
			}
		}
		close(ch)
	}()
	return ch
}

var dir4 = []point{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

func lowpointv(m *ascmap.Ascmap, x, y int) int {
	c := m.At(x, y)
	for _, e := range dir4 {
		v := m.At(x+e.x, y+e.y)
		if v != 0 && v <= c {
			return 0
		}
	}
	return int(c-'0') + 1
}

func basinsize(m *ascmap.Ascmap, x, y int) int {
	size := 0
	floodpts(x, y, func(x, y int) bool {
		c := m.At(x, y)
		if '0' <= c && c < '9' {
			size++
			return true
		}
		return false
	})
	return size
}

func floodpts(x, y int, cond func(x, y int) bool) {
	seen := make(map[point]struct{})
	var pts []point

	visit := func(x, y int) bool {
		p := point{x, y}
		if _, ok := seen[p]; ok || !cond(x, y) {
			return false
		}
		pts = append(pts, p)
		seen[p] = struct{}{}
		return true
	}

	if !visit(x, y) {
		return
	}

	for len(pts) > 0 {
		n := len(pts) - 1
		p := pts[n]
		pts = pts[:n]

		for _, delta := range dir4 {
			visit(p.x+delta.x, p.y+delta.y)
		}
	}
}
