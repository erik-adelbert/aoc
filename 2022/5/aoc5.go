package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	var ws worlds // part1&2

	input := bufio.NewScanner(os.Stdin)
	ws.load(input)

	// read moves
	for input.Scan() {
		// input text: ^move (\d+) from (\d+) to (\d+)$
		// args:          0    1     2    3    4   5
		args := bytes.Fields(input.Bytes())
		ws.move(args[5], args[3], args[1])
	}

	fmt.Println(ws)
}

// a world is a slice of adressable byte stacks
type world [][]byte

// for all byte stack methods:
// s is source stack,
// d is target stack,
// n is number of stack elements

func (w world) cut(s, n int) (x []byte) {
	i := len(w[s])                   // stack top index
	x, w[s] = w[s][i-n:], w[s][:i-n] // cut
	return
}

func (w world) pop(s, n int) (x []byte) {
	x = w.cut(s, n)
	// reverse x
	for l, r := 0, len(x)-1; l < r; {
		x[l], x[r] = x[r], x[l]
		l++
		r--
	}
	return
}

func (w world) push(d int, x []byte) {
	w[d] = append(w[d], x...)
}

func (w world) top(s int) byte {
	i := len(w[s]) - 1 // byte stack top index
	return w[s][i]
}

// worlds support muxed part 1&2 ops
type worlds [2]world

const (
	Part1 = iota
	Part2
)

func (ws worlds) String() string {
	var sb strings.Builder

	for i, s := range ws[Part1] {
		if len(s) > 0 {
			sb.WriteByte(ws[Part1].top(i))
		}
	}
	sb.WriteByte('\n')
	for i, s := range ws[Part2] {
		if len(s) > 0 {
			sb.WriteByte(ws[Part2].top(i))
		}
	}
	return sb.String()
}

// read initial state from input
// part 1&2 worlds start in the same state
func (ws *worlds) load(input *bufio.Scanner) {
	//                 [B] [L]     [J]
	//             [B] [Q] [R]     [D] [T]
	//             [G] [H] [H] [M] [N] [F]
	//         [J] [N] [D] [F] [J] [H] [B]
	//     [Q] [F] [W] [S] [V] [N] [F] [N]
	// [W] [N] [H] [M] [L] [B] [R] [T] [Q]
	// [L] [T] [C] [R] [R] [J] [W] [Z] [L]
	// [S] [J] [S] [T] [T] [M] [D] [B] [H]
	//  1   2   3   4   5   6   7   8   9

	lines := make([][]byte, 0, 8)
	for input.Scan() {
		raw := input.Bytes()
		if len(raw) == 0 {
			// done reading world state
			break
		}

		line := make([]byte, 0, 16)
		// copy letters only
		// see ex. input above
		for i := 1; i < len(raw); i += 4 {
			line = append(line, raw[i])
		}
		lines = append(lines, line)
	}

	// demux
	for len(ws[Part1]) <= len(lines[0]) {
		ws[Part1] = append(ws[Part1], make([]byte, 0, 16))
		ws[Part2] = append(ws[Part2], make([]byte, 0, 16))
	}

	// rotate lines to bstacks
	for i := len(lines) - 2; i >= 0; i-- { // discard last line
		for j, c := range lines[i] {
			if c != ' ' {
				// demux
				ws[Part1][j+1] = append(ws[Part1][j+1], c) // 1-indexed
				ws[Part2][j+1] = append(ws[Part2][j+1], c) // 1-indexed
			}
		}
	}
}

// muxed move for part1&2 worlds
func (ws worlds) move(dst, src, size []byte) {
	d, s, n := atoi(dst), atoi(src), atoi(size)

	// demux
	ws[Part1].push(d, ws[Part1].pop(s, n))
	ws[Part2].push(d, ws[Part2].cut(s, n))
}

// strconv.Atoi simplified core loop
// s is ^(\d+)$
func atoi(s []byte) int {
	var n int
	for _, c := range s {
		n = 10*n + int(c-'0')
	}
	return n
}
