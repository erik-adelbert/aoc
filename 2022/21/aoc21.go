// aoc21.go --
// advent of code 2022 day 21
//
// https://adventofcode.com/2022/day/21
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-21: initial commit

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

	fmt.Println(eval(PROG["root"]))
	fmt.Println(solve(PROG["root"], "humn"))
}

func eval(v *val) int {
	if v.op == "VAR" {
		v = PROG[v.s] // deref var
	}

	if v.op == "INT" {
		return v.n
	}

	a, b := v.a[0], v.a[1]
	switch v.op {
	case "+":
		return eval(a) + eval(b)
	case "-":
		return eval(a) - eval(b)
	case "*":
		return eval(a) * eval(b)
	case "/":
		return eval(a) / eval(b)
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
		n := eval(d.a[1])
		return force(d.a[0], n)
	case Right:
		n := eval(d.a[0])
		return force(d.a[1], n)
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
func force(v *val, n int) int {
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
		r := eval(d.a[1])
		switch d.op {
		case "+":
			return force(l, n-r)
		case "-":
			return force(l, n+r)
		case "*":
			return force(l, n/r)
		case "/":
			return force(l, n*r)
		}
	case Right:
		l := eval(d.a[0])
		r := d.a[1]
		switch d.op {
		case "+":
			return force(r, n-l)
		case "-":
			return force(r, l-n)
		case "*":
			return force(r, n/l)
		case "/":
			return force(r, l/n)
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
