package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type grid struct {
	d    [][]int
	w, h int
}

func newGrid() *grid {
	d := make([][]int, 128)
	for i := 0; i < 128; i++ {
		d[i] = make([]int, 128)
	}
	return &grid{d, 128, 128}
}

func (g *grid) set(v int, x, y int) {
	if x < 0 || y < 0 || g.d[y+1][x+1] == '9' {
		return
	}
	g.d[y+1][x+1] = v
}

func (g *grid) get(x, y int) int {
	if x < 0 || y < 0 {
		return 0
	}
	return g.d[y+1][x+1]
}

func (g *grid) components() map[int]int {
	label := newGrid() // reuse grid
	label.redim(g.w, g.h)

	labels := make([]int, 256+g.w*g.h) // labels (>256)
	for i := range labels {
		labels[i] = i
	}

	find := func(x int) int {
		for labels[x] != x {
			x, labels[x] = labels[x], labels[labels[x]] // path splitting
		}
		return x
	}

	union := func(x, y int) {
		if x > y {
			labels[y] = x
		} else {
			labels[x] = y
		}
	}

	id := 256                  // labels (>256)
	for y := 0; y < g.h; y++ { // Hoshen-Kopelman
		for x := 0; x < g.w; x++ {
			if g.get(x, y) == '9' {
				continue
			}
			nor, wes := label.get(x, y-1), label.get(x-1, y)
			switch {
			case nor == 0 && wes == 0:
				label.set(id, x, y)
				id++
			case nor != 0 && wes == 0:
				label.set(find(nor), x, y)
			case nor == 0 && wes != 0:
				label.set(find(wes), x, y)
			default:
				union(nor, wes)
				label.set(find(nor), x, y)
			}
		}
	}

	comps := make(map[int]int)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			v := label.get(x, y)
			if v != 0 {
				comps[find(v)]++
			}
		}
	}
	return comps
}

func (g *grid) redim(w, h int) {
	g.w, g.h = w, h
}

func (g *grid) copy(i int, data []byte) {
	buf := make([]int, len(data))
	for i := 0; i < len(data); i++ {
		buf[i] = int(data[i])
	}
	copy(g.d[i+1], append([]int{0}, buf...))
}

func (g *grid) String() string {
	var sb strings.Builder
	for i := 1; i <= g.h; i++ {
		for j := 1; j <= g.w; j++ {
			b := byte(' ')
			if g.d[i][j] != 0 {
				b = byte('0' + (g.d[i][j]-'0')%10) // works for data & labels
			}
			sb.WriteByte(b)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	g := newGrid() // data

	w, h, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := input.Bytes()
		g.copy(h, args) // data ('0'..'9')
		w, h = len(args), h+1
	}
	g.redim(w, h)

	comps := g.components()
	popcnt := make([]int, 0, len(comps))
	for _, pop := range comps {
		popcnt = append(popcnt, pop)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(popcnt)))

	fmt.Println(popcnt[0] * popcnt[1] * popcnt[2])
}
