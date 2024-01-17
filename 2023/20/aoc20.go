package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {

	mods := make(modules, 64)

	// parse network
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := split(input.Text(), " -> ")

		var and bool
		switch args[0][0] {
		case '&':
			and = true
			fallthrough
		case '%':
			args[0] = args[0][1:] // extract id
		}

		mods[args[0]] = newModule(and, split(args[1], ", "))
	}

	c := newCircuit(mods)
	fmt.Println(c.npulse(), c.rx1())

}

type module struct {
	outs []string
	and  bool
}

func newModule(and bool, outs []string) *module {
	return &module{outs, and}
}

func (m *module) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v", *m)
	return sb.String()
}

type modules map[string]*module

type circuit [4]uint32

func newCircuit(mods modules) circuit {
	type node struct {
		id   string
		v, b uint32
	}

	todo := make([]node, 0)
	push := func(id string, v, b uint32) {
		todo = append(todo, node{id, v, b})
	}
	pop := func() (id string, v, b uint32) {
		var top node
		top, todo = todo[len(todo)-1], todo[:len(todo)-1]
		return top.id, top.v, top.b
	}

	for _, start := range mods["broadcaster"].outs {
		push(start, 0, 1)
	}

	fflops := make([]uint32, 0, 4)
JOBS:
	for len(todo) > 0 {
		id, v, b := pop()

		for _, nxt := range mods[id].outs {
			if !mods[nxt].and {
				if len(mods[id].outs) == 2 {
					v |= b
				}
				push(nxt, v, b<<1)
				continue JOBS
			}
		}

		fflops = append(fflops, v|b)
	}

	return circuit(fflops)
}

func (c circuit) npulse() uint32 {
	type io struct {
		i, o uint32
	}

	ios := make([]io, len(c))
	for i := range c {
		ios[i] = io{c[i], 13 - uint32(popcnt32(c[i]))}
	}

	lo, hi := uint32(5000), uint32(0)

	for n := 0; n < 1000; n++ {
		rising := uint32(^n & (n + 1))
		hi += 4 * popcnt32(rising)

		falling := uint32(n & ^(n + 1))
		lo += 4 * popcnt32(falling)

		for i := range ios {
			i, o := ios[i].i, ios[i].o
			λ := popcnt32(rising & i)
			hi += λ * (o + 3)
			lo += λ

			λ = popcnt32(falling & i)
			hi += λ * (o + 2)
			lo += λ * 2
		}
	}

	return lo * hi
}

func (c circuit) rx1() (Π uint64) {
	Π = 1
	for i := range c {
		Π *= uint64(c[i])
	}
	return
}

var split = strings.Split

func popcnt32(u uint32) uint32 {
	return uint32(bits.OnesCount32(u))
}

const DEBUG = false

func debug(format string, a ...any) {
	if DEBUG {
		fmt.Printf(format, a...)
	}
}
