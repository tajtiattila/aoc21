package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(16, day16)
}

func day16() {
	pkt := parsepacket(aoc.MustString(16))

	fmt.Println("Day 16/1:", pkt.sumver())
}

type bitpacket struct {
	ver, typ uint

	literal uint

	sub []bitpacket
}

func (k bitpacket) sumver() uint {
	sumv := k.ver
	for _, sub := range k.sub {
		sumv += sub.sumver()
	}
	return sumv
}

func parsepacket(s string) bitpacket {
	bits, err := hex.DecodeString(strings.TrimSpace(s))
	check(err)

	br := bitreader{p: bits}
	return br.packet()
}

type bitreader struct {
	p []byte
	i int
}

func (br *bitreader) readbit() uint {
	ofs := br.i / 8
	if ofs >= len(br.p) {
		return 0
	}
	shift := 7 - br.i%8
	br.i++
	return uint(br.p[ofs]>>shift) & 1
}

func (br *bitreader) read(nbits int) uint {
	var v uint
	for i := 0; i < nbits; i++ {
		v = v<<1 | br.readbit()
	}
	return v
}

func (br *bitreader) literal() uint {
	var v uint
	for cont := true; cont; {
		cont = br.readbit() == 1
		v = v<<4 | br.read(4)
	}
	return v
}

func (br *bitreader) packet() bitpacket {
	ver, typ := br.read(3), br.read(3)
	pkt := bitpacket{ver: ver, typ: typ}
	if typ == 4 {
		pkt.literal = br.literal()
		return pkt
	}

	lengthtype := br.readbit()
	if lengthtype == 0 {
		nbits := br.read(15)
		end := br.i + int(nbits)
		for br.i != end {
			pkt.sub = append(pkt.sub, br.packet())
		}
	} else {
		nsub := int(br.read(11))
		for i := 0; i < nsub; i++ {
			pkt.sub = append(pkt.sub, br.packet())
		}
	}
	return pkt
}
