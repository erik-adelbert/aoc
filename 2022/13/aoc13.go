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

func main() {
	// parity, pair index, ordered pair count
	par, idx, popcnt := 0, 0, 0

	// part2 subkeys indices
	subkeys := []int{1, 2}

	// last 2 packets buffer
	packets := [2]string{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		packet := input.Text()

		// part1
		if len(packet) == 0 {
			idx++

			// last two packets
			a := packets[0]
			b := packets[1]

			if cmp(a, b) < 1 {
				// packets are ordered
				popcnt += idx
			}
			continue // process next pair
		}

		// part2
		switch {
		case cmp(packet, "2") < 1:
			// packet goes before subkeys
			subkeys[0]++
			subkeys[1]++
		case cmp(packet, "6") < 1:
			// packet goes between subkeys
			subkeys[1]++
		}

		// memoize for part1 according to parity
		packets[par], par = packet, 1-par
	}

	// part 1&2
	fmt.Println(popcnt, subkeys[0]*subkeys[1])
}

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
