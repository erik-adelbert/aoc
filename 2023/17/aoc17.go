// aoc17.go --
// advent of code 2022 day 17
//
// https://adventofcode.com/2022/day/17
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-17: initial commit

package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"os"
	"strings"
)

func main() {
	world := newGrid()

	h, w, input := 0, 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		input := input.Bytes()
		w = len(input)
		for i, b := range input {
			world.d[(h*w)+i] = b - '0'
		}
		w = len(input)
		h++
	}
	world.redim(h, w)

	g1 := NewGraph(world)
	g2 := g1.clone()

	fmt.Println(g1.shortest(1, 3), g2.shortest(4, 10))
}

type yx struct {
	y, x int
}

type grid struct {
	d    []byte
	h, w int
}

const MAXN = 141

func newGrid() *grid {
	H, W, g := MAXN, MAXN, new(grid)
	g.d = make([]byte, H*W)
	return g
}

func (g *grid) redim(h, w int) {
	g.h, g.w = max(g.h, h), max(g.w, w)
}

func (g *grid) String() string {
	var sb strings.Builder

	for j := 0; j < g.h; j++ {
		fmt.Fprintln(&sb, g.d[j*g.w:(j+1)*g.w])
	}

	return sb.String()
}

const (
	V = iota + 1
	H = 1 << iota
)

type node struct {
	yx
	θ int

	loss int

	move int
	best int

	prio int
}

type graph struct {
	nodes []node
	h, w  int
}

func (g *graph) nexts(u *node, min, max int) []*node {
	nexts := make([]*node, 0, 6)

	if u.θ&H > 0 {
		L, R := 0, 0
		for δy := 1; δy <= max; δy++ {
			if v := g.node(u.y-δy, u.x, V); v != nil {
				L += v.loss
				if δy >= min {
					v.move = L
					nexts = append(nexts, v)
				}
			}
			if v := g.node(u.y+δy, u.x, V); v != nil {
				R += v.loss
				if δy >= min {
					v.move = R
					nexts = append(nexts, v)
				}
			}
		}
	}

	if u.θ&V > 0 {
		L, R := 0, 0
		for δx := 1; δx <= max; δx++ {
			if v := g.node(u.y, u.x+δx, H); v != nil {
				L += v.loss
				if δx >= min {
					v.move = L
					nexts = append(nexts, v)
				}
			}
			if v := g.node(u.y, u.x-δx, H); v != nil {
				R += v.loss
				if δx >= min {
					v.move = R
					nexts = append(nexts, v)
				}
			}
		}
	}

	return nexts
}

func (g *graph) node(y, x, θ int) *node {
	if y < 0 || y >= g.h || x < 0 || x >= g.w/2 {
		return nil
	}

	return &g.nodes[y*g.w+x*2+(θ-1)]
}

func NewGraph(g *grid) graph {
	gr := graph{}
	gr.nodes = make([]node, g.h*g.w*2)
	gr.h, gr.w = g.h, g.w*2

	for i := range g.d[:g.h*g.w] {
		ii := 2 * i
		gr.nodes[ii] = node{
			yx:   yx{i / g.w, i % g.w},
			θ:    V,
			best: MaxInt,
			loss: int(g.d[i]),
		}

		gr.nodes[ii+1] = node{
			yx:   yx{i / g.w, i % g.w},
			θ:    H,
			best: MaxInt,
			loss: int(g.d[i]),
		}

	}

	return gr
}

func (g *graph) clone() (gc *graph) {
	gc = new(graph)

	*gc = *g
	gc.nodes = make([]node, len(g.nodes))
	copy(gc.nodes, g.nodes)

	return
}

func (g *graph) shortest(min, max int) int {
	nodes := g.nodes

	// start node
	nodes[0].best = 0
	nodes[0].θ = V | H

	// priority queue
	heap := make(heap, len(nodes))
	for i := 0; i < len(nodes); i++ {
		nodes[i].prio = i
		heap[i] = &nodes[i]
	}
	hp.Init(&heap)

	var end = &nodes[len(nodes)-1]
	var u *node
	for {
		u = hp.Pop(&heap).(*node)

		// stop short at end node
		if u.x == end.x && u.y == end.y {
			break
		}

		for _, v := range g.nexts(u, min, max) {
			if u.best+v.move < v.best {
				v.best = u.best + v.move
				heap.update(v)
			}
		}
	}

	return u.best
}

type heap []*node

func (h heap) Len() int { return len(h) }

func (h heap) Less(i, j int) bool {
	return h[i].best < h[j].best
}

func (h heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].prio, h[j].prio = i, j
}

func (h *heap) Push(x any) {
	node := x.(*node)
	node.prio = len(*h)
	*h = append(*h, node)
}

func (h *heap) Pop() any {
	q, i := *h, len(*h)-1
	pop := q[i]
	*h, q[i] = q[:i], nil
	return pop
}

func (h *heap) update(n *node) {
	hp.Fix(h, n.prio)
}

const MaxInt = int(^uint(0) >> 1)

var r = strings.NewReplacer(
	fmt.Sprint(MaxInt), "+∞",
)
