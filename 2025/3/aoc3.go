// aoc3.go --
// advent of code 2025 day 3
//
// https://adventofcode.com/2025/day/3
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2025-12-3: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Part1 = 2  // keep 2 digits for part 1
	Part2 = 12 // keep 12 digits for part 2
)

func main() {
	var sum1, sum2 int // sums for parts 1 and 2

	// process input lines
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		buf := input.Bytes()

		seq1 := newSeq(Part1, len(buf))
		seq2 := newSeq(Part2, len(buf))

		for _, c := range buf {
			seq1.push(c)
			seq2.push(c)
		}

		sum1 += seq1.val()
		sum2 += seq2.val()
	}

	fmt.Println(sum1, sum2)
}

// seq is a sequence of digits with greedy removal
// to keep it lexicographically largest
type seq struct {
	digits []byte
	size   int
	krem   int // remaining removals
}

// newSeq creates a new sequence of given size and input size
// authorizing inputSize - size removals to build the largest subsequence
func newSeq(size, inputSize int) *seq {
	return &seq{
		digits: make([]byte, 0, size),
		size:   size,
		krem:   inputSize - size,
	}
}

// push a new digit, removing larger trailing digits if possible
// to keep the sequence lexicographically largest
func (s *seq) push(d byte) {

	// remove larger trailing digits while we can
	for s.krem > 0 && !s.empty() && d > s.peek() {
		// remove last
		s.krem--
		s.digits = s.digits[:len(s.digits)-1] // pop
	}

	// add new digit
	s.digits = append(s.digits, d)
}

// val returns the integer value of the sequence
func (s *seq) val() (n int) {
	for i := range s.digits[:s.size] {
		n = 10*n + ctoi(s.digits[i])
	}
	return
}

// peek returns the last digit of the sequence
func (s *seq) peek() byte {
	return s.digits[len(s.digits)-1]
}

func (s *seq) empty() bool {
	return len(s.digits) == 0
}

func (s *seq) String() string {
	var sb strings.Builder
	for i := range s.digits[:s.size] {
		sb.WriteByte(s.digits[i])
	}
	return sb.String()
}

func ctoi(c byte) int {
	return int(c - '0')
}
