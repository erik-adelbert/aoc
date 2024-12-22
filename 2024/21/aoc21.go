// aoc21.go --
// advent of code 2024 day 21
//
// https://adventofcode.com/2024/day/21
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-21: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"slices"
)

const (
	MAXMSG = 12
	NLOOP1 = 2
	NLOOP2 = 25 - NLOOP1
)

func main() {

	count1, count2 := 0, 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		code := input.Text()
		m := Message{code: 1}.keytype()

		for i := 0; i < NLOOP1; i++ {
			m = m.cmdtype()
		}
		count1 += atoi(code[:3]) * m.len()

		for i := 0; i < NLOOP2; i++ {
			m = m.cmdtype()
		}
		count2 += atoi(code[:3]) * m.len()
	}

	fmt.Println(count1, count2) // part 1 & 2

}

type Cell struct {
	r, c int
}

// KEYPAD
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//     | 0 | A |
//     +---+---+

// CMDPAD
// 	   +---+---+
//     | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+

type Pad struct {
	k [][]byte
	c []Cell
}

var KEYPAD, CMDPAD Pad

func init() {
	KEYPAD.k = [][]byte{
		[]byte("789"),
		[]byte("456"),
		[]byte("123"),
		[]byte(".0A"),
	}

	CMDPAD.k = [][]byte{
		[]byte(".^A"),
		[]byte("<v>"),
	}

	setup := func(pad Pad) Pad {
		var keys []byte

		for _, row := range pad.k {
			keys = append(keys, row...)
		}
		n := slices.Max(keys) + 1

		pad.c = make([]Cell, n)
		for r, row := range pad.k {
			for c, key := range row {
				pad.c[key] = Cell{r, c}
			}
		}
		return pad
	}

	KEYPAD = setup(KEYPAD)
	CMDPAD = setup(CMDPAD)
}

func (p Pad) key(c Cell) byte {
	H, W := len(p.k), len(p.k[0])
	if c.r < 0 || c.r >= H || c.c < 0 || c.c >= W {
		return '.'
	}
	return p.k[c.r][c.c]
}

func (p Pad) rc(k byte) Cell {
	return p.c[k]
}

func (a Cell) sub(b Cell) Cell {
	return Cell{a.r - b.r, a.c - b.c}
}

func (p Pad) move(a, b byte) []byte {
	src, dst := p.rc(a), p.rc(b)
	δ := dst.sub(src)

	v := append(
		bytes.Repeat([]byte{'^'}, max(-δ.r, 0)),
		bytes.Repeat([]byte{'v'}, max(δ.r, 0))...,
	)

	h := append(
		bytes.Repeat([]byte{'<'}, max(-δ.c, 0)),
		bytes.Repeat([]byte{'>'}, max(δ.c, 0))...,
	)

	var buf bytes.Buffer

	write := func(shards ...[]byte) {
		for _, shard := range shards {
			buf.Write(shard)
		}
	}

	A := []byte{'A'}

	switch {
	case δ.c > 0 && p.key(Cell{dst.r, src.c}) != '.':
		write(v, h, A)
	case p.key(Cell{src.r, dst.c}) != '.':
		write(h, v, A)
	default:
		write(v, h, A)
	}

	return buf.Bytes()
}

func (m Message) cmdtype() Message {
	msg := make(Message, MAXMSG)

	for path, cnt := range m {
		cur := byte('A')
		for _, nxt := range []byte(path) {
			mv := string(CMDPAD.move(cur, nxt))
			msg[mv] += cnt
			cur = nxt
		}
	}

	return msg
}

func (m Message) keytype() Message {
	var buf bytes.Buffer

	for path := range m {
		cur := byte('A')
		for _, nxt := range []byte(path) {
			buf.Write(KEYPAD.move(cur, nxt))
			cur = nxt
		}
	}

	return Message{
		buf.String(): 1,
	}
}

// UP, RIGHT, DOWN, LEFT
var DIRS = []Cell{
	{-1, 0}, {0, 1}, {1, 0}, {0, -1},
}

type Message map[string]int

const (
	KEY = iota
	CMD
)

func (m Message) len() (n int) {
	for k, v := range m {
		n += len(k) * v
	}
	return
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
