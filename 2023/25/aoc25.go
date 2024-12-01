package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	g := newGraph(bufio.NewScanner(os.Stdin))

	s := g.furthest(0)
	e := g.furthest(s)

	size := g.flow(s, e)

	fmt.Println(size * (g.len() - size))
}

type graph struct {
	edges []int
	nodes [][2]int
}

type adjmat [][]int

func newGraph(input *bufio.Scanner) *graph {

	// perfect hashtable for ^[a..z]{3}$
	table := make([]int, 26*26*26)
	for i := range table {
		table[i] = MaxInt
	}

	// prealloc adjacency matrix
	A := make(adjmat, 1435) // from previous run
	for i := range A {
		A[i] = make([]int, 0, 8) // from previous run
	}
	A = A[:0:len(A)] // cover and protect rows

	hash := func(s string) int {
		h := 0
		for i := range s {
			h = 26*h + int(s[i]-'a')
		}

		i := table[h]
		if i == MaxInt {
			i = len(A)
			table[h], A = i, A[:i+1] // uncover new row
		}

		return i
	}

	for input.Scan() {
		args := split(input.Text(), ": ")
		src := hash(args[0])

		for _, s := range fields(args[1]) {
			dst := hash(s)
			A[src] = append(A[src], dst)
			A[dst] = append(A[dst], src)
		}
	}

	edges := make([]int, 0, 6416) // from previous run
	nodes := make([][2]int, 0, len(A))

	for _, row := range A {
		s := len(edges)
		e := s + len(row)
		edges = append(edges, row...)
		nodes = append(nodes, [2]int{s, e})
	}

	return &graph{edges, nodes}
}

func (g *graph) len() int {
	return len(g.nodes)
}

func (g *graph) nexts(node int) [][2]int {
	s, e := g.nodes[node][0], g.nodes[node][1]

	nexts := make([][2]int, 0, e-s)
	for i := s; i < e; i++ {
		nexts = append(nexts, [2]int{i, g.edges[i]})
	}

	return nexts
}

func (g *graph) furthest(start int) int {
	seen := make([]bool, len(g.nodes))

	todo := make([]int, 0, 331) // from previous run
	push := func(n int) {
		todo = append(todo, n)
		seen[n] = true
	}
	fpop := func() (n int) {
		n, todo = todo[0], todo[1:]
		return
	}

	push(start)

	far := start
	for len(todo) > 0 {
		cur := fpop()
		far = cur
		for _, x := range g.nexts(cur) {
			nxt := x[1]
			if !seen[nxt] {
				push(nxt)
			}
		}
	}

	return far
}

func (g *graph) flow(s, e int) (size int) {
	type state [2]int

	seen := make([]bool, g.len())

	todo := make([]state, 0, 319) // from previous run
	push := func(x state) {
		todo = append(todo, x)
		seen[x[0]] = true
	}
	fpop := func() (x state) {
		x, todo = todo[0], todo[1:]
		return
	}
	path := make([]state, 0)
	used := make([]bool, len(g.edges))

	reset := func() {
		todo = todo[:0]
		path = path[:0]
		clear(seen)
	}

	for i := 0; i < 5; i++ {
		size = 0
		push(state{s, MaxInt})

	BFS:
		for len(todo) > 0 {
			x := fpop()
			cur, head := x[0], x[1]
			size++

			if cur == e {
				i := head

				for i != MaxInt {
					edge, nxt := path[i][0], path[i][1]
					used[edge] = true
					i = nxt
				}
				break BFS
			}

			for _, x := range g.nexts(cur) {
				edge, nxt := x[0], x[1]
				if !used[edge] && !seen[nxt] {
					push(state{nxt, len(path)})
					path = append(path, state{edge, head})
				}
			}
		}

		reset()
	}

	return size
}

func hash(table *[]int, A *adjmat, s string) int {
	h := 0
	for i := range s {
		h = 26*h + int(s[i]-'a')
	}

	i := (*table)[h]
	if i == MaxInt {
		i = len(*A)
		(*table)[h] = i
		*A = (*A)[:i+1] // uncover row
	}

	return i
}

var split, fields = strings.Split, strings.Fields

const MaxInt = int(^uint(0) >> 1)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
