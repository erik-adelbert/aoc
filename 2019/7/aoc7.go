// aoc7.go --
// advent of code 2019 day 7
//
// https://adventofcode.com/2019/day/7
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

var TRACE bool // Global boolean variable

func main() {
	var code IntCode

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code = newIC(input.Text())

	ncpu := 5
	cpus := make([]*IntCodeCPU, ncpu)

	amplify := func(sigs []int) int {
		ins := make([]chan int, ncpu)

		for i := 0; i < ncpu; i++ {
			code := code.clone()
			ins[i] = make(chan int, 1)
			cpus[i] = newCPU(i, ins[i])

			go cpus[i].run(code) // run one amp per cpu
			ins[i] <- sigs[i]
		}

		var res int
		ins[0] <- 0

		done := 0
		for done != ncpu { // some running cpus
			for i, cpu := range cpus {
				// check output
				select {
				case r, ok := <-cpu.out:
					if !ok {
						// cpu is done
						cpu.out = nil
						done++
						continue
					}
					res = r // last output
					if ins[(i+1)%ncpu] != nil {
						ins[(i+1)%ncpu] <- r
					}
				default:
				}
			}
		}

		for _, ch := range ins {
			close(ch)
		}
		return res
	}

	sigmax := 0
	cur := []int{5, 6, 7, 8, 9}
	for cur != nil {
		sigmax = max(sigmax, amplify(cur))
		cur = next(cur)
	}
	fmt.Println(sigmax)
}

func next(arr []int) []int {
	// Find the rightmost element that is smaller than the element to its right
	i := len(arr) - 2
	for i >= 0 && arr[i] >= arr[i+1] {
		i--
	}

	// If no such element exists, we are at the last permutation
	if i < 0 {
		return nil
	}

	// Find the smallest element to the right of arr[i] that is larger than arr[i]
	j := len(arr) - 1
	for arr[j] <= arr[i] {
		j--
	}

	// Swap the two elements
	arr[i], arr[j] = arr[j], arr[i]

	// Reverse the sequence to the right of arr[i]
	reverse(arr[i+1:])

	return arr
}

// reverse reverses the elements of a slice in place
func reverse(arr []int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func init() {
	traceEnv := os.Getenv("TRACE")

	if traceEnv != "" {
		if val, err := strconv.ParseBool(traceEnv); err == nil {
			TRACE = val
		} else {
			TRACE = false
		}
	} else {
		TRACE = false
	}
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

// func (ic IntCode) patch(p map[int]int) IntCode {
// 	for i, v := range p {
// 		ic[i] = v
// 	}
// 	return ic
// }

type IntCodeCPU struct {
	id  int
	pc  int
	ins []<-chan int
	out chan int
}

func newCPU(id int, inputs ...<-chan int) *IntCodeCPU {
	cpu := &IntCodeCPU{id: id, out: make(chan int, 1)}
	cpu.connect(inputs...)
	return cpu
}

func (cpu *IntCodeCPU) connect(inputs ...<-chan int) {
	cpu.ins = append(cpu.ins, inputs...)
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

// func (cpu *IntCodeCPU) input(in int) *IntCodeCPU {
// 	cpu.in = in
// 	return cpu
// }

// func (cpu *IntCodeCPU) output() {
// 	if TRACE {
// 		fmt.Println("output:", cpu.out)
// 	} else {
// 		// fmt.Println(cpu.out)
// 	}
// }

func (cpu *IntCodeCPU) run(ic IntCode) IntCode {
	defer close(cpu.out)

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
			trace(cpu.id, opcode, ic[0], cpu.pc)
			// close(cpu.out)
			return ic // halt
		case JIT, JIF: // conditional jump
			var a [2]int
			for i, m := range pmods[:2] {
				a[i] = ic[cpu.pc+i+1]
				if m == 0 {
					a[i] = ic[a[i]]
				}
			}
			trace(cpu.id, opcode, pmods, a[:], ic[cpu.pc+1:cpu.pc+3])

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
			trace(cpu.id, opcode, pmods, a[0:2], ic[cpu.pc+1:cpu.pc+4])

			ic[a[2]] = μIC[opcode](a[0], a[1]) // execute
			cpu.pc += 4
		case INP:
			a := ic[cpu.pc+1]

			done := 0
		INSCAN:
			for done != len(cpu.ins) {
				for _, ch := range cpu.ins {
					select {
					case v, ok := <-ch:
						if !ok {
							done++
							continue
						}
						ic[a] = v
						trace(cpu.id, opcode, pmods, a, ic[a])
						break INSCAN
					default:
					}
				}
			}

			cpu.pc += 2
		case OUT: // unary op
			a := ic[cpu.pc+1]
			if pmods[0] == 0 {
				a = ic[a]
			}
			cpu.out <- a
			trace(cpu.id, opcode, pmods, ic[cpu.pc+1], a)

			// cpu.output()
			cpu.pc += 2
		}
	}
}

func trace(id, opcode int, args ...interface{}) (int, error) {
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
		return fmt.Printf("cpu%d: %s %d, pc: %d\n", id, ops[opcode], args[0], args[1])
	case INP, OUT:
		// trace(opcode, pmods, a, cpu.in)
		// trace(opcode, pmods, a, ic[a])
		pmods, a, ic := args[0].([]int), args[1].(int), args[2].(int)

		if pmods[0] == 0 {
			return fmt.Printf("cpu%d: %s $%d %d\t <- %d\n", id, ops[opcode], a, ic, ic)
		}
		return fmt.Printf("cpu%d: %s %d %d\t <- %d\n", id, ops[opcode], a, ic, ic)

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

		jmpfmt := "cpu%d: %s %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(jmpfmt, id, ops[opcode], dols[1], regs[1], dols[0], regs[0], a[0], a[1])

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

		binfmt := "cpu%d: %s $%d %s%d %s%d\t <- %d, %d\n"
		return fmt.Printf(binfmt, id, ops[opcode], dst, dols[0], regs[0], dols[1], regs[1], a[0], a[1])
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
