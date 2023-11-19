// aoc16.go --
// advent of code 2022 day 16
//
// https://adventofcode.com/2022/day/16
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-16: initial commit

package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Inf is max int
const Inf = int(^uint(0) >> 1)

func main() {
	w := mkworld()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		w.add(
			mknode(strings.Fields(r.Replace(input.Text()))),
		)
	}
	w.mkdist()
	w.mkedge()
	w.mkbest()

	fmt.Println(w.astar("AA", 0, 30))
	fmt.Println(w.astar("AA", 26, 26))
}

type (
	graph map[string]*node
	index map[string]int
)

type world struct {
	cave graph
	cidx index
	dist []int
	edge [][]edge
	best [][]wedge
}

func (w *world) astar(start string, t0, t1 int) int {
	best := 0
	cave, cidx, edge := w.cave, w.cidx, w.edge

	// initial state
	i := w.cidx[start]
	s0 := &state{
		prio: Inf,
		bmap: bmp(1 << i), // add start
		edge: [2]int{i, i},
		time: [2]int{t0, t1},
	}

	seen := make(set, 1024)

	heap := make(heap, 0, 1024)
	heap = heap[:0]
	hp.Init(&heap)

	// push initial state
	hp.Push(&heap, s0)
	for heap.Len() > 0 {
		cur := hp.Pop(&heap).(*state)
		if cur.prio <= best {
			return best
		}

		if _, ok := seen[cur]; ok {
			continue
		}

		seen.add(cur)
		for _, x := range edge[cur.edge[0]] {

			i, flow := cidx[x.to], cave[x.to].flo
			if cur.time[1] > x.δt && !cur.bmap.get(i) {
				tn := cur.time[1] - x.δt
				nxt := &state{
					flow: cur.flow + flow*tn,
					bmap: cur.bmap.set(i),
					edge: [2]int{i, cur.edge[1]},
					time: [2]int{cur.time[0], tn},
				}

				if nxt.time[0] > nxt.time[1] {
					nxt.edge[0], nxt.edge[1] = nxt.edge[1], nxt.edge[0]
					nxt.time[0], nxt.time[1] = nxt.time[1], nxt.time[0]
				}

				best = max(best, nxt.flow)
				if prio := nxt.hprio(w); prio > best {
					nxt.prio = prio
					hp.Push(&heap, nxt)
				}
			}
		}
	}

	return best
}

func mkworld() *world {
	w := new(world)
	w.cave = make(graph, 64)
	w.cidx = make(index, 64)
	return w
}

func (w *world) add(s string, n *node) {
	w.cidx[s] = len(w.cave)
	w.cave[s] = n
}

type node struct {
	nxt []string
	flo int
}

func mknode(line []string) (string, *node) {
	x := new(node)
	x.flo = atoi(line[5])
	x.nxt = make([]string, len(line[10:]))
	copy(x.nxt, line[10:])
	return line[1], x
}

// compute all pairs travel times
// floyd-warshall flooding
func (w *world) mkdist() {
	N := len(w.cave)

	dist := make([]int, N*N)
	for j := 0; j < N; j++ {
		for i := 0; i < N; i++ {
			dist[j*N+i] = Inf
		}
	}

	cave, cidx := w.cave, w.cidx
	for x, i := range cidx {
		dist[i*N+i] = 0
		for _, y := range cave[x].nxt {
			dist[cidx[y]*N+i] = 1
		}
	}
	for k := 0; k < N; k++ {
		for j := 0; j < N; j++ {
			for i := 0; i < N; i++ {
				if v := dist[k*N+i] + dist[j*N+k]; v > 0 {
					dist[j*N+i] = min(dist[j*N+i], v)
				}
			}
		}
	}

	w.dist = dist
}

type edge struct {
	to string
	δt int
}

func (w *world) mkedge() {
	cave, cidx, dist := w.cave, w.cidx, w.dist
	N := len(cave)

	edges := make([][]edge, N)
	for a, i := range cidx {
		edges[i] = make([]edge, 0, N)
		for b, j := range cidx {
			d := dist[j*N+i]
			if a == "AA" || dist[j*N+i] != Inf &&
				cave[a].flo*cave[b].flo > 0 {
				edges[i] = append(edges[i], edge{b, d + 1})
			}
		}
		sort.Slice(edges[i], func(x, y int) bool {
			f := func(x int) int {
				return cave[edges[i][x].to].flo
			}
			return f(x) > f(y)
		})
	}

	w.edge = edges
}

// weighted edge for heuristic search
type wedge struct {
	i, δt, v int // index, weight, value
}

func (w *world) mkbest() {
	cave, cidx, dist := w.cave, w.cidx, w.dist
	N := len(cave)

	best := make([][]wedge, 31)
	for t := range best {
		var (
			a, b string
			i, j int
		)
		for a, i = range cidx {
			w := Inf
			for b, j = range cidx {
				if i == j || cave[b].flo == 0 {
					continue
				}
				w = min(w, dist[j*N+i]+1)
			}
			if w < t {
				best[t] = append(best[t], wedge{i, w, cave[a].flo})
			}
		}
		sort.Slice(best[t], func(i, j int) bool {
			f := func(i int) int {
				return best[t][i].v * (t - best[t][i].δt)
			}
			return f(i) > f(j)
		})
	}
	w.best = best
}

// path head
type state struct {
	prio int

	bmap bmp
	flow int
	// time range
	time [2]int
	edge [2]int
}

// estimate maxflow
// heuristic is to go all the way down to clock = 0
// no matter if nodes are available or not
func (x *state) hprio(w *world) int {
	// unpack namespace
	bmap := x.bmap
	prio := x.flow
	best := w.best

	t0, t1 := x.time[0], x.time[1]

walk:
	for {
		edges := best[t1]
		for _, edge := range edges {
			// walk biggest value node at t1
			if !bmap.get(edge.i) {
				// update flow path and elapsed time
				t1 -= edge.δt
				prio += edge.v * t1
				if t0 > t1 {
					t0, t1 = t1, t0
				}
				bmap = bmap.set(edge.i)
				continue walk
			}
		}
		// done walking
		return prio
	}
}

type set map[*state]struct{}

func (s set) add(x *state) {
	s[x] = struct{}{}
}

// standard lib binary heap
type heap []*state

func (h heap) Len() int           { return len(h) }
func (h heap) Less(i, j int) bool { return h[i].prio < h[j].prio }
func (h heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *heap) Push(x any) {
	c := x.(*state)
	*h = append(*h, c)
}

func (h *heap) Pop() any {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}

// node setup bitmap for heuristics
type bmp uint64

func (b bmp) get(i int) bool {
	x := bmp(1) << i
	return b&x == x
}

func (b bmp) set(i int) bmp {
	return b | bmp(1)<<i
}

func (b bmp) clr(i int) bmp {
	return b & ^(bmp(1) << i)
}

var r = strings.NewReplacer(
	"=", " ",
	";", "",
	",", "",
)

// strconv.Atoi modified core loop
// s is ^\d+.*$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

var DEBUG = false

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
