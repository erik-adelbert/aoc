package main

import (
	"bufio"
	"bytes"
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
	// synchronous scan display to fb
	var fb bytes.Buffer
	crt := func() {
		// sync beamer xpos
		bmx := clk%40 - 1
		if bmx < 0 {
			bmx += 40
		}

		// CR/LF beamer
		if bmx == 0 && clk > 1 {
			fb.WriteByte('\n')
		}

		// beam-on window
		min, max := X-1, 39
		if X < max {
			max = X + 1
		}

		// beam
		pix := Black
		if min <= bmx && bmx <= max {
			// in range, light on!
			pix = White
		}
		fb.WriteByte(byte(pix))
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

	fmt.Println(sig)         // part1
	fmt.Println(fb.String()) // part2
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
