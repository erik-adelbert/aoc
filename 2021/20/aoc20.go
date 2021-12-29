package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	min = 0
	max = 1
)

type pixel struct {
	x, y int
}

type bitmap struct {
	kern [512]byte      // kernel filter
	bbox [2]pixel       // bounding box
	data map[pixel]byte // raw data
}

var par = false // par(ity)

func (b bitmap) inf(p pixel) bool {
	switch {
	case p.x < b.bbox[min].x || p.x > b.bbox[max].x:
		fallthrough
	case p.y < b.bbox[min].y || p.y > b.bbox[max].y:
		return true
	}
	return false
}

func (b bitmap) get(p pixel) int {
	if b.inf(p) { // p is infinite
		if par {
			return 1
		}
		return 0
	}
	return int(b.data[p])
}

func (b bitmap) count() int {
	return len(b.data)
}

func (b bitmap) enhance() *bitmap {
	nxt := bitmap{ // double buffer
		kern: b.kern,
		bbox: [2]pixel{
			{b.bbox[min].x - 1, b.bbox[min].y - 1},
			{b.bbox[max].x + 1, b.bbox[max].y + 1},
		},
		data: make(map[pixel]byte, len(b.data)),
	}

	kern9 := func(p pixel) byte { // apply filter
		δx := []int{-1, 0, 1, -1, 0, 1, -1, 0, 1}
		δy := []int{-1, -1, -1, 0, 0, 0, 1, 1, 1}

		n := 0
		for i := 0; i < len(δx); i++ {
			x, y := p.x+δx[i], p.y+δy[i]
			n = (n << 1) | b.get(pixel{x, y})
		}
		return b.kern[n]
	}

	for y := b.bbox[min].y - 1; y <= b.bbox[max].y+1; y++ {
		for x := b.bbox[min].x - 1; x <= b.bbox[max].x+1; x++ {
			if kern9(pixel{x, y}) == 1 {
				nxt.data[pixel{x, y}] = 1 // it's a map!
			}
		}
	}

	par = !par
	return &nxt
}

func (b bitmap) String() string {
	var sb strings.Builder

	for y := b.bbox[min].y - 1; y <= b.bbox[max].y+1; y++ {
		for x := b.bbox[min].x - 1; x <= b.bbox[max].x+1; x++ {
			if b.get(pixel{x, y}) == 1 {
				sb.WriteByte('@')
			} else {
				sb.WriteByte('.')
			}
		}
		if y != b.bbox[max].y+1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

func newBitmap(ker []byte, raw [][]byte) *bitmap {
	b := bitmap{
		kern: [512]byte{},
		bbox: [...]pixel{
			{0, 0},
			{len(raw[0]) - 1, len(raw) - 1},
		},
		data: make(map[pixel]byte, 64*64),
	}

	for i, c := range ker {
		if c == '#' {
			b.kern[i] = 1
		}
	}

	for y, row := range raw {
		for x, v := range row {
			if v == '#' {
				b.data[pixel{x, y}] = 1
			}
		}
	}
	return &b
}

func main() {
	var raw [][]byte
	kern := make([]byte, 512) // kern(el filter)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		switch len(input.Bytes()) {
		case 0: // continue
		case 512:
			copy(kern, input.Bytes())
		default:
			bytes := make([]byte, len(input.Bytes()))
			copy(bytes, input.Bytes())
			raw = append(raw, bytes)
		}
	}
	img := newBitmap(kern, raw)

	count := func(lim int) {
		tmp := img
		for i := 0; i < lim; i++ {
			tmp = tmp.enhance()
		}
		fmt.Println(tmp.count())
	}

	count(2)  // part1
	count(50) // part2
}
