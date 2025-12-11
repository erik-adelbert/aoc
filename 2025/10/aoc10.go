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
	"math"
	"math/bits"
	"os"
	"slices"
	"time"
)

func main() {
	t0 := time.Now()

	var acc1, acc2 uint16 // part 1 & 2 accumulators

	switches, lights, jolts := parseInput(bufio.NewScanner(os.Stdin))
	_ = jolts

	acc1 = part1(switches, lights)
	acc2 = part2(switches, jolts)

	fmt.Println(acc1, acc2, time.Since(t0))
}

// Part 1: bitmask BFS
func part1(switches [][]uint16, lights []uint16) uint16 {
	var acc1 uint16

	for i, row := range switches {
		w := 0 // width of bitmask
		for _, s := range row {
			w = max(w, bits.Len16(s))
		}

		// Build bitmask BFS table
		bmasks := slices.Repeat([]int{-1}, 1<<w)
		bmasks[0] = 0

		// BFS over bitmasks
		q := []uint16{0}
		for i := 0; i < len(q); i++ {
			u := q[i]
			for _, v := range row {
				nxt := u ^ v

				if bmasks[nxt] != -1 {
					continue
				}

				bmasks[nxt] = bmasks[u] + 1
				q = append(q, nxt)
			}
		}

		acc1 += uint16(bmasks[lights[i]])
	}

	return acc1
}

// Part 2: ILP branch-and-bound
func part2(switches [][]uint16, jolts [][]int) uint16 {
	var acc2 uint16

	for i := range switches {
		sw, jo := switches[i], jolts[i]

		n := len(sw)

		w := 0
		for _, s := range sw {
			w = max(w, bits.Len16(s))
		}

		// build constraint matrix
		M := make([][]float64, 2*w+n)
		for r := range M {
			M[r] = make([]float64, n+1)
		}

		// left-hand side
		for j, s := range sw {
			i := (2*w + len(sw)) - 1 - j
			M[i][j] = -1

			for b := 0; b < w; b++ {
				if s&(1<<b) != 0 {
					M[b][j] = 1
					M[b+w][j] = -1
				}
			}
		}

		// right-hand side
		for i := range w {
			M[i][n] = float64(jo[i])
			M[i+w][n] = -float64(jo[i])
		}

		acc2 += solve(M, slices.Repeat([]float64{1}, n))
	}

	return acc2
}

// solve solves the ILP defined by matrix m and coefficients c
func solve(M [][]float64, C []float64) uint16 {
	best := math.Inf(1) // best value found

	n := len(M[0]) - 1 // number of variables

	var rebranch func(M [][]float64)

	rebranch = func(M [][]float64) {
		val, x := simplex(M, C)

		if val+ε >= best || math.IsInf(val, -1) {
			return // infeasible
		}

		k, v := -1, 0
		for i, e := range x {
			if math.Abs(e-math.Round(e)) > ε {
				k = i
				v = int(e)
				break
			}
		}

		if k == -1 { // all integer
			if val+ε < best {
				best = val
			}
		} else {
			// first branch: x_k >= ceil(v)
			s := make([]float64, n+1)
			s[n] = float64(v)
			s[k] = 1
			rebranch(append(M, s))

			// second branch: x_k <= floor(v)
			s = make([]float64, n+1)
			s[n] = float64(^v) // bitwise NOT, same as -v-1 for integers
			s[k] = -1
			rebranch(append(M, s))
		}
	}

	rebranch(M)

	return uint16(math.Round(best))
}

