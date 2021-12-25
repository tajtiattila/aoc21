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
	ab.npod = 8
	path, cost := astar.FindPath(st, func(p0 astar.Point, dst []astar.State) []astar.State {
		st := p0.(amphistate)
		return ab.next(st, dst)
	})
	if aoc.Verbose {
		for _, st := range path {
			fmt.Println(st.(amphistate).image(ab.npod))
		}
	}
	fmt.Println("Day 23/1:", cost)
}

// positions of pods AABBCCDD:
//   0..10 hallway
//   13, 15, 17, 19 upper room
//   24, 26, 28, 30 lower room
type amphistate [8]byte

const amphiystride = 11

func parseamphistate(s string) amphistate {
	v := strings.Split(s, "\n")[1:]
	var podp [4][]byte
	for r := 0; r < 8; r++ {
		x, y := 2+(r/2)*2, 1+r%2
		pod := v[y][x+1]
		i := pod - 'A'
		podp[i] = append(podp[i], byte(x+amphiystride*y))
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

var amphimovecost = []int{1, 10, 100, 1000}

func amphixy(p int) (x, y int) {
	x, y = p%amphiystride, p/amphiystride
	return
}

func (st amphistate) image(npod int) string {
	rrows := npod / 4
	b := [][]byte{
		[]byte("#############\n"),
		[]byte("#...........#\n"),
		[]byte("###.#.#.#.###\n"),
	}
	for i := 1; i < rrows; i++ {
		b = append(b, []byte("  #.#.#.#.#\n"))
	}
	b = append(b, []byte("  #########\n"))

	for i, p := range st[:npod] {
		pod := i / rrows
		x, y := amphixy(int(p))
		b[y+1][x+1] = 'A' + byte(pod)
	}

	var sb strings.Builder
	for _, row := range b {
		sb.Write(row)
	}

	return sb.String()
}

func amphipodroomx(pod int) int {
	return 2 + 2*pod
}

func (st amphistate) costleft() int {
	n := 0
	for i, p := range st[:] {
		x, y := amphixy(int(p))
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
	im   []byte
	mbuf []amphimove

	npod int // 8 or 16
}

type amphimove struct {
	pos   byte
	nstep int
}

var logamphinext bool

func (b *amphibuf) next(st amphistate, dst []astar.State) []astar.State {
	// refresh map
	bottom := b.npod / 4
	nlen := amphiystride * (bottom + 1)
	if len(b.im) < nlen {
		b.im = make([]byte, nlen)
	}
	for i := range b.im[:nlen] {
		b.im[i] = 0
	}
	for i, p := range st[:] {
		pod := i / 2
		b.im[int(p)] = 'A' + byte(pod)
	}

	if logamphinext {
		fmt.Print(st.image(b.npod))
	}

	// find possible moves
	for i, p := range st[:b.npod] {
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

func (b *amphibuf) validmoves(pod int, start0 byte) []amphimove {
	start := int(start0)
	sx, sy := amphixy((start))
	bottom := b.npod / 4

	podc := 'A' + byte(pod)

	const yst = amphiystride
	v := b.mbuf[:0]
	if sy > 0 {
		// move out of the room
		rx := amphipodroomx(pod)
		if rx == sx {
			// pod in correct room

			// check positions below
			done := true
			for y, p := sy+1, start+yst; y <= bottom; y, p = y+1, p+yst {
				if b.im[p] != podc {
					// blocking another pod
					done = false
					break
				}
			}
			if done {
				// already at final room position
				return nil
			}
		}

		for p := start - yst; p > 0; p -= yst {
			if b.im[p] != 0 {
				// leave blocked
				return nil
			}
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
		var nsame, nother int
		for y, p := 0, dx; y <= bottom; y, p = y+1, p+yst {
			if dc := b.im[p]; dc != 0 {
				if podc == dc {
					nsame++
				} else {
					nother++
					break
				}
			}
		}
		if nother != 0 {
			// other pods blocking
			return nil
		}

		dy := bottom - nsame

		nstep := dy
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
					fmt.Printf("%c: %d → %d path to room blocked\n",
						'A'+byte(pod), start, dx)
				}
				return nil
			}
		}
		v = append(v, amphimove{pos: byte(dx + dy*yst), nstep: nstep})
	}

	b.mbuf = v
	return v
}
