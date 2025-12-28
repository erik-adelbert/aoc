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
	"time"
)

const (
	N          = 1000        // parts threshold
	CutoffDist = 196_000_000 // edge squared distance cutoff from prior runs
)

func main() {
	t0 := time.Now()

	var acc1, acc2 int // parts 1 and 2 accumulators

	points := make([]point, 0, N)

	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ {
		fields := bytes.SplitN(input.Bytes(), []byte(","), 3)

		points = append(points, point{
			X: atoi(fields[0]),
			Y: atoi(fields[1]),
			Z: atoi(fields[2]),
		})
	}

	n := len(points)
	edges := make([]edge, 0, 6*n) // pre-allocate edge slice

	// collect all edges below cutoff
	for i := range n - 1 {
		for j := i + 1; j < n; j++ {
			a, b := points[i], points[j]

			if d := dist2(a, b); d < CutoffDist { // cutoff heuristic
				edges = append(edges, edge{dist: d, i: u32(i), j: u32(j)})
			}
		}
	}

	// do minimal ordering of edges by distance
	// garantees ∀ i < n, ∀ j ≥ n : edges[i].dist ≤ edges[j].dist
	// this is sufficient for Kruskal's algorithm to produce correct MST
	qselect3(edges, len(edges)-1) // partition around k=len(edges)-1
	qselect3(edges, n)            // partition around k=n

	// initialize disjoint set union structure
	dsu := newDSU(n)

	unions := 0 // count of unions performed
	for i, e := range edges {

		if dsu.find(e.i) != dsu.find(e.j) {
			dsu.union(e.i, e.j)
			unions++
		}

		switch {
		case i == N-1: // part 1: after 1000 edges
			// stop as soon as the MST completes, so:
			// no total ordering is required
			// no tie-breaking is required
			// no instability among equal distances is relevant
			seen := make([]bool, n)

			var max1, max2, max3 u32 // sliding top 3 sizes

			for i := range n {
				root := dsu.find(u32(i))

				if !seen[root] {
					seen[root] = true

					switch sz := dsu.size[root]; {
					case sz > max1:
						max3, max2, max1 = max2, max1, sz
					case sz > max2:
						max3, max2 = max2, sz
					case sz > max3:
						max3 = sz
					}
				}
			}

			acc1 = int(max1 * max2 * max3) // product of 3 largest components

		case unions == N-1: // part 2: after 1000 unions, spanning tree is complete
			acc2 = points[e.i].X * points[e.j].X // product of X coords of last edge

			fmt.Println(acc1, acc2, time.Since(t0)) // output results

			return
		}
	}
}

// qselect3 performs a quickselect using median-of-three partitioning
func qselect3(A []edge, k int) {
	l, r := 0, len(A)-1

	for l < r {
		mid := l + (r-l)/2

		// median-of-three
		if A[r].dist < A[l].dist {
			A[l], A[r] = A[r], A[l]
		}
		if A[mid].dist < A[l].dist {
			A[mid], A[l] = A[l], A[mid]
		}
		if A[r].dist < A[mid].dist {
			A[mid], A[r] = A[r], A[mid]
		}

		piv := A[mid].dist
		A[mid], A[r-1] = A[r-1], A[mid]

		i, j := l, r-1
		for {
			for i++; A[i].dist < piv; i++ {
			}
			for j--; A[j].dist > piv; j-- {
			}
			if i >= j {
				break
			}
			A[i], A[j] = A[j], A[i]
		}
		A[i], A[r-1] = A[r-1], A[i]

		// pivot is now at i
		switch {
		case i == k:
			return
		case i < k:
			l = i + 1
		default:
			r = i - 1
		}
	}
}

type u32 = uint32

// edge represents an edge between two points with a squared distance
type edge struct {
	dist int // squared distance
	i, j u32 // point indices
}

// dsu is a disjoint set union (union-find) data structure
type dsu struct {
	parent []u32
	size   []u32
}

// newDSU creates a new DSU with n elements
func newDSU(n int) *dsu {
	p := make([]u32, n)
	for i := range n {
		p[i] = u32(i)
	}

	sz := make([]u32, n)
	for i := range n {
		sz[i] = 1
	}

	return &dsu{parent: p, size: sz}
}

// find returns the root of the set containing x, with path compression
func (d *dsu) find(x u32) u32 {
	root := x

	// find the root
	for d.parent[root] != root {
		root = d.parent[root]
	}

	// path compression: make all nodes on path point directly to root
	for d.parent[x] != x {
		nxt := d.parent[x]
		d.parent[x] = root
		x = nxt
	}

	return root
}

// union merges the sets containing a and b
func (d *dsu) union(a, b u32) {
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
