package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

const MaxDigit = 10 // maximum digit count for our inputs

func main() {
	var (
		acc1, acc2 int                // parts 1 and 2
		ss         [2 * MaxDigit]byte // buffer for rotation check. 20 is enough for our inputs.
	)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	line := bytes.TrimSpace(input.Bytes()) // single line input
	for _, span := range bytes.Split(line, []byte(",")) {
		// parse range
		bufA, bufB, _ := bytes.Cut(span, []byte("-"))

		a, b := atoi(bufA), atoi(bufB)

		for i := a; i <= b; i++ {
			s := itoa(i) // convert to byte slice
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

func itoa(n int) []byte {
	if n == 0 {
		return []byte("0")
	}

	var buf [MaxDigit]byte // 10 is enough for our inputs
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	return buf[i:]
}
