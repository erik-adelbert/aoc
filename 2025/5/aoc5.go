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

func main() {
	var acc1, acc2 int // parts 1 and 2 counts

	// read spans and queries
	input := bufio.NewScanner(os.Stdin)

	spans := make([]span, 0, SpanCountHint)

	// state machine parser
	state := ReadSpans
	for input.Scan() {
		buf := input.Bytes()

		switch {
		case len(buf) == 0:
			// blank line separates spans and queries

			acc2, spans = merge(spans) // merge intervals and calculate total coverage

			state = ReadQueries
		case state == ReadSpans:
			start, end, _ := bytes.Cut(buf, []byte("-")) // parse range

			spans = append(spans, span{atoi(start), atoi(end)})
		case state == ReadQueries:
			v := atoi(buf) // parse query point

			if query(spans, v) {
				acc1++
			}
		}
	}

	fmt.Println(acc1, acc2)
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

	// check spans that could contain v (working backwards from found position)
	for j := i - 1; j >= 0; j-- {
		switch {
		case spans[j].end < v:
			return false // spans are sorted by start, so no earlier spans can contain v
		case spans[j].start <= v && v <= spans[j].end:
			return true
		}
	}

	return false
}

// span is an interval [start, end]
type span struct {
	start, end int
}

const (
	ReadSpans = iota
	ReadQueries
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
