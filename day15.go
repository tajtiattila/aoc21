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

	fmt.Println("Day 15/1:", path15cost(am, 1))
	fmt.Println("Day 15/2:", path15cost(am, 5))
}

func path15cost(am *ascmap.Ascmap, rep int) int {
	dx, dy := rep*am.Dx, rep*am.Dy

	inm := func(p point) bool {
		return 0 <= p.x && 0 <= p.y && p.x < dx && p.y < dy
	}
	costm := func(p point) int {
		add := p.x/am.Dx + p.y/am.Dy
		v := int(am.At(p.x%am.Dx, p.y%am.Dy) - '0')
		return 1 + (v+add-1)%9
	}

	endp := point{dx - 1, dy - 1}
	_, cost := astar.FindPath(point{0, 0}, func(p0 astar.Point, dst []astar.State) []astar.State {
		p := p0.(point)
		for _, d := range dir4 {
			q := point{p.x + d.x, p.y + d.y}
			if inm(q) {
				dst = append(dst, astar.State{
					Point:        q,
					Cost:         costm(q),
					EstimateLeft: iabs(q.x-endp.x) + iabs(q.y-endp.y),
				})
			}
		}
		return dst
	})
	return cost
}
