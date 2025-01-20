// aoc21.go --
// advent of code 2019 day 21
//
// https://adventofcode.com/2019/day/21
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
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code := newIC(input.Text())
	in := make(chan int, 1)

	writeln := func(s string) {
		for _, c := range s {
			in <- int(c)
		}
		in <- '\n'
	}

	spring := func(script []string) int {
		cpu := newCPU(0, in)
		go cpu.run(code)

		// input script
		fmt.Println(cpu.readline())
		for _, s := range script {
			fmt.Println(s)
			writeln(s)
		}

		// display ack
		_ = cpu.readline()
		fmt.Println(cpu.readline())
		_ = cpu.readline()

		// return damage
		return <-cpu.out
	}

	script1 := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"WALK",
	}

	script2 := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"NOT E T",
		"NOT T T",
		"OR H T",
		"AND T J",
		"RUN",
	}

	fmt.Println(spring(script1), spring(script2))
}

var TRACE bool

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

// func (ic IntCode) clone() IntCode {
// 	return append(IntCode(nil), ic...)
// }

// func (ic IntCode) patch(p map[int]int) IntCode {
// 	for i, v := range p {
// 		ic[i] = v
// 	}
// 	return ic
// }

type IntCodeCPU struct {
	id, pc, rbo, ram int

	in  <-chan int
	out chan int
}

func newCPU(id int, in <-chan int) *IntCodeCPU {
	cpu := &IntCodeCPU{id: id, in: in, out: make(chan int, 1)}
	return cpu
}

func (cpu *IntCodeCPU) readline() string {
	line := make([]byte, 0, 32)
LINESCAN:
	for v := range cpu.out {
		v := byte(v)
		switch v {
		case '\n':
			break LINESCAN
		default:
			line = append(line, v)
		}
	}
	return string(line)
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
			// reallocate memory
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

	const (
		ABS = "$"
		IMM = ""
		REL = "."
	)

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

	printf := func(format string, args ...interface{}) (int, error) {
		format = fmt.Sprintf("cpu%d: %s ", id, ops[opcode]) + format
		return fmt.Printf(format, args...)
	}

	switch opcode {
	case HLT:
		return printf("%d, pc: %d\n", args[0], args[1])

	case RBO:
		pmods, a, ic, rbo := args[0].([]int), args[1].([]int), args[2].(IntCode), args[3].(int)
		sym := []string{ABS, IMM, REL}[pmods[0]]
		src := []int{ic[0], a[0], ic[0]}[pmods[0]]

		return printf("%s%d\t <- %d (%d)\n", sym, src, a[0], rbo)

	case INP, OUT:
		pmods, a, ic := args[0].([]int), args[1].(int), args[2].(int)
		sym := []string{ABS, IMM, ABS}[pmods[0]]

		return printf("%s%d %d\t <- %d\n", sym, a, ic, ic)

	case JIF, JIT:
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)
		sym := [2]string{ABS, ABS}

		var src [2]int
		for i, m := range pmods[:2] {
			src[i] = ic[i]
			switch m {
			case 1:
				sym[i] = IMM
				src[i] = a[i]
			case 2:
				sym[i] = REL
			}
		}

		return printf("%s%d %s%d\t <- %d, %d\n", sym[0], src[0], sym[1], src[1], a[0], a[1])

	case ADD, MUL, LT, EQ:
		pmods, a, ic := args[0].([]int), args[1].([]int), args[2].(IntCode)
		sym := [3]string{ABS, ABS, ABS}

		var src [2]int
		for i, m := range pmods[:2] {
			src[i] = ic[i]
			switch m {
			case 1:
				sym[i] = IMM
				src[i] = a[i]
			case 2:
				sym[i] = REL
			}
		}

		dst := ic[2]
		if pmods[2] == 2 {
			sym[2] = REL
		}

		return printf("%s%d %s%d %s%d\t <- %d, %d\n", sym[2], dst, sym[0], src[0], sym[1], src[1], a[0], a[1])
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
