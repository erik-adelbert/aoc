// aoc2.go --
// advent of code 2025 day 2
//
// https://adventofcode.com/2025/day/2
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-2: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

const MaxDigits = 10 // maximum digit count for our inputs

func main() {
	var (
		acc1, acc2 int                 // parts 1 and 2
		ss         [2 * MaxDigits]byte // for part 2 rotation check
	)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	line := bytes.TrimSpace(input.Bytes()) // single line input
	for _, span := range bytes.Split(line, []byte(",")) {
		// parse range
		bufA, bufB, _ := bytes.Cut(span, []byte("-"))

		a, b := atoi(bufA), atoi(bufB)

		for i := a; i <= b; i++ {
			s := itoa(i)
			slen := len(s)

			// part1: check if s is a repetition
			half := slen >> 1
			if (slen&1 == 0) && slices.Equal(s[:half], s[half:]) {
				acc1 += i
			}

			// part2: check if s is a rotation of itself
			// double the slice in the buffer
			copy(ss[:slen], s)
			copy(ss[slen:], ss[:slen])

			// check inside excluding full matches at ends
			if bytes.Contains(ss[1:slen<<1-1], s) {
				acc2 += i
			}
		}
	}

	fmt.Println(acc1, acc2)
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

// itoa converts n to its byte slice representation.
// It is taylored for our inputs (n >= 0) of at most MaxDigits digits.
// It prevents allocations by using a fixed-size array internally
func itoa(n int) []byte {
	if n == 0 {
		return []byte("0")
	}

	var buf [MaxDigits]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	return buf[i:]
}
