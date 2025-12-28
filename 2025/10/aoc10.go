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
	"math/bits"
	"os"
	"sync"
	"time"
	"unsafe"
)

const MaxSolver = 16 // number of parallel solvers -- sweet spot on M1 instead of 8 (why!?)

// use this for dynamic max procs
// var MaxSolver = runtime.GOMAXPROCS(0)

func main() {
	t0 := time.Now()

	var wg sync.WaitGroup // wait group for solvers

	in := make(chan mach, 4*MaxSolver)   // input machines into solvers
	out := make(chan parts, 4*MaxSolver) // output part results from solvers

	// each solver processes machines from the input channel for parts 1 and 2
	solver := func() {
		for m := range in { // for each machine
			out <- parts{
				part1(m.flips, m.light), // min presses to reach light pattern
				part2(m.flips, m.jolts), // min presses to reach joltage targets
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
	for range MaxSolver {
		wg.Go(solver) // Go 1.25+
	}

	go closer() // launch watchdog

	var acc1, acc2 i32 // part 1 & 2 accumulators

	// collect results
	for r := range out {
		acc1 += r.p1
		acc2 += r.p2
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

type i32 = int32
type u16 = uint16

// Part 1: bitmask BFS
func part1(flips []i32, light i32) i32 {
	w := flips[len(flips)-1]     // width of bitmask
	flips = flips[:len(flips)-1] // remove bounding switch

	// build bitmask BFS table
	N := 1 << w

	bmasks := make([]i32, N)
	for i := range N {
		bmasks[i] = math.MaxInt16
	}
	bmasks[0] = 0

	// BFS over bitmasks
	q := make([]i32, 0, N)

	q = append(q, 0)
	for i := 0; i < len(q); i++ {
		u := q[i]
		for _, v := range flips {
			nxt := u ^ v

			if bmasks[nxt] != math.MaxInt16 {
				continue
			}

			if bmasks[nxt] = bmasks[u] + 1; nxt == light {
				return bmasks[light]
			}

			q = append(q, nxt)
		}
	}

	panic("unreached")
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
func part2(flips []i32, jolts []i32) i32 {
	flips = flips[:len(flips)-1] // remove bounding switch

	m := len(jolts) // equations (cols)
	n := len(flips) // variables (rows)

	// M is the transpose of the usual matrix because we perform
	// row-HNF. Each row corresponds to a variable (flips),
	var M mat
	for i := range n {
		for j := range m {
			M[x(i, j)] = (flips[i] >> j) & 1
		}
	}

	var K ker3

	// row-Hermite Normal Form
	x, kdim := hnf(&K, &M, jolts, m, n) // base solution + kernel dimension

	// base sum
	var sumX i32
	for i := range n {
		sumX += x[i]
	}

	if kdim == 0 {
		// no free variables, return base sum
		return sumX
	}

	k0 := ρ(K[:], 0) // kernel basis vectors
	sum0 := Σplus(k0)

	if kdim == 1 {
		return sumX + min1D(x, k0, n)
	}

	k1 := ρ(K[:], 1) // kernel basis vectors
	sum1 := Σplus(k1)

	if sum0 < sum1 {
		sum0, sum1 = sum1, sum0
		swap(k0, k1)
	}

	if kdim == 2 {
		return sumX + min2D(x, sum0, k0, sum1, k1, n, math.MaxInt16)
	}

	k2 := ρ(K[:], 2) // kernel basis vectors
	sum2 := Σplus(k2)

	if sum0 < sum2 {
		sum0, sum2 = sum2, sum0
		swap(k0, k2)
	}

	// no more than 3 free variables in the input
	return sumX + min3D(x, sum0, k0, sum1, k1, sum2, k2, n)
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

			var light i32

			for i, c := range lfield {
				if c == On {
					light |= 1 << i
				}
			}

			// flips
			sfields := fields[1 : len(fields)-1]

			flips := make([]i32, len(sfields)+1)

			var w i32 // width of bitmask

			for i, f := range sfields {
				f = f[1 : len(f)-1] // trim brackets

				for b := range bytes.SplitSeq(f, []byte(",")) {
					n := atoi(b)

					w = max(w, n)
					flips[i] |= 1 << n
				}
			}
			// add dummy switch for bounding
			flips[len(sfields)] = w + 1

			// joltages
			jfield := fields[len(fields)-1]
			jfield = jfield[1 : len(jfield)-1] // trim brackets
			jsets := bytes.Split(jfield, []byte(","))

			jolts := make([]i32, len(jsets))
			for i := range jolts {
				jolts[i] = atoi(jsets[i])
			}

			if !yield(mach{flips, jolts, light}) {
				return
			}
		}
	}
}

// utility functions and types
type parts struct{ p1, p2 i32 }
type mach struct {
	flips []i32
	jolts []i32
	light i32
}

// ker3 type has a maximum of 3 free variables
type ker3 = [3 * 16]i32

// flat 16x16 matrix type
type mat = [16 * 16]i32

// v16 is a 16-dimensional integer vector
type v16 [16]i32

// min1D finds the minimum of the objective over a 1D integer affine subspace
// subject to non-negativity and implicit feasibility constraints.
func min1D(x *v16, v0 *v16, n int) i32 {
	var sum i32

	for i := range n {
		sum += v0[i]
	}

	if sum > 0 {
		sum = -sum

		for i := range n {
			v0[i] *= -1
		}
	}

	var negi u16 // negative entries in v0
	for i := range n {
		if v0[i] < 0 {
			negi |= 1 << i
		}
	}

	i := bits.TrailingZeros16(negi)
	a, b := x[i], -v0[i] // first negative entry
	negi &= negi - 1     // clear lowest set bit

	for negi != 0 {
		i = bits.TrailingZeros16(negi) // next negative entry

		v0i, xi := -v0[i], x[i]

		if a*v0i > b*xi {
			a, b = xi, v0i
		}

		negi &= negi - 1
	}

	return sum * fdiv(a, b)
}

// min2D performs a bounded search in a 2D kernel space.
// One dimension is chosen as primary based on the smallest feasible range.
// This is not a general ILP solver; it relies on small kernel dimension.
func min2D(x *v16, sum0 i32, v0 *v16, sum1 i32, v1 *v16, n int, best i32) i32 {
	var posMask0, negMask1 u16

	for i := range n {
		if v0[i] > 0 {
			posMask0 |= 1 << i
		}

		if v1[i] < 0 {
			negMask1 |= 1 << i
		}
	}

	addMask := posMask0 &^ negMask1
	lostMask := ^posMask0 &^ negMask1

	min0 := fmbounds2D(x, v0, v1, n)
	min1 := fmbounds2D(x, v1, v0, n)

	var min10, min01 i32 = math.MinInt16, math.MinInt16
	for i := range n {
		v0i, v1i := v0[i], v1[i]

		if v0i > 0 {
			r := x[i] + min1*v1i
			if min01*v0i < -r {
				min01 = cdiv(-r, v0i)
			}
		}

		if v1i > 0 {
			r := x[i] + min0*v0i
			if min10*v1i < -r {
				min10 = cdiv(-r, v1i)
			}
		}
	}

	sm0 := min0*sum0 + min10*sum1
	sm1 := min01*sum0 + min1*sum1

	if min(sm0, sm1) >= best {
		return best
	}

	if sm1 < sm0 {
		v0, v1 = v1, v0
		sum0, sum1 = sum1, sum0
		min0, min1 = min1, min0
		min01, min10 = min10, min01
	}

	var xx v16 = *x // working copy
	x = &xx

	for i := range n {
		x[i] += min0*v0[i] + min10*v1[i]
	}

	sum := min0*sum0 + min10*sum1
	k0 := min0

	hasNeg := false
	for i := range n {
		if x[i] < 0 {
			hasNeg = true
			break
		}
	}

	if !hasNeg {
		if sum < best {
			best = sum
		}

		k0++

		for i := range n {
			x[i] += v0[i] - v1[i]
		}

		sum += sum0 - sum1
	}

	for {
		if k0*sum0+min1*sum1 >= best {
			return best
		}

		var negMask u16
		for i := range n {
			if x[i] < 0 {
				negMask |= 1 << i
			}
		}

		if negMask != 0 {
			switch {
			case negMask&lostMask != 0:
				return best
			case negMask&addMask != 0:
				k0++

				for i := range n {
					x[i] += v0[i]
				}

				sum += sum0
			default:
				for i := range n {
					x[i] -= v1[i]
				}

				sum -= sum1
			}

			continue
		}

		// subtract v1 until minimal
		hasNeg = false
		for {
			for i := range n {
				x[i] -= v1[i]

				if x[i] < 0 {
					hasNeg = true
				}
			}

			if !hasNeg {
				sum -= sum1
				continue
			}

			for i := range n {
				x[i] += v1[i] // undo last subtraction
			}

			break
		}

		if sum < best {
			best = sum
		}

		k0++

		for i := range n {
			x[i] += v0[i] - v1[i]
		}

		sum += sum0 - sum1
	}
}

// min3D reduces the 3D kernel problem to nested 2D searches.
// This is terminal and allowed to modify x in-place.
func min3D(x *v16, sum0 i32, v0 *v16, sum1 i32, v1 *v16, sum2 i32, v2 *v16, n int) i32 {
	min0, max0 := fmbounds3D(x, v0, v1, v2, n)

	for i := range n {
		x[i] += min0 * v0[i]
	}

	var xx v16 = *x // working copy
	sum := min0 * sum0

	best := sum + min2D(&xx, sum1, v1, sum2, v2, n, math.MaxInt16)

	for k := min0 + 1; k <= max0; k++ {
		for i := range n {
			x[i] += v0[i]
		}
		sum += sum0

		xx, cur := *x, best-sum
		if nxt := min2D(&xx, sum1, v1, sum2, v2, n, cur); nxt < cur {
			best = nxt + sum
		}
	}

	return best
}

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
func hnf(K *ker3, M *mat, rhs []i32, m, n int) (*v16, int) {
	// U starts as the identity and accumulates the same row operations as M.
	// Because U is unimodular (det = ±1), it preserves the set of integer solutions.
	// Lower rows of U span ker(M), upper rows map pivot variables to solution space.
	U := new(mat)
	for i := range n {
		U[x(i, i)] = 1
	}

	// perform HNF on M
	rank := 0
	var pivs [16]int // pivot columns

	// Invariant:
	//   - Rows [0:rank) form a triangular system on pivot columns pivs[0:rank)
	//   - Rows [rank:] are unconstrained
	//   - All transformations are unimodular row operations
	for c := 0; c < m && rank < n; c++ {
		// search for pivot row
		rp := rank

		for rp < n && M[x(rp, c)] == 0 {
			rp++
		}

		if rp == n {
			continue
		}

		pivs[rank] = c

		mr := ρ(M[:], rank) // current rank row
		ur := ρ(U[:], rank) // pivot row in U

		if rp != rank {
			// swap pivot row with current rank row
			swap(mr, ρ(M[:], rp))
			swap(ur, ρ(U[:], rp))
		}

		// Eliminate M[i][c] using an extended-GCD-based unimodular transform.
		// This keeps entries small and guarantees exact integer arithmetic.
		for i := rank + 1; i < n; i++ {
			var mi, ui *v16

			if mi, ui = ρ(M[:], i), ρ(U[:], i); mi[c] == 0 {
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

	// copy kernel basis vectors
	for i := range kdim {
		*ρ(K[:], i) = *ρ(U[:], rank+i)
	}

	// Reconstruct a particular integer solution.
	// The triangular structure guarantees exact divisibility here for valid inputs.
	// s[r] is the coefficient for the r-th pivot variable in the transformed system.
	// NOTE: No consistency check is performed, solvability is expected.
	x0 := new(v16)

	var s [16]i32 // solution for rank variables
	for r := 0; r < rank; r++ {
		c := pivs[r] // pivot column

		n := rhs[c]
		for i := range r {
			n -= M[x(i, c)] * s[i]
		}

		s[r] = n / M[x(r, c)]

		axpy(s[r], ρ(U[:], r), x0) // x0 += s[r] * U[r]
	}

	return x0, kdim // return base solution and kernel dimension
}

// fmbouds2D computes Fourier-Motzkin bounds for 2D kernel
func fmbounds2D(x, v0, v1 *v16, n int) i32 {
	type k1 struct{ x, v0, v1 i32 }

	var pos1, p1 = [16]k1{}, 0
	var neg1, n1 = [16]k1{}, 0

	var min0 i32 = math.MinInt16
	for i := range n {
		xi, v0i, v1i := x[i], v0[i], v1[i]

		switch {
		case v1i == 0 && v0i > 0:
			b := -xi
			if min0*v0i < b {
				min0 = cdiv(b, v0i)
			}
		case v1i > 0:
			pos1[p1] = k1{xi, v0i, v1i}
			p1++
		case v1i < 0:
			neg1[n1] = k1{xi, v0i, v1i}
			n1++
		}
	}

	for _, pc := range pos1[:p1] {
		for _, nc := range neg1[:n1] {
			if a := nc.v0*pc.v1 - pc.v0*nc.v1; a > 0 {
				if b := pc.x*nc.v1 - nc.x*pc.v1; a*min0 < b {
					min0 = cdiv(b, a)
				}
			}
		}
	}

	return min0
}

// Fourier-Motzkin bounds for 3D kernel
func fmbounds3D(x, v0, v1, v2 *v16, n int) (i32, i32) {
	type k2 struct{ x, v0, v1, v2 i32 }
	type k1 struct{ x, v0, v1 i32 }

	var pos2, p2 = [8]k2{}, 0
	var neg2, n2 = [8]k2{}, 0

	var pos1, p1 = [32]k1{}, 0
	var neg1, n1 = [32]k1{}, 0

	var min0, max0 i32 = math.MinInt16, math.MaxInt16

	limit := func(n, d i32) {
		switch {
		case d > 0:
			if min0*d < n {
				min0 = cdiv(n, d)
			}
		case d < 0:
			if max0*d < n {
				max0 = fdiv(-n, -d)
			}
		}
	}

	for i := range n {
		v2i := v2[i]

		switch {
		case v2i == 0:
			v1i := v1[i]
			switch {
			case v1i == 0:
				limit(-x[i], v0[i])
			case v1i < 0:
				neg1[n1] = k1{x[i], v0[i], v1i}
				n1++
			default:
				pos1[p1] = k1{x[i], v0[i], v1i}
				p1++
			}
		case v2i > 0:
			pos2[p2] = k2{x[i], v0[i], v1[i], v2i}
			p2++
		default:
			neg2[n2] = k2{x[i], v0[i], v1[i], v2i}
			n2++
		}
	}

	for _, p := range pos2 {
		pv2 := p.v2

		for _, n := range neg2 {
			nv2 := n.v2

			x := n.x*pv2 - p.x*nv2
			v0 := n.v0*pv2 - p.v0*nv2
			v1 := n.v1*pv2 - p.v1*nv2

			switch {
			case v1 == 0:
				limit(-x, v0)
			case v1 < 0:
				neg1[n1] = k1{x, v0, v1}
				n1++
			default:
				pos1[p1] = k1{x, v0, v1}
				p1++
			}
		}
	}

	for _, p := range pos1 {
		for _, n := range neg1 {
			b := p.x*n.v1 - n.x*p.v1
			a := n.v0*p.v1 - p.v0*n.v1

			limit(b, a)
		}
	}

	return min0, max0
}

// Σplus computes the sum of entries in v, making the sum non-negative
func Σplus(v *v16) (s i32) {
	for i := range v {
		s += v[i]
	}

	if s < 0 {
		s = -s
		for i := range v {
			v[i] = -v[i]
		}
	}

	return
}

// linear index for row-major 16x16 matrix
func x(r, c int) int {
	return r*16 + c
}

// swap exchanges the contents of two v16 vectors
func swap(a, b *v16) {
	*a, *b = *b, *a
}

// ρ returns a row as *v16 for a flat matrix
func ρ(s []i32, i int) *v16 {
	return (*v16)(unsafe.Pointer(&s[i*16]))
}

// linc performs a linear combination:
// x <- a*x + b*y
// y <- c*x + d*y
func linc(x, y *v16, a, b, c, d i32) {
	for i := range x {
		xi, yi := x[i], y[i]

		x[i] = a*xi + b*yi
		y[i] = c*xi + d*yi
	}
}

// axpy performs the operation y <- y + a*x
func axpy(a i32, x, y *v16) {
	for i := range y {
		y[i] += x[i] * a
	}
}

// Helper functions

// cdiv: ceiling division a/b
func cdiv(a, b i32) i32 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

// fdiv: floor division a/b
func fdiv(a, b i32) i32 {
	if a >= 0 {
		return a / b
	}
	return (a - (b - 1)) / b
}

// extended gcd returns gcd(a,b) and x,y such that ax+by=gcd(a,b)
func egcd(a, b i32) (i32, i32, i32) {
	var x0, x1 i32 = 1, 0
	var y0, y1 i32 = 0, 1

	for b != 0 {
		q := a / b

		a, b = b, a%b

		x0, x1 = x1, x0-q*x1
		y0, y1 = y1, y0-q*y1
	}

	return a, x0, y0
}

// atoi: convert byte slice to integer
func atoi(s []byte) (n i32) {
	for _, c := range s {
		n = 10*n + i32(c-'0')
	}
	return
}
