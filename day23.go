package main

import (
	"fmt"
	"strings"

	"github.com/tajtiattila/aoc"
	"github.com/tajtiattila/aoc/astar"
)

func init() {
	aoc.Register(23, day23)
}

func day23() {
	st := parseamphistate(aoc.MustString(23))

	var ab amphibuf
	path, cost := astar.FindPath(st, func(p0 astar.Point, dst []astar.State) []astar.State {
		st := p0.(amphistate)
		return ab.next(st, dst)
	})
	if aoc.Verbose {
		for _, st := range path {
			fmt.Println(st.(amphistate))
		}
	}
	fmt.Println("Day 23/1:", cost)
}

// positions of pods AABBCCDD:
//   0..10 hallway
//   13, 15, 17, 19 upper room
//   24, 26, 28, 30 lower room
type amphistate [8]byte

func parseamphistate(s string) amphistate {
	v := strings.Split(s, "\n")[1:]
	var podp [4][]byte
	for r := 0; r < 8; r++ {
		x, y := 2+(r/2)*2, 1+r%2
		pod := v[y][x+1]
		i := pod - 'A'
		podp[i] = append(podp[i], byte(x+11*y))
	}
	var w []byte
	for _, vv := range podp[:] {
		if len(vv) != 2 {
			panic("parse error")
		}
		w = append(w, vv...)
	}
	var st amphistate
	copy(st[:], w)
	return st
}

var amphitpl = `#############
#...........#
###.#.#.#.###
  #.#.#.#.#
  #########
`

var amphivalidpos = []byte{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	13, 15, 17, 19,
	24, 26, 28, 30,
}

var amphimovecost = []int{1, 10, 100, 1000}

func (st amphistate) String() string {
	b := []byte(amphitpl)
	for i, p := range st[:] {
		pod := i / 2
		x, y := p%11, p/11
		ofs := x + 1 + (y+1)*14
		b[ofs] = 'A' + byte(pod)
	}
	return string(b)
}

func amphipodroomx(pod int) int {
	return 2 + 2*pod
}

func (st amphistate) costleft() int {
	n := 0
	for i, p := range st[:] {
		x, y := int(p%11), int(p/11)
		pod := i / 2
		roomx := amphipodroomx(pod)
		if x == roomx {
			// pod in correct room
			continue
		}
		var nstep int
		switch y {
		case 0:
			// pod in hallway
			nstep = 1 + iabs(x-roomx)
		case 1:
			nstep = 4
		case 2:
			nstep = 5
		}
		n += nstep * amphimovecost[pod]
	}
	return n
}

type amphibuf struct {
	im   [33]byte
	mbuf []amphimove
}

type amphimove struct {
	pos   byte
	nstep int
}

var logamphinext bool

func (b *amphibuf) next(st amphistate, dst []astar.State) []astar.State {
	// refresh map
	for i := range b.im[:] {
		b.im[i] = 0
	}
	for i, p := range st[:] {
		pod := i / 2
		b.im[int(p)] = 'A' + byte(pod)
	}

	if logamphinext {
		fmt.Print(st)
	}

	// find possible moves
	for i, p := range st[:] {
		pod := i / 2
		stepcost := amphimovecost[pod]
		nstate := st
		for _, m := range b.validmoves(pod, p) {
			nstate[i] = m.pos
			if logamphinext {
				fmt.Printf("%c: %d → %d (%d)\n", 'A'+byte(pod), p, m.pos, m.nstep)
			}
			dst = append(dst, astar.State{
				Point: nstate,
				Cost:  stepcost * m.nstep,

				EstimateLeft: nstate.costleft(),
			})
		}
	}

	return dst
}

func (b *amphibuf) validmoves(pod int, start byte) []amphimove {
	sx, sy := int(start%11), int(start/11)

	v := b.mbuf[:0]
	if sy > 0 {
		// move out of the room
		rx := amphipodroomx(pod)
		if int(start) == rx+22 {
			// already at inner room pos
			return nil
		}

		pc := 'A' + byte(pod)
		ro, ri := b.im[rx+11], b.im[rx+22]
		if ro == pc && ri == pc {
			// room done
			return nil
		}

		if sy == 2 && b.im[int(start)-11] != 0 {
			// blocked
			return nil
		}

		c0 := sy

		// hwxvalid reports if the hallway position is valid
		hwxvalid := func(x int) bool {
			return x%2 == 1 || x == 0 || x == 10
		}

		// move left in hallway
		for ex := sx - 1; ex >= 0; ex-- {
			if !hwxvalid(ex) {
				continue
			}
			if b.im[ex] != 0 {
				break // movement blocked
			}
			c1 := sx - ex
			v = append(v, amphimove{pos: byte(ex), nstep: c0 + c1})
		}

		// move right in hallway
		for ex := sx + 1; ex <= 10; ex++ {
			if !hwxvalid(ex) {
				continue
			}
			if b.im[ex] != 0 {
				break // movement blocked
			}
			c1 := ex - sx
			v = append(v, amphimove{pos: byte(ex), nstep: c0 + c1})
		}
	} else {
		// move to destination room
		dx := amphipodroomx(pod)
		if b.im[11+dx] != 0 {
			return nil // room full
		}
		dpos := dx + 11
		nstep := 1
		if b.im[22+dx] == 0 {
			// move to inner slot
			dpos += 11
			nstep++
		}

		var xadd int
		if dx < sx {
			xadd = -1
			nstep += sx - dx
		} else {
			xadd = 1
			nstep += dx - sx
		}

		for ex := sx + xadd; ex != dx; ex += xadd {
			if b.im[ex] != 0 {
				if logamphinext {
					fmt.Printf("%c: %d → %d room blocked\n",
						'A'+byte(pod), start, dx)
				}
				return nil
			}
		}
		v = append(v, amphimove{pos: byte(dpos), nstep: nstep})
	}

	b.mbuf = v
	return v
}
