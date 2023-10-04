// aoc22.go --
// advent of code 2021 day 22
//
// https://adventofcode.com/2021/day/22
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-22: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type axis int

// Axis sugars
const (
	X axis = iota
	Y
	Z
)

// Axis sugar
var Axis = []axis{X, Y, Z}

type vec [3]int

var null vec

type cubo struct {
	min, max vec
}

var empty cubo

func (c cubo) ok() bool {
	λ := func(a axis) bool {
		return c.min[a] <= c.max[a]
	}
	return λ(X) && λ(Y) && λ(Z)
}

func (c cubo) aabb(d cubo) bool {
	λ := func(a axis) bool {
		return c.min[a] <= d.max[a] && c.max[a] >= d.min[a]
	}
	return λ(X) && λ(Y) && λ(Z)
}

func (c cubo) contains(d cubo) bool {
	λ := func(a axis) bool {
		return c.min[a] <= d.min[a] && c.max[a] >= d.max[a]
	}
	return λ(X) && λ(Y) && λ(Z)
}

func (c cubo) trim(r cubo) cubo { // r(egion)
	return cubo{
		max(c.min, r.min),
		min(c.max, r.max),
	}
}

func (c cubo) vol() int64 {
	λ := func(d axis) int64 {
		return int64(1 + c.max[d] - c.min[d])
	}
	return λ(X) * λ(Y) * λ(Z)
}

// https://en.wikipedia.org/wiki/K-d_tree
type node struct {
	c     cubo // leaf only
	a     axis
	v     int
	left  *node
	right *node
}

func add(n *node, c cubo) *node {
	switch {
	case n.c.contains(c):
		return n
	case c.contains(n.c):
		return &node{c: c}
	}

	for _, a := range Axis { // for X, Y, Z
		switch {
		case c.max[a] < n.c.min[a]:
			return &node{a: a, v: c.max[a], left: &node{c: c}, right: n} // left leaf
		case c.min[a] > n.c.max[a]:
			return &node{a: a, v: c.min[a] - 1, left: n, right: &node{c: c}} // right leaf
		default:
			if v := n.c.min[a]; c.min[a] < v && c.max[a] >= v {
				in, out := c, c
				in.min[a], out.max[a] = v, v-1

				return &node{a: a, v: v - 1, left: &node{c: out}, right: add(n, in)} // left leaf recurse right
			}
			if v := n.c.max[a]; c.max[a] > v && c.min[a] <= v {
				in, out := c, c
				in.max[a], out.min[a] = v, v+1

				return &node{a: a, v: v, left: add(n, in), right: &node{c: out}} // recurse left right leaf
			}
		}
	}
	panic("add: unreachable")
}

func del(n *node, c cubo) *node {
	switch {
	case c.contains(n.c):
		return nil
	case !c.aabb(n.c):
		return n
	}

	for _, a := range Axis { // for X, Y, Z
		if v := c.min[a]; n.c.min[a] < v && n.c.max[a] >= v {
			in, out := n.c, n.c
			in.min[a], out.max[a] = v, v-1
			left := &node{c: out}

			right := del(&node{c: in}, c)
			if right == nil {
				return left
			}
			return &node{a: a, v: v - 1, left: left, right: right}
		}
		if v := c.max[a]; n.c.max[a] > v && n.c.min[a] <= v {
			in, out := n.c, n.c
			in.max[a], out.min[a] = v, v+1
			right := &node{c: out}

			left := del(&node{c: in}, c)
			if left == nil {
				return right
			}
			return &node{a: a, v: v, left: left, right: right}
		}
	}
	panic("del: unreachable")
}

func leaf(n *node) bool {
	return n.left == nil && n.right == nil
}

func reboot(n *node, on bool, c cubo) *node {
	if n == nil {
		if on {
			n = &node{c: c}
		}
		return n
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	if !leaf(n) {
		left, right := c, c
		left.max[n.a] = min(left.max[n.a], n.v)
		right.min[n.a] = max(right.min[n.a], n.v+1)
		if left.ok() {
			n.left = reboot(n.left, on, left)
		}
		if right.ok() {
			n.right = reboot(n.right, on, right)
		}
		if n.left == nil {
			return n.right
		}
		if n.right == nil {
			return n.left
		}
		return n
	}

	if on {
		return add(n, c)
	}
	return del(n, c)
}

func recount(n *node) int64 {
	switch {
	case n == nil:
		return 0
	case leaf(n):
		return n.c.vol()
	default:
		return recount(n.left) + recount(n.right)
	}
}

func main() {
	var p1, p2 *node
	r := regexp.MustCompile(
		`(on|off) x=(-?\d+)..(-?\d+),y=(-?\d+)..(-?\d+),z=(-?\d+)..(-?\d+)`,
	)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		matches := r.FindAllStringSubmatch(input.Text(), -1)[0]
		args := make([]int, 6)
		for i := 0; i < 6; i++ {
			args[i], _ = strconv.Atoi(matches[i+2])
		}
		on := (matches[1] == "on")
		c := cubo{
			vec{args[0], args[2], args[4]},
			vec{args[1], args[3], args[5]},
		}

		part1 := cubo{vec{-50, -50, -50}, vec{50, 50, 50}}
		if sub := c.trim(part1); sub.ok() {
			p1 = reboot(p1, on, sub)
		}

		p2 = reboot(p2, on, c)
	}

	fmt.Println(recount(p1)) // part1
	fmt.Println(recount(p2)) // part2
}

func min(a, b vec) vec {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	return vec{
		min(a[X], b[X]),
		min(a[Y], b[Y]),
		min(a[Z], b[Z]),
	}
}

func max(a, b vec) vec {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	return vec{
		max(a[X], b[X]),
		max(a[Y], b[Y]),
		max(a[Z], b[Z]),
	}
}
