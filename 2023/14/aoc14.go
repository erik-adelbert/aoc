package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {

	g := newGrid[byte](0)
	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		g.load(j, input.Bytes())
	}

	b := newBoard(g.widen())

	fmt.Println(
		b.tiltNorth(),
		b.tiltCycle(1_000_000_000),
	)
}

const MAXN = 100

type ints interface {
	byte | int16
}

type grid[T ints] struct {
	φ func(y, x int) int
	d []T
	w int
}

func newGrid[T ints](w int) *grid[T] {
	if w == 0 {
		w = MAXN + 2
	}

	g := grid[T]{
		d: make([]T, w*w),
		w: w,
	}
	g.φ = func(y, x int) int {
		return y*g.w + x
	}

	return &g
}

func (g *grid[T]) load(j int, s []byte) {
	g.w = len(s)
	row := g.d[j*g.w:]
	for i := range s {
		row[i] = T(s[i])
	}
}

func (g *grid[T]) widen() *grid[T] {
	w := g.w
	buf := make([]T, (w+2)*(w+2))
	for i := range buf {
		buf[i] = '#'
	}

	for j := 0; j < g.w; j++ {
		copy(buf[(j+1)*(w+2)+1:], g.d[j*w:(j+1)*w])
	}
	g.w = w + 2
	g.d = buf
	return g
}

// func (g *grid[T]) clone() *grid[T] {
// 	c := *g
// 	c.d = make([]T, c.w*c.w)
// 	copy(c.d, g.d)
// 	return &c
// }

func (g *grid[T]) String() string {
	var sb strings.Builder
	for j := 0; j < g.w; j++ {
		fmt.Fprintln(&sb, g.d[j*g.w:(j+1)*g.w])
	}
	return sb.String()
}

const (
	North = iota
	West
	South
	East
)

type board struct {
	fixes [4]*grid[int16]
	rolls [4][]int16
	rocks []int16
	w     int
}

func newBoard(g *grid[byte]) (b *board) {
	b = &board{w: g.w}

	b.rocks = make([]int16, 0, 2037)
	for i := range g.d {
		if g.d[i] == 'O' {
			b.rocks = append(b.rocks, int16(i))
		}
	}

	b.fixes = [4]*grid[int16]{
		North: newGrid[int16](g.w),
		West:  newGrid[int16](g.w),
		South: newGrid[int16](g.w),
		East:  newGrid[int16](g.w),
	}

	b.rolls = [4][]int16{
		North: make([]int16, 0, 2061),
		West:  make([]int16, 0, 2061),
		South: make([]int16, 0, 2061),
		East:  make([]int16, 0, 2061),
	}

	for y := 0; y < g.w; y++ {
		for x := 0; x < g.w; x++ {
			for θ, i := range []int{
				North: g.φ(x, y),
				West:  g.φ(y, x),
				South: g.φ(g.w-1-x, y),
				East:  g.φ(y, g.w-1-x),
			} {
				if g.d[i] == '#' {
					b.rolls[θ] = append(b.rolls[θ], int16(i))
				}
				b.fixes[θ].d[i] = int16(len(b.rolls[θ]) - 1)
			}
		}
	}

	return
}

func (b *board) tilt(θ int) []int16 {
	var clone = slices.Clone[[]int16]

	w := b.w
	state := clone(b.rolls[θ])

	offs := []int{
		North: +w,
		West:  +1,
		South: -w,
		East:  -1,
	}

	for i, r := range b.rocks {
		ii := b.fixes[θ].d[r]
		state[ii] += int16(offs[θ])
		b.rocks[i] = state[ii]
	}

	return state
}

func (b *board) tiltNorth() (load int) {
	var clone = slices.Clone[[]int16]

	rocks := clone(b.rocks)
	state := b.tilt(North)
	b.rocks = rocks

	for i, x := range b.rolls[North] {
		y := state[i]

		for i := x; i < y; i += int16(b.w) {
			y := int(i) / b.w
			load += b.w - 2 - y
		}
	}

	return
}

type hashkey []int16

func (h hashkey) hash() string {
	var n int

	if n = len(h); n > 99 {
		// /!\ tune cropping if needed
		n = n * 54 / 100
	}

	var sb strings.Builder
	for i := range h[:n] {
		fmt.Fprintf(&sb, "%x", h[i])
	}
	return sb.String()
}

func (b *board) tiltCycle(n int) (load int) {

	type state struct {
		s []int16
		i int
	}

	seen := make(map[string]state, len(b.rocks))

	s, e := 0, 0
CYCLE:
	for {
		for _, θ := range []int{North, West, South} {
			b.tilt(θ)
		}
		cycle := b.tilt(East)

		h := hashkey(cycle).hash()
		if x, ok := seen[h]; ok {
			s, e = x.i, len(seen)
			break CYCLE
		}
		seen[h] = state{cycle, len(seen)}
	}

	size, off := e-s, n-1-s
	i := s + off%size

	var last []int16
	for _, v := range seen {
		if v.i == i {
			last = v.s
			break
		}
	}

	w := b.w
	for i, α := range b.rolls[East] {
		β := last[i]

		n := int(α - β)
		y := int(α) / w

		load += n * (w - 1 - y)
	}

	return
}
