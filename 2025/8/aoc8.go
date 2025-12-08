// aoc8.go --
// advent of code 2025 day 8
//
// https://adventofcode.com/2025/day/8
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-8: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
	"time"
)

const CutoffDist = 195_000_000 // edge squared distance cutoff from prior runs

func main() {
	t0 := time.Now()

	var acc1, acc2 int // parts 1 and 2 accumulators

	points := make([]point, 0, 1000)

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		b := bytes.Split(input.Bytes(), []byte(","))
		points = append(points, point{
			X: atoi(b[0]),
			Y: atoi(b[1]),
			Z: atoi(b[2]),
		})
	}

	n := len(points)
	edges := make([]edge, 0, 6*n) // pre-allocate edge slice
	// collect all edges below cutoff
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			a, b := points[i], points[j]
			if d := dist2(a, b); d < CutoffDist {
				edges = append(edges, edge{dist: d, a: i, b: j})
			}
		}
	}
	// sort edges by distance
	slices.SortFunc(edges, func(a, b edge) int { return a.dist - b.dist })

	unions := 0
	dsu := newDSU(n)
	for i, e := range edges {

		if dsu.find(e.a) != dsu.find(e.b) {
			dsu.union(e.a, e.b)
			unions++
		}

		switch {
		case i == n-1: // part 1: after 1000 edges
			seen := make([]bool, n)
			sizes := make([]int, 0, n)

			for i := range n {
				root := dsu.find(i)
				if !seen[root] {
					seen[root] = true
					sizes = append(sizes, dsu.size[root])
				}
			}

			slices.SortFunc(sizes, func(a, b int) int { return b - a }) // reverse sort sizes
			acc1 = sizes[0] * sizes[1] * sizes[2]                       // product of 3 largest components

		case unions == n-1: // part 2: after 1000 unions, spanning tree is complete
			acc2 = int(points[e.a].X) * int(points[e.b].X) // product of X coords of last edge
			fmt.Println(acc1, acc2, time.Since(t0))        // output results

			return
		}
	}
}

// edge represents an edge between two points with a squared distance
type edge struct {
	dist int // squared distance
	a, b int
}

// dsu is a disjoint set union (union-find) data structure
type dsu struct {
	parent []int
	size   []int
}

// newDSU creates a new DSU with n elements
func newDSU(n int) *dsu {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &dsu{p, sz}
}

// find returns the root of the set containing x, with path compression
func (d *dsu) find(x int) int {
	root := x
	// Find the root
	for d.parent[root] != root {
		root = d.parent[root]
	}
	// Path compression: make all nodes on path point directly to root
	for d.parent[x] != x {
		next := d.parent[x]
		d.parent[x] = root
		x = next
	}
	return root
}

// union merges the sets containing a and b
func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

// point represents a 3D point
type point struct{ X, Y, Z int }

// dist2 returns the squared distance between points a and b
func dist2(a, b point) int {
	dx, dy, dz := a.X-b.X, a.Y-b.Y, a.Z-b.Z
	return dx*dx + dy*dy + dz*dz
}

// atoi converts a byte slice representing a non-negative integer to int
func atoi(s []byte) (n int) {
	for _, c := range s {
		n = 10*n + int(c-'0')
	}

	return
}
