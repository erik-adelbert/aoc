package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	w := new(world)

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		w.readline(j, input.Text())
	}

	//fmt.Println(w)

	path, area := w.findpath()
	fmt.Println(len(path)/2, area)
}

const MAXN = 140

type world struct {
	O, H, W int
	maze    [MAXN * MAXN]byte
}

func (w *world) readline(j int, line string) {
	w.H = max(w.H, j+1)
	w.W = max(w.W, len(line))

	i := φ(j, 0)
	copy(w.maze[i:], line)

	if i := strings.Index(line, "S"); i > 0 {
		w.O = φ(j, i)
	}
}

func (w *world) findpath() ([]int, int) {
	var old, cur, nxt int

	area := 0
	path := make([]int, 0, 8192)

	// match S with its next pipe
	align := func(x int) int {
		var matchers = []struct {
			dir func(int) int // north | west | south | east
			pat string
		}{
			// ex. go north if north == ('7' | '|' | 'F' )
			{north, "7|F"}, {west, "L-F"}, {south, "J|L"},
		}

		for _, m := range matchers {
			for i := range m.pat {
				if m.pat[i] == w.maze[m.dir(x)] {
					return m.dir(x)
				}
			}
		}
		return east(x) // default to east
	}

	// go fwd without turning around
	fwd := func(x int) int {
		var a, b func(int) int
		switch w.maze[x] {
		case 'J':
			a, b = north, west
		case 'L':
			a, b = north, east
		case '|':
			a, b = north, south
		case '-':
			a, b = east, west
		case '7':
			a, b = south, west
		case 'F':
			a, b = south, east
		}

		// select between a and b
		// prevent turning around
		if a(x) != old {
			return a(x)
		}
		return b(x)
	}

	old, cur, nxt = 0, 0, w.O // 3 window
	for {
		old, cur = cur, nxt
		path = append(path, cur)

		if cur == w.O {
			nxt = align(cur)
		} else {
			nxt = fwd(cur)
		}

		// shoelace formula
		cj, ci := ji(cur)
		nj, ni := ji(nxt)
		area += ci*nj - cj*ni

		if nxt == w.O {
			break
		}
	}

	area = 1 + (area-len(path))/2

	return path, area
}

func (w *world) String() string {
	var sb strings.Builder

	const (
		NS = "│" // pipe connecting north and south
		EW = "─" // pipe connecting east and west
		NE = "╰" // 90-degree connecting north and east
		NW = "╯" // 90-degree connecting north and west
		SW = "╮" // 90-degree connecting south and west
		SE = "╭" // 90-degree connecting south and east
	)
	r := strings.NewReplacer(
		string(byte(0)), ".",
		"-", EW, "|", NS, "L", NE, "J", NW, "7", SW, "F", SE,
	)

	buf := make([]byte, MAXN*MAXN)
	path, _ := w.findpath()
	for i := range path {
		buf[path[i]] = w.maze[path[i]]
	}

	for j := 0; j < w.H; j++ {
		fmt.Fprintln(&sb, r.Replace(string(buf[φ(j, 0):φ(j, w.W)])))
	}

	return sb.String()
}

func φ(j, i int) int {
	return j*MAXN + i
}

func ji(i int) (int, int) {
	return i / MAXN, i % MAXN
}

func north(i int) int {
	return i - MAXN
}

func south(i int) int {
	return i + MAXN
}

func east(i int) int {
	return i + 1
}

func west(i int) int {
	return i - 1
}
