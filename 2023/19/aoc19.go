package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	workflows := make(workflows, 1024)
	parts := make([]cub4, 0, 256)

	var parse func(string)

	// parsing state machine
	parseCube := func(s string) {
		// input is "{x=(\d+),m=(\d+),a=(\d+),s=(\d+)}"
		part := cub4{}
		// tokenize
		for i, x := range split(s[1:len(s)-1], ",") {
			// parse token
			n := atoi(x[2:])
			part[i] = span{n, n + 1}
		}
		parts = append(parts, part)
	}

	parseWorkflow := func(s string) {
		if len(s) == 0 { // transition on empty string
			parse = parseCube
			return
		}

		// ex.
		// ktc{a<1998:R,m>3286:A,x<1292:R,A}
		// cb{a>1858:R,s>3028:R,zt}
		i := index(s, "{") // seek start

		// tokenize
		id, rules := s[:i], split(s[i+1:len(s)-1], ",")
		for i := range rules {
			// parse token
			workflows[id] = append(workflows[id], parseRule(rules[i]))
		}
	}

	parse = parseWorkflow // start state
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		parse(input.Text()) // self transition
	}

	sum1 := 0
	for _, p := range workflows.process(parts) {
		for _, x := range p {
			sum1 += x.lo
		}
	}

	sum2 := 0
	for _, p := range workflows.process(all()) {
		prod := 1
		for _, x := range p {
			prod *= x.hi - x.lo + 1
		}
		sum2 += prod
	}

	fmt.Println(sum1, sum2)
}

func parseRule(s string) rule {
	var r rule

	// ex. "m>3286:A", "R", "A"
	switch index(s, ":") {
	case -1:
		r = rule{jmp: s}
	default:
		args := split(s, ":")
		r = rule{
			cut{
				a: axis(index("xmas", args[0][:1])),
				x: atoi(args[0][2:]),
			},
			args[0][1],
			args[1],
		}
	}
	return r
}

type axis int

// Axis sugars
const (
	X axis = iota
	M
	A
	S
)

const (
	LO = iota
	HI
)

type span struct{ lo, hi int }

type cub4 [4]span // hyperrectangle-4

var null cub4 // sugar

const (
	MIN = 1
	MAX = 4000
)

func all() []cub4 {
	var x cub4
	for i := range x {
		x[i] = span{MIN, MAX}
	}
	return []cub4{x}
}

// divide tesseract (hypercube) by splitting hyperplan (a, x)
func (hc cub4) split(hp cut, off int) (lo, hi cub4) {
	spaces := make([]cub4, 2)

	// half space ranges
	a, x := hp.a, hp.x+off // (LO, HI) offsets when LO: (-1, 0) | HI: (0, 1)
	halves := []span{
		{ // low
			hc[a].lo,
			min(hc[a].hi, x-1),
		},
		{ // high
			max(hc[a].lo, x),
			hc[a].hi,
		},
	}

	// validate halves
	for i, h := range halves {
		if h.lo < h.hi {
			// valid range
			cut := hc       // clone
			cut[a] = h      // update
			spaces[i] = cut // store
		} else {
			spaces[i] = null
		}
	}

	return spaces[0], spaces[1]
}

// addressable stack
type astack map[string][]cub4

func (r astack) pop() (k string, c []cub4) {
	for k, c = range r {
		delete(r, k)
		return
	}
	return "", []cub4{}
}

// splitting hyperplan
type cut struct {
	a axis
	x int // coordinate
}

type rule struct {
	cut
	tst byte
	jmp string
}

type workflow []rule

func (w workflow) filter(parts []cub4) astack {
	done := make(astack, len(parts))

	// match rules
	cur, nxt := parts, make([]cub4, len(parts))
	for _, r := range w { // each rule
		nxt = nxt[:0] // reset
		for _, p := range cur {
			switch r.tst {
			case 0:
				done[r.jmp] = append(done[r.jmp], p)
			case '<':
				lo, hi := p.split(r.cut, LO)

				if lo != null {
					done[r.jmp] = append(done[r.jmp], lo)
				}
				if hi != null {
					nxt = append(nxt, hi) // no match yet
				}
			case '>':
				lo, hi := p.split(r.cut, HI)

				if lo != null {
					nxt = append(nxt, lo) // no match yet
				}
				if hi != null {
					done[r.jmp] = append(done[r.jmp], hi)
				}
			}
		}
		cur, nxt = nxt, cur
	}

	return done
}

type workflows map[string]workflow

func (ws workflows) process(parts []cub4) []cub4 {
	A := make([]cub4, 0, len(parts))

	stack := make(astack)
	stack["in"] = parts
	for len(stack) > 0 {
		k, v := stack.pop()

		for k, v := range ws[k].filter(v) {
			switch k {
			case "R": // reject
				// discard
			case "A": // accept
				A = append(A, v...) // push to A
			default:
				stack[k] = append(stack[k], v...) // push back
			}
		}
	}

	return A
}

var index, split = strings.Index, strings.Split

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
