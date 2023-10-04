// aoc13.go --
// advent of code 2022 day 13
//
// https://adventofcode.com/2022/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-13: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

type packet struct {
	list []packet
	val  int
}

func main() {
	popcnt := 0
	packets := []packet{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bytes := input.Bytes()

		// part1
		if len(bytes) == 0 {
			a := packets[len(packets)-2]
			b := packets[len(packets)-1]

			if cmp(a, b) < 1 {
				popcnt += len(packets) / 2
			}
			continue
		}

		// part2
		packets = append(packets, mkpacket(bytes))
	}

	// part1
	fmt.Println(popcnt)

	// part2
	keys := []int{1, 2}
	markers := []packet{
		// from u/Elavid on reddit
		mkint(2),
		mkint(6),
		// mkpacket([]byte("[[2]]")),
		// mkpacket([]byte("[[6]]")),
	}

	for i := range packets {
		if cmp(packets[i], markers[0]) <= 0 {
			keys[0]++
		}
		if cmp(packets[i], markers[1]) <= 0 {
			keys[1]++
		}
	}
	fmt.Println(keys[0] * (keys[1]))
}

func (p packet) isint() bool {
	return p.val != -1
}

// If both values are integers, the lower integer should come first.
// If the left integer is lower, the inputs are right.
// If the left integer is higher,the inputs are not right.
// Otherwise, the inputs are the same integer; continue
//
// If both values are lists, compare the first value of each list, then the second value, and so on.
// If the left list runs out of items first, the inputs are right.
// If the right list runs out of items first, the inputs are not right.
// If the lists are the same length and no comparison makes a decision about the order, continue.
//
// If exactly one value is an integer, convert the integer to a list, then retry.
func cmp(a, b packet) int {
	switch {
	case a.isint() && b.isint():
		switch {
		case a.val < b.val:
			return -1
		case a.val > b.val:
			return 1
		}
	case !(a.isint() || b.isint()):
		for i := range a.list {
			if i >= len(b.list) {
				return 1
			}
			if r := cmp(a.list[i], b.list[i]); r != 0 {
				return r
			}
		}
		if len(b.list) > len(a.list) {
			return -1
		}
	case a.isint():
		return cmp(mklist([]packet{a}), b)
	case b.isint():
		return cmp(a, mklist([]packet{b}))
	}
	return 0
}

func mkint(v int) packet {
	return packet{val: v}
}

func mklist(l []packet) packet {
	return packet{l, -1}
}

func mkpacket(s []byte) packet {
	var rec func(int) (packet, int)

	rec = func(i int) (packet, int) {
		a := packet{val: -1}

		for ; i < len(s); i++ {
			switch s[i] {
			case '[':
				var b packet
				b, i = rec(i + 1)
				a.list = append(a.list, b)
				fallthrough
			case ',':
				continue
			case ']':
				return a, i
			}

			a.list = append(
				a.list, mkint(atoi(s[i:])))
		}
		return a, i
	}

	a, _ := rec(0)
	return a
}

// strconv.Atoi modified core loop
// s is ^\d+.*
// capture breaks at first non digit
func atoi(s []byte) int {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = 10*n + int(c-'0')
	}
	return n
}
