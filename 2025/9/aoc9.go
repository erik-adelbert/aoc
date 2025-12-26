package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
	"time"
)

type (
	// tile represents a point (x, y)
	tile = [2]u32

	// span represents a 1D interval [l, r]
	span struct {
		l, r u32
	}

	// rect represents a candidate rectangle generator with top-left corner (x, y) and horizontal span
	rect struct {
		x, y u32
		span
	}
)

// inter returns the intersection of two spans
func (a span) inter(b span) span {
	return span{
		l: max(a.l, b.l),
		r: min(a.r, b.r),
	}
}

// contains checks if the span contains the point i
func (a span) contains(i u32) bool {
	return a.l <= i && i <= a.r
}

func main() {
	t0 := time.Now()

	var acc1, acc2 u64 // parts 1 and 2 accumulators

	tiles := make([]tile, 0, SizeHint)

	// read input points
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		bufX, bufY, _ := bytes.Cut(input.Bytes(), []byte(","))
		X, Y := atoi(bufX), atoi(bufY)

		tiles = append(tiles, tile{X, Y})
	}

	// sort tiles by y, then x
	slices.SortFunc(tiles, func(a, b tile) int {
		ya, yb := int(a[1]), int(b[1])
		if ya != yb {
			return ya - yb
		}

		xa, xb := int(a[0]), int(b[0])
		return xa - xb
	})

	acc1 = part1(tiles)
	acc2 = part2(tiles)

	fmt.Println(acc1, acc2, time.Since(t0))
}

// part1 computes the maximum rectangle area using left and right tops and bottoms
func part1(tiles []tile) u64 {
	selit := reverse(tiles)

	ltops, rtops := limits(tiles) // left and right tops
	lbots, rbots := limits(selit) // left and right bottoms

	return max(area(ltops, rbots), area(rtops, lbots))
}

// part2 computes the maximum rectangle area using a sweep line algorithm
func part2(tiles []tile) u64 {
	var area u64 = 0

	edges := make([]u32, 0, 4)  // descending edges
	spans := make([]span, 0, 8) // current intervals

	cur := make([]rect, 0, 512) // current candidates

	for i := 0; i < len(tiles); i += 2 {
		// Tiles are assumed to come in pairs on the same y line
		x0, y, x1 := tiles[i][0], tiles[i][1], tiles[i+1][0]

		// Toggle x values in the descending edge list
		edges = toggle(edges, x0)
		edges = toggle(edges, x1)

		// Compute intervals from descending edges
		spans = toSpans(spans, edges)

		// Check rectangles with current candidates
		for _, c := range cur {
			if c.contains(x0) {
				w := adiff(c.x, x0) + 1
				h := adiff(c.y, y) + 1
				area = max(area, w*h)
			}

			if c.contains(x1) {
				w := adiff(c.x, x1) + 1
				h := adiff(c.y, y) + 1
				area = max(area, w*h)
			}
		}

		// Shrink or remove candidates
		nxt := cur[:0]
		for _, c := range cur {
			ok := false
			for _, s := range spans {
				if ok = s.contains(c.x); ok {
					c.span = c.inter(s)
					break
				}
			}

			if ok {
				nxt = append(nxt, c)
			}
		}
		cur = nxt

		// Add new candidates
		for _, x := range []u32{x0, x1} {
			for _, s := range spans {
				if s.contains(x) {
					cur = append(cur, rect{x: x, y: y, span: s})
					break
				}
			}
		}
	}

	return area
}

// limits returns leftmost and rightmost tiles for each y level
func limits(tiles []tile) (left, right []tile) {
	left = make([]tile, 0, len(tiles)/2)
	right = make([]tile, 0, len(tiles)/2)

	last := func(s []tile) tile {
		return s[len(s)-1]
	}

	for i := 0; i < len(tiles); {
		y := tiles[i][1]
		xmin, xmax := tiles[i][0], tiles[i][0]

		j := i + 1
		for ; j < len(tiles) && tiles[j][1] == y; j++ {
			xmin = min(xmin, tiles[j][0])
			xmax = max(xmax, tiles[j][0])
		}

		if len(left) == 0 || xmin <= last(left)[0] {
			left = append(left, tile{xmin, y})
		}

		if len(right) == 0 || xmax >= last(right)[0] {
			right = append(right, tile{xmax, y})
		}

		i = j
	}
	return
}

// reverse returns a reversed copy of the input slice
func reverse(tiles []tile) []tile {
	rev := slices.Clone(tiles)
	slices.Reverse(rev)

	return rev
}

// area computes the maximum area from pairs of tiles in a and b
func area(a, b []tile) (best u64) {
	adiff := func(a, b u32) u64 {
		if a > b {
			return u64(a - b)
		}
		return u64(b - a)
	}

	for _, p := range a {
		for _, q := range b {
			w := adiff(p[0], q[0]) + 1
			h := adiff(p[1], q[1]) + 1
			best = max(best, w*h)
		}
	}

	return
}

// toggle inserts or removes v from sorted slice A
func toggle(A []u32, v u32) []u32 {
	i, ok := slices.BinarySearch(A, v)

	if ok {
		// Remove v at index i
		copy(A[i:], A[i+1:])
		return A[:len(A)-1]
	}

	// Insert v at index i
	A = append(A, 0)
	copy(A[i+1:], A[i:]) // shift right
	A[i] = v

	return A
}

// toSpans converts descending edges [l0, r0, l1, r1, ...] into intervals
func toSpans(buf []span, edges []u32) []span {
	buf = buf[:0]
	for i := 0; i < len(edges); i += 2 {
		buf = append(buf, span{l: edges[i], r: edges[i+1]})
	}
	return buf
}

// adiff computes the absolute difference between two u32 values
func adiff(a, b u32) u64 {
	if a > b {
		return u64(a - b)
	}
	return u64(b - a)
}

// atoi converts a byte slice representing a non-negative integer to uint32
func atoi(s []byte) (n u32) {
	for _, c := range s {
		n = 10*n + u32(c-'0')
	}

	return
}

const SizeHint = 497

type (
	u64 = uint64
	u32 = uint32
)
