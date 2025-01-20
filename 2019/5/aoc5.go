// aoc5.go --
// advent of code 2019 day 5
//
// https://adventofcode.com/2019/day/5
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
	"strings"
)

const TRACE = false

func main() {
	var code IntCode

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code = newIC(input.Text())

	cpu := newCPU()

	cpu.input(1).run(code.clone())
	p1 := cpu.out

	cpu.input(5).run(code)
	p2 := cpu.out

	fmt.Println(p1, p2)
}

type IntCode []int

func newIC(s string) IntCode {
	words := strings.Split(s, ",")

	ic := make(IntCode, 0, len(words))
	for _, w := range words {
		ic = append(ic, atoi(w))
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
	pc, in, out int
}

func newCPU() *IntCodeCPU {
	return &IntCodeCPU{}
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

func (cpu *IntCodeCPU) input(in int) *IntCodeCPU {
	cpu.in = in
	return cpu
}

func (cpu *IntCodeCPU) output() {
	if TRACE {
		fmt.Println("output:", cpu.out)
	} else {
		// fmt.Println(cpu.out)
	}
}

func (cpu *IntCodeCPU) run(ic IntCode) IntCode {

	pmods := make([]int, 0, 3)

	// busy loop
	cpu.pc = 0
	for {
		// fetch

		opcode := ic[cpu.pc] % 100
		pmods = pmods[:0]
		mode := ic[cpu.pc] / 100
		for i := 0; i < 3; i++ {
			pmods = append(pmods, mode%10)
			mode /= 10
		}

		switch opcode { // decode
		case HLT:
			trace(opcode, ic[0], cpu.pc)
			return ic // halt
		case JIT, JIF: // conditional jump
			var a [2]int
			for i, m := range pmods[:2] {
				a[i] = ic[cpu.pc+i+1]
				if m == 0 {
					a[i] = ic[a[i]]
				}
			}
			trace(opcode, pmods, a[:], ic[cpu.pc+1:cpu.pc+3])

			switch {
			case opcode == JIT && a[0] != 0, opcode == JIF && a[0] == 0:
				cpu.pc = a[1] // jump
			default:
				cpu.pc += 3
			}

		case ADD, MUL, LT, EQ: // binary op
			var a [3]int
			for i, m := range pmods[:2] {
				a[i] = ic[cpu.pc+i+1]
				if m == 0 {
					a[i] = ic[a[i]]
				}
			}
			a[2] = ic[cpu.pc+3]
			trace(opcode, pmods, a[0:2], ic[cpu.pc+1:cpu.pc+4])

			ic[a[2]] = μIC[opcode](a[0], a[1]) // execute
			cpu.pc += 4
		case INP:
			a := ic[cpu.pc+1]
			ic[a] = cpu.in
			trace(opcode, pmods, a, cpu.in)

			cpu.pc += 2
		case OUT: // unary op
			a := ic[cpu.pc+1]
			if pmods[0] == 0 {
				a = ic[a]
			}
			cpu.out = a
			trace(opcode, pmods, ic[cpu.pc+1], a)

			cpu.output()
			cpu.pc += 2
		}
	}
}

func trace(opcode int, args ...interface{}) (int, error) {
	ops := []string{
		ADD: "ADD",
		MUL: "MUL",
		HLT: "HLT",
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
		return fmt.Printf("%s %d, pc: %d\n", ops[opcode], args[0], args[1])
	case INP, OUT:
		// trace(opcode, pmods, a, cpu.in)
		// trace(opcode, pmods, a, ic[a])
		pmods, a, ic := args[0].([]int), args[1].(int), args[2].(int)

		if pmods[0] == 0 {
			return fmt.Printf("%s $%d %d\t <- %d\n", ops[opcode], a, ic, ic)
		}
		return fmt.Printf("%s %d %d\t <- %d\n", ops[opcode], a, ic, ic)

	case JIF, JIT:
		// trace(opcode, pmods, a, ic[cpu.pc+1:cpu.pc+3])
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)

		dols := [2]string{"$", "$"}
		var regs [2]int

		for i, m := range pmods[:2] {
			regs[i] = ic[i]
			if m == 1 {
				dols[i] = ""
				regs[i] = a[i]
			}
		}

		jmpfmt := "%s %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(jmpfmt, ops[opcode], dols[1], regs[1], dols[0], regs[0], a[0], a[1])

	case ADD, MUL, LT, EQ:
		// trace(opcode, 0: pmods, 1: a[0:2], 2: ic[cpu.pc+1:cpu.pc+3])
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)

		dols := [2]string{"$", "$"}
		var regs [2]int

		for i, m := range pmods[:2] {
			regs[i] = ic[i]
			if m == 1 {
				dols[i] = ""
				regs[i] = a[i]
			}
		}
		dst := ic[2]

		binfmt := "%s $%d %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(binfmt, ops[opcode], dst, dols[0], regs[0], dols[1], regs[1], a[0], a[1])
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
