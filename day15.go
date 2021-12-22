package main

import (
	"fmt"

	"github.com/tajtiattila/aoc"
	"github.com/tajtiattila/aoc/ascmap"
	"github.com/tajtiattila/aoc/astar"
)

func init() {
	aoc.Register(15, day15)
}

func day15() {
	am, err := ascmap.FromReader(aoc.Reader(15))
	check(err)

	fmt.Println("Day 15/1:", path15cost(am))
}

func path15cost(am *ascmap.Ascmap) int {
	endp := point{am.Dx - 1, am.Dy - 1}
	_, cost := astar.FindPath(point{0, 0}, func(p0 astar.Point, dst []astar.State) []astar.State {
		p := p0.(point)
		for _, d := range dir4 {
			q := point{p.x + d.x, p.y + d.y}
			if am.In(q.x, q.y) {
				dst = append(dst, astar.State{
					Point:        q,
					Cost:         int(am.At(q.x, q.y) - '0'),
					EstimateLeft: iabs(q.x-endp.x) + iabs(q.y-endp.y),
				})
			}
		}
		return dst
	})
	return cost
}
