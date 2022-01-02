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
	y, x int
}

type bitmap struct {
	kern *[512]byte     // kernel filter
	bbox [2]pixel       // bounding box
	data map[pixel]byte // raw data
}

var par = false // par(ity)

func (b bitmap) inf(y, x int) bool {
	switch {
	case y < b.bbox[min].y || y > b.bbox[max].y:
		fallthrough
	case x < b.bbox[min].x || x > b.bbox[max].x:
		return true
	}
	return false
}

func (b bitmap) get(y, x int) int {
	if b.inf(y, x) { // p is infinite
		if par {
			return 1
		}
		return 0
	}
	return int(b.data[pixel{y, x}])
}

func (b bitmap) count() int {
	return len(b.data)
}

func (b bitmap) enhance() *bitmap {
	nxt := bitmap{ // double buffer
		kern: b.kern,
		bbox: [2]pixel{
			{b.bbox[min].y - 1, b.bbox[min].x - 1},
			{b.bbox[max].y + 1, b.bbox[max].x + 1},
		},
		data: make(map[pixel]byte, 1<<15),
	}

	kern9 := func(y, x int) byte { // apply filter
		δy := []int{-1, -1, -1, +0, 0, 0, +1, 1, 1}
		δx := []int{-1, +0, +1, -1, 0, 1, -1, 0, 1}

		n := 0
		for i := 0; i < len(δx); i++ {
			n = (n << 1) | b.get(y+δy[i], x+δx[i])
		}
		return b.kern[n]
	}

	for y := b.bbox[min].y - 1; y <= b.bbox[max].y+1; y++ {
		for x := b.bbox[min].x - 1; x <= b.bbox[max].x+1; x++ {
			if kern9(y, x) == 1 {
				nxt.data[pixel{y, x}] = 1 // it's a map!
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
			if b.get(y, x) == 1 {
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
		kern: &[512]byte{},
		bbox: [...]pixel{
			{0, 0},
			{len(raw) - 1, len(raw[0]) - 1},
		},
		data: make(map[pixel]byte, 1<<15),
	}

	for i, c := range ker {
		if c == '#' {
			b.kern[i] = 1
		}
	}

	for y, row := range raw {
		for x, v := range row {
			if v == '#' {
				b.data[pixel{y, x}] = 1
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

	for i := 0; i < 50; i++ {
		if i == 2 {
			fmt.Println(img.count()) // part1
		}
		img = img.enhance()
	}
	fmt.Println(img.count()) // part2
}
