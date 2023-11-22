// aoc24.go --
// advent of code 2022 day 24
//
// https://adventofcode.com/2022/day/24
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-24: initial commit
// 2023-11-22: use u128 as maneatingape does

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	w := make(waze, 0, 32)

	input := bufio.NewScanner(os.Stdin)

	var row string
	input.Scan() // discard first row
	for input.Scan() {
		if len(row) > 0 {
			// discard first and last col
			w = w.append(row[1 : len(row)-1])
		}
		row = input.Text()
	}
	// discard last row

	// build maze over time
	m := newTaze(w)

	const (
		FWD = true
		BCK = !FWD
	)

	lap1 := m.solve(0, FWD)
	lap2 := m.solve(lap1, BCK)
	lap3 := m.solve(lap2, FWD)

	fmt.Println(lap1, lap3)
}

// timed waze
type taze struct {
	H, W       int
	rows, cols []uint128
}

func newTaze(m waze) *taze {
	H, W := len(m), len(m[0])

	encode := func(b byte) []uint128 {
		pack := func(s string) uint128 {
			var acc uint128

			for i := range s {
				acc = acc.lsh(1)
				if s[i] != b {
					acc = acc.or(one)
				}
			}
			return acc
		}

		rows := make([]uint128, H)
		for i := range m {
			rows[i] = pack(m[i])
		}
		return rows
	}

	// Left, Right, Up, Down
	L, R, U, D := encode('<'), encode('>'), encode('^'), encode('v')

	rows := make([]uint128, 0, H*W)
	for t := 0; t < W; t++ {
		for j := 0; j < H; j++ {
			l := (L[j].lsh(t)).or(L[j].rsh(W - t))
			r := (R[j].rsh(t)).or(R[j].lsh(W - t))
			rows = append(rows, l.and(r))
		}
	}

	cols := make([]uint128, 0, W*H)
	for t := 0; t < H; t++ {
		for i := 0; i < W; i++ {
			u := U[mod(i+t, H)]
			d := D[mod(i-t, H)]
			cols = append(cols, u.and(d))
		}
	}

	return &taze{H, W, rows, cols}
}

func (m *taze) solve(t0 int, fwd bool) int {
	H, W, rows, cols := m.H, m.W, m.rows, m.cols
	state := make([]uint128, H+1)
	for t := t0; ; {
		var old, cur, nxt uint128

		t++
		nxt = state[0]
		for i := 0; i < H; i++ {
			old, cur, nxt = cur, nxt, state[i+1]
			state[i] = cur.or(cur.rsh(1)).or(cur.lsh(1).or(old).or(nxt))
			state[i] = state[i].and(rows[H*mod(t, W)+i])
			state[i] = state[i].and(cols[W*mod(t, H)+i])
		}

		if fwd {
			state[0] = state[0].or(one.lsh(W - 1))
			if state[H-1].and(one) != zero {
				return t + 1
			}
		} else {
			state[H-1] = state[H-1].or(one)
			if state[0].and(one.lsh(W-1)) != zero {
				return t + 1
			}
		}
	}
}

// windy maze
type waze []string

func (w waze) append(elems ...string) waze {
	return waze(append([]string(w), elems...))
}

const uint128size = 128

type uint128 struct {
	hi, lo uint64
}

var (
	zero = uint128{0, 0}
	one  = uint128{0, 1}
)

func (u uint128) lsh(n int) uint128 {
	if n >= 64 {
		return uint128{u.lo << (n - 64), 0}
	}
	return uint128{u.hi<<n | u.lo>>(64-n), u.lo << n}
}

func (u uint128) rsh(n int) uint128 {
	if n >= 64 {
		return uint128{0, u.hi >> (n - 64)}
	}
	return uint128{u.hi >> n, u.lo>>n | u.hi<<(64-n)}
}

func (u uint128) and(m uint128) uint128 {
	return uint128{u.hi & m.hi, u.lo & m.lo}
}

func (u uint128) or(m uint128) uint128 {
	return uint128{u.hi | m.hi, u.lo | m.lo}
}

func (u uint128) String() string {
	var sb strings.Builder
	if u.hi != 0 {
		fmt.Fprintf(&sb, "%x%016x", u.hi, u.lo)
	} else {
		fmt.Fprintf(&sb, "%x", u.lo)
	}
	return sb.String()
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

var DEBUG = false

func debug(a ...any) (int, error) {
	if DEBUG {
		return fmt.Println(a...)
	}
	return 0, nil
}
