package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(13, day13)
}

func day13() {
	fp := parsefoldpaper(aoc.Reader(13))

	handlefold(fp.pt, fp.fold[0])
	fmt.Println("Day 13/1:", len(fp.pt))

	for _, f := range fp.fold[1:] {
		handlefold(fp.pt, f)
	}

	var xmin, xmax, ymin, ymax int
	first := true
	for pt := range fp.pt {
		if first {
			xmin, xmax = pt.x, pt.x
			ymin, ymax = pt.y, pt.y
			first = false
			continue
		}
		if pt.x < xmin {
			xmin = pt.x
		}
		if pt.y < ymin {
			ymin = pt.y
		}
		if xmax < pt.x {
			xmax = pt.x
		}
		if ymax < pt.y {
			ymax = pt.y
		}
	}

	aoc.Logln(xmin, ymin, xmax, ymax)

	line := make([][]byte, ymax-ymin+1)
	for i := range line {
		line[i] = bytes.Repeat([]byte{' '}, xmax-xmin+1)
	}
	for pt := range fp.pt {
		i, j := pt.y-ymin, pt.x-xmin
		line[i][j] = '#'
	}

	fmt.Println("Day 13/2:")

	for _, l := range line {
		fmt.Println(string(bytes.TrimRight(l, " ")))
	}
}

type paperfold struct {
	x bool
	v int
}

type foldpaper struct {
	pt   map[point]struct{}
	fold []paperfold
}

func parsefoldpaper(r io.Reader) foldpaper {
	p := foldpaper{
		pt: make(map[point]struct{}),
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "fold") {
			p.fold = append(p.fold, parsefold(t))
		} else {
			p.pt[parsept(t)] = struct{}{}
		}
	}
	check(scanner.Err())
	return p
}

func parsept(pt string) point {
	var p point
	_, err := fmt.Sscanf(pt, "%d,%d", &p.x, &p.y)
	check(err)
	return p
}

func parsefold(fold string) paperfold {
	var axis rune
	var value int
	_, err := fmt.Sscanf(fold, "fold along %c=%d", &axis, &value)
	check(err)
	return paperfold{x: axis == 'x', v: value}
}

func handlefold(m map[point]struct{}, fold paperfold) {
	for p := range m {
		if fold.x {
			if p.x > fold.v {
				delete(m, p)
				p.x = 2*fold.v - p.x
				m[p] = struct{}{}
			}
		} else {
			if p.y > fold.v {
				delete(m, p)
				p.y = 2*fold.v - p.y
				m[p] = struct{}{}
			}
		}
	}
}
