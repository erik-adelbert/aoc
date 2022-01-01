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

func (g *grid) set(v byte, y, x int) {
	g.d[y+1][x+1] = v
}

func (g *grid) get(y, x int) byte {
	return g.d[y+1][x+1]
}

func (g *grid) filter(y, x int) int {
	v := g.get(y, x)
	for _, x := range g.neigh(y, x) {
		if x > 0 && v >= x {
			return 0
		}
	}
	return 1 + int(v-'0')
}

func (g *grid) neigh(y, x int) []byte {
	return []byte{
		g.d[y][x+1], g.d[y+1][x], g.d[y+2][x+1], g.d[y+1][x+2],
	}
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

func (g *grid) copy(i int, data []byte) {
	copy(g.d[i+1], append([]byte{0}, data...))
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
		args := input.Bytes()
		g.copy(h, args)
		h, w = h+1, len(args)
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
