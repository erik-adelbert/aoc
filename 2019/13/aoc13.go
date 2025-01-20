// aoc13.go --
// advent of code 2019 day 13
//
// https://adventofcode.com/2019/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-1: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var TRACE bool // Global boolean variable

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code0 := newIC(input.Text())
	code1 := code0.clone().patch(map[int]int{0: 2})

	ncpu := 2

	ins := make([]chan int, ncpu)
	for i := range ins {
		ins[i] = make(chan int, 1)
	}
	defer func() {
		for _, in := range ins {
			close(in)
		}
	}()

	cpus := make([]*IntCodeCPU, ncpu)
	for i := range cpus {
		cpus[i] = newCPU(i, ins[i])
	}

	go cpus[0].run(code0)
	go cpus[1].run(code1)

	nblock, score := 0, 0

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for x := range cpus[0].out {
			y, t := <-cpus[0].out, <-cpus[0].out
			if t == 2 {
				nblock++
			}
			_, _ = x, y
		}
	}()

	go func() {
		defer wg.Done()

		paddle := [2]int{0, 0}

		cmd := func(x int) int {
			switch {
			case x < paddle[0]:
				return -1
			case x > paddle[0]:
				return 1
			}
			return 0
		}

		for x := range cpus[1].out {
			y, t := <-cpus[1].out, <-cpus[1].out
			switch {
			case x < 0:
				score = t
			case t == 3:
				paddle = [2]int{x, y}
			case t == 4:
				ins[1] <- cmd(x)
			}
		}
	}()

	wg.Wait()

	fmt.Println(nblock, score)
}

func init() {
	env := os.Getenv("TRACE")

	TRACE = false
	if env != "" {
		if val, err := strconv.ParseBool(env); err == nil {
			TRACE = val
		}
	}
}

type IntCode []int

func newIC(s string) IntCode {
	words := strings.Split(s, ",")

	ic := make(IntCode, len(words))
	for i, w := range words {
		ic[i] = atoi(w)
	}
	return ic
}

func (ic IntCode) clone() IntCode {
	return append(IntCode(nil), ic...)
}

func (ic IntCode) patch(p map[int]int) IntCode {
	for i, v := range p {
		ic[i] = v
	}
	return ic
}

type IntCodeCPU struct {
	id, pc, rbo, ram int

	in  <-chan int
	out chan int
}

func newCPU(id int, in <-chan int) *IntCodeCPU {
	cpu := &IntCodeCPU{id: id, in: in, out: make(chan int, 1)}
	return cpu
}

const (
	ADD = iota + 1
	MUL
	INP
	OUT
	JIT
	JIF
	LT
	EQ
	RBO
	HLT = 99
)

type μcmd func(...int) int

var μIC = []μcmd{
	ADD: func(args ...int) int { return args[0] + args[1] }, // add
	MUL: func(args ...int) int { return args[0] * args[1] }, // mul
	LT: func(args ...int) int {
		if args[0] < args[1] {
			return 1
		}
		return 0
	}, // less than
	EQ: func(args ...int) int {
		if args[0] == args[1] {
			return 1
		}
		return 0
	},
}

