// aoc10.go --
// advent of code 2025 day 10
//
// https://adventofcode.com/2025/day/10
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-10: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"math"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"
)

var nsolver = runtime.NumCPU() // number of parallel solvers

func main() {
	t0 := time.Now()

	var wg sync.WaitGroup // wait group for solvers

	in := make(chan mach, 2*nsolver)   // input machines into solvers
	out := make(chan parts, 2*nsolver) // output part results from solvers

	// each solver processes machines from the input channel for parts 1 and 2
	solver := func() {
		for m := range in { // for each machine
			out <- parts{
				part1(m.switches, m.light), // min presses to reach light pattern
				part2(m.switches, m.jolts), // min presses to reach joltage targets
			}
		}
	}

	// watchdog to close output channel when all solvers are done
	closer := func() {
		wg.Wait()
		close(out)
	}

	// iterate over parsed machines and send them to input channel
	parser := func() {
		for m := range parse(bufio.NewScanner(os.Stdin)) {
			in <- m
		}
		close(in)
	}

	go parser() // launch feeder

	// start worker pool
	for range nsolver {
		wg.Go(solver) // Go 1.25+
	}

	go closer() // launch watchdog

	var acc1, acc2 uint16 // part 1 & 2 accumulators

	// collect results
	for r := range out {
		acc1 += r.p1
		acc2 += r.p2
	}

	fmt.Printf("%d %d %dμs\n", acc1, acc2, time.Since(t0).Microseconds())
}

type parts struct{ p1, p2 uint16 }
type mach struct {
	switches []uint16
	jolts    []int32
	light    uint16
}

// Part 1: bitmask BFS
func part1(switches []uint16, light uint16) uint16 {
	w := switches[len(switches)-1]        // width of bitmask
	switches = switches[:len(switches)-1] // remove bounding switch

	// build bitmask BFS table
	bmasks := slices.Repeat([]uint16{math.MaxUint16}, 1<<w)
	bmasks[0] = 0

	// BFS over bitmasks
	q := []uint16{0}
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range switches {
			nxt := u ^ v

			if bmasks[nxt] != math.MaxUint16 {
				continue
			}

			bmasks[nxt] = bmasks[u] + 1
			q = append(q, nxt)
		}
	}

	return uint16(bmasks[light])
}

// Part 2 solves an integer linear system M x = rhs using a row-Hermite Normal Form.
// The matrix is small and structured, so we:
//  1. compute a particular integer solution,
//  2. extract a small integer kernel (≤ 3 dimensions),
//  3. search the affine solution space for the minimal feasible solution.
//
// This is exact arithmetic; no floating point or Gaussian elimination is used.
// This approach is exact and fast for our inputs. It is not intended as a general-purpose
// integer linear optimizer.
func part2(switches []uint16, jolts []int32) uint16 {
	switches = switches[:len(switches)-1] // remove bounding switch

	m := len(jolts)    // equations (columns)
	n := len(switches) // variables (rows)

	// M is built with variables as rows and constraints as columns.
	// Row operations therefore preserve the integer solution set of M x = rhs
	// and allow kernel vectors to be read directly from the unimodular transform.
	M := make([][]int32, n)
	for i := range M {
		M[i] = make([]int32, m)

		for j := range M[i] {
			M[i][j] = int32((switches[i] >> j) & 1) // variable i affects joltage j
		}
	}

	var K kern // kernel vectors
	for i := range K {
		K[i] = make([]int32, n)
	}

	// row-Hermite Normal Form
	x0, kdim := hnf(K, M, jolts) // base solution + kernel dimension

	// base sum
	var sum0 int32
	for i := range x0 {
		sum0 += x0[i]
	}

	if kdim == 0 {
		// no free variables, return base sum
		return uint16(sum0)
	}

	// Each kernel vector shifts the solution within the affine space.
	// δi is the change in the total sum when moving one step along kernel vector i.
	// Signs are normalized so that δi ≥ 0 for minimization.
	δs := make([]int32, kdim)
	for i := range δs {
		k := K[i]

		var δ int32
		for j := range k {
			δ += k[j]
		}

		// pretend all step are positive (no negative presses)
		if δ < 0 {
			δ = -δ
			for j := range k {
				k[j] = -k[j]
			}
		}

		δs[i] = δ
	}

	// Variable upper bounds induced by joltage constraints.
	// These bounds restrict feasible movement along kernel directions.
	lims := make([]int32, n)
	for i := range n {
		var lim int32 = math.MaxInt32 // upper limit for variable i

		for j := range m {
			if (switches[i] & (1 << j)) != 0 {
				lim = min(lim, jolts[j]) // variable i can be used to increment joltage j
			}
		}

		lims[i] = lim
	}

	switch kdim { // count of free variables
	case 1:
		return min1D(x0, sum0, δs[0], K[0])
	case 2:
		return min2D(x0, sum0, δs[:2], K[:2], lims, false)
	case 3:
		return min3D(x0, sum0, δs[:3], K[:3], lims)
	}

	// no more than 3 free variables in the input
	panic("unreached")
}

