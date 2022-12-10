package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// machine register X, clock and signal power:
	X := 1
	clk, pwr := 0, 0

	// part2
	const (
		Black = ' '
		// undefined is very bright
		White = '\uFFFD'
	)

	// scan display
	crt := func() {
		// column from clock
		c := clk%40 - 1
		if c < 0 {
			c += 40
		}

		// wrap beam
		if c == 0 && clk > 1 {
			fmt.Println()
		}

		// window
		min, max := X-1, 39
		if X < max {
			max = X + 1
		}

		// beam
		x := Black
		if min <= c && c <= max {
			x = White
		}
		fmt.Printf("%c", x)
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// split incoming instruction:
		// fields:  0   1
		// values: cmd arg
		inst := strings.Fields(input.Text())

		clk++ // tick
		crt() // beam CRT

		// execute, monitor power, beam CRT
		switch inst[0][0] {
		case 'a': // addx
			// part1 clock sync
			switch (clk + 21) % 40 {
			default:
				clk++
			case 1:
				pwr += clk * X
				clk++
			case 0:
				clk++
				pwr += clk * X
			}

			crt()
			X += atoi(inst[1])
		case 'n': // noop
			// part1 clock sync
			if (clk+20)%40 == 0 {
				pwr += clk * X
			}
		}

	}

	// part1
	fmt.Println(pwr)
}

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(s string) int {
	sign := 1
	if s[0] == '-' {
		sign, s = -1, s[1:]
	}
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return sign * n
}
