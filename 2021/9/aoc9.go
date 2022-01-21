package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type grid struct {
	d    [128][128]int
	h, w int
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
}

func (g *grid) copy(i int, data []byte) int {
	t := g.d[i+1][1:]

	for i, b := range data {
		t[i] = int(b)
	}
	return len(data)
}

func (g *grid) set(y, x, v int) {
	if !(y >= 0 && x >= 0 && g.d[y+1][x+1] != '9') {
		return
	}
	g.d[y+1][x+1] = v
}

func (g *grid) get(y, x int) int {
	if !(y >= 0 && x >= 0) {
		return 0
	}
	return g.d[y+1][x+1]
}

func (g *grid) filter(y, x int) int {
	btoi := func(b int) int {
		return int(b - '0') // fast convert
	}

	neighbors := func(y, x int) []int {
		return []int{
			g.d[y][x+1], g.d[y+1][x], g.d[y+2][x+1], g.d[y+1][x+2],
		}
	}

	v := g.get(y, x)
	for _, n := range neighbors(y, x) {
		if 0 < n && n <= v {
			return 0
		}
	}
	return 1 + btoi(v)
}

func (g *grid) groups() map[int]int {
	label := new(grid)
	label.redim(g.h, g.w)

	labels := make([]int, g.w*g.h+256)
	for i := range labels {
		labels[i] = i
	}

	find := func(a int) int {
		for labels[a] != a {
			a, labels[a] = labels[a], labels[labels[a]] // path splitting
		}
		return a
	}

	union := func(a, b int) {
		if a < b {
			labels[a] = b
		} else {
			labels[b] = a
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

	groups := make(map[int]int)
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if v := label.get(y, x); v != 0 {
				groups[find(v)]++
			}
		}
	}
	return groups
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
	g := new(grid) // data

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		data := input.Bytes()
		w = g.copy(h, data) // data ('0'..'9')
		h++
	}
	g.redim(h, w)

	sum := 0
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			sum += g.filter(y, x)
		}
	}
	fmt.Println(sum) // part1

	popcnts := values(g.groups())
	sort.Sort(sort.Reverse(sort.IntSlice(popcnts)))
	fmt.Println(popcnts[0] * popcnts[1] * popcnts[2]) // part2
}

func values(m map[int]int) []int {
	vals := make([]int, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}
