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
		input := Fields(input.Text())
		h, b := input[Hand], atoi(input[Bid])

		games1 = append(games1, game{hand: newHand(h, Jack), bid: b})
		games2 = append(games2, game{hand: newHand(h, Joker), bid: b})
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
	*hand
	bid int
}

func cmp(a, b game) int {
	return int(*a.hand - *b.hand)
}

// hand
//
// see day write-up:
//
//	5   4   3   2   1   XKKK  fields
//	0123456789abcdef01234567  bit index (24bits)
type hand int

const (
	Jack  = false
	Joker = !Jack
)

func newHand(s string, mode bool) (h *hand) {
	var counts [14]int // hand card counts ex: "A23AA" -> [3, 1, 1 ..., 0]

	h = new(hand)
	h.set(K, 1)
	for i := range s {
		n := ctoi(s[i], mode)
		h.set(field{16 - (i << 2), 0x7}, n) // store reversed, see write-uo
		counts[n]++
	}

	nread, J := 0, ctoi('J', mode)
	for i := range counts {
		if mode == Joker && i == J {
			continue
		}

		nread += counts[i]

		// rank hand
		switch counts[i] {
		// default to Hand
		case 5: // Five
			h.set(X, 1)
			fallthrough
		case 4: // Four
			h.set(K, 4)
		case 3: // Three
			if h.get(K) == 2 { // Full
				h.set(X, 1)
			}
			h.set(K, 3)
		case 2: // One
			switch {
			case h.get(K) >= 2: // Two or Full
				h.set(X, 1)
			default:
				h.set(K, 2)
			}
		}

		if nread == 5 || (mode == Joker && nread == 5-counts[J]) {
			break
		}
	}

	if mode == Joker {
		k := h.get(K) + counts[J]
		switch {
		case counts[J] == 0:
			// no joker nothing to do
		case k == 6 || k == 5:
			// JJJJJ or XXXXJ
			h.set(K, 4)
			h.set(X, 1)
		default:
			h.set(K, k)
		}
	}

	return
}

type field struct {
	n, mask int
}

var (
	K = field{0x15, 0x7}
	X = field{0x14, 0x1}
)

func (h *hand) clr(f field) {
	n, mask := f.n, f.mask
	*h &= hand(^(mask << n))
}

func (h *hand) get(f field) int {
	n, mask := f.n, f.mask
	return int(*h>>n) & mask
}

func (h *hand) set(f field, k int) {
	n, mask := f.n, f.mask
	*h &= hand(^(mask << n))
	*h |= hand(k << n)
}

func ctoi(c byte, mode bool) int {
	s := "?23456789TJQKA"
	if mode == Joker {
		s = "J23456789T?QKA"
	}

	for i := range s {
		if s[i] == c {
			return i
		}
	}
	return -1
}

var Fields = strings.Fields

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
