// aoc6.go --
// advent of code 2025 day 6
//
// https://adventofcode.com/2025/day/6
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-6: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

const BufSizeHint = 1 << 10 // 1 K buffer size

func main() {
	t0 := time.Now() // start timer

	var acc1, acc2 int // parts 1 and 2 accumulators

	input := bufio.NewScanner(os.Stdin)

	var lines [][]byte

	// read all input lines
	for input.Scan() {
		lines = append(lines, bytes.Clone(input.Bytes()))
	}

	// parse layout from the last line
	layout := lines[len(lines)-1]

	// determine column positions and extract operations
	cols := make([]int, 0, BufSizeHint) // pre-allocate for typical columns
	ops := make([]byte, 0, BufSizeHint) // pre-allocate for typical operations

	for i, op := range layout {
		if op != ' ' {
			cols = append(cols, i)
			ops = append(ops, op)
		}
	}

	// get the column widths
	widths := make([][2]int, len(cols))

	for i := range len(cols) - 1 {
		widths[i] = [2]int{cols[i], cols[i+1]}
	}
	// add last column
	last, size := cols[len(cols)-1], len(lines[0])+1
	widths[len(widths)-1] = [2]int{last, size}

	// get the column tokens
	h, w := len(lines)-1, len(widths)
	tokens := make([][][]byte, w)
	for c := range tokens {
		tokens[c] = make([][]byte, h)
	}

	for r := range h {
		for c, w := range widths {
			α, ω := w[0], w[1]
			tokens[c][r] = lines[r][α : ω-1]
		}
	}

	// pre-allocate number transpose matrix
	nekot := make([][]byte, len(tokens[0]))

	for i := range nekot {
		nekot[i] = make([]byte, len(tokens))
	}

	// compute results
	for c, token := range tokens {
		h, w := len(token), len(token[0])

		// transpose for part 2
		for i := range w {
			for j := range h {
				nekot[i][j] = token[j][i]
			}
		}

		switch ops[c] {
		case '+':
			for r := range h {
				acc1 += atoi(bytes.TrimSpace(token[r][:w]))
			}

			for c := range w {
				acc2 += atoi(bytes.TrimSpace(nekot[c][:h]))
			}
		case '*':
			prod := 1
			for r := range h {
				prod *= atoi(bytes.TrimSpace(token[r][:w]))
			}
			acc1 += prod

			prod = 1
			for c := range w {
				prod *= atoi(bytes.TrimSpace(nekot[c][:h]))
			}
			acc2 += prod
		}
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}

	return
}
