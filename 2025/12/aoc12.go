// aoc12.go --
// advent of code 2025 day 12
//
// https://adventofcode.com/2025/day/12
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-12: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

const (
	PolyCount = 6
	PolyDim   = 3
)

func main() {
	t0 := time.Now() // start timer

	var acc1 int // part 1 accumulator

	block := make([]byte, 0, sq(PolyDim)) // current polyomino block (flattened)
	cells := make([]int, 0, PolyCount)    // polyomino cell counts

	// flushBlock appends current block cell count to polys
	flushBlock := func() {
		if len(block) > 0 {
			cells = append(cells, bytes.Count(block, []byte("#")))
			block = block[:0]
		}
	}

	state := Poly // input state

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		buf := input.Bytes()

		// ---- polyomino processing phase ----
		// ex. buf format with ncells=6:
		// .#.
		// ###
		// ##.
		if state == Poly {
			switch {
			case bytes.Contains(buf, []byte("x")):
				// grid layout line indicates end of polyominoes
				flushBlock()
				state = Grid
				// proceed to grid processing phase

			case len(buf) == 0:
				// end of current block
				flushBlock()
				continue // proceed to next line

			case buf[0] == '#' || buf[0] == '.':
				// polyomino block line
				block = append(block, buf...) // accumulate block lines
				continue                      // proceed to next line
			}
		}

		// ---- grid processing phase ----
		// buf format: "WxH: n0 n1 n2 n3 n4 n5"
		lhs, rhs, _ := bytes.Cut(buf, []byte(": ")) // [WxH] [n0 n1 n2 n3 n4 n5]

		// parse grid dimensions and calculate area
		w, h, _ := bytes.Cut(lhs, []byte("x")) // [W] [H]
		area := atoi(w) * atoi(h)

		// parse polyomino counts and calculate total cell count
		ncell := 0
		for j, n := range bytes.Fields(rhs) { // enumerate [n0 n1 n2 n3 n4 n5]
			ncell += atoi(n) * cells[j] // total cells += count * poly cell size
		}

		// apply empirical 87% heuristic
		if ncell*23 < area*20 {
			acc1++
		}
	}

	fmt.Println(acc1, time.Since(t0))
}

// input states
const (
	Poly bool = (iota == 0)
	Grid
)

// sq returns the square of n
func sq(n int) int {
	return n * n
}

// atoi converts a byte slice representing a non-negative integer to int
func atoi(s []byte) (n int) {
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return
}
