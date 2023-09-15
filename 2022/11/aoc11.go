package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type part struct {
	ctx *context
	fun func(n int) int
	lim int
}

func main() {
	ctx1 := new(context)

	// load inputs, compute running lcm
	m, n := 1, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		m = lcm(m, ctx1[n].load(input))
		n++
	}

	// parameterized part 1 & 2
	parts := []part{
		{
			ctx1,
			func(n int) int { return n / 3 },
			20,
		},
		{
			ctx1.clone(),
			func(n int) int { return n % m }, // capture m
			10_000,
		},
	}

	// run simulations
	for _, p := range parts {
		fmt.Println(p.solve())
	}
}

func (p part) solve() int {
	ctx := p.ctx
	for i := 0; i < p.lim; i++ {
		for j := range ctx {
			// single state part update
			ctx[j].update(p)
		}
	}

	max := [2]int{0, 0}

	// maintain 2 highest
	max2 := func(n int) {
		switch {
		case n >= max[0]: // accept duplicate
			max[1], max[0] = max[0], n
		case n > max[1]:
			max[1] = n
		}
	}

	for i := range ctx {
		max2(ctx[i].count)
	}
	return max[0] * max[1]
}

type context [8]state

func (c context) clone() *context {
	new := c
	for i := range c {
		new[i].items = make([]int, len(c[i].items))
		copy(new[i].items, c[i].items)
	}
	return &new
}

// state is aligned by go/fieldalignment
type state struct {
	items []int
	cmd   arit
	links [2]int
	mod   int
	count int
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

func (s *state) update(p part) {
	ctx, f := p.ctx, p.fun
	op, m := s.cmd, s.mod

	for _, x := range s.items {
		r := op.eval(f, x)
		nxt := s.links[0]
		if r%m > 0 {
			nxt = s.links[1]
		}
		ctx[nxt].items = append(ctx[nxt].items, r)
	}
	s.count += len(s.items)
	s.items = s.items[:0]
}

type arit struct {
	op   byte
	args [2]int
}

func (a arit) eval(f func(int) int, n int) int {
	// local args[1] defaults to n...
	args := [2]int{n, n}
	//... unless defined
	if a.args[1] >= 0 {
		args[1] = a.args[1]
	}

	x := 0
	switch a.op {
	case '+':
		x = args[0] + args[1]
	case '-':
		x = args[0] - args[1]
	case '*':
		x = args[0] * args[1]
	}

	return f(x)
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