// simplex solves the LP defined by matrix M and coefficients C
func simplex(M [][]float64, C []float64) (float64, []float64) {
	h := len(M)
	w := len(M[0]) - 1

	// N: non-basic variable indices
	N := make([]int, w+1)
	for i := range w {
		N[i] = i
	}
	N[w] = -1

	// B: basic variable indices
	B := make([]int, h)
	for i := range h {
		B[i] = w + i
	}

	// T: tableau
	T := make([][]float64, h+2)

	for i := range h {
		T[i] = make([]float64, w+2)
		T[i][w+1] = -1
		copy(T[i], M[i])

	}

	T[h] = make([]float64, w+2)
	copy(T[h], C)

	// zeros
	T[h][w] = 0
	T[h][w+1] = 0
	T[h+1] = make([]float64, w+2)

	// Swap last two columns in constraint rows
	for i := range h {
		T[i][w], T[i][w+1] = T[i][w+1], T[i][w]
	}
	T[h+1][w] = 1

	// Gauss-Jordan pivoting
	pivot := func(r, c int) {
		k := 1.0 / T[r][c] // pivot element

		// eliminate column c in all other rows
		for i := range h + 2 {
			if i == r {
				continue
			}
			for j := range w + 2 {
				if j != c {
					T[i][j] -= T[r][j] * T[i][c] * k
				}
			}
		}

		// adjust pivot row
		for i := range w + 2 {
			T[r][i] *= k
		}

		// adjust pivot column
		for i := range h + 2 {
			T[i][c] *= -k
		}

		T[r][c] = k // set pivot element

		// swap basic and non-basic variables
		B[r], N[c] = N[c], B[r]
	}

	// find optimal
	find := func(p int) bool {
		for {
			// find s
			s := -1
			smin := math.Inf(1)
			for i := range w + 1 {
				if p != 0 || N[i] != -1 {
					val := T[h+p][i]
					if val < smin || (val == smin && N[i] < N[s]) {
						s = i
						smin = val
					}
				}
			}

			if T[h+p][s] > -ε {
				return true // optimal
			}

			// find r
			r := -1
			rmin := math.Inf(1)
			for i := range h {
				if T[i][s] > ε {
					ratio := T[i][w+1] / T[i][s]
					if ratio < rmin || (ratio == rmin && B[i] < B[r]) {
						r = i
						rmin = ratio
					}
				}
			}

			if r == -1 {
				return false // unbounded
			}

			pivot(r, s)
		}
	}

	// Initialization
	r := 0
	vmin := T[0][w+1]
	for i := 1; i < h; i++ {
		if T[i][w+1] < vmin {
			r = i
			vmin = T[i][w+1]
		}
	}

	if T[r][w+1] < -ε {
		pivot(r, w)
		if !find(1) || T[h+1][w+1] < -ε {
			return math.Inf(-1), nil
		}
	}

	for i := range h {
		if B[i] == -1 {
			// Find s
			s := 0
			vmin := T[i][0]
			for j := 1; j < w; j++ {
				if T[i][j] < vmin || (T[i][j] == vmin && N[j] < N[s]) {
					s = j
					vmin = T[i][j]
				}
			}
			pivot(i, s)
		}
	}

	if find(0) {
		x := make([]float64, w)
		for i := range h {
			if B[i] >= 0 && B[i] < w {
				x[B[i]] = T[i][w+1]
			}
		}

		sum := 0.0
		for i := range w {
			sum += C[i] * x[i]
		}
		return sum, x
	}
	return math.Inf(-1), nil
}

const (
	Off = '.'
	On  = '#'
)

func parseInput(input *bufio.Scanner) ([][]uint16, []uint16, [][]int) {
	var switches [][]uint16
	var lights []uint16
	var jolts [][]int

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

		lights = append(lights, light)

		// switches
		sfields := fields[1 : len(fields)-1]

		row := make([]uint16, len(sfields))
		for i, f := range sfields {
			f = f[1 : len(f)-1] // trim brackets

			for buf := range bytes.SplitSeq(f, []byte(",")) {
				row[i] |= uint16(1 << atoi(buf))
			}

		}
		switches = append(switches, row)

		// joltages
		jfield := fields[len(fields)-1]
		jfield = jfield[1 : len(jfield)-1] // trim brackets
		jsets := bytes.Split(jfield, []byte(","))

		jolt := make([]int, len(jsets))
		for i := range jolt {
			jolt[i] = atoi(jsets[i])
		}
		jolts = append(jolts, jolt)
	}

	return switches, lights, jolts
}

const ε = 1e-9

func atoi(s []byte) (n int) {
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return
}
