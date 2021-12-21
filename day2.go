package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(2, day2)
}

func day2() {
	input, err := parsesubcmds(aoc.Reader(2))
	if err != nil {
		log.Fatalf("Error parsing commands: %w", err)
	}

	day2_1(input)
	day2_2(input)
}

func day2_1(input []subcmd) {
	var x, y int
	for _, c := range input {
		x += c.dx
		y += c.dy
	}
	fmt.Println("Day 2/1:", x*y)
}

func day2_2(input []subcmd) {
	var x, y, aim int64
	for _, c := range input {
		f := int64(c.dx)
		x += f
		y += aim * f
		aim += int64(c.dy)
	}
	fmt.Println("Day 2/2:", x*y)
}

type subcmd struct {
	dx, dy int
}

func parsesubcmds(r io.Reader) ([]subcmd, error) {
	var cmds []subcmd
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var cmd string
		var v int
		if _, err := fmt.Sscan(scanner.Text(), &cmd, &v); err != nil {
			return nil, err
		}

		sc, err := decodesubcmd(cmd, v)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, sc)
	}
	return cmds, scanner.Err()
}

func decodesubcmd(cmd string, value int) (subcmd, error) {
	switch cmd {
	case "forward":
		return subcmd{value, 0}, nil
	case "down":
		return subcmd{0, value}, nil
	case "up":
		return subcmd{0, -value}, nil
	}

	return subcmd{}, fmt.Errorf("Invalid command %s", cmd)
}
