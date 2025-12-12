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
	PolyominoCount = 6
	PolyominoDim   = 3
)

func main() {
	t0 := time.Now() // start timer

	var acc1 int // part 1 accumulator

	input := bufio.NewScanner(os.Stdin)

	polys := make([]polyomino, 0, PolyominoCount)
	block := make([][]byte, 0, PolyominoDim)

	// read all input lines
	for i := 0; input.Scan(); i++ {
		buf := input.Bytes()
		maxPoly := PolyominoCount * (PolyominoDim + 2)

		switch {
		case len(buf) == 0:
			// end of polyomino block
			if i < maxPoly && len(block) > 0 {
				polys = append(polys, parse(block))
				block = block[:0]
			}

		case i < maxPoly && i%(PolyominoDim+2) == 0:
			// skip header line with index
		case i < maxPoly:
			// polyomino line
			block = append(block, buf)

		default:
			// processing grid line: "WxH: n1 n2 n3 ..."
			lhs, rhs, _ := bytes.Cut(buf, []byte(": "))  // ["WxH", " n1 n2 n3 ..."]
			wbuf, hbuf, _ := bytes.Cut(lhs, []byte("x")) // ["W", "H"]
			fields := bytes.Fields(rhs)

			w, h := atoi(wbuf), atoi(hbuf)
			area := w * h

			// compute total required cells
			size := 0
			for j, f := range fields {
				if j >= len(polys) {
					break
				}
				num := atoi(f)
				size += num * polys[j].area
			}

			// apply empirical 87% heuristic
			if size*115 < area*100 {
				// if total area used is less than 87% of grid area, count it
				acc1++
			}
		}
	}

	fmt.Println(acc1, time.Since(t0))
}

// polyomino represents a polyomino by its area
type polyomino struct {
	area int
}

// parse counts '#' in the block and returns a polyomino with that area
func parse(block [][]byte) polyomino {
	area := 0
	for _, row := range block {
		area += bytes.Count(row, []byte{'#'})
	}
	return polyomino{area}
}

func atoi(s []byte) (n int) {
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return
}
