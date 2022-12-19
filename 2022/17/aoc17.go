package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var tetros = [][][]byte{
	{
		[]byte{1, 1, 1, 1},
	},
	{
		[]byte{0, 1, 0},
		[]byte{1, 1, 1},
		[]byte{0, 1, 0},
	},
	{
		[]byte{0, 0, 1},
		[]byte{0, 0, 1},
		[]byte{1, 1, 1},
	},
	{
		[]byte{1},
		[]byte{1},
		[]byte{1},
		[]byte{1},
	},
	{
		[]byte{1, 1},
		[]byte{1, 1},
	},
}

// world
const (
	H    = 3192
	W    = 7
	XORG = 2
	YORG = 3
)

type board [H][W]byte

var jets []byte

func main() {
	b := board{}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		jets = input.Bytes()
	}

	// part 1&2
	b.play([]int{2022, 1_000_000_000_000})
}

type state struct {
	// tetro, jet, count
	t, j, n int

	// cur & last highests
	h [2]int
}

// play n tetrominoes with cycle detection
func (b *board) play(n []int) {
	// initial state
	s := state{0, 0, 0, [2]int{H, 0}}

	// cycle detection key type and map
	type key struct {
		t, j int
		// if it doesn't work for an input,
		// try to replace W/2 by W
		s [W / 2]int
	}
	seen := make(map[key]state, 2048)

	for {
		// simulate play
		s.drop(b)

		k := key{
			s.t, s.j, b.skyline(s.h[0]),
		}
		if _, ok := seen[k]; !ok {
			seen[k] = s
			continue
		}

		// cycle detected
		s0 := seen[k]    // starting state
		p0 := s.n - s0.n // period

		rh := make([][2]int, len(n))
		for i := range rh {
			q := (n[i] - s0.n) / p0
			r := (n[i] - s0.n) % p0

			// retrograd height!!
			h := (H - s0.h[1]) + (s0.h[1]-s.h[1])*q

			// sort/insert by remainders
			j := sort.Search(len(rh), func(i int) bool { return rh[i][0] >= r }) - 1
			rh[j][0] = r
			rh[j][1] = h
		}

		// fast forward play
		j, rmax := 0, rh[len(rh)-1][0]+1
		y := s.h[1]

		for i := 0; i < rmax; i++ {
			// simulate for remainders
			s.drop(b)
			if i == rh[j][0] {
				// tetro count is just over n:
				// account for cycle and added height and output
				fmt.Println(rh[j][1] + (y - s.h[1]))
				j++
			}
		}
		return
	}
}

// drop next tetromino and update board from initial state
func (s *state) drop(b *board) {
	// current tetro from state
	t, th, tw := tetros[s.t], len(tetros[s.t]), len(tetros[s.t][0])

	// collision helper
	collide := func(y, x int) bool {
		for j := range t {
			for i := range t[j] {
				if b[y+j][x+i]+t[j][i] > 1 {
					// collision!
					return true
				}
			}
		}
		return false
	}

	// drop and track tetro
	y, x := s.h[0]-th-YORG, XORG
	for ; ; s.j++ {
		m := jets[s.j%len(jets)]
		switch {
		case m == '<' && x > 0 && !collide(y, x-1):
			x--
		case m == '>' && x+tw < W && !collide(y, x+1):
			x++
		}

		if y+th < H && !collide(y+1, x) {
			y++
			continue
		}

		break
	}
	// update board / copy block
	for j := range t {
		for i := range t[j] {
			b[y+j][x+i] += t[j][i]
		}
	}

	// update state for:
	//   current and last highest
	//   next jet
	//   next tetro
	//   tetros current count
	s.h[0], s.h[1] = min(s.h[0], y), s.h[0]
	s.j = (s.j + 1) % len(jets)
	s.t = (s.t + 1) % len(tetros)
	s.n++
}

// compute half board top height per column
func (b board) skyline(h int) [W / 2]int {
	heights := [W / 2]int{}

	for i := range heights[:W/2] {
		heights[i] = h
		for heights[i] < H && b[heights[i]][i] == 0 {
			heights[i]++
		}
		heights[i] -= h
	}

	return heights
}

// classical min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
