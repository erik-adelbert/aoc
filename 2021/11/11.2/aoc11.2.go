package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// submarine cave abstraction
type cave struct {
	cells []byte // each cell contains a (0..9) value
	h, w  int
}

func newCave(h, w int) *cave {
	cells := make([]byte, h*w)
	return &cave{cells, h, w}
}

func (c *cave) popcount() int {
	return c.h * c.w
}

func (c *cave) inc(ji idx) byte {
	j, i, w := ji[0], ji[1], c.w
	b := (c.cells[j*w+i] + 1) % 10
	c.cells[j*w+i] = b
	return b
}

func (c *cave) String() string {
	var sb strings.Builder
	for j := 0; j < c.h; j++ {
		for i := 0; i < c.w; i++ {
			sb.WriteByte(c.cells[j*c.w+i] + '0')
		}
		if j != c.h {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

type idx [2]int // row major index

// Row and Column
const (
	R = iota
	C
)

// A blast is entirely made of flashing cells.
type blast map[idx]bool

func (b blast) popcount() int {
	return len(b)
}

// safe determines if a cave has steady (non flashing) cells (safe) or not (unsafe).
// It computes the global blast of one step taken from  a given cave state. It also
// asserts safeness by checking on steady cells.
func safe(c *cave) bool {
	// main reaction blast
	cur := make(blast)
	for j := 0; j < c.h; j++ {
		for i := 0; i < c.w; i++ {
			ji := idx{j, i}
			if c.inc(ji) == 0 { // flashing when 0
				cur[ji] = true
			}
		}
	}
	return c.cascade(make(blast), cur)
}

// cascade computes and aggregate chain reaction blasts & asserts global safeness
// it double buffers using two maps because updating a map while iterating it is
// undefined.
func (c *cave) cascade(glob, cur blast) bool {
	// chain reaction blast
	for {
		nxt := make(blast)
		for flash := range cur {
			δj := []int{-1, -1, -1, +0, +0, +1, +1, +1}
			δi := []int{-1, +0, +1, -1, +1, -1, +0, +1}

			for k := 0; k < len(δj); k++ {
				j, i := flash[R]+δj[k], flash[C]+δi[k]
				if j < 0 || j >= c.h || i < 0 || i >= c.w {
					continue
				}
				ji := idx{j, i}
				if !glob[ji] && !cur[ji] && !nxt[ji] { // new one!
					if c.inc(ji) == 0 { // flashing
						nxt[ji] = true // neighbor chain reacts
					}
				}
			}
		}

		for k, v := range cur { // add current blast to global one
			glob[k] = v
		}

		if len(nxt) == 0 { // no more flashing cell
			return glob.popcount() != c.popcount() // safe/true when some non flashing
		}

		cur, nxt = nxt, nil
	}
}

func main() {
	ctob := func(b byte) byte {
		return b - '0' // fast convert
	}

	cave := newCave(10, 10)
	j, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for i, c := range input.Bytes() {
			cave.cells[j*10+i] = ctob(c)
		}
		j++
	}

	n := 1
	for safe(cave) {
		n++
	}
	fmt.Println(n)
}
