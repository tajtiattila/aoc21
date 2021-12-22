package main

import (
	"fmt"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(11, day11)
}

func day11() {
	const w = 10
	d := parsedumboct(w, aoc.MustString(11))

	fmt.Println("Day 11/1:", d.run(100))

	n := 100
	for {
		n++
		nflash := d.step()
		if nflash == w*w {
			break
		}
	}

	fmt.Println("Day 11/2:", n)
}

type dumboct struct {
	m []uint8
	w int
}

func (d dumboct) String() string {
	var b strings.Builder
	for i, v := range d.m {
		if v < 10 {
			b.WriteByte('0' + v)
		} else {
			b.WriteByte('a' + v - 10)
		}
		if i%d.w == d.w-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func (d dumboct) in(x, y int) bool {
	return x >= 0 && y >= 0 && x < d.w && y < d.w
}

func (d dumboct) inc(x, y int) {
	i := x + y*d.w
	d.m[i]++
	if d.m[i] == 10 {
		for _, p := range dir8 {
			x2, y2 := x+p.x, y+p.y
			if d.in(x2, y2) {
				d.inc(x2, y2)
			}
		}
	}
}

func (d dumboct) step() (nflash int) {
	for i := range d.m {
		x, y := i%d.w, i/d.w
		d.inc(x, y)
	}

	for i := range d.m {
		if d.m[i] > 9 {
			nflash++
			d.m[i] = 0
		}
	}

	return nflash
}

func (d dumboct) run(nstep int) (nflash int) {
	for i := 0; i < nstep; i++ {
		nflash += d.step()
	}
	return nflash
}

var dir8 = []point{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

func parsedumboct(w int, src string) dumboct {
	m := make([]uint8, w*w)
	i := 0
	for _, r := range src {
		if '0' <= r && r <= '9' {
			m[i] = uint8(r - '0')
			i++
		}
	}
	return dumboct{m, w}
}
