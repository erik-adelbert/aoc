// aoc25.go --
// advent of code 2021 day 25
//
// https://adventofcode.com/2021/day/25
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-25: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	w = 139
	h = 137
)

type board [w * h]byte

func main() {
	var nxt board

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		low, max := slice(j)
		copy(nxt[low:max:max], input.Bytes())
	}

	s, n := 0, 1 // step, change count
	for n > 0 {
		s, n = s+1, 0 // advance step, reset change count

		cur := nxt
		// east scan
		for j := 0; j < h; j++ {
			crow := cur.row(j)
			nrow := nxt.row(j)
			for i := 0; i < w; i++ {
				ii := (i + 1) % w
				if crow[i] == '>' && crow[ii] == '.' {
					nrow[i], nrow[ii] = nrow[ii], nrow[i] // swap!
					n++
				}
			}
		}

		cur = nxt
		// south scan
		for j := 0; j < h; j++ {
			jj := (j + 1) % h
			chead, ctail := cur.row(j), cur.row(jj)
			nhead, ntail := nxt.row(j), nxt.row(jj)
			for i := 0; i < w; i++ {
				if chead[i] == 'v' && ctail[i] == '.' {
					nhead[i], ntail[i] = ntail[i], nhead[i] // swap!
					n++
				}
			}
		}
	}
	fmt.Println(s)
}

func at(j, i int) int {
	return j*w + i
}

func slice(j int) (low, max int) {
	low = j * w
	max = low + w
	return
}

func (b *board) row(j int) []byte {
	low, max := slice(j)
	return b[low:max:max]
}

func (b *board) String() string {
	var sb strings.Builder
	for j := 0; j < h; j++ {
		low, max := slice(j)
		sb.Write(b[low:max:max])
		sb.WriteByte('\n')
	}

	return sb.String()
}
