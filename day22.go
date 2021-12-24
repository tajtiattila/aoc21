package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(22, day22)
}

func day22() {
	ops := parsecuboidops(aoc.Reader(22))

	var p1 []cuboidop
	const a = 50

Outer:
	for _, o := range ops {
		for i := 0; i < 3; i++ {
			if iabs(o.min[i]) > a || iabs(o.max[i]) > a {
				continue Outer
			}
		}
		p1 = append(p1, o)
	}

	fmt.Println("Day 22/1:", mapcuboids(p1))
	fmt.Println("Day 22/2:", mapcuboids(ops))
}

func parsecuboidops(r io.Reader) []cuboidop {
	var v []cuboidop
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var s string
		var o cuboidop
		_, err := fmt.Sscanf(scanner.Text(), "%s x=%d..%d,y=%d..%d,z=%d..%d",
			&s,
			&o.min[0], &o.max[0],
			&o.min[1], &o.max[1],
			&o.min[2], &o.max[2])
		check(err)
		if s == "on" {
			o.on = true
		}
		v = append(v, o)
	}
	check(scanner.Err())
	return v
}

type cuboidop struct {
	min, max vec3
	on       bool
}

func mapcuboids(ops []cuboidop) int64 {
	type M map[int]struct{}
	mm := [3]M{make(M), make(M), make(M)}
	for _, o := range ops {
		for i := 0; i < 3; i++ {
			mm[i][o.min[i]] = struct{}{}
			mm[i][o.max[i]+1] = struct{}{}
		}
	}

	n := 1
	var r [3][]int
	for i := 0; i < 3; i++ {
		for c := range mm[i] {
			r[i] = append(r[i], c)
		}
		n *= len(r[i])
		sort.Ints(r[i])
	}

	cm := cuboidm{
		axv: r,
		pix: make([]byte, n),
	}
	cm.ystride = len(cm.axv[0])
	cm.zstride = cm.ystride * len(cm.axv[1])

	for _, o := range ops {
		lo, hi := o.min, o.max
		for i := 0; i < 3; i++ {
			hi[i]++
		}
		cm.set(cm.ofs3(lo), cm.ofs3(hi), o.on)
	}

	return cm.countones()
}

type cuboidm struct {
	axv [3][]int
	pix []byte

	ystride int
	zstride int
}

func (cm *cuboidm) ofs3(v vec3) vec3 {
	var o vec3
	for i := 0; i < 3; i++ {
		w := sort.SearchInts(cm.axv[i], v[i])
		if cm.axv[i][w] != v[i] {
			panic("impossible")
		}
		o[i] = w
	}
	return o
}

func (cm *cuboidm) set(lo, hi vec3, one bool) {
	var b byte
	if one {
		b = 1
	} else {
		b = 0
	}
	for z := lo[2]; z < hi[2]; z++ {
		for y := lo[1]; y < hi[1]; y++ {
			ofs := z*cm.zstride + y*cm.ystride + lo[0]
			for x := lo[0]; x < hi[0]; x++ {
				cm.pix[ofs] = b
				ofs++
			}
		}
	}
}

func (cm *cuboidm) countones() int64 {
	xa := cm.axv[0]
	ya := cm.axv[1]
	za := cm.axv[2]

	xhi := len(xa) - 1
	yhi := len(ya) - 1
	zhi := len(za) - 1

	var n int64
	for z := 0; z < zhi; z++ {
		zm := za[z+1] - za[z]
		for y := 0; y < yhi; y++ {
			ym := ya[y+1] - ya[y]
			ofs := z*cm.zstride + y*cm.ystride
			for x := 0; x < xhi; x++ {
				if cm.pix[ofs] > 0 {
					xm := xa[x+1] - xa[x]
					n += int64(xm * ym * zm)
				}
				ofs++
			}
		}
	}
	return n
}
