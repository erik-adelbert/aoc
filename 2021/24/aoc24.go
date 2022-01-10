//
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	code := make([]string, 0, 256)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		code = append(code, line)
	}

	type item struct {
		j, b int
	}
	stack := make([]item, 0, 8)

	push := func(j, b int) {
		stack = append(stack, item{j, b})
	}

	pop := func() (int, int) {
		i := len(stack) - 1
		pop := stack[i]
		stack, stack[i] = stack[:i], item{}
		return pop.j, pop.b
	}

	p := 99999999999999
	q := 11111111111111

	for i := 0; i < 14; i++ {
		arg := strings.Fields(code[18*i+5])
		a, _ := strconv.Atoi(arg[len(arg)-1])

		arg = strings.Fields(code[18*i+15])
		b, _ := strconv.Atoi(arg[len(arg)-1])

		if a > 0 {
			push(i, b)
		} else {
			j, b := pop()

			xp, xq := 13-i, 13-i // exponents
			switch {
			case a > -b:
				xp += i - j
			case a < -b:
				xq += i - j
			}

			p -= abs(a+b) * pow(10, xp)
			q += abs(a+b) * pow(10, xq)
		}
	}
	fmt.Println(p) // part1
	fmt.Println(q) // part2
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func pow(a, n int) int {
	if n == 0 {
		return 1
	}

	b := 1
	for n > 1 {
		if (n & 1) == 1 {
			b = a * b
			a = a * a
			n = (n - 1) / 2
		} else {
			a = a * a
			n /= 2
		}
	}
	return a * b
}
