package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(25, day25)
}

func day25() {
	puzzle := parsecucumap(aoc.Reader(25))

	m := puzzle.clone()
	var scratch cucumap
	nstep := 1
	for m.step(&scratch) {
		nstep++
	}
	fmt.Println("Day 25/1:", nstep)
}

type cucumap struct {
	w   int
	pix []byte
}

func parsecucumap(r io.Reader) cucumap {
	var cm cucumap
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Bytes()
		cm.w = len(l)
		cm.pix = append(cm.pix, l...)
	}
	check(scanner.Err())
	return cm
}

func (cm cucumap) String() string {
	var sb strings.Builder

	for i := 0; i < len(cm.pix); i += cm.w {
		sb.Write(cm.pix[i : i+cm.w])
		sb.WriteByte('\n')
	}

	return sb.String()
}

func (cm cucumap) clone() cucumap {
	dst := cucumap{
		w:   cm.w,
		pix: make([]byte, len(cm.pix)),
	}
	copy(dst.pix, cm.pix)
	return dst
}

func (cm cucumap) copy(dst *cucumap) {
	dst.w = cm.w
	if len(dst.pix) != len(cm.pix) {
		dst.pix = make([]byte, len(cm.pix))
	}
	copy(dst.pix, cm.pix)
}

func (cm cucumap) step(scratch *cucumap) bool {
	scratch.w = cm.w
	if len(scratch.pix) != len(cm.pix) {
		scratch.pix = make([]byte, len(cm.pix))
	}

	moved := false

	// move east
	for i := range scratch.pix {
		scratch.pix[i] = '.'
	}
	for i, p := range cm.pix {
		switch p {
		case 'v':
			scratch.pix[i] = 'v'
			fallthrough
		case '.':
			continue
		}
		var right int
		if x := i % cm.w; x == cm.w-1 {
			right = i - x
		} else {
			right = i + 1
		}
		if cm.pix[right] == '.' {
			moved = true
			scratch.pix[right] = '>'
		} else {
			scratch.pix[i] = '>'
		}
	}

	// move south
	for i := range cm.pix {
		cm.pix[i] = '.'
	}
	dy := len(cm.pix) / cm.w
	i := 0
	for y := 0; y < dy; y++ {
		for x := 0; x < cm.w; x, i = x+1, i+1 {
			p := scratch.pix[i]
			switch p {
			case '>':
				cm.pix[i] = '>'
				fallthrough
			case '.':
				continue
			}
			var down int
			if y == dy-1 {
				down = x
			} else {
				down = i + cm.w
			}
			if scratch.pix[down] == '.' {
				moved = true
				cm.pix[down] = 'v'
			} else {
				cm.pix[i] = 'v'
			}
		}
	}

	return moved
}
