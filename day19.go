package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(19, day19)
}

func day19() {
	sd := parse19(aoc.Reader(19))
	aoc.Logf("%d scanners\n", len(sd))

	beacon := beacons19(sd)
	fmt.Println("Day 10/1:", len(beacon))
}

func beacons19(sd []scannerdata) []vec3 {
	if len(sd) == 0 {
		return nil
	}

	matched := matchbeacons(sd)
	m := make(map[vec3]struct{})
	for _, d := range matched {
		for _, p := range d.beacon {
			p = p.mul(d.rot).add(d.tr)
			m[p] = struct{}{}
		}
	}
	var v []vec3
	for p := range m {
		v = append(v, p)
	}
	return v
}

func matchbeacons(src []scannerdata) []scannerdata {
	if len(src) == 0 {
		return nil
	}
	var done, left []scannerdata
	done = append(done, src[0])
	left = append(left, src[1:]...)

	for i := 0; len(left) != 0 && i < len(done); i++ {
		base := done[i]

		var nleft []scannerdata
		for _, next := range left {
			if beaconmatch(&base, &next) {
				done = append(done, next)
			} else {
				nleft = append(nleft, next)
			}
		}
		left = nleft
	}

	if len(left) != 0 {
		panic("matchbeacons: no progress")
	}
	return done
}

func beaconmatch(a, b *scannerdata) bool {
	for _, rot := range rot3 {
		m := make(map[vec3]int)
		for _, ka := range a.beacon {
			ka = ka.mul(a.rot).add(a.tr)
			for _, kb := range b.beacon {
				kb = kb.mul(rot)
				tr := ka.sub(kb)
				m[tr]++
			}
		}

		bestn := 0
		var besttr vec3
		for tr, n := range m {
			if n > bestn {
				besttr, bestn = tr, n
			}
		}
		if bestn >= 12 {
			b.rot = rot
			b.tr = besttr
			return true
		}
	}
	return false
}

type scannerdata struct {
	idx int // scanner index

	beacon []vec3

	rot mat3
	tr  vec3
}

func parse19(r io.Reader) []scannerdata {
	var d []scannerdata
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "---") {
			k := scannerdata{
				rot: mat3id,
			}
			_, err := fmt.Sscanf(t, "--- scanner %d ---", &k.idx)
			check(err)
			d = append(d, k)
		} else {
			var x, y, z int
			_, err := fmt.Sscanf(t, "%d,%d,%d", &x, &y, &z)
			check(err)
			n := len(d) - 1
			d[n].beacon = append(d[n].beacon, vec3{x, y, z})
		}
	}
	check(scanner.Err())
	return d
}

type vec3 [3]int

func (a vec3) add(b vec3) vec3 {
	return vec3{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
	}
}

func (a vec3) sub(b vec3) vec3 {
	return vec3{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}

func (a vec3) cross(b vec3) vec3 {
	return vec3{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

type mat3 [9]int

var mat3id = mat3{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1,
}

func mat3fromv(x, y, z vec3) mat3 {
	return mat3{
		x[0], x[1], x[2],
		y[0], y[1], y[2],
		z[0], z[1], z[2],
	}
}

func (v vec3) mul(m mat3) vec3 {
	return vec3{
		v[0]*m[0] + v[1]*m[3] + v[2]*m[6],
		v[0]*m[1] + v[1]*m[4] + v[2]*m[7],
		v[0]*m[2] + v[1]*m[5] + v[2]*m[8],
	}
}

var rot3 []mat3

func init() {
	rm := func(xa int, v int) {
		var x vec3
		x[xa] = v

		i := (xa + 1) % 3
		j := (xa + 2) % 3
		for _, d := range dir4 {
			var y vec3
			y[i] = d.x
			y[j] = d.y

			z := x.cross(y)
			rot3 = append(rot3, mat3fromv(x, y, z))
		}
	}

	for xa := 0; xa < 3; xa++ {
		rm(xa, -1)
		rm(xa, 1)
	}
}
