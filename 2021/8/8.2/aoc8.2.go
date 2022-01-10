package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type seg map[rune]bool // 8 segments display abstraction

// Seg constructs a seg object
func Seg(s string) seg {
	x := make(seg)
	for _, r := range s {
		x[r] = true
	}
	return x
}

func (s seg) inter(t seg) int { // common segments
	if len(s) > len(t) {
		s, t = t, s
	}

	n := 0
	for r := range s {
		if t[r] {
			n++
		}
	}
	return n
}

func (s seg) String() string {
	var sb strings.Builder
	for r := range s {
		sb.WriteRune(r)
	}
	return sb.String()
}

var digs = []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6} //  segment counts for 0..9

func match(segs []seg, sigs [][]seg) int {
	sig1 := sigs[digs[1]][0] // segment signal for 1
	sig4 := sigs[digs[4]][0] // segment signal for 4

	n := 0
	for _, s := range segs {
		n *= 10
		switch len(s) { // segments
		case 5:
			switch { // discriminate 2, 3, 5
			case s.inter(sig1) == 2: // only 1 & 3 have 2 common segments
				n += 3
			case s.inter(sig4) == 3: // only 4 & 5 have 3 common segments
				n += 5
			default: // it isn't 3 or 5: it's 2
				n += 2
			}
		case 6:
			switch { // discriminate 0, 6, 9
			case s.inter(sig1) == 1:
				n += 6
			case s.inter(sig4) == 4:
				n += 9
			}
		default: // for 2, 3, 4, 7 segment counts are 1, 7, 4, 8
			known := []int{2: 1, 3: 7, 4: 4, 7: 8}
			n += known[len(s)]
		}
	}
	// fmt.Println(segs, n)
	return n
}

func main() {
	var sigs [][][]seg
	var outs [][]seg

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args := strings.Split(input.Text(), "|")

		tokens := strings.Fields(strings.TrimSpace(args[0]))
		sig := make([][]seg, 8)
		for _, t := range tokens {
			s := Seg(t)
			sig[len(s)] = append(sig[len(s)], s)
		}
		sigs = append(sigs, sig)

		tokens = strings.Fields(strings.TrimSpace(args[1]))
		out := make([]seg, 0, 4)
		for _, t := range tokens {
			out = append(out, Seg(t))
		}
		outs = append(outs, out)
	}

	sum := 0
	for i, segs := range outs {
		sum += match(segs, sigs[i])
	}
	fmt.Println(sum)
}
