// aoc22.go --
// advent of code 2024 day 22
//
// https://adventofcode.com/2024/day/22
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-22: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
	VOFF   = 9
	NLOOP  = 2000
	MAXDIM = 0xFFFFF
)

func main() {
	seqs := make([]int, MAXDIM)

	sum1 := 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n := atoi(input.Text())
		sum1 += rehash(n, NLOOP, seqs)
	}

	count2 := slices.Max(seqs)

	fmt.Println(sum1, count2) // part 1 & 2
}

const MinInt = -1 << 31

var SEEN = make([]int, MAXDIM)

func rehash(a, n int, seqs []int) int {
	color := a

	// loop on digits
	key, cur := 0, a%10
	for i := 0; i < n; i++ {
		a = hash(a)   // rehash
		nxt := a % 10 // next digit

		δ := nxt - cur + VOFF           // vertical offset to avoid negative keys
		key = ((key << 5) & MAXDIM) + δ // running key

		if i > 3 && SEEN[key] != color { // color as seen
			SEEN[key] = color
			seqs[key] += nxt
		}

		cur = nxt
	}

	return a
}

func hash(a int) int {
	a ^= (a << 6) & 0xFFFFFF
	a ^= (a >> 5) & 0xFFFFFF
	a ^= (a << 11) & 0xFFFFFF
	return a
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
