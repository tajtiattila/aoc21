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

	cost := runday23(8, st)
	fmt.Println("Day 23/1:", cost)

	st = st.unfold()
	cost = runday23(16, st)
	fmt.Println("Day 23/2:", cost)
}

func runday23(npod int, st amphistate) int {
	var ab amphibuf
	ab.npod = npod
	path, cost := astar.FindPath(st, func(p0 astar.Point, dst []astar.State) []astar.State {
		st := p0.(amphistate)
		return ab.next(st, dst)
	})
	if aoc.Verbose {
		for _, st := range path {
			fmt.Println(st.(amphistate).image(ab.npod))
		}
	}
	if len(path) == 0 {
		return -1
	}
	return cost
}

// positions of pods AABBCCDD:
//   0..10 hallway
//   13, 15, 17, 19 upper room
//   24, 26, 28, 30 lower room
//   35, 37, 39, 41 upper unfolded room
//   46, 48, 50, 52 lower unfolded room
type amphistate [16]byte

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

func (st amphistate) unfold() amphistate {
	const yst = amphiystride
	var uf amphistate
	for i := 0; i < 8; i++ {
		j := (i/2)*4 + i%2
		p := st[i]
		if p >= 2*yst {
			p += 2 * yst
		}
		uf[j] = p
	}
	//     2 4 6 8
	// y1 #D#C#B#A#
	// y2 #D#B#A#C#
	const y1 = 2 * yst
	const y2 = 3 * yst
	uf[2] = y1 + 8
	uf[3] = y2 + 6
	uf[6] = y1 + 6
	uf[7] = y2 + 4
	uf[10] = y1 + 4
	uf[11] = y2 + 8
	uf[14] = y1 + 2
	uf[15] = y2 + 2
	return uf
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

func (st amphistate) costleft(npod int) int {
	n := 0
	nkind := npod / 4
	for i, p := range st[:npod] {
		x, y := amphixy(int(p))
		pod := i / nkind
		roomx := amphipodroomx(pod)
		if x == roomx {
			// pod in correct room
			continue
		}
		nstep := iabs(x - roomx) // x movement
		if y == 0 {
			// pod in hallway, must move down at least once
			nstep++
		} else {
			// pod in wrong room, must exit and then down at least once
			nstep += y + 1
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
	nkind := b.npod / 4
	nlen := amphiystride * (bottom + 1)
	if len(b.im) < nlen {
		b.im = make([]byte, nlen)
	}
	for i := range b.im[:nlen] {
		b.im[i] = 0
	}
	for i, p := range st[:b.npod] {
		pod := i / nkind
		b.im[int(p)] = 'A' + byte(pod)
	}

	if logamphinext {
		fmt.Print(st.image(b.npod))
	}

	// find possible moves
	hasvalid := false
	for i, p := range st[:b.npod] {
		pod := i / nkind
		stepcost := amphimovecost[pod]
		nstate := st
		for _, m := range b.validmoves(pod, p) {
			nstate[i] = m.pos
			if logamphinext {
				fmt.Printf("%c: %d → %d (%d)\n", 'A'+byte(pod), p, m.pos, m.nstep)
			}
			hasvalid = true
			dst = append(dst, astar.State{
				Point: nstate,
				Cost:  stepcost * m.nstep,

				EstimateLeft: nstate.costleft(b.npod),
			})
		}
	}
	if !hasvalid && logamphinext {
		fmt.Println("no valid moves")
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
		for y, p := 1, dx+yst; y <= bottom; y, p = y+1, p+yst {
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
			if logamphinext {
				fmt.Printf("%c: %d → %d other pods in room\n",
					'A'+byte(pod), start, dx)
			}
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
