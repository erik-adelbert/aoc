package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type grid struct {
	d    [][]byte
	w, h int
}

func newGrid() *grid {
	d := make([][]byte, 128)
	for i := 0; i < 128; i++ {
		d[i] = make([]byte, 128)
	}

	return &grid{d, 128, 128}
}

func (g *grid) set(v byte, x, y int) {
	g.d[y+1][x+1] = v
}

func (g *grid) get(x, y int) byte {
	return g.d[y+1][x+1]
}

func (g *grid) neigh(x, y int) []byte {
	neigh := []byte{
		g.d[y][x+1], g.d[y+1][x], g.d[y+2][x+1], g.d[y+1][x+2],
	}
	return neigh
}

func (g *grid) redim(w, h int) {
	g.w, g.h = w, h
}

func (g *grid) copy(i int, data []byte) {
	copy(g.d[i+1], append([]byte{0}, data...))
}

func (g *grid) String() string {
	var sb strings.Builder
	for i := 1; i <= g.h; i++ {
		for j := 1; j <= g.w; j++ {
			sb.WriteByte(g.d[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	g := newGrid()

	w, h, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := input.Bytes()
		g.copy(h, args)
		w = len(args)
		h++
	}
	g.redim(w, h)

	n := 0
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			v := g.get(x, y)

			min := v
			for _, x := range g.neigh(x, y) {
				if x > 0 && v >= x {
					min = 0
				}
			}
			if min == v {
				n += 1 + int(v-'0')
			}
		}
	}
	fmt.Println(n)
}
