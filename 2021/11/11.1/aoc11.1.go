package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// submarine cave abstraction
type cave struct {
	cells []byte
	h, w  int
}

func newCave(h, w int) *cave {
	cells := make([]byte, h*w)
	return &cave{cells, h, w}
}

func (ca *cave) inc(r, c int) byte {
	w := ca.w
	b := (ca.cells[r*w+c] + 1) % 10
	ca.cells[r*w+c] = b
	return b
}

func (c *cave) String() string {
	var sb strings.Builder
	for j := 0; j < c.h; j++ {
		for i := 0; i < c.w; i++ {
			sb.WriteByte(c.cells[j*c.w+i] + '0')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type coo [2]int

const (
	R = iota
	C
)

type blast map[coo]bool

func flash(c *cave) int {
	cur := make(blast)
	for j := 0; j < c.h; j++ {
		for i := 0; i < c.w; i++ {
			if c.inc(j, i) == 0 { // flashing
				cur[coo{j, i}] = true // record in current blast
			}
		}
	}
	_, n := c.cascade(make(blast), cur) // discard global blast
	return n
}

func (c *cave) cascade(glob, cur blast) (blast, int) {
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
				if !glob[coo{j, i}] && !cur[coo{j, i}] && !nxt[coo{j, i}] { // new one!
					if c.inc(j, i) == 0 { // flashing
						nxt[coo{j, i}] = true // neighbor chain reacts
					}
				}
			}
		}
		for k, v := range cur {
			glob[k] = v // merge current blast to the global one
		}
		cur = nxt

		if len(nxt) == 0 { // no more flashing cell
			return glob, len(glob)
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
		for i, c := range input.Bytes() {
			cave.cells[j*10+i] = ctob(c)
		}
		j++
	}

	n := 0
	for i := 0; i < 100; i++ {
		n += flash(cave)
	}
	fmt.Println(n)
}
