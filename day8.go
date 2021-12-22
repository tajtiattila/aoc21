package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/bits"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(8, day8)
}

func day8() {
	vs := parsesevensegs(aoc.Reader(8))
	n := 0
	for _, x := range vs {
		for _, d := range x.rawvalue {
			switch bits.OnesCount8(d) {
			case 2, 3, 4, 7:
				n++
			}
		}
	}

	fmt.Println("Day 8/1:", n)

	s := 0
	for _, x := range vs {
		s += x.value()
	}
	fmt.Println("Day 8/2:", s)
}

type sevenseg struct {
	digit    []uint8
	rawvalue []uint8
}

func (s sevenseg) value() int {
	m := make(map[uint8]int)
	for i, x := range s.digit {
		m[x] = i
	}

	v := 0
	for _, x := range s.rawvalue {
		v = v*10 + m[x]
	}
	return v
}

func parsesevensegs(r io.Reader) []sevenseg {
	var v []sevenseg
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		v = append(v, parsesevenseg(scanner.Text()))
	}
	check(scanner.Err())
	return v
}

func parsesevenseg(s string) sevenseg {
	var v []uint8
	for _, s := range strings.Fields(s) {
		if s == "|" {
			continue
		}
		v = append(v, parsealphabyte(s))
	}
	return sevenseg{
		digit:    sortDigits(v[:10]),
		rawvalue: v[10:],
	}
}

func parsealphabyte(s string) uint8 {
	var v uint8
	for _, r := range s {
		if r < 'a' || 'g' < r {
			panic("invalid rune")
		}
		v |= 1 << (r - 'a')
	}
	return v
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func sortDigits(src []uint8) []uint8 {
	db := make([]uint8, 10)
	set := func(digit int, vseg uint8) {
		if db[digit] != 0 {
			log.Panicf("digit %d can't be set to %07b, it was set to %07b", digit, vseg, db[digit])
		}
		db[digit] = vseg
	}

	nseg := bits.OnesCount8

	// Find 1, 4, 7, and 8.
	var v, w []uint8
	for _, d := range src {
		switch nseg(d) {
		case 2:
			set(1, d)
		case 3:
			set(7, d)
		case 4:
			set(4, d)
		case 7:
			set(8, d)
		default:
			v = append(v, d)
		}
	}

	if len(v) != 6 {
		panic("logic error after 1, 4, 7, 8")
	}

	// Find 0, 6 and 9.
	v, w = w[:0], v
	for _, d := range w {
		if nseg(d) != 6 {
			v = append(v, d)
			continue
		}
		d4, d7 := db[4], db[7]
		switch d {
		case d | d4 | d7:
			set(9, d)
		case d | d7:
			set(0, d)
		default:
			set(6, d)
		}
	}

	if len(v) != 3 {
		panic("logic error after 0, 6, 9")
	}

	// Find 2, 3 and 5.
	w = v
	for _, d := range w {
		switch d {
		case d & db[6]:
			set(5, d)
		case d | db[7]:
			set(3, d)
		default:
			set(2, d)
		}
	}

	return db
}

func fail(args ...interface{}) {
	panic(fmt.Sprint(args))
}
