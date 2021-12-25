package main

import (
	"strings"
	"testing"
)

func TestDay25(t *testing.T) {
	src := `v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>
`
	m := parsecucumap(strings.NewReader(src))

	var scratch cucumap
	nstep := 1
	for m.step(&scratch) {
		nstep++
	}

	t.Log(nstep)
	t.Log(m)
}
