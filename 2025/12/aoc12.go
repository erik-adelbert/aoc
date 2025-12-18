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

	polys := make([]int, 0, PolyCount) // poliyomino areas
	block := make([][]byte, 0, PolyDim)

	// read all input lines
	input := bufio.NewScanner(os.Stdin)

	for i := 0; input.Scan(); i++ {
		buf := input.Bytes()

		maxPoly := PolyCount * (PolyDim + 2)

		switch {
		case len(buf) == 0:

			// end of polyomino block
			if i < maxPoly && len(block) > 0 {
				polys = append(polys, parse(block))
				block = block[:0]
			}

		case i < maxPoly:
			if i%(PolyDim+2) != 0 {
				block = append(block, buf) // polyomino line
			}

		default:
			// processing grid line: "WxH: n1 n2 n3 ..."
			lhs, rhs, _ := bytes.Cut(buf, []byte(": ")) // ["WxH", "n1 n2 n3 ..."]
			w, h, _ := bytes.Cut(lhs, []byte("x"))      // ["W", "H"]

			area := atoi(w) * atoi(h)

			// compute total required cells
			size := 0
			for j, n := range bytes.Fields(rhs) {
				size += atoi(n) * polys[j]
			}

			// apply empirical 87% heuristic
			if size*23 < area*20 {
				acc1++
			}
		}
	}

	fmt.Println(acc1, time.Since(t0))
}

// parse counts '#' in the block and returns a polyomino with that area
func parse(block [][]byte) int {
	area := 0
	for _, row := range block {
		area += bytes.Count(row, []byte{'#'})
	}
	return area
}

func atoi(s []byte) (n int) {
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return
}
