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
	h, w int
}

func newGrid() *grid {
	d := make([][]int, 128)
	for i := 0; i < 128; i++ {
		d[i] = make([]int, 128)
	}
	return &grid{d, 128, 128}
}

func (g *grid) set(y, x, v int) {
	if y < 0 || x < 0 || g.d[y+1][x+1] == '9' {
		return
	}
	g.d[y+1][x+1] = v
}

func (g *grid) get(y, x int) int {
	if x < 0 || y < 0 {
		return 0
	}
	return g.d[y+1][x+1]
}

func (g *grid) comps() map[int]int {
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

	union := func(y, x int) {
		if x > y {
			labels[y] = x
		} else {
			labels[x] = y
		}
	}

	id := 256                  // labels (>256)
	for y := 0; y < g.h; y++ { // Hoshen-Kopelman
		for x := 0; x < g.w; x++ {
			if g.get(y, x) == '9' {
				continue
			}
			nor, wes := label.get(y-1, x), label.get(y, x-1)
			switch {
			case nor == 0 && wes == 0:
				label.set(y, x, id)
				id++
			case nor != 0 && wes == 0:
				label.set(y, x, find(nor))
			case nor == 0 && wes != 0:
				label.set(y, x, find(wes))
			default:
				union(nor, wes)
				label.set(y, x, find(nor))
			}
		}
	}

	comps := make(map[int]int)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if v := label.get(y, x); v != 0 {
				comps[find(v)]++
			}
		}
	}
	return comps
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
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
	for j := 1; j <= g.h; j++ {
		for i := 1; i <= g.w; i++ {
			b := byte(' ')
			if g.d[j][i] != 0 {
				b = byte('0' + (g.d[j][i]-'0')%10) // works for data & labels
			}
			sb.WriteByte(b)
		}
		if j != g.h {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	g := newGrid() // data

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := input.Bytes()
		g.copy(h, args) // data ('0'..'9')
		h, w = h+1, len(args)
	}
	g.redim(h, w)

	comps := g.comps()
	popcnt := make([]int, 0, len(comps))
	for _, pop := range comps {
		popcnt = append(popcnt, pop)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(popcnt)))

	fmt.Println(popcnt[0] * popcnt[1] * popcnt[2])
}
