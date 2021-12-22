package main

import (
	"fmt"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(17, day17)
}

func day17() {
	ta := parsetargetarea(aoc.MustString(17))

	fmt.Println("Day 17/1:", maxshotheight(ta))
}

type box struct {
	min point
	max point
}

func (b box) in(p point) bool {
	return b.min.x <= p.x && p.x <= b.max.x &&
		b.min.y <= p.y && p.y <= b.max.y
}

func parsetargetarea(s string) box {
	var b box
	_, err := fmt.Sscanf(s, "target area: x=%d..%d, y=%d..%d",
		&b.min.x, &b.max.x,
		&b.min.y, &b.max.y)
	check(err)
	return b
}

func maxshotheight(ta box) int {
	maxh := 0
	for x := 1; x < ta.min.x/2; x++ {
		maxx := x * (x + 1) / 2
		if maxx < ta.min.x {
			continue // shot can't reach minx
		}
		maxy := -ta.min.y // shot would "skip" target
		for y := 1; y <= maxy; y++ {
			h, r := shot(ta, x, y)
			if r == shotabove {
				break
			}
			if r == shothit && h > maxh {
				maxh = h
			}
		}
	}
	return maxh
}

type shotres int

const (
	shothit shotres = iota
	shotleft
	shotabove
)

func shot(ta box, vx, vy int) (height int, res shotres) {
	x, y := 0, 0
	for {
		if y > height {
			height = y
		}
		if ta.in(point{x, y}) {
			return height, shothit
		}
		if y < ta.min.y {
			return 0, shotleft
		}
		if x > ta.max.x {
			return 0, shotabove
		}
		x += vx
		y += vy
		if vx > 0 {
			vx--
		}
		vy--
	}
}
