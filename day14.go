package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(14, day14)
}

func day14() {
	templ, rule := parsepolimeripuzzle(aoc.Reader(14))

	fmt.Println("Day 14/1:", rule.common(templ, 10))
	fmt.Println("Day 14/2:", rule.common(templ, 40))
}

type polimerirule map[[2]byte]byte

func (pr polimerirule) common(s string, depth uint) int {
	m := pr.maps(s, depth)
	type runecount struct {
		r byte
		n int
	}
	v := make([]runecount, 0, len(m))
	for r, n := range m {
		v = append(v, runecount{r, n})
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].n > v[j].n
	})
	n := len(v) - 1
	return v[0].n - v[n].n
}

func (pr polimerirule) maps(s string, depth uint) map[byte]int {
	if s == "" {
		return nil
	}
	sc := stepctx{
		m:     make(map[byte]int),
		cache: make(map[polimeridepthkey]map[byte]int),
	}
	left := s[0]
	sc.m[left]++
	for i := 1; i < len(s); i++ {
		right := s[i]
		pr.mapstep(sc, left, right, depth)
		left = right
	}
	return sc.m
}

type polimeridepthkey struct {
	key   [2]byte
	depth uint
}

type stepctx struct {
	m     map[byte]int
	cache map[polimeridepthkey]map[byte]int
}

func (pr polimerirule) mapstep(sc stepctx, left, right byte, depth uint) {
	if depth == 0 {
		sc.m[right]++
		return
	}

	dk := polimeridepthkey{
		key:   [2]byte{left, right},
		depth: depth,
	}

	addm := func(m, add map[byte]int) {
		for r, c := range add {
			m[r] += c
		}
		return
	}

	if xm, ok := sc.cache[dk]; ok {
		addm(sc.m, xm)
		return
	}

	parentm := sc.m
	sc.m = make(map[byte]int)
	pr.mapstepimpl(sc, left, right, depth)
	sc.cache[dk] = sc.m
	addm(parentm, sc.m)
}

func (pr polimerirule) mapstepimpl(sc stepctx, left, right byte, depth uint) {
	k := [2]byte{left, right}
	mid, ok := pr[k]
	if !ok {
		sc.m[right]++
		return
	}

	pr.mapstep(sc, left, mid, depth-1)
	pr.mapstep(sc, mid, right, depth-1)
}

func (r polimerirule) run(s string, n int) string {
	for i := 0; i < n; i++ {
		s = r.step(s)
	}
	return s
}

func (pr polimerirule) step(s string) string {
	if len(s) == 0 {
		return s
	}

	var p []byte
	l := s[0]
	p = append(p, l)
	for i := 1; i < len(s); i++ {
		r := s[i]
		k := [2]byte{l, r}
		if ins, ok := pr[k]; ok {
			p = append(p, ins)
		}
		p = append(p, r)
		l = r
	}

	return string(p)
}

func parsepolimeripuzzle(r io.Reader) (templ string, rule polimerirule) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	templ = scanner.Text()

	rule = make(polimerirule)
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			continue
		}
		var betw [2]byte
		var ins byte
		_, err := fmt.Sscanf(t, "%c%c -> %c", &betw[0], &betw[1], &ins)
		check(err)
		rule[betw] = ins
	}
	return templ, rule
}
