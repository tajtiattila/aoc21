package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tajtiattila/aoc"
)

func init() {
	aoc.Register(20, day20)
}

func day20() {
	algo, im := parse20(aoc.Reader(20))

	fmt.Println("Day 20/1:", im.enhance(algo, 2).nlit())
	fmt.Println("Day 20/2:", im.enhance(algo, 50).nlit())
}

func parse20(r io.Reader) (algo []byte, im eimage) {
	scanner := bufio.NewScanner(r)

	scanner.Scan()
	algo = make([]byte, 512)
	copy(algo, scanner.Bytes())
	if len(scanner.Bytes()) != 512 {
		log.Panic("invalid algo")
	}

	scanner.Scan()

	var w int
	var pix []byte
	for scanner.Scan() {
		b := scanner.Bytes()
		w = len(b)
		pix = append(pix, b...)
	}

	check(scanner.Err())

	return algo, eimage{
		dim:    point{w, (len(pix) / w)},
		stride: w,
		pix:    pix,
		def:    '.',
	}
}

type eimage struct {
	org point
	dim point

	stride int

	pix []byte

	def byte
}

func (im eimage) String() string {
	var sb strings.Builder
	for i := 0; i != len(im.pix); i += im.stride {
		sb.Write(im.pix[i : i+im.stride])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (im eimage) nlit() int {
	if im.def == '#' {
		panic("eimage.nlit: infinite")
	}
	n := 0
	for _, p := range im.pix {
		if p == '#' {
			n++
		}
	}
	return n
}

func (im eimage) ofs(x, y int) int {
	xo := x - im.org.x
	yo := y - im.org.y
	if xo < 0 || yo < 0 || xo >= im.dim.x || yo >= im.dim.y {
		return -1
	}
	return xo + yo*im.stride
}

func (im eimage) at(x, y int) byte {
	o := im.ofs(x, y)
	if o < 0 {
		return im.def
	}
	return im.pix[o]
}

func (im eimage) set(x, y int, b byte) {
	o := im.ofs(x, y)
	if o < 0 {
		panic("eimage.set")
	}
	im.pix[o] = b
}

var logimageenhance bool

func (im eimage) enhance(algo []byte, num int) eimage {
	if logimageenhance {
		fmt.Printf("Before enhance:\n%s\n", im)
	}
	for i := 0; i < num; i++ {
		im = im.enhanceonce(algo)
		if logimageenhance {
			fmt.Printf("After enhance %d:\n%s\n", i+1, im)
		}
	}
	return im
}

func (im eimage) enhanceonce(algo []byte) eimage {
	w := eimage{
		org: point{
			im.org.x - 1,
			im.org.y - 1,
		},
		dim: point{
			im.dim.x + 2,
			im.dim.y + 2,
		},
	}
	w.stride = w.dim.x
	w.pix = make([]byte, w.dim.x*w.dim.y)

	var defidx int
	if im.def == '#' {
		defidx = 511
	}
	w.def = algo[defidx]

	for y, ye := w.org.y, w.org.y+w.dim.y; y < ye; y++ {
		for x, xe := w.org.x, w.org.x+w.dim.x; x < xe; x++ {
			var bits int
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					c := im.at(x+dx, y+dy)

					bits <<= 1
					if c == '#' {
						bits |= 1
					}
				}
			}
			w.set(x, y, algo[bits])
		}
	}
	return w
}
