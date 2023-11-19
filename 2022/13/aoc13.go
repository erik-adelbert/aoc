// aoc13.go --
// advent of code 2022 day 13
//
// https://adventofcode.com/2022/day/13
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-13: initial commit
// 2023-11-19: adapt https://github.com/maneatingape/advent-of-code-rust beautiful analysis

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type packet struct {
	data  []byte
	index int
	extra []byte
}

func newPacket(s string) *packet {
	p := new(packet)
	p.data = []byte(s)
	p.extra = make([]byte, 16)
	return p
}

func (p *packet) next() byte {
	if i := len(p.extra); i > 0 {
		pop := p.extra[i-1]
		p.extra = p.extra[:i-1]
		return pop
	}
	data, i := p.data, p.index
	if data[i] == '1' && data[i+1] == '0' {
		p.index += 2
		return 'X' // single-digit 10
	}
	p.index++
	return data[i]
}

func (p *packet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{ s: %s, i: %d, x: '%s' }", p.data, p.index, string(p.extra))
	return sb.String()
}

func cmp(a, b string) int {
	left, right := newPacket(a), newPacket(b)
	for {
		a, b := left.next(), right.next()
		switch {
		case a == b:
			continue
		case a == ']':
			return -1
		case b == ']':
			return 1
		case a == '[':
			right.extra = append(right.extra, ']')
			right.extra = append(right.extra, b)
		case b == '[':
			left.extra = append(left.extra, ']')
			left.extra = append(left.extra, a)
		case a < b:
			return -1
		default:
			return 1
		}
	}
}

func main() {
	popcnt := 0
	packets := make([]string, 0)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		packet := input.Text()

		// part1
		if len(packet) == 0 {
			// last two fifo packets
			a := packets[len(packets)-2]
			b := packets[len(packets)-1]

			if cmp(a, b) < 1 {
				popcnt += len(packets) / 2
			}
			continue
		}

		// memoize for part2
		packets = append(packets, packet)
	}

	// part1
	fmt.Println(popcnt)

	// part2
	keys := []int{1, 2}
	for i := range packets {
		switch {
		case cmp(packets[i], "[[2]]") < 1:
			keys[0]++
			keys[1]++
		case cmp(packets[i], "[[6]]") < 1:
			keys[1]++
		}
	}
	fmt.Println(keys[0] * (keys[1]))
}
