// aoc5.go --
// advent of code 2025 day 5
//
// https://adventofcode.com/2025/day/5
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-5: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

const (
	MaxSpanTreeNodes = 187 // maximum span tree nodes from previous runs
)

func main() {
	var acc1, acc2 int // parts 1 and 2 counts

	// read spans and queries
	input := bufio.NewScanner(os.Stdin)

	tree := newSpanTree(MaxSpanTreeNodes)
	var spans []span

	// state machine parser
	state := ReadSpans
	for input.Scan() {
		buf := input.Bytes()

		switch {
		case len(buf) == 0:
			// blank line separates spans and queries

			// calculate total coverage from merged intervals
			// and populate the span tree
			acc2 = cover(tree, spans)

			state = ReadQueries
		case state == ReadSpans:
			// parse span
			start, end, _ := bytes.Cut(buf, []byte("-")) // parse range

			spans = append(spans, span{atoi(start), atoi(end)})
		case state == ReadQueries:
			// parse query point
			v := atoi(buf)

			// query spans containing v
			if len(tree.query(v)) > 0 {
				acc1++
			}
		}
	}

	fmt.Println(acc1, acc2)
}

// cover merges overlapping intervals and calculates total coverage
func cover(tree *spanTree, spans []span) int {
	if len(spans) == 0 {
		return 0
	}

	// sort intervals by start position
	slices.SortFunc(spans, func(a, b span) int {
		return a.start - b.start
	})

	// merge overlapping intervals, count cover and populate tree
	cover := 0
	cur := spans[0]

	for i := 1; i < len(spans); i++ {
		if spans[i].start <= cur.end+1 {
			// overlapping or adjacent intervals - merge
			if spans[i].end > cur.end {
				cur.end = spans[i].end
			}
		} else {
			// non-overlapping interval - add current coverage and start new interval
			cover += cur.end - cur.start + 1
			tree.insert(cur.start, cur.end)

			cur = spans[i]
		}
	}

	// add the last interval coverage
	cover += cur.end - cur.start + 1
	tree.insert(cur.start, cur.end)

	return cover
}

// span is an interval [start, end]
type span struct {
	start, end int
}

// spanTree is an interval tree for efficient querying of overlapping intervals
type spanTree struct {
	start  []int
	end    []int
	maxEnd []int
	left   []int
	right  []int
	root   int
}

// newSpanTree creates a new span tree with given capacity
func newSpanTree(cap int) *spanTree {
	return &spanTree{
		start:  make([]int, 0, cap),
		end:    make([]int, 0, cap),
		maxEnd: make([]int, 0, cap),
		left:   make([]int, 0, cap),
		right:  make([]int, 0, cap),
		root:   -1,
	}
}

// addNode adds a new node to the span tree and returns its index
func (t *spanTree) addNode(s, e int) int {
	idx := len(t.start)

	t.start = append(t.start, s)
	t.end = append(t.end, e)
	t.maxEnd = append(t.maxEnd, e)
	t.left = append(t.left, -1)
	t.right = append(t.right, -1)

	return idx
}

// insert adds a new span [s, e] to the span tree
func (t *spanTree) insert(s, e int) {
	i := t.addNode(s, e)

	// first node becomes root
	if t.root == -1 {
		t.root = i
		return
	}

	// first step: BST insert
	cur := t.root

	for {
		if s < t.start[cur] {
			// go left
			if t.left[cur] == -1 {
				t.left[cur] = i
				break
			}

			cur = t.left[cur]
		} else {
			// go right
			if t.right[cur] == -1 {
				t.right[cur] = i
				break
			}

			cur = t.right[cur]
		}
	}

	// 2nd step: update maxEnd
	cur = t.root
	for cur != -1 {
		if t.maxEnd[cur] < e {
			t.maxEnd[cur] = e
		}

		if s < t.start[cur] {
			cur = t.left[cur]
		} else {
			cur = t.right[cur]
		}
	}
}

// query returns all spans containing v
func (t *spanTree) query(v int) []span {
	var result []span

	stack := []int{t.root}

	for len(stack) > 0 {
		n := stack[len(stack)-1] // pop
		stack = stack[:len(stack)-1]

		if n == -1 {
			continue
		}

		// check overlap
		if t.start[n] <= v && v <= t.end[n] {
			result = append(result, span{t.start[n], t.end[n]})
		}

		// prune left subtree via maxEnd
		if t.left[n] != -1 && t.maxEnd[t.left[n]] >= v {
			stack = append(stack, t.left[n])
		}

		// right subtree if start <= v
		if t.right[n] != -1 && t.start[n] <= v {
			stack = append(stack, t.right[n])
		}
	}

	return result
}

const (
	ReadSpans = iota
	ReadQueries
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
