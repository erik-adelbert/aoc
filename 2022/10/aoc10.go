package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	X := 1
	clk, pwr := 0, 0

	y := 0
	frame := [6][40]byte{}

	beam := func() {
		x := clk%40 - 1
		if x < 0 {
			x += 40
		}

		if x == 0 && clk > 1 {
			y++
		}

		min := X - 1
		max := 39
		if X < 39 {
			max = X + 1
		}

		if min <= x && x <= max {
			frame[y][x] = 1
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fields := strings.Fields(input.Text())

		clk++
		beam()

		switch fields[0][0] {
		case 'a':
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

			beam()
			X += atoi(fields[1])
		case 'n':
			if (clk+20)%40 == 0 {
				pwr += clk * X
			}
		}

	}

	// part1
	fmt.Println(pwr)

	// part2
	for _, r := range frame {
		for _, c := range r {
			x := ' '
			if c > 0 {
				x = '�'
			}
			fmt.Printf("%c", x)
		}
		fmt.Println()
	}
}

// strconv.Atoi simplified core loop
// s is ^(-)?\d+$
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
