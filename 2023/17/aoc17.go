// aoc17.go --
// advent of code 2022 day 17
//
// https://adventofcode.com/2022/day/17
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-17: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	world := newGrid()

	h, w := 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bytes := input.Bytes()

		w = len(bytes)
		copy(world.d[h*w:], bytes)

		h++
	}
	world.redim(h, w)

	fmt.Println(world.astar(pot{1, 3}), world.astar(pot{4, 10}))
}

type grid struct {
	φ func(y, x int) (i int)
	d []byte
	h int
	w int
}

const MAXN = 141

func newGrid() *grid {
	H, W := MAXN, MAXN
	g := grid{
		func(y, x int) (i int) { return y*MAXN + x },
		make([]byte, H*W),
		H, W,
	}
	return &g
}

func (g *grid) redim(h, w int) {
	g.h, g.w = h, w
	g.φ = func(y, x int) (i int) { return y*w + x }
}

func (g *grid) loss(i int) int {
	return int(g.d[i] - '0')
}

type pot struct {
	lo, hi int // crucible lowest, highest move length
}

const NBUCKET = 99

func (g *grid) astar(p pot) int {
	h, w, φ := g.h, g.w, g.φ
	N := NBUCKET

	const (
		V = iota
		H
	)

	const (
		L = iota
		R
		U
		D
	)

	type state struct {
		y, x, θ int
	}

	type states []state

	todo := make([][]state, N)
	for i := range todo {
		todo[i] = make([]state, 0, 221) // 221 tuned from previous run
	}
	push := func(ss []state, y, x, θ int) []state { return append(ss, state{y, x, θ}) }
	pop := func(ss []state) (int, int, int, []state) {
		var top state
		top, ss = ss[len(ss)-1], ss[:len(ss)-1]
		return top.y, top.x, top.θ, ss
	}

	type parm struct {
		valid func(i, y, x int) bool
		step  func(i, y, x int) (int, int)
		seek  func(x int) int
		θ     int
	}

	parms := []*parm{
		L: {
			func(i, _, x int) bool { return i <= x },
			func(i, y, x int) (int, int) { return y, x - i },
			func(i int) int { return i - 1 },
			H,
		},
		R: {
			func(i, _, x int) bool { return i < w-x },
			func(i, y, x int) (int, int) { return y, x + i },
			func(i int) int { return i + 1 },
			H,
		},
		U: {
			func(i, y, _ int) bool { return i <= y },
			func(i, y, x int) (int, int) { return y - i, x },
			func(i int) int { return i - w },
			V,
		},
		D: {
			func(i, y, _ int) bool { return i < h-y },
			func(i, y, x int) (int, int) { return y + i, x },
			func(i int) int { return i + w },
			V,
		},
	}

	losses := make([][2]int, h*w)

	move := func(m *parm, i, y, x, loss int) {
		χ := func(y, x, c int) int { return (h - y + w - x + c) % N }

		for ii := 1; ii <= p.hi && m.valid(ii, y, x); ii++ {
			i = m.seek(i)
			loss += g.loss(i)

			if ii >= p.lo && (losses[i][m.θ] == 0 || loss < losses[i][m.θ]) {
				y, x := m.step(ii, y, x)
				h := χ(y, x, loss)
				todo[h] = push(todo[h], y, x, m.θ)
				losses[i][m.θ] = loss
			}
		}
	}

	i := 0
	todo[0] = push(todo[0], 0, 0, V)
	todo[0] = push(todo[0], 0, 0, H)

	for {
		for len(todo[i%N]) > 0 {
			var y, x, θ int

			y, x, θ, todo[i%N] = pop(todo[i%N])
			i := φ(y, x)
			loss := losses[i][θ]

			if y == h-1 && x == w-1 {
				return loss
			}

			if θ == V {
				move(parms[L], i, y, x, loss)
				move(parms[R], i, y, x, loss)
			} else {
				move(parms[U], i, y, x, loss)
				move(parms[D], i, y, x, loss)
			}
		}
		i++
	}
}

func (g *grid) String() string {
	var sb strings.Builder

	for j := 0; j < g.h; j++ {
		fmt.Fprintln(&sb, string(g.d[j*g.w:(j+1)*g.w]))
	}

	return sb.String()
}
