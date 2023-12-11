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

	path, area := w.path()
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

func (w *world) path() ([]int, int) {
	area := 0
	path := make([]int, 0, 1<<13)

	old, cur, nxt := 0, 0, w.O
	for {
		fwd := func(x int, a, b func(int) int) int {
			p, q := a(x), b(x)
			if p != old { // prevent turning back
				return p
			}
			return q
		}

		old, cur = cur, nxt
		path = append(path, cur)
		switch w.maze[cur] {
		case 'S':
			type matcher struct {
				fun func(int) int
				pat string
			}

			align := func(x int) int {
				var ms = []matcher{
					{north, "7|F"},
					{west, "L-F"},
					{south, "J|L"},
				}

				// ex. go north if north(cur) == '7' or '|' or 'F'
				for _, m := range ms {
					for i := range m.pat {
						if w.maze[m.fun(x)] == m.pat[i] {
							return m.fun(x)
						}
					}
				}
				return east(x) // default to east
			}

			nxt = align(cur)
		case 'J':
			nxt = fwd(cur, north, west)
		case 'L':
			nxt = fwd(cur, north, east)
		case '|':
			nxt = fwd(cur, north, south)
		case '-':
			nxt = fwd(cur, east, west)
		case '7':
			nxt = fwd(cur, south, west)
		case 'F':
			nxt = fwd(cur, south, east)
		}

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
	r := strings.NewReplacer("-", EW, "|", NS, "L", NE, "J", NW, "7", SW, "F", SE)

	for j := 0; j < w.H; j++ {
		fmt.Fprintln(&sb, r.Replace(string(w.maze[φ(j, 0):φ(j, w.W)])))
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
