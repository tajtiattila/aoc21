package main

import (
	"bufio"
	"fmt"
	"sort"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(10, day10)
}

func day10() {
	scanner := bufio.NewScanner(aoc.Reader(10))
	errscore := 0
	var compscore []int
	for scanner.Scan() {
		es, cs := bracer(scanner.Text())
		errscore += es
		if cs != 0 {
			compscore = append(compscore, cs)
		}
	}
	check(scanner.Err())

	fmt.Println("Day 10/1:", errscore)

	sort.Ints(compscore)
	mid := len(compscore) / 2
	fmt.Println("Day 10/2:", compscore[mid])
}

var closebrace = map[rune]rune{
	'(': ')',
	'[': ']',
	'{': '}',
	'<': '>',
}

var openbrace = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
	'>': '<',
}

var bracescore = map[rune]int{
	// completion scores
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,

	// error scores
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func bracer(s string) (errscore, compscore int) {
	var stack []rune
	for _, r := range s {
		switch r {
		case '(', '[', '{', '<':
			stack = append(stack, r)
		case ')', ']', '}', '>':
			n := len(stack) - 1
			if n < 0 {
				return 0, 0
			}
			if stack[n] != openbrace[r] {
				return bracescore[r], 0
			}
			stack = stack[:n]
		default:
			panic("invalid rune")
		}
	}

	m := 1
	for _, r := range stack {
		compscore += m * bracescore[r]
		m *= 5
	}
	return 0, compscore
}
