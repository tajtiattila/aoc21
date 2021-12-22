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

	maxh, nhit := simshot(ta)
	fmt.Println("Day 17/1:", maxh)
	fmt.Println("Day 17/2:", nhit)
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

func simshot(ta box) (maxheight, nhit int) {
	for x := 1; x <= ta.max.x; x++ {
		maxx := x * (x + 1) / 2
		if maxx < ta.min.x {
			continue // shot can't reach minx
		}
		d := -ta.min.y
		for y := -d; y <= d; y++ {
			if h, ok := shot(ta, x, y); ok {
				nhit++
				if h > maxheight {
					maxheight = h
				}
			}
		}
	}
	return maxheight, nhit
}

func shot(ta box, vx, vy int) (height int, hit bool) {
	x, y := 0, 0
	for {
		if y > height {
			height = y
		}
		if ta.in(point{x, y}) {
			return height, true
		}
		if y < ta.min.y {
			return 0, false
		}
		if x > ta.max.x {
			return 0, false
		}
		x += vx
		y += vy
		if vx > 0 {
			vx--
		}
		vy--
	}
}
