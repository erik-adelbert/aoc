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
)

func main() {
	var acc1, acc2 int // parts 1 and 2 accumulators

	input := bufio.NewScanner(os.Stdin)

	var lines [][]byte

	// read all input lines
	for input.Scan() {
		lines = append(lines, bytes.Clone(input.Bytes()))
	}

	// parse operations from last line
	ops := lines[len(lines)-1]

	// determine column positions and extract operations
	cols := make([]int, 0, 1000)      // pre-allocate for typical 1000 columns
	cleanOps := make([]byte, 0, 1000) // pre-allocate for typical 1000 operations

	for i, op := range ops {
		if op != ' ' {
			cols = append(cols, i)
			cleanOps = append(cleanOps, op)
		}
	}

	// get the column widths
	widths := make([][2]int, len(cols))

	for i := range len(cols) - 1 {
		widths[i] = [2]int{cols[i], cols[i+1]}
	}
	// add last column width
	last, size := cols[len(cols)-1], len(lines[0])+1
	widths[len(widths)-1] = [2]int{last, size}

	// split lines into columns
	numRows := len(lines) - 1
	splits := make([][][]byte, numRows)

	for row := range splits {
		splits[row] = make([][]byte, len(widths))
		for col, width := range widths {
			a, b := width[0], width[1]
			splits[row][col] = lines[row][a : b-1]
		}
	}

	// transpose columns for easier processing
	tokens := make([][][]byte, len(splits[0]))

	for col := range tokens {
		tokens[col] = make([][]byte, numRows)
		for row := range splits {
			tokens[col][row] = splits[row][col]
		}
	}

	// pre-allocate transpose matrix
	trans := make([][]byte, len(tokens[0]))

	for row := range trans {
		trans[row] = make([]byte, len(tokens[0]))
	}

	// compute results
	for col, token := range tokens {
		transLen := len(token[0])

		// transpose for part 2
		for row := range transLen {
			for j := range token {
				trans[row][j] = token[j][row]
			}
		}

		switch cleanOps[col] {
		case '+':
			for row := range token {
				acc1 += atoi(bytes.TrimSpace(token[row]))
			}

			for row := range transLen {
				acc2 += atoi(bytes.TrimSpace(trans[row][:len(token)]))
			}
		case '*':
			prod := 1
			for row := range token {
				prod *= atoi(bytes.TrimSpace(token[row]))
			}
			acc1 += prod

			prod = 1
			for row := range transLen {
				prod *= atoi(bytes.TrimSpace(trans[row][:len(token)]))
			}
			acc2 += prod
		}
	}

	fmt.Println(acc1, acc2)
}

// atoi converts a byte slice representing a decimal integer to int
func atoi(s []byte) int {
	n := 0
	for i := range s {
		n = n*10 + int(s[i]-'0')
	}
	return n
}
