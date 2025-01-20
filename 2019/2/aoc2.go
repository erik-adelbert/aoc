// aoc2.go --
// advent of code 2019 day 2
//
// https://adventofcode.com/2019/day/2
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
	var p1, p2 int

	const (
		NOUN = iota + 1
		VERB
	)

	const (
		V    = 2
		N    = 12
		GOAL = 19690720
	)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code = newIC(input.Text())
	cpu := newCPU()

	noun := cpu.run(code.clone().patch(map[int]int{NOUN: 1}))[0]
	verb := cpu.run(code.clone().patch(map[int]int{VERB: 1}))[0]

	fg := GOAL
	f0 := cpu.run(code)[0]
	δn, δv := noun-f0, verb-f0
	Δ := fg - f0

	p1 = f0 + N*δn + V*δv
	p2 = 100*Δ/δn + (Δ%δn)/δv

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
	pc int
}

func newCPU() *IntCodeCPU {
	return &IntCodeCPU{}
}

const (
	ADD = iota + 1
	MUL
	HLT = 99
)

type μcmd func(...int) int

var μIC = []μcmd{
	ADD: func(args ...int) int { return args[0] + args[1] }, // add
	MUL: func(args ...int) int { return args[0] * args[1] }, // mul
}

func (cpu *IntCodeCPU) run(ic IntCode) IntCode {
	// busy loop
	cpu.pc = 0
	for {
		opcode := ic[cpu.pc] // fetch

		switch opcode { // decode
		case HLT:
			trace(opcode, ic[0], cpu.pc)
			return ic // halt
		case ADD, MUL: // binary op
			a := ic[cpu.pc+1 : cpu.pc+4] // addresses
			trace(opcode, a[2], a[0], a[1], ic[a[0]], ic[a[1]])
			ic[a[2]] = μIC[opcode](ic[a[0]], ic[a[1]]) // execute
			cpu.pc += 4
		}
	}
}

func trace(opcode int, args ...interface{}) (int, error) {
	ops := []string{
		ADD: "ADD",
		MUL: "MUL",
		HLT: "HLT",
	}

	BINFMT := "%s $%d $%d $%d\t <- %d, %d\n"

	fmts := []string{
		ADD: BINFMT,
		MUL: BINFMT,
		HLT: "%s %d, pc: %d\n",
	}

	if TRACE {
		format := fmts[opcode]
		args = append([]interface{}{ops[opcode]}, args...)
		return fmt.Printf(format, args...)
	}
	return 0, nil
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
