package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	const (
		Hand = iota
		Bid
	)

	games1 := make([]game, 0, 1024)
	games2 := make([]game, 0, 1024)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		input := fields(input.Text())
		h, b := input[Hand], atoi(input[Bid])

		games1 = append(games1, game{hand: mkHand(h, Jack), bid: b})
		games2 = append(games2, game{hand: mkHand(h, Wild), bid: b})
	}

	slices.SortFunc(games1, cmp)
	slices.SortFunc(games2, cmp)

	sum1, sum2 := 0, 0
	for i := range games1 {
		sum1 += (i + 1) * games1[i].bid
		sum2 += (i + 1) * games2[i].bid
	}
	fmt.Println(sum1, sum2)
}

type game struct {
	hand
	bid int
}

func cmp(a, b game) int {
	return int(a.hand - b.hand)
}

func (g game) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "{%v %d}", g.hand, g.bid)
	return sb.String()
}

// see day write-up:
//
//	5   4   3   2   1   XRRR  fields X: special RRR: base rank
//	0123456789abcdef01234567  bit index (24bits)
type hand int

const (
	Jack = false
	Wild = !Jack
)

// rank sugars
const (
	High  = iota + 1
	One   // special One is Two
	Three // special Three is Full
	Four  // special Four is Five
	Five
)

func mkHand(s string, mode bool) (h hand) {
	var counts [14]int // card counts ex: "A23AA" -> [0, 1, 1, ..., 3]

	// encode and count hand cards
	h = h.set(R, High) // default rank
	for i := range s {
		n := ctoi(s[i], mode)
		h = h.set(card(i), n) // store reversed, see write-uo
		counts[n]++
	}

	// rank hand
	nread, J := 0, ctoi('J', mode)
	for i := range counts {
		nread += counts[i]

		// only One pair or more contribute to the rank
		// in wildcard mode, do not rank jacks
		if counts[i] < One || (mode == Wild && i == J) {
			continue
		}

		if h.get(R) > High || counts[i] == Five {
			// special hand:
			// One + One, One + Three, Three + One, Five
			counts[i] = min(counts[i], Four)
			h = h.set(X, On)
		}
		h = h.set(R, max(h.get(R), counts[i])) // update rank

		if nread == 5 || (mode == Wild && nread == 5-counts[J]) {
			break
		}
	}

	if mode == Wild {
		// max out base rank
		h = h.set(R, h.get(R)+counts[J])
		if counts[J] == Five {
			// JJJJJ
			h = h.set(R, Four).set(X, On)
		}
	}

	return
}

type field struct {
	n, mask int
}

const On = 1 // X field sugar

var (
	R = field{0x15, 0x7}
	X = field{0x14, 0x1}
)

// card# to bit field
func card(i int) field {
	return field{16 - (i << 2), 0xf} // reverse
}

func (h hand) get(f field) int {
	n, mask := f.n, f.mask
	return int(h>>n) & mask
}

func (h hand) set(f field, k int) hand {
	n, mask := f.n, f.mask
	h &= hand(^(mask << n))
	h |= hand(k << n)
	return h
}

func (h hand) String() string {
	var sb strings.Builder

	ranks := [...][]string{
		{}, {"High"}, {"One", "Two"}, {"Three", "Full"}, {"Four", "Five"},
	}

	itoc := func(i int) byte {
		return "J23456789TJQKA"[i]
	}

	// ex. "{AAJAA, Four}"
	sb.WriteString("{")
	for i := 0; i < 5; i++ {
		sb.WriteByte(itoc(h.get(card(i))))
	}
	sb.WriteString(", ")
	sb.WriteString(ranks[h.get(R)][h.get(X)])
	sb.WriteString("}")

	return sb.String()
}

func ctoi(c byte, mode bool) int {
	s := "?23456789TJQKA"
	if mode == Wild {
		s = "J23456789T?QKA"
	}

	for i := range s {
		if s[i] == c {
			return i
		}
	}
	return -1
}

var fields = strings.Fields

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
