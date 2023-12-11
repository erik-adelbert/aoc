package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	w := newWorld()

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

func newWorld() *world {
	w := new(world)
	return w
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
		choose := func(p, q int) int {
			if p != old {
				return p
			}
			return q
		}

		old, cur = cur, nxt
		path = append(path, cur)
		switch w.maze[cur] {
		case 'S':
			branch := func(x int, s string) bool {
				for i := range s {
					if w.maze[x] == s[i] {
						return true
					}
				}
				return false
			}

			switch {
			case branch(north(cur), "7|F"):
				nxt = north(cur)
			case branch(west(cur), "L-F"):
				nxt = west(cur)
			case branch(south(cur), "J|L"):
				nxt = south(cur)
			default:
				nxt = east(cur)
			}
		case 'F':
			nxt = choose(east(cur), south(cur))
		case '|':
			nxt = choose(north(cur), south(cur))
		case 'L':
			nxt = choose(north(cur), east(cur))
		case '-':
			nxt = choose(west(cur), east(cur))
		case 'J':
			nxt = choose(west(cur), north(cur))
		case '7':
			nxt = choose(west(cur), south(cur))
		}

		cj, ci := ji(cur)
		nj, ni := ji(nxt)
		area += ci*nj - cj*ni

		if nxt == w.O {
			break
		}
	}

	area = 1 + (abs(area)-len(path))/2

	return path, area
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
