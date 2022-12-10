package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// machine initial state
	// register X, clock and signal strength
	X := 1
	clk, sig := 0, 0

	// part2
	const (
		Black = ' '
		// undefined is very bright
		White = '\uFFFD'
	)

	// part2 decode/display message
	// synchronous scan display
	crt := func() {
		// sync column from clock
		c := clk%40 - 1
		if c < 0 {
			c += 40
		}

		// wrap beamer
		if c == 0 && clk > 1 {
			fmt.Println()
		}

		// beam-on window
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
		// fetch and tokenize instruction
		// fields:  0   1
		// values: cmd arg
		ins := strings.Fields(input.Text())

		clk++ // tick
		crt() // beam CRT

		// decode, monitor power, beam CRT, execute
		switch ins[0][0] {
		case 'a':
			// part1 sync signal monitoring
			switch (clk + 21) % 40 {
			default:
				clk++
			case 1:
				sig += clk * X
				clk++
			case 0:
				clk++
				sig += clk * X
			}

			crt() // beam CRT

			// addx
			X += atoi(ins[1])
		case 'n':
			// part1 sync signal monitoring
			if (clk+20)%40 == 0 {
				sig += clk * X
			}
			// noop
		}

	}

	// part1
	fmt.Println(sig)
}

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(s string) int {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return neg * n
}
