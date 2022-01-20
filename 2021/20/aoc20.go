package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type bitmap struct {
	data   [][]byte
	h, w   int
	popcnt int
}

var (
	cur, nxt = 0, 1    // parity, ^parity
	kern     []byte    // kern(el filter)
	bufs     [2]bitmap // double buffering
)

func init() {
	N := 200
	bufs[0].data = make([][]byte, N)
	bufs[1].data = make([][]byte, N)
	for j := range bufs[0].data {
		bufs[0].data[j] = make([]byte, N)
		bufs[1].data[j] = make([]byte, N)
	}
}

func (b *bitmap) redim(h, w int) {
	b.h, b.w = h, w
	b.popcnt = 0
}

func (b bitmap) get(y, x int) int {
	if y < 0 || y >= b.h || x < 0 || x >= b.w { // p is infinite
		return cur
	}
	return int(b.data[y][x])
}

func enhance() {
	h, w := bufs[cur].h+2, bufs[cur].w+2
	bufs[nxt].redim(h, w)

	kern9 := func(y, x int) byte { // apply filter
		δy := []int{-1, -1, -1, +0, 0, 0, +1, 1, 1}
		δx := []int{-1, +0, +1, -1, 0, 1, -1, 0, 1}

		n := 0
		for i := range δx {
			n = (n << 1) | bufs[cur].get(y+δy[i], x+δx[i])
		}

		if kern[n] == '#' {
			return 1
		}
		return 0
	}

	data := bufs[nxt].data
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			data[y][x] = kern9(y-1, x-1)
			if data[y][x] == 1 {
				bufs[nxt].popcnt++
			}
		}
	}

	cur, nxt = nxt, cur // swap buffers (and switch parity)
}

func (b bitmap) String() string {
	var sb strings.Builder

	for y := -1; y < b.h+1; y++ {
		for x := -1; x < b.w+1; x++ {
			if b.get(y, x) == 1 {
				sb.WriteByte('@')
			} else {
				sb.WriteByte('.')
			}
		}
		if y != b.h+1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func main() {
	var raw []string

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch len(input.Bytes()) {
		case 0: // continue
		case 512:
			kern = []byte(input.Text())
		default:
			line := input.Text()
			raw = append(raw, line)
			h, w = h+1, len(line)
		}
	}

	for j := range raw {
		for i, v := range raw[j] {
			if v == '#' {
				bufs[cur].data[j][i] = 1
				bufs[cur].popcnt++
			}
		}
	}
	bufs[cur].redim(h, w)

	for i := 0; i < 50; i++ {
		if i == 2 {
			fmt.Println(bufs[cur].popcnt) // part1
		}
		enhance()
	}
	fmt.Println(bufs[cur].popcnt) // part2
}
