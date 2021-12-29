package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var h, w int

type board []byte

func (b board) String() string {
	var sb strings.Builder
	for i := 0; i < h; i++ {
		sb.WriteString(string(b[i*w : (i+1)*w]))
		sb.WriteByte('\n')
	}
	sb.WriteString(string(b[(h-2)*w : (h-1)*w]))

	return sb.String()
}

func clone(b board) board {
	return append(b[:0:0], b...)
}

func main() {
	var cur, nxt board
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		bytes := input.Bytes()
		nxt = append(nxt, bytes...)
		h, w = h+1, len(bytes)
	}

	s, n := 0, 1 // step, change counts
	for n > 0 {
		s, n, cur = s+1, 0, clone(nxt) // reset
		for j := 0; j < h; j++ {       // east scan
			for i := 0; i < w; i++ {
				x := (i + 1) % w
				if cur[j*w+i] == '>' && cur[j*w+x] == '.' {
					nxt[j*w+i], nxt[j*w+x] = '.', '>'
					n++
				}
			}
		}
		cur = clone(nxt)
		for j := 0; j < h; j++ { // south scan
			y := (j + 1) % h
			for i := 0; i < w; i++ {
				if cur[j*w+i] == 'v' && cur[y*w+i] == '.' {
					nxt[j*w+i], nxt[y*w+i] = '.', 'v'
					n++
				}
			}
		}
	}
	fmt.Println(s)
}
