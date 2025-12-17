package main

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"os"
	"slices"
	"time"
)

func main() {
	t0 := time.Now()

	var acc1, acc2 int // parts 1 and 2 accumulators

	path := make([]point, 0, SizeHint)

	Xs := make([]uint32, 0, 2*SizeHint)
	Ys := make([]uint32, 0, 2*SizeHint)

	// read input points
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bufX, bufY, _ := bytes.Cut(input.Bytes(), []byte(","))
		X, Y := atoi(bufX), atoi(bufY)

		path = append(path, point{X: X, Y: Y})

		Xs = append(Xs, X, X+1) // add +1 for edge handling
		Ys = append(Ys, Y, Y+1)
	}

	// compute coordinate sets
	slices.Sort(Xs)
	slices.Sort(Ys)

	Xs = slices.Compact(Xs)
	Ys = slices.Compact(Ys)

	// compact coordinates mapping
	xmax, ymax := Xs[len(Xs)-1], Ys[len(Ys)-1]

	xmap := make([]uint32, xmax+1)
	ymap := make([]uint32, ymax+1)

	for i, x := range Xs {
		xmap[x] = uint32(i)
	}

	for i, y := range Ys {
		ymap[y] = uint32(i)
	}

	// scanline fill
	R, C := uint32(len(Xs)), uint32(len(Ys))
	edges := make([]uint32, R*C)

	for i1 := range path {
		x1, y1 := path[i1].X, path[i1].Y

		i2 := (i1 + 1) % len(path)
		x2, y2 := path[i2].X, path[i2].Y

		x1, x2 = xmap[x1], xmap[x2]
		if x1 > x2 {
			x1, x2 = x2, x1
		}

		y1, y2 = ymap[y1], ymap[y2]

		edges[x1*C+y1] |= Left
		edges[x2*C+y1] |= Right

		for x := x1 + 1; x < x2; x++ {
			edges[x*C+y1] |= OnEdge
		}
	}

	// build prefix sums
	bmask := make([]uint8, R*C)
	for i := range R {
		var inside uint8

		for j := range C {
			cell := uint8(edges[i*C+j])

			if inside > 0 || cell > 0 {
				bmask[i*C+j] = 1
			}

			inside ^= cell
		}
	}

	θ := func(i, j uint32) uint32 { return uint32(i)*(C+1) + uint32(j) } // linear index

	sums := make([]uint32, (R+1)*(C+1))
	for i := range R {
		for j := range C {
			n := uint32(bmask[i*C+j])

			sums[θ(i+1, j+1)] = n + sums[θ(i, j+1)] + sums[θ(i+1, j)] - sums[θ(i, j)]
		}
	}

	// evaluate all rectangles
	for i, j := range allIndexPairs(path) {
		var x1, x2, y1, y2 uint32

		if x1, x2 = path[i].X, path[j].X; x1 > x2 {
			x1, x2 = x2, x1
		}

		if y1, y2 = path[i].Y, path[j].Y; y1 > y2 {
			y1, y2 = y2, y1
		}

		δx, δy := int(x2-x1), int(y2-y1)
		area := (δx + 1) * (δy + 1)

		// remap to compressed coordinates
		x1, x2 = xmap[x1], xmap[x2]+1
		y1, y2 = ymap[y1], ymap[y2]+1

		// part 2
		all := (x2 - x1) * (y2 - y1) // total area

		// count insiders using prefix sums
		insiders := sums[θ(x2, y2)] - sums[θ(x2, y1)] - sums[θ(x1, y2)] + sums[θ(x1, y1)]

		// part 1
		acc1 = max(acc1, area)

		// part 2
		if all == insiders {
			acc2 = max(acc2, area)
		}
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

const SizeHint = 497

const (
	Left = 1 + iota
	Right
	OnEdge
)

type point struct {
	X, Y uint32
}

// allIndexPairs yields all unique pairs of indices (i, j) with i < j
func allIndexPairs(pts []point) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		n := len(pts)

		for i := range pts[:n-1] {
			for j := i + 1; j < n; j++ {
				if !yield(i, j) {
					return
				}
			}
		}
	}
}

// atoi converts a byte slice representing a non-negative integer to int
func atoi(s []byte) (n uint32) {
	for _, c := range s {
		n = 10*n + uint32(c-'0')
	}

	return
}
