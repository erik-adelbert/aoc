package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// idx abstracts an index of a 2D matrix
type idx struct {
	x, y int
}

// idxmap precomputes all 8-neighbor indices in a w*h matrix
func idxmap(w, h int) [][]idx {
	idxmap := make([][]idx, w*h)
	for y := 0; y < 10; y++ { // rc shit
		for x := 0; x < 10; x++ {
			idxmap[y*10+x] = mask8(idx{x, y}, w, h)
		}
	}
	return idxmap
}

// mask8 builds adaptative indices for 8-neighbor matrices
// it handles well corners and borders
func mask8(cell idx, w, h int) []idx {
	mask := make([]idx, 0, 8)
	for δy := -1; δy < 2; δy++ {
		for δx := -1; δx < 2; δx++ {
			if δx != 0 || δy != 0 {
				i, j := cell.x+δx, cell.y+δy            // all of this is so painful
				if i < 0 || j < 0 || i >= w || j >= h { // out of bounds
					continue // reject
				}
				mask = append(mask, idx{i, j})
			}
		}
	}
	return mask
}

// submarine cave abstraction
type cave struct {
	cells  [][]byte // each cell contains a (0..9) value
	w, h   int
	neighs [][]idx // neighbors
}

func newCave(w, h int) *cave {
	cells := make([][]byte, h)
	for i := 0; i < h; i++ {
		cells[i] = make([]byte, w)
	}
	return &cave{cells, w, h, idxmap(w, h)}
}

func (c *cave) popcount() int {
	return c.w * c.h
}

func (c *cave) set(cl idx, b byte) {
	c.cells[cl.y][cl.x] = b
}

func (c *cave) inc(cl idx) byte {
	b := (c.cells[cl.y][cl.x] + 1) % 10
	c.cells[cl.y][cl.x] = b
	return b
}

func (c *cave) neighbors(cl idx) []idx {
	return c.neighs[cl.y*c.w+cl.x]
}

func (c *cave) String() string {
	var sb strings.Builder
	for j := 0; j < c.h; j++ {
		for _, b := range c.cells[j][:c.w] {
			sb.WriteByte(b + '0')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// A blast is entirely made of flashing cells.
type blast map[idx]bool

// safe determines if a cave has steady (non flashing) cells (safe) or not (unsafe).
// It computes the global blast of one step taken from  a given cave state. It also
// asserts safeness by checking on steady cells.
func safe(c *cave) bool {
	// main reaction blast
	cur := make(blast)
	for j := 0; j < c.h; j++ {
		for i := 0; i < c.w; i++ {
			cl := idx{i, j}     // cave cell
			if c.inc(cl) == 0 { // flashing when 0
				cur[cl] = true
			}
		}
	}
	_, safe := c.cascade(make(blast), cur) // discard global blast
	return safe
}

// cascade computes and aggregate chain reaction blasts & asserts global safeness
// it double buffers using two maps because updating a map while iterating it is
// undefined.
func (c *cave) cascade(glob, cur blast) (blast, bool) {
	// chain reaction blast
	for {
		nxt := make(blast)
		for flash := range cur {
			for _, ne := range c.neighbors(flash) {
				if !glob[ne] && !cur[ne] && !nxt[ne] { // new one!
					if c.inc(ne) == 0 { // flashing
						nxt[ne] = true // neighbor chain reacts
					}
				}
			}
		}

		for k, v := range cur { // add current blast to global one
			glob[k] = v
		}
		cur = nxt

		if len(nxt) == 0 { // no more flashing cell
			return glob, len(glob) != c.popcount() // safe/true when some non flashing
		}
	}
}

func main() {
	ctob := func(b byte) byte {
		return b - '0' // fast convert
	}

	cave := newCave(10, 10)
	j, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		for i, b := range input.Bytes() {
			cave.set(idx{i, j}, ctob(b))
		}
		j++
	}

	n := 1
	for safe(cave) {
		n++
	}
	fmt.Println(n)
}
