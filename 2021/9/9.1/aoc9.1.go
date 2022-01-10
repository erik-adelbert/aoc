package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type grid struct {
	d    [][]byte
	h, w int
}

func newGrid() *grid {
	d := make([][]byte, 128)
	for i := 0; i < 128; i++ {
		d[i] = make([]byte, 128)
	}

	return &grid{d, 128, 128}
}

func (g *grid) get(y, x int) byte {
	return g.d[y+1][x+1]
}

func (g *grid) filter(y, x int) int {
	btoi := func(b byte) int {
		return int(b - '0') // fast convert
	}

	v := g.get(y, x)
	for _, x := range g.neigh(y, x) {
		if 0 < x && x <= v {
			return 0
		}
	}
	return 1 + btoi(v)
}

func (g *grid) neigh(y, x int) []byte {
	return []byte{
		g.d[y][x+1], g.d[y+1][x], g.d[y+2][x+1], g.d[y+1][x+2],
	}
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

func (g *grid) copy(i int, data []byte) int {
	t := g.d[i+1]

	t, t[0] = t[1:], 0
	return copy(t, data)
}

func (g *grid) String() string {
	var sb strings.Builder
	for j := 1; j <= g.h; j++ {
		for i := 1; i <= g.w; i++ {
			sb.WriteByte(g.d[j][i])
		}
		if j != g.h {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	g := newGrid()

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		data := input.Bytes()
		w = g.copy(h, data)
		h++
	}
	g.redim(h, w)

	sum := 0
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			sum += g.filter(y, x)
		}
	}
	fmt.Println(sum)
}
