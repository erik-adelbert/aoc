package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	// strconv.Atoi simplified core loop
	// s is ^\d+$
	atoi := func(s []byte) int {
		n := 0
		for _, c := range s {
			n = 10*n + int(c-'0')
		}
		return n
	}

	const (
		Part1 = iota
		Part2
	)

	// 	               [B] [L]     [J]
	// 	           [B] [Q] [R]     [D] [T]
	// 	           [G] [H] [H] [M] [N] [F]
	//         [J] [N] [D] [F] [J] [H] [B]
	//     [Q] [F] [W] [S] [V] [N] [F] [N]
	// [W] [N] [H] [M] [L] [B] [R] [T] [Q]
	// [L] [T] [C] [R] [R] [J] [W] [Z] [L]
	// [S] [J] [S] [T] [T] [M] [D] [B] [H]
	//  1   2   3   4   5   6   7   8   9

	STACKS := [2][][]byte{
		{ // part1
			{},
			[]byte("SLW"),
			[]byte("JTNQ"),
			[]byte("SCHFJ"),
			[]byte("TRMWNGB"),
			[]byte("TRLSDHQB"),
			[]byte("MJBVFHRL"),
			[]byte("DWRNJM"),
			[]byte("BZTFHNDJ"),
			[]byte("HLQNBFT"),
		},
		{ // part2
			{},
			[]byte("SLW"),
			[]byte("JTNQ"),
			[]byte("SCHFJ"),
			[]byte("TRMWNGB"),
			[]byte("TRLSDHQB"),
			[]byte("MJBVFHRL"),
			[]byte("DWRNJM"),
			[]byte("BZTFHNDJ"),
			[]byte("HLQNBFT"),
		},
	}

	rev := func(x []byte) []byte {
		for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
			x[i], x[j] = x[j], x[i]
		}
		return x
	}

	// for top, pop, cut and push:
	//   p: Part1|Part2,
	//   s: source stack,
	//   d: target stack

	top := func(p, s int) byte {
		i := len(STACKS[p][s]) - 1 // stack top index
		return STACKS[p][s][i]
	}

	cut := func(p, s, n int) []byte {
		i := len(STACKS[p][s]) // stack top index
		x := STACKS[p][s][i-n:]
		STACKS[p][s] = STACKS[p][s][:i-n] // cut
		return x
	}

	pop := func(p, s, n int) []byte {
		return rev(cut(p, s, n))
	}

	push := func(p, d int, x []byte) {
		STACKS[p][d] = append(STACKS[p][d], x...)
	}

	move := func(nel, src, dst []byte) {
		n, s, d := atoi(nel), atoi(src), atoi(dst)

		push(Part1, d, pop(Part1, s, n))
		push(Part2, d, cut(Part2, s, n))
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// input text: ^move (\d+) from (\d+) to (\d+)$
		// args:          0    1     2    3    4   5
		args := bytes.Fields(input.Bytes())
		move(args[1], args[3], args[5])
	}

	display := func(p int) {
		var out bytes.Buffer
		for i, s := range STACKS[p] {
			if len(s) > 0 {
				out.WriteByte(top(p, i))
			}
		}
		fmt.Println(out.String())
	}

	display(Part1)
	display(Part2)
}
