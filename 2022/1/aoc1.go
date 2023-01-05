package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var m1, m2, m3, sum int

	max3 := func() {
		switch {
		case sum > m1:
			m1, m2, m3 = sum, m1, m2
		case sum > m2:
			m2, m3 = sum, m2
		case sum > m3:
			m3 = sum
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Bytes()
		if len(line) > 1 {
			sum += atoi(line)
		} else {
			max3()
			sum = 0
		}
	}
	max3()

	fmt.Println(m1, m1+m2+m3) // part 1 & 2
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) int {
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}