// kernel type has a maximum of 3 free variables
type kern = [3][]int32

// hnf computes a *row* Hermite Normal Form of M using integer row operations.
// U tracks the unimodular transformation such that:
//
//	U * M = H
//
// where H is in row-HNF. Rows below 'rank' in H are zero rows, and the
// corresponding rows in U form an integer basis of ker(M).
//
// A particular solution to M x = rhs is reconstructed by applying the same
// transformations to rhs via U.
func hnf(K kern, M [][]int32, rhs []int32) ([]int32, int) {
	m, n := len(M[0]), len(M)

	// U starts as the identity and accumulates the same row operations as M.
	// Because U is unimodular (det = ±1), it preserves the set of integer solutions.
	// Lower rows of U span ker(M), upper rows map pivot variables to solution space.
	U := make([][]int32, n)
	for i := range n {
		U[i] = make([]int32, n)
		U[i][i] = 1
	}

	// perform HNF on M
	rank := 0
	var pivs [16]int // pivot columns

	// Invariant:
	//   - Rows [0:rank) form a triangular system on pivot columns pivs[0:rank)
	//   - Rows [rank:] are unconstrained
	//   - All transformations are unimodular row operations
	for c := 0; c < m && rank < n; c++ {
		r := rank

		for r < n && M[r][c] == 0 {
			r++
		}

		if r == n {
			continue
		}

		pivs[rank] = c

		if r != rank {
			M[rank], M[r] = M[r], M[rank]
			U[rank], U[r] = U[r], U[rank]
		}

		mr := M[rank]
		ur := U[rank]

		// Eliminate M[i][c] using an extended-GCD-based unimodular transform.
		// This keeps entries small and guarantees exact integer arithmetic.
		for i := rank + 1; i < n; i++ {
			var mi, ui []int32

			if mi, ui = M[i], U[i]; mi[c] == 0 {
				continue
			}

			u, v := mr[c], mi[c]

			gcd, x, y := egcd(u, v)

			u /= gcd
			v /= gcd

			linc(mr, mi, x, y, -v, u)
			linc(ur, ui, x, y, -v, u)
		}

		rank++
	}

	// Rows below 'rank' are zero rows in H.
	// Their corresponding rows in U form an integer basis of ker(M).
	// Kernel dimension is guaranteed small (≤ 3) for this input.
	kdim := n - rank

	for i := range kdim {
		copy(K[i], U[rank+i])
	}

	// Reconstruct a particular integer solution.
	// The triangular structure guarantees exact divisibility here for valid inputs.
	// s[r] is the coefficient for the r-th pivot variable in the transformed system.
	// NOTE: No consistency check is performed, solvability is expected.
	x0 := make([]int32, n)

	var s [16]int32 // solution for rank variables
	for r := 0; r < rank; r++ {
		c := pivs[r] // pivot column

		n := rhs[c]
		for i := range r {
			n -= M[i][c] * s[i]
		}

		s[r] = n / M[r][c]

		axpy(s[r], U[r], x0) // x0 += s[r] * U[r]
	}

	return x0, kdim // return base solution and kernel dimension
}

