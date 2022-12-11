package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type arit struct {
	op   byte
	args [2]int
}

func (a arit) eval(f func(int) int, x int) int {
	// when a.arg[1] < 0 binds to x
	args := [2]int{x, x}
	if a.args[1] >= 0 {
		args[1] = a.args[1]
	}

	r := 0
	switch a.op {
	case '+':
		r = args[0] + args[1]
	case '-':
		r = args[0] - args[1]
	case '*':
		r = args[0] * args[1]
	}

	return f(r)
}

type state struct {
	cmd arit
	mod int

	count int
	links [2]int
	items []int
}

func (s *state) load(input *bufio.Scanner) int {
	for input.Scan() {
		var line string
		if line = strings.Trim(input.Text(), " "); len(line) == 0 {
			return s.mod
		}

		switch line[0] {
		case 'M':
			// discard name
		case 'S':
			items := strings.Split(line[16:], ", ")
			for _, v := range items {
				s.items = append(s.items, atoi(v))
			}
		case 'O':
			cmd := strings.Fields(line[17:])
			s.cmd.op = cmd[1][0]

			s.cmd.args[0], s.cmd.args[1] = -1, -1
			if cmd[0][0] != 'o' {
				s.cmd.args[0] = atoi(cmd[0])
			}
			if cmd[2][0] != 'o' {
				s.cmd.args[1] = atoi(cmd[2])
			}
		case 'T':
			s.mod = atoi(line[19:])
		case 'I':
			if line[3] == 't' {
				s.links[0] = atoi(line[25:])
			} else {
				s.links[1] = atoi(line[26:])
			}
		}
	}
	return s.mod
}

func (s *state) update(f func(int) int) {
	op, m := s.cmd, s.mod

	for _, x := range s.items {
		r := op.eval(f, x)
		nxt := s.links[0]
		if r%m > 0 {
			nxt = s.links[1]
		}
		states[nxt].items = append(states[nxt].items, r)
	}
	s.count += len(s.items)
	s.items = s.items[:0]
}

var states [8]state

func main() {
	m, n := 1, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		m = lcm(m, states[n].load(input))
		n++
	}

	part1 := func(n int) int { return n / 3 }
	part2 := func(n int) int { return n % m }

	max := [2]int{0, 0}

	max2 := func(n int) {
		switch {
		case n >= max[0]:
			max[1], max[0] = max[0], n
		case n >= max[1]:
			max[1] = n
		}
	}

	backup := states
	for i := range states {
		backup[i].items = make([]int, len(states[i].items))
		copy(backup[i].items, states[i].items)
	}

	for i := 0; i < 20; i++ {
		for j := range states {
			states[j].update(part1)
		}
	}

	for i := range states {
		max2(states[i].count)
	}

	// part 1
	fmt.Println(max[0] * max[1])

	max[0], max[1] = 0, 0

	for i := range states {
		states[i] = backup[i]
		states[i].items = make([]int, len(backup[i].items))
		copy(states[i].items, backup[i].items)
	}

	for i := 0; i < 10_000; i++ {
		for j := range states {
			states[j].update(part2)
		}
	}

	for i := range states {
		max2(states[i].count)
	}

	// part2
	fmt.Println(max[0] * max[1])
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
