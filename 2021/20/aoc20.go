// aoc20.go --
// advent of code 2021 day 20
//
// https://adventofcode.com/2021/day/20
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-20: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MAXLEN = 1 << 8

type (
	kernel [2 * MAXLEN]byte
	bitmap [MAXLEN * MAXLEN]byte
)

type image struct {
	bmap   *bitmap
	h, w   int
	popcnt int
}

var (
	kern     kernel    // kernel filter
	bufs     [2]*image // double buffer
	cur, nxt = 0, 1    // parity, ^parity
)

func main() {
	for i := range bufs {
		bufs[i] = newImage()
	}

	j, i := 0, MAXLEN
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bit := strings.NewReplacer(".", "\x00", "#", "\x01")
		switch len(input.Bytes()) {
		case 0:
			// continue
		case len(kern):
			copy(kern[:], input.Bytes())
		default:
			low, max := slice(j, i)
			i = copy(bufs[cur].bmap[low:max:max], []byte(bit.Replace(input.Text())))
			j++
		}
	}
	bufs[cur].redim(j, i)

	for i := 0; i < 50; i++ {
		if i == 2 {
			fmt.Println(bufs[cur].popcnt) // part1
		}
		enhance()
	}
	fmt.Println(bufs[cur].popcnt) // part2
}

func enhance() {
	h, w := bufs[cur].h+2, bufs[cur].w+2
	bufs[nxt].redim(h, w)

	bitmap := bufs[nxt].bmap
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			δy := []int{-1, -1, -1, +0, 0, 0, +1, 1, 1}
			δx := []int{-1, +0, +1, -1, 0, 1, -1, 0, 1}

			n := 0
			for i := range δx {
				n = (n << 1) | bufs[cur].get(y-1+δy[i], x-1+δx[i])
			}

			bitmap[y*MAXLEN+x] = 0
			if kern[n] == '#' {
				bitmap[y*MAXLEN+x] = 1
				bufs[nxt].popcnt++
			}
		}
	}

	cur, nxt = nxt, cur // swap buffers (and switch parity)
}

func newBitmap() *bitmap {
	return new(bitmap)
}

func newImage() *image {
	return &image{
		newBitmap(), 0, 0, 0,
	}
}

func (a *image) redim(h, w int) {
	a.h, a.w, a.popcnt = h, w, 0
}

func slice(j, w int) (low, max int) {
	low = j * MAXLEN
	max = low + w
	return
}

func (a *image) get(y, x int) int {
	if y < 0 || y >= a.h || x < 0 || x >= a.w { // p is infinite
		return cur
	}
	return int(a.bmap[y*MAXLEN+x])
}

func (a *image) String() string {
	unbit := strings.NewReplacer("\x00", ".", "\x01", "#")
	var sb strings.Builder
	for j := 0; j < a.h; j++ {
		low, max := slice(j, a.w)
		sb.Write(append(a.bmap[low:max:max], '\n'))
	}
	return unbit.Replace(sb.String())
}