// min1D finds the minimum of the objective over a 1D integer affine subspace
// subject to non-negativity and implicit feasibility constraints.
func min1D(x []int32, sum0, δ0 int32, K []int32) uint16 {
	var min0, max0 int32 = math.MinInt32, math.MaxInt32

	// minimize over k0*x[i] + ... constraints
	for i := range x {
		k0 := K[i]

		switch {
		case k0 > 0:
			min0 = max(min0, cdiv(-x[i], k0))
		case k0 < 0:
			max0 = min(max0, fdiv(x[i], -k0))
		}
	}

	if min0 > max0 {
		// infeasible
		return math.MaxUint16
	}

	// optimal at min0
	return uint16(sum0 + min0*δ0)
}

// min2D performs a bounded search in a 2D kernel space.
// One dimension is chosen as primary based on the smallest feasible range.
// This is not a general ILP solver; it relies on small kernel dimension.
func min2D(x []int32, sum0 int32, δs []int32, K [][]int32, lims []int32, safe bool) uint16 {
	var min0, max0 int32 = math.MinInt32 / 2, math.MaxInt32 / 2
	var min1, max1 int32 = math.MinInt32 / 2, math.MaxInt32 / 2

	for i := range x {
		k0, k1 := K[0][i], K[1][i]

		lim := lims[i]

		switch {
		case k0 != 0 && k1 == 0:
			limit(&min0, &max0, k0, x[i], lim)
		case k0 == 0 && k1 != 0:
			limit(&min1, &max1, k1, x[i], lim)
		}
	}

	range0 := max0 - min0
	range1 := max1 - min1

	// use smaller range as primary
	if range0 > range1 {
		// use range1 as primary
		min0, max0 = min1, max1

		// swap 0,1 -> 1,0
		δs[0], δs[1] = δs[1], δs[0]
		K[0], K[1] = K[1], K[0]

		defer func() {
			// restore original order for caller

			// reswap 0,1 -> 1,0
			δs[0], δs[1] = δs[1], δs[0]
			K[0], K[1] = K[1], K[0]
		}()
	}

	if safe {
		// work on a copy to avoid modifying caller's x
		// not needed if min2D is terminal
		x = slices.Clone(x)
	}

	for i := range x {
		x[i] += min0 * K[0][i]
	}

	// move along primary
	best := min1D(x, sum0+min0*δs[0], δs[1], K[1])
	for i := min0 + 1; i <= max0; i++ {
		for j := range x {
			x[j] += K[0][j]
		}

		// minimize along secondary
		m := min1D(x, sum0+i*δs[0], δs[1], K[1])
		if m > best {
			break
		}

		best = m
	}

	return best
}

// min3D reduces the 3D kernel problem to nested 2D searches.
// This is terminal and allowed to modify x in-place.
func min3D(x []int32, sum0 int32, δs []int32, K [][]int32, limits []int32) uint16 {
	var min0, max0 int32 = math.MinInt32 / 2, math.MaxInt32 / 2
	var min1, max1 int32 = math.MinInt32 / 2, math.MaxInt32 / 2
	var min2, max2 int32 = math.MinInt32 / 2, math.MaxInt32 / 2

	for i := range x {
		k0, k1, k2 := K[0][i], K[1][i], K[2][i]

		lim := limits[i]

		switch {
		case k0 != 0 && k1 == 0 && k2 == 0:
			limit(&min0, &max0, k0, x[i], lim)
		case k0 == 0 && k1 != 0 && k2 == 0:
			limit(&min1, &max1, k1, x[i], lim)
		case k0 == 0 && k1 == 0 && k2 != 0:
			limit(&min2, &max2, k2, x[i], lim)
		}
	}

	range0 := max0 - min0
	range1 := max1 - min1
	range2 := max2 - min2

	// use smallest range as primary
	switch min(range0, range1, range2) {
	case range0:
		// do nothing more

	case range1:
		min0, max0 = min1, max1

		// swap 0, 1, 2 -> 1, 0, 2
		δs[0], δs[1] = δs[1], δs[0]
		K[0], K[1] = K[1], K[0]

		// don't restore order because min3D is terminal

	case range2:
		min0, max0 = min2, max2

		// swap 0, 1, 2 -> 2, 1, 0
		δs[0], δs[2] = δs[2], δs[0]
		K[0], K[2] = K[2], K[0]

		// don't restore order because min3D is terminal
	}

	// primary is in 0
	for i := range x {
		x[i] += min0 * K[0][i] // in-place because min3D is terminal
	}

	// move along primary hyperplane
	best := min2D(x, sum0+min0*δs[0], δs[1:3], K[1:3], limits, true)
	for i := min0 + 1; i <= max0; i++ {
		for j := range x {
			x[j] += K[0][j]
		}

		// minimize along secondary hyperplane
		best = min(best, min2D(x, sum0+i*δs[0], δs[1:3], K[1:3], limits, true))
	}

	return best
}

