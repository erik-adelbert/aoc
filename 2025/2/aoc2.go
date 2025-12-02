package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

func main() {
	var (
		acc1, acc2 int      // parts 1 and 2
		s2         [19]byte // buffer for rotation check
	)

	buf, _ := os.ReadFile("input.txt") // read entire input file
	line := bytes.TrimSpace(buf)       // single line input

	for _, span := range bytes.Split(line, []byte(",")) {
		// parse range
		bounds := bytes.Split(span, []byte("-"))

		a, b := atoi(bounds[0]), atoi(bounds[1])

		for i := a; i <= b; i++ {
			s := itoa(i) // convert to byte slice

			// create doubled slice
			copy(s2[:len(s)], s)
			copy(s2[len(s):], s)

			if slices.Equal(s[len(s)/2:], s[:len(s)/2]) {
				// part1: s is a repetition
				acc1 += i
			}

			if bytes.Contains(s2[1:2*len(s)-1], s) {
				// part2: s has a rotation
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

func itoa(n int) []byte {
	if n == 0 {
		return []byte("0")
	}

	var buf [19]byte // enough for 64-bit integer
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	return buf[i:]
}
