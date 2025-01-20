// aoc23.go --
// advent of code 2019 day 23
//
// https://adventofcode.com/2019/day/23
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
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	code := newIC(input.Text())

	ncpu := 50

	// build network
	done := make(chan struct{})
	cpus := make([]*cpu, ncpu)

	// network select cases to read from each cpu
	conns := make([]reflect.SelectCase, ncpu)

	// create cpus and output connectors
	for i := 0; i < ncpu; i++ {
		cpus[i], conns[i] = newCPU(i, code.clone())
		go cpus[i].run(done)
	}
	conns = append(conns, reflect.SelectCase{Dir: reflect.SelectDefault})
	IDLE := len(conns) - 1

	nat := make(chan packet, 1)
	defer close(nat)
	go func() {
		// NAT anonymous logic
		defer close(done)

		var p, p0 packet
		var active bool

		for {
			runtime.Gosched() // throttle the NAT

			select {
			case p = <-nat:
				fmt.Println("NAT: receiving", p)
				p.src = 255
				active = true
			default:
				nidle := 0
				for i := 0; i < ncpu; i++ {
					if cpus[i].ready() {
						nidle++
					}
				}
				if nidle == ncpu {
					if active && p0 == p {
						fmt.Println("NAT: repeating", p0)
						return
					}
					p0, active = p, false

					fmt.Println("NAT: waking up cpu00", p)
					cpus[0].recv(p)
				}
			}
		}
	}() // launch the NAT

NETWORK: // loop
	for {
		select {
		case <-done:
			break NETWORK // exit ordered from NAT
		default:
		}

		id, val, _ := reflect.Select(conns)

		switch id {
		case IDLE:
			// nothing to do, yield!
			runtime.Gosched()
		default:
			// route incoming packet
			dst := int(val.Int())
			tgt := fmt.Sprintf("cpu%02d", dst)
			if dst == 255 {
				tgt = "NAT"
			}
			p := packet{id, <-cpus[id].nic.out, <-cpus[id].nic.out}
			fmt.Printf("net: routing %v to %s\n", p, tgt)
			switch dst {
			case 255:
				nat <- p
			default:
				cpus[dst].recv(p)
			}
		}
	}
}

type packet struct {
	src, x, y int
}

func (p packet) String() string {
	if p.src == 255 {
		return fmt.Sprintf("NAT{%d, %d}", p.x, p.y)
	}
	return fmt.Sprintf("cpu%02d{%d, %d}", p.src, p.x, p.y)
}

type queue []packet

func mkqueue() queue {
	return make(queue, 0, 150)
}

func (q queue) pop() (packet, queue) {
	return q[0], q[1:]
}

type cpu struct {
	nic *IntCodeCPU
	in  chan int
	q   queue
	sync.Mutex
	idle bool
}

func newCPU(id int, code IntCode) (*cpu, reflect.SelectCase) {
	in := make(chan int, 2)
	nic := newIntCodeCPU(id, in)
	go nic.run(code)

	// network select ready case for this cpu nic
	sc := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(nic.out),
	}

	in <- id // send address

	return &cpu{q: mkqueue(), nic: nic, in: in}, sc
}

func (c *cpu) run(quit chan struct{}) {
	fmt.Printf("cpu%02d: running\n", c.nic.id)
	for {
		runtime.Gosched()

		select {
		case <-quit:
			close(c.in)
			return
		default:
			if c.nic.blocking.Load() {
				c.Lock()
				if len(c.q) == 0 {
					c.idle = true
					c.Unlock()

					c.in <- -1
				} else {
					var p packet

					c.idle = false
					p, c.q = c.q.pop()
					c.Unlock()

					c.in <- p.x
					c.in <- p.y
				}
			}
		}
	}
}

func (c *cpu) ready() bool {
	if c.TryLock() {
		defer c.Unlock()
		return c.idle
	}
	return false
}

func (c *cpu) recv(p packet) {
	c.Lock()
	defer c.Unlock()

	c.idle = false
	c.q = append(c.q, p)
	fmt.Printf("cpu%02d: received %v\n", c.nic.id, p)
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
	blocking *atomic.Bool
	in       <-chan int
	out      chan int
	id       int
	pc       int
	rbo      int
	ram      int
}

func newIntCodeCPU(id int, in <-chan int) *IntCodeCPU {
	cpu := &IntCodeCPU{id: id, in: in, out: make(chan int, 1), blocking: new(atomic.Bool)}
	return cpu
}

// func (cpu *IntCodeCPU) readline() string {
// 	line := make([]byte, 0, 32)
// LINESCAN:
// 	for v := range cpu.out {
// 		v := byte(v)
// 		switch v {
// 		case '\n':
// 			break LINESCAN
// 		default:
// 			line = append(line, v)
// 		}
// 	}
// 	return string(line)
// }

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
	EOT = 100
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
			var ok bool

			cpu.blocking.Store(true)
			a := getdst(ic[cpu.pc+1], pmods[0])
			if ic[a], ok = <-cpu.in; !ok {
				trace(cpu.id, EOT, ic[0], cpu.pc)
				return ic
			}
			trace(cpu.id, opcode, pmods, a, ic[a])

			cpu.blocking.Store(false)
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
		EOT: "EOT",
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
		format = fmt.Sprintf("cpu%02d: %s ", id, ops[opcode]) + format
		return fmt.Printf(format, args...)
	}

	switch opcode {
	case EOT, HLT:
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