func (cpu *IntCodeCPU) run(ic IntCode) IntCode {
	defer close(cpu.out)
	cpu.ram = len(ic)

	pmods := make([]int, 0, 3)

	decode := func(code int) int {
		opcode := code % 100
		pmods = pmods[:0]
		mode := code / 100
		for i := 0; i < 3; i++ {
			pmods = append(pmods, mode%10)
			mode /= 10
		}
		return opcode
	}

	getdst := func(pc, mode int) int {
		if mode == 2 {
			pc += cpu.rbo
		}
		for len(ic) <= pc {
			ic = append(ic, make(IntCode, len(ic))...)
		}
		return pc
	}

	getargs := func(n int) []int {
		args := make([]int, n)
		for i := 0; i < n; i++ {
			args[i] = ic[cpu.pc+i+1]
			if pmods[i]&1 == 0 {
				args[i] = ic[getdst(args[i], pmods[i])]
			}
		}
		return args
	}

	// busy loop
	cpu.pc = 0
	for {
		// fetch
		opcode := decode(ic[cpu.pc])

		switch opcode {
		case HLT:
			trace(cpu.id, opcode, ic[0], cpu.pc)
			return ic // halt
		case RBO:
			a := getargs(1)
			cpu.rbo += a[0]
			trace(cpu.id, opcode, pmods, a, ic[cpu.pc+1:cpu.pc+2], cpu.rbo)

			cpu.pc += 2
		case JIT, JIF: // conditional jump
			a := getargs(2)
			trace(cpu.id, opcode, pmods, a, ic[cpu.pc+1:cpu.pc+4])

			switch {
			case opcode == JIT && a[0] != 0, opcode == JIF && a[0] == 0:
				cpu.pc = a[1] // jump
			default:
				cpu.pc += 3
			}

		case ADD, MUL, LT, EQ: // binary op
			a := getargs(2)
			dst := getdst(ic[cpu.pc+3], pmods[2])

			trace(cpu.id, opcode, pmods, a, ic[cpu.pc+1:cpu.pc+4])

			ic[dst] = μIC[opcode](a[0], a[1]) // execute
			cpu.pc += 4
		case INP:
			a := getdst(ic[cpu.pc+1], pmods[0])
			ic[a] = <-cpu.in
			trace(cpu.id, opcode, pmods, a, ic[a])
			cpu.pc += 2
		case OUT: // unary op
			a := getargs(1)
			trace(cpu.id, opcode, pmods, ic[cpu.pc+1], a[0])

			cpu.out <- a[0]
			cpu.pc += 2
		}
	}
}

func trace(id, opcode int, args ...interface{}) (int, error) {
	ops := []string{
		ADD: "ADD",
		MUL: "MUL",
		HLT: "HLT",
		RBO: "RBO",
		INP: "INP",
		OUT: "OUT",
		JIT: "JIT",
		JIF: "JIF",
		LT:  "LT",
		EQ:  "EQ",
	}

	if !TRACE {
		return 0, nil
	}

	switch opcode {
	case HLT:
		// trace(opcode, ic[0], cpu.pc)
		return fmt.Printf("cpu%d: %s %d, pc: %d\n", id, ops[opcode], args[0], args[1])
	case RBO:
		//trace(cpu.id, opcode, pmods, a, ic[cpu.pc+1:cpu.pc+2])
		pmods, a, ic, rbo := args[0].([]int), args[1].([]int), args[2].(IntCode), args[3].(int)
		sym := []string{"$", "", "@"}[pmods[0]]
		src := []int{ic[0], a[0], ic[0]}[pmods[0]]

		return fmt.Printf("cpu%d: %s %s%d\t <- %d (%d)\n", id, ops[opcode], sym, src, a[0], rbo)

	case INP, OUT:
		// trace(opcode, pmods, a, cpu.in)
		// trace(opcode, pmods, a, ic[a])
		pmods, a, ic := args[0].([]int), args[1].(int), args[2].(int)
		sym := []string{"$", "", "$"}[pmods[0]]

		return fmt.Printf("cpu%d: %s %s%d %d\t <- %d\n", id, ops[opcode], sym, a, ic, ic)

	case JIF, JIT:
		// trace(opcode, pmods, a, ic[cpu.pc+1:cpu.pc+3])
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)

		dols := [2]string{"$", "$"}
		var regs [2]int

		for i, m := range pmods[:2] {
			regs[i] = ic[i]
			switch m {
			case 1:
				dols[i] = ""
				regs[i] = a[i]
			case 2:
				dols[i] = "@"
			}
		}

		jmpfmt := "cpu%d: %s %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(jmpfmt, id, ops[opcode], dols[0], regs[0], dols[1], regs[1], a[0], a[1])

	case ADD, MUL, LT, EQ:
		// trace(opcode, 0: pmods, 1: a[0:2], 2: ic[cpu.pc+1:cpu.pc+3])
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)

		dols := [3]string{"$", "$", "$"}
		var regs [2]int

		for i, m := range pmods[:2] {
			regs[i] = ic[i]
			switch m {
			case 1:
				dols[i] = ""
				regs[i] = a[i]
			case 2:
				dols[i] = "@"
			}
		}
		dst := ic[2]
		if pmods[2] == 2 {
			dols[2] = "@"
		}

		binfmt := "cpu%d: %s %s%d %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(binfmt, id, ops[opcode], dols[2], dst, dols[0], regs[0], dols[1], regs[1], a[0], a[1])
	}

	return 0, nil
}

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(s string) (n int) {
	neg := 1
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return neg * n
}
