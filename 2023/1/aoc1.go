package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	sum1, sum2 := 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		input := input.Text()
		sum1 += stoi(input, DigitsOnly)    // part1
		sum2 += stoi(input, TextAndDigits) // part2
	}
	fmt.Println(sum1, sum2) // part 1 & 2
}

const (
	DigitsOnly    = false
	TextAndDigits = !DigitsOnly
)

func stoi(s string, mode bool) int {
	trie := digits
	if mode == TextAndDigits {
		trie = texts
	}
	l, r := 0, 0 // left, right
	for i := range s {
		// factorizing both scans into a closure is consistently slower

		// LR scan
		if nodes := trie[s[i]]; l == 0 && len(nodes) > 0 { // once for l
			for k := range nodes {
				if nodes[k].match(s[i+1:], LR) {
					l = nodes[k].atoi() // get leftmost number
					break
				}
			}
		}

		// RL scan
		j := (len(s) - 1) - i
		if nodes := trie[s[j]]; r == 0 && len(nodes) > 0 { // once for r
			for k := range nodes {
				if nodes[k].match(s[:j], RL) {
					r = nodes[k].atoi() // get rightmost number
					break
				}
			}
		}

		// eventually terminate with:
		//  l == r when there is a single number in s
		//  l != r when there are at least two
		if l*r > 0 {
			return 10*l + r
		}
	}

	panic("unreachable")
}

// see:
// https://web.stanford.edu/class/archive/cs/cs166/cs166.1146/lectures/09/Small09.pdf
type trie [][]node

const (
	LR = false
	RL = !LR
)

// node is a trie node
type node struct {
	trans string // transition
	state int    // terminal state ex. "three" or "3" -> 3
}

func (n node) match(s string, dir bool) bool {
	match := strings.HasPrefix
	if dir == RL {
		match = strings.HasSuffix
	}
	return match(s, n.trans)
}

func (n node) atoi() int {
	return n.state
}

// ε is inconditional transition (ie. any follow-up char is ok)
const ε = ""

// digits only trie
var digits = trie{
	'1': {{ε, 1}}, '2': {{ε, 2}}, '3': {{ε, 3}}, '4': {{ε, 4}}, '5': {{ε, 5}},
	'6': {{ε, 6}}, '7': {{ε, 7}}, '8': {{ε, 8}}, '9': {{ε, 9}}, 'z': {},
}

// digits and text texts trie
//
// Wether in digits or texts trie there is no proper suffix to link back to
// because there is no common suffix at all.
// Hence walking a trie branch is searching for a full (pre|suf)fix at a time
// and iterating (~backtracking) to the next (pre|suf)fix on failure.
// This is why transitions are merged in texts trie:
//
//	ex. o -> n -> e -> 1 => o -> ne -> 1
var texts = trie{
	'1': {{ε, 1}}, '2': {{ε, 2}}, '3': {{ε, 3}}, '4': {{ε, 4}}, '5': {{ε, 5}},
	'6': {{ε, 6}}, '7': {{ε, 7}}, '8': {{ε, 8}}, '9': {{ε, 9}}, 'z': {},

	'e': {{"ight", 8}, {"on", 1}, {"thre", 3}, {"fiv", 5}, {"nin", 9}},
	'f': {{"ive", 5}, {"our", 4}},
	'n': {{"ine", 9}, {"seve", 7}},
	'o': {{"ne", 1}, {"tw", 2}},
	'r': {{"fou", 4}},
	's': {{"even", 7}, {"ix", 6}},
	't': {{"hree", 3}, {"wo", 2}, {"eigh", 8}},
	'x': {{"si", 6}},
}
