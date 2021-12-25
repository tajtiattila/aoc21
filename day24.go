package main

import (
	"fmt"
	"io"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(24, day24)
}

func day24() {
	prg := parse24(aoc.Reader(24))
	fmt.Println(prg)

	cf := codefind24{
		prg: prg,
	}
	cf.digit(0, 0)

	fmt.Println("Day 24/1:", cf.bestc)
}

type codefind24 struct {
	prg []coef24
	num []byte

	bestc uint64
}

func (cf *codefind24) digit(z, ndigit int) {
	if len(cf.num) != 14 {
		cf.num = make([]byte, 14)
		for i := 0; i < 14; i++ {
			cf.num[i] = '1'
		}
	}

	if ndigit == 14 {
		if z == 0 {
			cf.havenum()
		}
		return
	}

	coef := cf.prg[ndigit]
	loz := coef.divz == 26 && coef.addx < 0
	if loz {
		ww := (z % 26) + coef.addx
		if ww < 1 || 9 < ww {
			return
		}
		cf.num[ndigit] = '0' + byte(ww)
		_, z = coef.apply(z, ww)
		cf.digit(z, ndigit+1)
		return
	}

	// check all digits 1..9
	for digit := 1; digit <= 9; digit++ {
		cf.num[ndigit] = '0' + byte(digit)
		_, znext := coef.apply(z, digit)
		if !loz || znext < z {
			cf.digit(znext, ndigit+1)
		}
	}
}

func (cf *codefind24) havenum() {
	var c uint64
	for _, b := range cf.num {
		aoc.Logf("%c", b)
		c = c*10 + uint64(b-'0')
	}
	aoc.Logln()
	if c > cf.bestc {
		cf.bestc = c
	}
}

func runcoef24(prg []coef24, w []byte) int {
	z := 0
	for _, x := range w {
		aoc.Logf("%c", x)
	}
	aoc.Log(":")
	for i, coef := range prg {
		var x int
		x, z = coef.apply(z, int(w[i]-'0'))
		aoc.Logf(" (%d)%d", x, z)
	}
	aoc.Logln()
	return z
}

type coef24 struct {
	divz int
	addx int
	addy int
}

func (coef coef24) apply(z, w int) (xcmp, zo int) {
	var x, y int

	// mul x 0
	// add x z
	// mod x 26
	// div z coef₀
	// add x coef₁
	x = (z % 26) + coef.addx
	z /= coef.divz

	// eql x w
	// eql x 0
	xcmp = x
	if x == w {
		x = 0
	} else {
		x = 1
	}

	// mul y 0
	// add y 25
	// mul y x
	// add y 1
	// mul z y
	y = 25*x + 1
	z *= y

	// mul y 0
	// add y w
	// add y coef₂
	// mul y x
	// add z y
	y = (w + coef.addy) * x
	z += y

	return xcmp, z
}

// CHRISTMAS

func parse24(r io.Reader) []coef24 {
	const digitf = `inp w
mul x 0
add x z
mod x 26
div z %d
add x %d
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y %d
mul y x
add z y
`

	var v []coef24
	for i := 0; i < 14; i++ {
		var coef coef24
		_, err := fmt.Fscanf(r, digitf, &coef.divz, &coef.addx, &coef.addy)
		check(err)
		v = append(v, coef)
	}

	return v
}
