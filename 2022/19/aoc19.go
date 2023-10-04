// aoc19.go --
// advent of code 2022 day 19
//
// https://adventofcode.com/2022/day/19
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-19: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

// robot/resource type
const (
	A = 1 // ore
	B = 2 // clay
	C = 3 // obsidian ore
	// implicit 4 obsidian clay
	D = 5 // geode ore
	// implicit 6 geode obsidian
)

type world struct {
	rules []int
	robot [8]int
	stock [8]int
	timer int
	ruleA int
}

func main() {
	worlds := make([]world, 0, 32)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Bytes()

		blue := make([]int, 0, 8)
		for i := 0; i < len(line); i++ {
			if line[i] < '0' || '9' < line[i] {
				continue
			}

			n, j := atoi(line[i:])
			blue = append(blue, n)
			i += j
		}
		worlds = append(worlds, mkworld(blue))
	}

	var bestD int

	// part1
	sum := 0
	for _, w := range worlds {
		bestD, w.timer = 0, 24
		w.maxout(&bestD)
		sum += w.id() * bestD
	}

	// part2
	prd := 1
	if len(worlds) > 3 {
		worlds = worlds[:3]
	}
	for _, w := range worlds {
		bestD, w.timer = 0, 32
		w.maxout(&bestD)
		prd *= bestD
	}

	fmt.Println(sum, prd)
}

func (w world) maxout(best *int) {
	*best = max(w.stock[D], *best)
	for _, x := range w.moves() {
		if x.hcost() > *best {
			x.maxout(best)
		}
	}
	return
}

// building cost heuristics
func (w world) hcost() int {
	// build as many C&Ds as fast as possible
	// consider A&B freely available
	c, r, d := w.stock[C], w.robot[C], w.stock[D]
	for left := w.timer - 1; left >= 0; left-- {
		if c >= w.rules[D+1] {
			c += r - w.rules[D+1]
			d += left
		} else {
			c += r
			r++
		}
	}
	return d
}

func (w world) build(m int) world {
	for left := w.timer - 1; left > 0; left-- {
		past := w.timer - left - 1

		// fast forward
		x := w
		for _, k := range []int{A, B, C} {
			x.stock[k] += x.robot[k] * past
		}

		// check build ability
		switch m {
		case A, B:
			if x.stock[A] < x.rules[m] {
				continue
			}
		case C:
			if x.stock[A] < x.rules[C] ||
				x.stock[B] < x.rules[C+1] {
				continue
			}
		case D:
			if x.stock[A] < x.rules[D] ||
				x.stock[C] < x.rules[D+1] {
				continue
			}
		}

		// build!
		x.timer = left
		x.stock[A] -= w.rules[m]
		switch m {
		case C:
			x.stock[B] -= x.rules[C+1]
		case D:
			x.stock[C] -= x.rules[D+1]
			x.stock[D] += left
		}
		for _, k := range []int{A, B, C} {
			x.stock[k] += x.robot[k]
		}
		x.robot[m]++
		return x
	}

	return world{} // none
}

func (w world) moves() []world {
	next := make([]world, 0, 4)

	want := []bool{
		A: w.robot[A] < w.ruleA,
		B: w.robot[B] < w.rules[C+1],
		C: w.robot[C] < w.rules[D+1] && w.robot[B] > 0,
		D: w.robot[C] > 0,
	}

	for _, i := range []int{A, B, C, D} {
		if want[i] {
			if x := w.build(i); !x.isnull() {
				next = append(next, x)
			}
		}
	}

	return next
}

func mkworld(p []int) world {
	var w world
	w.rules = p

	w.robot[A] = 1

	for _, k := range []int{B, C, D} {
		w.ruleA = max(w.ruleA, w.rules[k])
	}

	return w
}

func (w world) id() int {
	if len(w.rules) > 1 {
		return w.rules[0]
	}
	return 0
}

func (w world) isnull() bool {
	return w.id() == 0
}

// strconv.Atoi modified core loop
// s is ^\d+.*$
func atoi(s []byte) (int, int) {
	var n, i int
	for i = 0; i < len(s); i++ {
		if s[i] < '0' || '9' < s[i] {
			break
		}
		n = 10*n + int(s[i]-'0')
	}
	return n, i
}

// maximum of two ints
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const DEBUG = true

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
