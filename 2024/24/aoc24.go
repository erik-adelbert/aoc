// aoc24.go --
// advent of code 2024 day 24
//
// https://adventofcode.com/2024/day/24
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-24: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	JTAG   = false
	TRACE  = false
	NLATCH = 312
	NCOMP  = 222
	NBIT   = 45
	NFAULT = 8
)

func main() {
	setup := make(map[string]bool, 2*NBIT)
	latches := make(map[string]*conn, NLATCH)
	circuit := make(map[string]LC, NCOMP)
	logics := make([]Logic, 0, NCOMP)

	mode := INIT
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case line == "":
			mode = WIRE

		case mode == INIT:
			args := strings.Split(line, ": ")

			a, b := args[0], args[1]
			setup[a] = b == "1"
			latches[a] = newConn(nil, a, false)

			trace("Setup %s to %t\n", a, setup[a])

		case mode == WIRE:
			args := strings.Fields(line)

			a, b, c, op := args[0], args[2], args[4], args[1]
			for _, arg := range []string{a, b, c} {
				if _, ok := latches[arg]; !ok {
					latches[arg] = newConn(nil, arg, false)
				}
			}

			// connect
			gate := mkgate(c, op)
			circuit[c] = &gate
			latches[a].wire(gate.A)
			latches[b].wire(gate.B)
			gate.C.wire(latches[c])

			// store
			logics = append(logics, Logic{op: op, A: a, B: b, C: c})
			trace("Store %s %s %s -> %s\n", a, op, b, c)
		}
	}

	// setup & evaluate
	for reg, value := range setup {
		latches[reg].set(value, true)
	}

	z := out('z', latches)

	faults := inspect(logics)
	slices.Sort(faults)

	fmt.Println(z, strings.Join(faults, ",")) // part 1 & 2
}

/* Our 1bit full adder:
A ---+---------+
     |          \ (2)
	 |		    XOR --- + -------+
     |          /       |         \ (1)
B ------+------+        |         XOR -------------> S
     |  |               |         /
C_in ----------------------------+
     |  |               |         \
	 |	|				|	      AND -----+
	 |	|				|         /         \
	 |	|				+--------+           \
	 |	|				              (3)    OR ---> C_out
	 +-------------------------+		     /
		|                       \		    /
		|					     AND ------+
		|					    /
        +----------------------+
*/
// statically inspect the circuit for faults
func inspect(logics []Logic) []string {

	is_swapped :=
		func(a Logic) bool {
			op := "XOR"
			if a.op == "AND" {
				op = "OR"
			}
			for _, b := range logics {
				if a == b {
					continue
				}
				if b.op == op && (b.A == a.C || b.B == a.C) {
					return false
				}
			}
			return true
		}

	faults := make([]string, 0, NFAULT)
	for _, a := range logics {
		op, in0, in1, out := a.op, a.A, a.B, a.C
		switch {
		case out[0] == 'z' && out != "z45" && op != "XOR":
			// fault in the output stage (1)
			fallthrough
		case out[0] != 'z' && op == "XOR" && in0[0] != 'x' && in0[0] != 'y' && in1[0] != 'x' && in1[0] != 'y':
			// fault in the input stage (2)
			faults = append(faults, a.C)
		case (op == "XOR" || op == "AND") && (in0[0] == 'x' || in0[0] == 'y') && (in1[0] == 'x' || in1[0] == 'y'):
			// fault in the carry stage (3)
			if in0[1:] != "00" && in1[1:] != "00" {
				if is_swapped(a) {
					faults = append(faults, a.C)
				}
			}
		}
	}
	return faults
}

func out(varname byte, latches map[string]*conn) int {
	type reg struct {
		id  string
		val bool
	}

	digits := make([]reg, 0, len(latches))
	for id, c := range latches {
		if id[0] == varname {
			digits = append(digits, reg{id: id, val: c.val})
		}
	}
	slices.SortFunc(digits, func(a, b reg) int {
		return -strings.Compare(a.id, b.id)
	})

	n := 0
	for _, x := range digits {
		n <<= 1
		if x.val {
			n++
		}
	}
	return n
}

const (
	INIT = iota
	WIRE
)

type conn struct {
	own  LC
	name string
	outs []*conn
	val  bool
	jtag bool
	wkup bool
}

func newConn(owner LC, name string, wakeup bool) *conn {
	return &conn{val: false, own: owner, name: name, jtag: JTAG, wkup: wakeup}
}

func (c *conn) wire(inputs ...*conn) {
	c.outs = append(c.outs, inputs...)
}

func (c *conn) set(val, force bool) {
	if !force && c.val == val {
		return
	}
	c.val = val
	if c.wkup {
		c.own.evaluate()
	}
	if c.jtag {
		owner := fmt.Sprintf("m%s", c.name)
		if c.own != nil {
			owner = c.own.name()
		}
		fmt.Printf("Connector %s-%s set to %t\n", owner, c.name, c.val)
	}
	for _, con := range c.outs {
		(*con).set(val, false)
	}
}

type LC interface {
	name() string
	evaluate()
}

type Gate struct {
	op func(bool, bool) bool
	A  *conn
	B  *conn
	C  *conn
	id string
}

func (g *Gate) name() string {
	return g.id
}

func (g *Gate) evaluate() {
	g.C.set(g.op(g.A.val, g.B.val), false)
}

func mkgate(id string, opcode string) Gate {

	var op func(bool, bool) bool
	switch opcode {
	case "AND":
		op = func(a, b bool) bool { return a && b }
	case "OR":
		op = func(a, b bool) bool { return a || b }
	case "XOR":
		op = func(a, b bool) bool { return a != b }
	}

	gate := Gate{id: id, op: op}
	gate.A = newConn(&gate, "A", true)
	gate.B = newConn(&gate, "B", true)
	gate.C = newConn(&gate, "C", false)
	return gate
}

type Logic struct {
	op      string
	A, B, C string
}

func trace(format string, args ...interface{}) {
	if TRACE {
		fmt.Printf(format, args...)
	}
}
