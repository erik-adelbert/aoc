package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
	"sync"
)

func main() {
	grid := newGrid()

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		grid.load(input.Bytes())
	}

	graph := grid.graph()
	fmt.Println(graph.walk(), graph.hike())
}

type grid struct {
	φ func(j, i int) int
	d []byte
	w int
}

const WMAX = 141

func newGrid() (g *grid) {
	g = new(grid)
	g.w = WMAX
	g.d = make([]byte, 0, WMAX*WMAX)
	g.φ = func(j, i int) int { return j*g.w + i }
	return
}

func (g *grid) load(s []byte) {
	g.w = len(s)
	g.d = append(g.d, s...)
}

func (g *grid) get(y, x int) byte { return g.d[g.φ(y, x)] }

func (g *grid) String() string {
	φ := g.φ
	var sb strings.Builder
	for j := 0; j < g.w; j++ {
		fmt.Fprintln(&sb, string(g.d[φ(j, 0):φ(j+1, 0)]))
	}

	return sb.String()
}

const (
	Y = iota
	X
)

type graph struct {
	φ    func(j, i int) int
	und  []uint64
	dir  []uint64
	cs   []int
	s    int
	e    int
	off  int
	size int
}

func (g *grid) graph() *graph {
	φ, w := g.φ, g.w

	s, e := φ(1, 1), φ(w-2, w-2)
	g.d[s-w], g.d[s] = '#', 'X'
	g.d[e+w], g.d[e] = '#', 'X'

	// junctions
	const size = 36
	J := make(map[int]int, size)
	J[s], J[e] = 0, 1

	Δ := [][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

	for y := 1; y < w-1; y++ {
		for x := 1; x < w-1; x++ {
			if g.d[φ(y, x)] != '#' {
				nlink := 0
				for _, δ := range Δ {
					if g.d[φ(y+δ[Y], x+δ[X])] != '#' {
						nlink++
					}
				}
				if nlink > 2 { // junction
					g.d[φ(y, x)] = 'X'
					J[φ(y, x)] = len(J)
				}
			}
		}
	}

	type state struct {
		i, c int
		fwd  bool
	}

	dir := make([]uint64, size)
	und := make([]uint64, size)
	cs := make([]int, size*size)

	todo := make([]state, 0, 4)
	for start, from := range J {
		todo = append(todo, state{start, 0, true})
		g.d[start] = '#'

		for len(todo) > 0 {
			var s state
			s, todo = todo[0], todo[1:]
			i, w, fwd := s.i, s.c, s.fwd
			y, x := i/g.w, i%g.w
			for θ, δ := range Δ {
				y, x := y+δ[Y], x+δ[X]
				i := φ(y, x)
				switch g.get(y, x) {
				case '#':
				case 'X':
					to := J[i]

					if fwd {
						dir[from] |= 1 << uint64(to)
					} else {
						dir[to] |= 1 << uint64(from)
					}

					und[from] |= 1 << uint64(to)
					und[to] |= 1 << uint64(from)

					cs[from*size+to] = w + 1
					cs[to*size+from] = w + 1
				case '.':
					todo = append(todo, state{i, w + 1, fwd})
					g.d[i] = '#'
				default:
					ds := []byte{'^', '<', '>', 'v'}
					todo = append(todo, state{i, w + 1, fwd && (ds[θ] == g.d[i])})
					g.d[i] = '#'
				}
			}
		}
	}

	s = trail0(und[0])
	e = trail0(und[1])
	off := 2 + cs[s] + cs[size+e]

	mask := 0
	for i, e := range und {
		if popcnt(e) < 4 {
			mask |= 1 << i
		}
	}

	for i, e := range und {
		if popcnt(e) < 4 {
			und[i] = (e & ^uint64(mask)) | dir[i]
		}
	}

	φ = func(j, i int) int { return size*j + i }

	return &graph{φ, und, dir, cs, s, e, off, size}
}

func (g *graph) walk() int {
	cost := make([]int, g.size)
	todo := []int{g.s}

	for len(todo) > 0 {
		var from int
		from, todo = todo[0], todo[1:]
		nodes := g.dir[from]

		for nodes > 0 {
			to := trail0(nodes)
			mask := 1 << to
			nodes ^= uint64(mask)

			cost[to] = max(cost[to], cost[from]+g.cs[from*36+to])
			todo = append(todo, to)
		}
	}
	return cost[g.e] + g.off
}

type seed struct {
	from, cost int
	seen       uint64
}

func (g graph) hike() int {
	cost := 0

	const MAXPROCS = 64

	seeds := []seed{{g.s, 0, 1 << g.s}}
	for len(seeds) < MAXPROCS {
		var x seed

		x, seeds = seeds[0], seeds[1:]
		f, c, s := x.from, x.cost, x.seen
		if f == g.e {
			cost = max(cost, c)
			continue
		}

		nodes := g.und[f] & ^s
		for nodes > 0 {
			to := trail0(nodes)
			mask := uint64(1 << to)
			nodes ^= mask
			seeds = append(seeds, seed{to, c + g.cs[g.φ(f, to)], s | mask})
		}
	}

	var wg sync.WaitGroup
	costs := make(chan int, MAXPROCS)
	for _, s := range seeds {
		wg.Add(1)
		go func(s seed) {
			defer wg.Done()
			f, c, seen := s.from, s.cost, s.seen
			costs <- g.dfs(f, seen) + c
		}(s)
	}

	go func() {
		wg.Wait()
		close(costs)
	}()

	for c := range costs {
		cost = max(cost, c)
	}

	return cost + g.off
}

func (g *graph) dfs(from int, seen uint64) int {
	if from == g.e {
		return 0
	}

	cost := 0
	nodes := g.und[from] & ^seen

	for nodes > 0 {
		to := trail0(nodes)
		mask := uint64(1 << to)
		nodes ^= mask

		cost = max(cost, g.cs[g.φ(from, to)]+g.dfs(to, seen|mask))
	}

	return cost
}

var popcnt, trail0 = bits.OnesCount64, bits.TrailingZeros64
