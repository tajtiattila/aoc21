package main

import (
	"fmt"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(21, day21)
}

func day21() {
	p1, p2 := parse21(aoc.MustString(21))

	fmt.Println("Day 21/1:", sim21(detdice(100), p1, p2))
	fmt.Println("Day 21/2:", diracwinmost(p1, p2))
}

func parse21(s string) (p1, p2 int) {
	_, err := fmt.Sscanf(s, `Player 1 starting position: %d
Player 2 starting position: %d`, &p1, &p2)
	check(err)
	return p1, p2
}

func sim21(d dice, q1, q2 int) int {
	p1 := player21{pos: q1}
	p2 := player21{pos: q2}

	nroll := 0
	stepf := func(p *player21) {
		adv := d.roll() + d.roll() + d.roll()
		nroll += 3

		p.pos = 1 + (p.pos+adv-1)%10
		p.score += p.pos
	}

	for {
		stepf(&p1)
		if p1.score >= 1000 {
			return nroll * p2.score
		}
		stepf(&p2)
		if p2.score >= 1000 {
			return nroll * p1.score
		}
	}
}

type player21 struct {
	pos   int
	score int
}

type dice interface {
	roll() int
}

type deterministicdice struct {
	next int
	max  int
}

func detdice(sides int) dice {
	return &deterministicdice{max: sides}
}

func (d *deterministicdice) roll() int {
	i := d.next + 1
	d.next = (d.next + 1) % d.max
	return i
}

type diracpl struct {
	p [2]int

	score [2]int

	pno int // next player num (0/1)
	rno int // next roll idx (0..2)
}

type diracwin [2]int64

func diracwinmost(p1, p2 int) int64 {
	m := make(map[diracpl]diracwin)
	state := diracpl{
		p: [2]int{p1, p2},
	}
	dw := diracwinsub(m, state, 21)
	if dw[0] > dw[1] {
		return dw[0]
	} else {
		return dw[1]
	}
}

/*
10@19 10@19 (1) → 6,2
r1: 1@20 10@19 (2) → 3,2
	r1: 1@20 1@20 (1) → 3,0
		p1win*3
	r2: 1@20 2@21 p2win
	r3: 1@20 3@22 p2win
r2: 2@21 10@19 p1win
r3: 3@22 10@19 p1win
*/
func diracwinsub(m map[diracpl]diracwin, state diracpl, winscore int) diracwin {
	if r, ok := m[state]; ok {
		return r
	}

	aoc.Logf("%d@%d %d@%d (%d:%d)\n",
		state.p[0], state.score[0], state.p[1], state.score[1],
		state.pno+1, state.rno+1)

	var dw diracwin

	i := state.pno
	nstate := state
	nstate.rno = (state.rno + 1) % 3
	if nstate.rno == 0 {
		nstate.pno = 1 - state.pno
	}
	for roll := 1; roll <= 3; roll++ {
		nstate.p[i] = 1 + (state.p[i]+roll-1)%10
		nstate.score[i] = state.score[i] + nstate.p[i]
		if nstate.rno == 0 && nstate.score[i] >= winscore {
			dw[i]++
		} else {
			r := diracwinsub(m, nstate, winscore)
			dw[0] += r[0]
			dw[1] += r[1]
		}
	}

	m[state] = dw

	aoc.Logf(" %d@%d %d@%d (%d:%d) → %d,%d\n",
		state.p[0], state.score[0], state.p[1], state.score[1],
		state.pno+1, state.rno+1,
		dw[0], dw[1])
	return dw
}
