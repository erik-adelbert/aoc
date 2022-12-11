package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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
	// synchronous scan display to fb
	var fb bytes.Buffer
	crt := func() {
		// sync beamer xpos
		bmx := clk%40 - 1
		if bmx < 0 {
			bmx += 40
		}

		// CR/LF beamer, slightly incorrect but ok
		if bmx == 0 {
			fb.WriteByte('\n')
		}

		// beam
		pix := Black
		if abs(X-bmx) <= 1 {
			// in range, light on!
			pix = White
		}
		fb.WriteByte(byte(pix))
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// fetch and tokenize instruction
		ins := input.Text()

		clk++ // tick
		crt() // beam CRT

		// decode, monitor power, beam CRT, execute
		switch ins[0] {
		case 'a':
			// part1 sync signal monitoring
			switch clk % 40 {
			default:
				clk++
			case 19:
				clk++
				sig += clk * X
			case 20:
				sig += clk * X
				clk++
			}

			crt() // beam CRT

			// addx
			X += atoi(ins[5:])
		case 'n':
			// part1 sync signal monitoring
			if clk%40 == 20 {
				sig += clk * X
			}
			// noop
		}

	}

	fmt.Print(sig)           // part1
	fmt.Println(fb.String()) // part2
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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
