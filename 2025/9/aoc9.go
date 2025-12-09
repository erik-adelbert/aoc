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

	var acc1, acc2 uint32 // parts 1 and 2 accumulators

	points := make([]point, 0, SizeHint)

	xraw := make([]uint32, 2*SizeHint)
	yraw := make([]uint32, 2*SizeHint)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bufX, bufY, _ := bytes.Cut(input.Bytes(), []byte(","))
		X, Y := atoi(bufX), atoi(bufY)

		points = append(points, point{X: X, Y: Y})

		xraw = append(xraw, X, X+1) // add +1 for edge handling
		yraw = append(yraw, Y, Y+1)
	}

	// fmt.Println("Parsed input...", time.Since(t0))

	slices.Sort(xraw)
	slices.Sort(yraw)

	// fmt.Println("Sorted coordinates...", time.Since(t0))

	// coordinate sets
	Xs := slices.Compact(xraw)
	Ys := slices.Compact(yraw)

	xmax, ymax := Xs[len(Xs)-1], Ys[len(Ys)-1]

	// compact coordinates mapping
	xmap := make([]uint32, xmax+1)
	ymap := make([]uint32, ymax+1)

	for i, x := range Xs {
		xmap[x] = uint32(i)
	}

	for i, y := range Ys {
		ymap[y] = uint32(i)
	}

	// fmt.Println("Built coordinate maps...", time.Since(t0))

	// scanline fill
	R, C := uint32(len(Xs)), uint32(len(Ys))
	edges := make([]uint32, R*C)

	for i := range points {
		x1, y1 := points[i].X, points[i].Y
		x2, y2 := points[(i+1)%len(points)].X, points[(i+1)%len(points)].Y

		x1, x2 = xmap[x1], xmap[x2]
		x1, x2 = min(x1, x2), max(x1, x2)

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

	// fmt.Println("Completed scanline fill...", time.Since(t0))

	θ := func(i, j uint32) int { return int(i*(C+1) + j) } // linear index

	sums := make([]uint32, (R+1)*(C+1))
	for i := range R {
		for j := range C {
			n := uint32(bmask[i*C+j])

			sums[θ(i+1, j+1)] = n + sums[θ(i, j+1)] + sums[θ(i+1, j)] - sums[θ(i, j)]
		}
	}

	// fmt.Println("Built prefix sum table...", time.Since(t0))

	// evaluate all rectangles
	for i, j := range allIndexPairs(points) {
		x1, x2 := points[i].X, points[j].X
		x1, x2 = min(x1, x2), max(x1, x2)

		y1, y2 := points[i].Y, points[j].Y
		y1, y2 = min(y1, y2), max(y1, y2)

		area := (x2 - x1 + 1) * (y2 - y1 + 1)

		// remap to compressed coordinates
		xr1, xr2 := xmap[x1], xmap[x2]+1
		yr1, yr2 := ymap[y1], ymap[y2]+1

		// part 2
		all := (xr2 - xr1) * (yr2 - yr1) // total area

		// count insiders using prefix sums
		insiders := sums[θ(xr2, yr2)] - sums[θ(xr2, yr1)] - sums[θ(xr1, yr2)] + sums[θ(xr1, yr1)]

		// part 1
		if area > acc1 {
			acc1 = area
		}

		// part 2
		if all == insiders && area > acc2 {
			acc2 = area
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
