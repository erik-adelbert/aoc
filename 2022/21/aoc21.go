package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type val struct {
	op, s string
	a     []*val
	n, x  int
}

// PROG stores instructions
var PROG map[string]*val

func main() {
	PROG = make(map[string]*val, 2048)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cmd := strings.Fields(r.Replace(input.Text()))
		load(cmd)
	}

	fmt.Println(PROG["root"].eval())
	fmt.Println(solve(PROG["root"], "humn"))
}

func (v *val) eval() int {
	if v.op == "VAR" {
		v = PROG[v.s] // deref var
	}

	if v.op == "INT" {
		return v.n
	}

	a, b := v.a[0], v.a[1]
	switch v.op {
	case "+":
		return a.eval() + b.eval()
	case "-":
		return a.eval() - b.eval()
	case "*":
		return a.eval() * b.eval()
	case "/":
		return a.eval() / b.eval()
	}

	panic("unreachable")
}

func solve(v *val, k string) int {
	d := v
	if v.op == "VAR" {
		d = PROG[v.s] // deref var
	}

	// mark symbolic path
	mksym(d, k)

	// walk along and force symbolic values
	// eval the rest
	switch d.x {
	case Left:
		n := d.a[1].eval()
		return d.a[0].force(n)
	case Right:
		n := d.a[0].eval()
		return d.a[1].force(n)
	}

	panic("unreachable")
}

// symbolic left/right hand side
const (
	_ = iota
	Left
	Right
)

// mark symbolic code path
func mksym(v *val, k string) bool {
	d := v
	if d.op == "VAR" {
		d = PROG[v.s] // deref var
	}

	if d.op == "INT" {
		return v.s == k
	}

	if mksym(d.a[0], k) {
		d.x = Left
		return true
	}

	if mksym(d.a[1], k) {
		d.x = Right
		return true
	}

	return false
}

// force values along symbolic path
func (v *val) force(n int) int {
	d := v
	if d.op == "VAR" {
		d = PROG[v.s] // deref var
	}

	if d.op == "INT" {
		return n
	}

	// reverse symbolic side ops
	switch d.x {
	case Left:
		l := d.a[0]
		r := d.a[1].eval()
		switch d.op {
		case "+":
			return l.force(n - r)
		case "-":
			return l.force(n + r)
		case "*":
			return l.force(n / r)
		case "/":
			return l.force(n * r)
		}
	case Right:
		l := d.a[0].eval()
		r := d.a[1]
		switch d.op {
		case "+":
			return r.force(n - l)
		case "-":
			return r.force(l - n)
		case "*":
			return r.force(n / l)
		case "/":
			return r.force(l / n)
		}
	}

	panic("unreachable")
}

func load(args []string) {
	v := new(val)
	switch len(args) {
	case 2:
		// var
		v.op = "INT"
		v.n = atoi(args[1])
	default:
		// binary op
		v.op = args[2]
		v.a = append(
			v.a,
			&val{op: "VAR", s: args[1]},
			&val{op: "VAR", s: args[3]},
		)
	}

	PROG[args[0]] = v
	return
}

func (v *val) int() int {
	return v.n
}

func (v *val) String() string {
	var sb strings.Builder
	sb.WriteString(v.op + "(")
	switch v.op {
	default:
		sb.WriteString(
			fmt.Sprintf("%v, %v", v.a[0], v.a[1]),
		)
	case "INT":
		sb.WriteString(
			fmt.Sprint(v.n),
		)
	case "VAR":
		sb.WriteString(
			fmt.Sprint(v.s),
		)
	}

	switch v.x {
	case Left:
		sb.WriteRune('L')
	case Right:
		sb.WriteRune('R')
	}
	sb.WriteString(")")
	return sb.String()
}

var r = strings.NewReplacer(
	":", "",
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) int {
	var n int

	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}