const (
	Off = '.'
	On  = '#'
)

// parse input into machines
func parse(input *bufio.Scanner) iter.Seq[mach] {
	return func(yield func(mach) bool) {

		for input.Scan() {
			buf := input.Bytes()
			fields := bytes.Split(buf, []byte(" "))

			// lights
			lfield := fields[0]
			lfield = lfield[1 : len(lfield)-1] // trim brackets

			light := uint16(0)
			for i, c := range lfield {
				if c == On {
					light |= 1 << i
				}
			}

			// switches
			sfields := fields[1 : len(fields)-1]

			switches := make([]uint16, len(sfields)+1)

			var w int32 // width of bitmask
			for i, f := range sfields {
				f = f[1 : len(f)-1] // trim brackets

				for b := range bytes.SplitSeq(f, []byte(",")) {
					n := atoi(b)
					w = max(w, n)
					switches[i] |= uint16(1 << n)
				}
			}
			// add dummy switch for bounding
			switches[len(sfields)] = uint16(w + 1)

			// joltages
			jfield := fields[len(fields)-1]
			jfield = jfield[1 : len(jfield)-1] // trim brackets
			jsets := bytes.Split(jfield, []byte(","))

			jolts := make([]int32, len(jsets))
			for i := range jolts {
				jolts[i] = atoi(jsets[i])
			}

			if !yield(mach{switches, jolts, light}) {
				return
			}
		}
	}
}

// Helper functions

// cdiv: ceiling division a/b
func cdiv(a, b int32) int32 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

// fdiv: floor division a/b
func fdiv(a, b int32) int32 {
	if a >= 0 {
		return a / b
	}
	return (a - (b - 1)) / b
}

// limit min/max bounds based on k*x + ... constraints
func limit(mini, maxi *int32, k, x, lim int32) {
	switch {
	case k > 0:
		// k*min >= -x  and  k*max <= lim - x
		*mini = max(*mini, cdiv(-x, k))
		*maxi = min(*maxi, fdiv(lim-x, k))
	case k < 0:
		// k*min <= lim - x  and  k*max >= -x
		*mini = max(*mini, cdiv(x-lim, -k))
		*maxi = min(*maxi, fdiv(x, -k))
	}
}

// extended gcd returns gcd(a,b) and x,y such that ax+by=gcd(a,b)
func egcd(a, b int32) (int32, int32, int32) {
	var x0, x1 int32 = 1, 0
	var y0, y1 int32 = 0, 1

	for b != 0 {
		q := a / b

		a, b = b, a%b

		x0, x1 = x1, x0-q*x1
		y0, y1 = y1, y0-q*y1
	}

	return a, x0, y0
}

// linc performs a linear combination:
// x <- a*x + b*y
// y <- c*x + d*y
func linc(x, y []int32, a, b, c, d int32) {
	for i := range x {
		xi, yi := x[i], y[i]

		x[i] = a*xi + b*yi
		y[i] = c*xi + d*yi
	}
}

// axpy performs the operation y <- y + a*x
func axpy(a int32, x, y []int32) {
	for i := range y {
		y[i] += x[i] * a
	}
}

// atoi: convert byte slice to integer
func atoi(s []byte) (n int32) {
	for _, c := range s {
		n = 10*n + int32(c-'0')
	}
	return
}
