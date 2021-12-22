package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tajtiattila/aoc"

	"log"
)

func init() {
	aoc.Register(18, day18)
}

func day18() {
	p := parsesnailpuzzle(aoc.Reader(18))
	var sum *snailnum
	for _, v := range p {
		sum = sum.add(v)
	}
	fmt.Println("Day 18/1:", sum.mag())

	maxmag := 0
	for i, vi := range p {
		for j, vj := range p {
			if i == j {
				continue
			}
			mij := vi.add(vj).mag()
			mji := vj.add(vi).mag()
			maxmag = imax(maxmag, mij, mji)
		}
	}
	fmt.Println("Day 18/2:", maxmag)
}

func parsesnailpuzzle(r io.Reader) []*snailnum {
	scanner := bufio.NewScanner(r)
	var v []*snailnum
	for scanner.Scan() {
		sn, err := parsesnailnum(scanner.Text())
		check(err)
		v = append(v, sn)
	}
	check(scanner.Err())
	return v
}

type snailnum struct {
	child [2]*snailnum
	value int
}

type snailnumparser struct {
	s string
	i int

	err error
}

func (p *snailnumparser) ch() byte {
	if p.i >= len(p.s) {
		return 0
	}
	i := p.i
	p.i++
	return p.s[i]
}

func (p *snailnumparser) num() *snailnum {
	if p.err != nil {
		return nil
	}

	c := p.ch()
	if '0' <= c && c <= '9' {
		return &snailnum{
			value: int(c - '0'),
		}
	}

	if c != '[' {
		return p.errch(c, '[')
	}

	l := p.num()

	if !p.eat(',') {
		return nil
	}

	r := p.num()

	if !p.eat(']') {
		return nil
	}

	return &snailnum{
		child: [2]*snailnum{l, r},
	}
}

func (p *snailnumparser) eat(want byte) bool {
	if p.err != nil {
		return false
	}
	if c := p.ch(); c != want {
		p.errch(c, want)
		return false
	}
	return true
}

func (p *snailnumparser) errch(c, want byte) *snailnum {
	p.err = fmt.Errorf("%q: unexpected rune %q at position %d, want %q", p.s, c, p.i, want)
	return nil
}

func parsesnailnum(s string) (*snailnum, error) {
	p := snailnumparser{s: s}
	return p.num(), p.err
}

func (sn *snailnum) clone() *snailnum {
	if sn == nil {
		return nil
	}

	return &snailnum{
		child: [2]*snailnum{
			sn.child[0].clone(),
			sn.child[1].clone(),
		},
		value: sn.value,
	}
}

func (a *snailnum) add(b *snailnum) *snailnum {
	if a == nil {
		return b.clone()
	}

	if b == nil {
		return a.clone()
	}

	sn := &snailnum{
		child: [2]*snailnum{a.clone(), b.clone()},
	}

	sn.reduce()
	return sn
}

func (sn *snailnum) mag() int {
	if sn.child[0] == nil && sn.child[1] == nil {
		return sn.value
	}
	return 3*sn.child[0].mag() + 2*sn.child[1].mag()
}

func (sn *snailnum) addv(v int) {
	if sn == nil {
		return
	}
	if sn.child[0] != nil || sn.child[1] != nil {
		log.Fatalln("invalid snailnum addv")
	}

	sn.value += v
}

func (sn *snailnum) regpair() (left, right int, ok bool) {
	if sn == nil {
		return 0, 0, false // invalid
	}
	if sn.child[0] == nil || sn.child[1] == nil {
		return 0, 0, false // regular number, not a pair
	}
	return sn.child[0].value, sn.child[1].value, true
}

func (sn *snailnum) reduce() {
	v := make([]*snailnum, 0, 64)
	for sn.explode(v) || sn.split() {
	}
}

func (sn *snailnum) explode(stk []*snailnum) bool {
	if sn == nil {
		return false
	}

	n := len(stk)
	if n >= 4 {
		if l, r, ok := sn.regpair(); ok {
			sn.sibling(stk, 0).addv(l)
			sn.sibling(stk, 1).addv(r)

			sn.child[0] = nil
			sn.child[1] = nil
			sn.value = 0
			return true
		}
	}

	stk = append(stk, sn)
	defer func() {
		stk = stk[:n]
	}()

	expl := sn.child[0].explode(stk)
	if !expl {
		expl = sn.child[1].explode(stk)
	}
	return expl
}

func (sn *snailnum) split() bool {
	if sn == nil {
		return false
	}

	if sn.child[0] == nil && sn.child[1] == nil && sn.value >= 10 {
		l := sn.value / 2
		r := sn.value - l
		sn.child[0] = &snailnum{value: l}
		sn.child[1] = &snailnum{value: r}
		sn.value = 0
		return true
	}

	return sn.child[0].split() || sn.child[1].split()
}

// dir: 0: left, 1: right
func (sn *snailnum) sibling(stk []*snailnum, dir int) *snailnum {
	p := sn
	for i := len(stk) - 1; i >= 0; i-- {
		z := stk[i]
		if z.child[1-dir] == p {
			return z.child[dir].deepregular(1 - dir)
		}
		p = z
	}
	return nil
}

func (sn *snailnum) deepregular(dir int) *snailnum {
	p := sn
	for {
		if c := p.child[dir]; c != nil {
			p = c
		} else {
			break
		}
	}
	return p
}
