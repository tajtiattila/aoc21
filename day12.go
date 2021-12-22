package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(12, day12)
}

func day12() {
	g := parsecavegraph(aoc.Reader(12))
	fmt.Println("Day 12/1:", paths11(g, false))
	fmt.Println("Day 12/2:", paths11(g, true))
}

type cavegraph map[string][]string

func parsecavegraph(r io.Reader) cavegraph {
	scanner := bufio.NewScanner(r)
	g := make(map[string][]string)

	edge := func(a, b string) {
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}

	for scanner.Scan() {
		v := strings.Split(scanner.Text(), "-")
		edge(v[0], v[1])
	}
	check(scanner.Err())

	return g
}

func minorcave(s string) bool {
	return len(s) > 0 && 'a' < s[0] && s[0] < 'z'
}

type path11inf struct {
	seen      map[string]struct{}
	twice     bool
	seentwice string
}

func (i *path11inf) enter(c string) bool {
	if c == "start" {
		return false
	}

	if !minorcave(c) {
		return true
	}

	if _, ok := i.seen[c]; !ok {
		i.seen[c] = struct{}{}
		return true
	}
	if i.twice && i.seentwice == "" {
		i.seentwice = c
		return true
	}
	return false
}

func (i *path11inf) leave(c string) {
	if i.twice && i.seentwice == c {
		i.seentwice = ""
		return
	}

	delete(i.seen, c)
}

func paths11(g cavegraph, twice bool) int {
	i := &path11inf{
		seen:  make(map[string]struct{}),
		twice: twice,
	}
	const from = "start"
	i.enter(from)
	return paths11x(g, i, from)
}

func paths11x(g cavegraph, inf *path11inf, from string) int {
	n := 0

	for _, next := range g[from] {
		switch next {
		case "end":
			n++
		default:
			if inf.enter(next) {
				n += paths11x(g, inf, next)
				inf.leave(next)
			}
		}
	}

	return n
}
