// aoc19.go --
// advent of code 2024 day 19
//
// https://adventofcode.com/2024/day/19
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-19: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	RULES = iota
	WORDS
)

func main() {
	rules := make([]string, 0, 450) // arbitrary size
	words := make([]string, 0, 400) // arbitrary size

	state := RULES
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		switch {
		case line == "":
			state = WORDS
		case state == RULES:
			rules = strings.Split(line, ", ")
		case state == WORDS:
			words = append(words, line)
		}
	}

	trie := build(rules)

	count1, count2 := 0, 0
	for _, w := range words {
		if n := match(w, trie); n > 0 {
			count1 += 1
			count2 += n
		}
	}

	fmt.Println(count1, count2) // part 1 & 2
}

type TrieNode struct {
	next map[byte]*TrieNode
	stop bool
}

func newNode() *TrieNode {
	return &TrieNode{next: make(map[byte]*TrieNode)}
}

// Build a trie from a list of words
func build(words []string) *TrieNode {
	root := newNode()
	for _, word := range words {
		cur := root
		for _, x := range word {
			car := byte(x)
			if _, ok := cur.next[car]; !ok {
				cur.next[car] = newNode()
			}
			cur = cur.next[car]
		}
		cur.stop = true
	}
	return root
}

// Count all possible ways to fully match a string using words in the trie without overlaps
func match(line string, trie *TrieNode) int {
	end := len(line)
	memo := make(map[int]int, 58) // arbitrary size

	// DFS with memoization
	var recount func(int) int
	recount = func(start int) (count int) {
		if start == end {
			return 1 // success on the entire line!
		}

		if cnt, ok := memo[start]; ok {
			return cnt // use cached value
		}

		cur := trie
		for i := start; i < end; i++ {
			var ok bool
			var nxt *TrieNode

			car := line[i]
			if nxt, ok = cur.next[car]; !ok {
				break
			}

			cur = nxt
			if cur.stop {
				count += recount(i + 1) // add all ways from the next position
			}
		}

		memo[start] = count
		return
	}

	return recount(0)
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
