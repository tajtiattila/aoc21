package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/tajtiattila/aoc"
	"github.com/tajtiattila/aoc/input"
)

func init() {
	aoc.Register(4, day4)
}

func day4() {
	puzzle, err := parsebingo(aoc.Reader(4))
	if err != nil {
		log.Fatal(err)
	}

	var firstscore, lastscore int

	won := make([]bool, len(puzzle.card))
	n := len(puzzle.card)
Outer:
	for _, d := range puzzle.draw {
		for i, c := range puzzle.card {
			if won[i] {
				continue
			}
			c.mark(d)
			if !c.win() {
				continue
			}

			switch n {
			case len(puzzle.card):
				firstscore = c.score(d)
			case 1:
				lastscore = c.score(d)
				break Outer
			}

			won[i] = true
			n--
		}
	}

	fmt.Println("Day 4/1:", firstscore)
	fmt.Println("Day 4/2:", lastscore)

}

type bingocell struct {
	value int
	mark  bool
}

type bingocard []bingocell

func (c bingocard) resetmarks() {
	for i := range c {
		c[i].mark = false
	}
}

func (c bingocard) mark(n int) {
	for i := range c {
		if c[i].value == n {
			c[i].mark = true
		}
	}
}

func (c bingocard) win() bool {
	for i := 0; i < 5; i++ {
		if c.winrow(i) || c.wincol(i) {
			return true
		}
	}
	return false
}

func (c bingocard) winrow(row int) bool {
	ofs := 5 * row
	for _, x := range c[ofs : ofs+5] {
		if !x.mark {
			return false
		}
	}
	return true
}

func (c bingocard) wincol(col int) bool {
	for i := col; i < 25; i += 5 {
		if !c[i].mark {
			return false
		}
	}
	return true
}

func (c bingocard) score(lastdraw int) int {
	var u int
	for _, x := range c {
		if !x.mark {
			u += x.value
		}
	}
	return u * lastdraw
}

type bingopuzzle struct {
	draw []int
	card []bingocard
}

func parsebingo(r io.Reader) (*bingopuzzle, error) {
	var draw []int
	var card []bingocard
	scanner := input.NewBlockScanner(r)
	for scanner.Scan() {
		var err error
		if draw == nil {
			draw, err = parsedraw(scanner.Text())
		} else {
			var c bingocard
			c, err = parsebingocard(scanner.Text())
			card = append(card, c)
		}
		if err != nil {
			return nil, err
		}
	}
	return &bingopuzzle{
		draw: draw,
		card: card,
	}, scanner.Err()
}

func parsedraw(s string) ([]int, error) {
	var v []int
	for _, x := range strings.Split(s, ",") {
		n, err := strconv.Atoi(x)
		if err != nil {
			return nil, err
		}
		v = append(v, n)
	}
	return v, nil
}

func parsebingocard(s string) (bingocard, error) {
	var v []bingocell
	for _, x := range strings.Fields(s) {
		n, err := strconv.Atoi(x)
		if err != nil {
			return nil, err
		}
		v = append(v, bingocell{value: n})
	}

	if len(v) != 25 {
		return nil, fmt.Errorf("Invalid bingo card: %s", s)
	}

	return bingocard(v), nil
}
