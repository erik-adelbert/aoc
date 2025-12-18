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
	"time"
)

func main() {
	t0 := time.Now() // start timer

	var acc1, acc2 int // parts 1 and 2 counts

	// read spans and queries
	input := bufio.NewScanner(os.Stdin)

	spans := make([]span, 0, SpanCountHint)

	// state machine parser
	state := Span
	for input.Scan() {
		buf := input.Bytes()

		switch {
		case len(buf) == 0:
			// blank line separates spans and queries

			acc2, spans = merge(spans) // merge intervals and calculate total coverage

			state = Query

		case state == Span:
			start, end, _ := bytes.Cut(buf, []byte("-")) // parse range

			spans = append(spans, span{atoi(start), atoi(end)})

		case state == Query:
			qp := atoi(buf) // parse query point

			if query(spans, qp) {
				acc1++
			}
		}
	}

	fmt.Println(acc1, acc2, time.Since(t0))
}

// merge merges overlapping intervals and calculates total coverage
func merge(spans []span) (int, []span) {
	// sort intervals by start position
	slices.SortFunc(spans, func(a, b span) int {
		return a.start - b.start
	})

	merged := make([]span, 0, MergedSpanCountHint) // could also be len(spans)

	// merge overlapping intervals, count cover and populate tree
	cover, cur := 0, spans[0]

	for i := range spans[1:] {
		if spans[i].start <= cur.end+1 {
			// overlapping or adjacent intervals - merge
			if spans[i].end > cur.end {
				cur.end = spans[i].end
			}
		} else {
			// non-overlapping interval - add current coverage and start new interval
			cover += cur.end - cur.start + 1
			merged = append(merged, cur)

			cur = spans[i]
		}
	}

	// add the last interval
	cover += cur.end - cur.start + 1
	merged = append(merged, cur)

	return cover, merged
}

// query returns true if v is contained in any of the spans
func query(spans []span, v int) bool {
	// find the insertion point for a span that would start at v
	i, _ := slices.BinarySearchFunc(spans, span{v, v}, func(a, b span) int {
		return a.start - b.start
	})

	// check if v is contained in the span just before the insertion point
	return i > 0 && spans[i-1].start <= v && v <= spans[i-1].end
}

// span is an interval [start, end]
type span struct {
	start, end int
}

const (
	Span = iota
	Query
)

const (
	// hints for pre-allocations from prior runs
	SpanCountHint       = 187
	MergedSpanCountHint = 78
)

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s []byte) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
