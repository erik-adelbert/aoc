package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	g := newGrid()

	input := bufio.NewScanner(os.Stdin)
	var j int
	for j = 0; input.Scan(); j++ {
		input := input.Bytes()
		copy(g.d[j*len(input):], input)
	}
	g.redim(j)

	// part1 trace from first cell heading right
	part1 := g.trace(cell{0, R})

	// part2 init
	part2, n := part1, 4*g.w-2

	// build rays
	rays := make([]cell, 0, n)
	rays = append(rays, []cell{{0, D}, {j*j - 1, U}}...) // first and last column
	for i := 1; i < g.w; i++ {
		// borders
		rays = append(rays, []cell{{i * j, R}, {i*(j+1) - 1, L}, {i, D}, {j*j - 1 - i, U}}...)
	}

	// parallel raytracing
	var wg sync.WaitGroup
	wg.Add(n)

	traces := make(chan int, n)

	// distribute
	for _, r := range rays {
		go func(r cell) {
			traces <- g.trace(r)
			wg.Done()
		}(r)
	}

	// barrier
	go func() {
		wg.Wait()
		close(traces)
	}()

	// collect
	for t := range traces {
		part2 = max(part2, t)
	}

	fmt.Println(part1, part2)
}

const (
	L = 1 << iota
	D
	R
	U
)

type cell struct {
	i int
	θ byte
}

type move func(int, int) (int, bool)
type test func(int, int) bool

var inbounds = []test{
	L: func(i, w int) bool { return (i-1) >= 0 && (i-1)/w == i/w },
	D: func(i, w int) bool { return (i + w) <= w*w-1 },
	R: func(i, w int) bool { return (i+1) <= w*w-1 && (i+1)/w == i/w },
	U: func(i, w int) bool { return (i - w) >= 0 },
}

var moves = []move{
	L: func(i, w int) (int, bool) { return i - 1, inbounds[L](i, w) },
	D: func(i, w int) (int, bool) { return i + w, inbounds[D](i, w) },
	R: func(i, w int) (int, bool) { return i + 1, inbounds[R](i, w) },
	U: func(i, w int) (int, bool) { return i - w, inbounds[U](i, w) },
}

func (c cell) deflect(g *grid) cell {
	i, θ, gi, w := c.i, c.θ, g.d[c.i], g.w

	deflect := [][]byte{
		'/':  {D: L, L: D, U: R, R: U},
		'\\': {U: L, R: D, D: R, L: U},
	}

	move := func(θ byte) cell {
		if i, ok := moves[θ](i, w); ok {
			return cell{i, θ}
		}
		return cell{}
	}

	return move(deflect[gi][θ])
}

func (c cell) split(g *grid) []cell {
	cells := make([]cell, 0, 2)

	i, θ, gi, w := c.i, c.θ, g.d[c.i], g.w

	nexts := []byte{θ}
	switch {
	case (θ == L || θ == R) && gi == '|':
		nexts = []byte{U, D}
	case (θ == U || θ == D) && gi == '-':
		nexts = []byte{L, R}
	}

	for _, θ := range nexts {
		if i, ok := moves[θ](i, w); ok {
			cells = append(cells, cell{i, θ})
		}
	}

	return cells
}

const MAX = 110

type grid struct {
	d []byte
	w int
}

func newGrid() *grid {
	var g grid
	g.d = make([]byte, MAX*MAX)
	g.w = MAX
	return &g
}

func (g *grid) get(c cell) byte {
	return g.d[c.i]
}

func (g *grid) redim(w int) {
	g.d = g.d[:w*w]
	g.w = w
}

func (g *grid) trace(start cell) int {
	// cycle detection presence map
	T := make([]byte, g.w*g.w)

	seen := func(c cell) bool {
		i, θ := c.i, c.θ
		if T[i]&θ == θ {
			return true
		}
		T[i] |= θ
		return false
	}

	// raytrace segment
	ray := func(c cell) cell {
		var (
			cur, nxt, w = c.i, c.i, g.w
			step        = moves[c.θ]
			ok          = false
		)
		for {
			if nxt, ok = step(cur, w); !ok {
				return cell{cur, c.θ}
			}

			if g.d[nxt] != '.' || seen(cell{nxt, c.θ}) {
				break
			}

			cur = nxt
		}

		return cell{nxt, c.θ}
	}

	// stack-based raytracing over grid
	stack := make([]cell, 0, 1024)

	push := func(x ...cell) { stack = append(stack, x...) }

	pop := func() cell {
		i := len(stack) - 1
		pop := stack[i]
		stack, stack[i] = stack[:i], cell{}
		return pop
	}

	push(start)
	for len(stack) > 0 {
		if x := pop(); !seen(x) {
			switch g.get(x) {
			case '/', '\\':
				if x := x.deflect(g); x != (cell{}) {
					push(x)
				}
			case '-', '|':
				if x := x.split(g); len(x) > 0 {
					push(x...)
				}
			default:
				push(ray(x))
			}
		}
	}

	n := 0
	for i := range T {
		if T[i] > 0 {
			n++
		}
	}

	return n
}

func (g *grid) String() string {
	var sb strings.Builder
	for j := 0; j < g.w; j++ {
		fmt.Fprintln(&sb, string(g.d[j*g.w:(j+1)*g.w]))
	}
	return sb.String()
}
