package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(5, day5)
}

func day5() {
	vseg, err := readsegments(aoc.Reader(5))
	if err != nil {
		log.Fatal(err)
	}
	ma := make(map[point]int)
	mb := make(map[point]int)
	for _, s := range vseg {
		if s.a.x == s.b.x || s.a.y == s.b.y {
			mapseg(ma, s)
		}
		mapseg(mb, s)
	}

	n := 0
	for _, v := range ma {
		if v > 1 {
			n++
		}
	}

	m := 0
	for _, v := range mb {
		if v > 1 {
			m++
		}
	}

	fmt.Println("Day 4/1:", n)
	fmt.Println("Day 4/2:", m)
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

type segment struct {
	a, b point
}

func (s segment) String() string {
	return fmt.Sprintf("%s -> %s", s.a, s.b)
}

func readsegments(r io.Reader) ([]segment, error) {
	var v []segment
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var s segment
		_, err := fmt.Sscanf(scanner.Text(), "%d,%d -> %d,%d", &s.a.x, &s.a.y, &s.b.x, &s.b.y)
		if err != nil {
			return nil, err
		}
		v = append(v, s)
	}
	return v, scanner.Err()
}

func mapseg(m map[point]int, s segment) {
	p := s.a

	for {
		m[p]++
		if p == s.b {
			return
		}
		adji(&p.x, s.b.x)
		adji(&p.y, s.b.y)
	}
}

func adji(v *int, goal int) {
	if *v == goal {
		return
	}
	if *v < goal {
		*v++
	} else {
		*v--
	}
}
