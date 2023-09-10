package main

import (
	"bufio"
	"fmt"
	"os"
)

// windy maze
type waze [2][]rc

// timed waze
type taze []waze

func main() {
	w := new(waze)

	input := bufio.NewScanner(os.Stdin)

	row := make([]byte, 0, 128)
	input.Scan() // discard first row
	for input.Scan() {
		if len(row) > 0 {
			// discard first and last col
			w = w.append(row[1 : len(row)-1])
		}
		row = input.Bytes()
	}
	// discard last row

	// build maze over time
	m := mkmaze(w)

	t0 := 0
	t1 := m.flood(t0, Fwd) // part 1
	t2 := m.flood(t1, Bck) // part 1.5
	t3 := m.flood(t2, Fwd) // part 2

	// part 1 & 2
	fmt.Println(t1, t3)
}

const (
	J = iota
	I
)

type JI [2]int

func (a JI) add(b JI) JI {
	return JI{a[J] + b[J], a[I] + b[I]}
}

func (a JI) eq(b JI) bool {
	return a[J] == b[J] && a[I] == b[I]
}

type JISet map[JI]struct{}

func (s JISet) add(x JI) {
	s[x] = struct{}{}
}

func clear(s JISet) JISet {
	// optimized at compile time
	for k := range s {
		delete(s, k)
	}
	return s
}

const (
	Fwd = iota
	Bck
)

func (m *taze) flood(t0, dir int) int {
	w0 := (*m)[0]
	H, W := len(w0[R]), len(w0[C])

	a, z := JI{0, 0}, JI{H - 1, W - 1}
	if dir == Bck {
		a, z = z, a
	}

	out := func(p JI) bool {
		return p[J] < 0 || p[J] >= H || p[I] < 0 || p[I] >= W
	}

	free := func(t int, x JI) bool {
		tJ, tI := t%W, t%H
		return (*m)[tJ][R][x[J]][x[I]] == 0 &&
			(*m)[tI][C][x[I]][x[J]] == 0
	}

	// enter asap
	for !free(t0, a) {
		t0++
	}
	t := t0
	cur := make(JISet, H*W/3)
	nxt := make(JISet, H*W/3)
	nxt.add(a)
	for { // t, passing time
		Δ := []JI{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

		cur, nxt = nxt, clear(cur) // avoid map allocs
		if len(cur) == 0 {
			// no way through, restart later
			return m.flood(t0+1, dir)
		}

		// generate paths (moves in time)
		for x := range cur {
			if free(t+1, x) {
				// staying put is a valid next move
				nxt.add(x)
			}
			for _, δ := range Δ {
				x := x.add(δ)
				switch {
				case z.eq(x):
					return t + 2 // goal!
				case !out(x) && free(t+1, x):
					nxt.add(x) // en route
				}
			}
		}
		t++
	}
}

func mkmaze(w *waze) *taze {
	w.split()

	m := make(taze, len(w[C]))
	// m := make(taze, max(len(w[R]), len(w[C])))

	H := len(w[R])
	rows, cols := w[R], w[C]
	for t := range m {

		// gen row
		m[t][R] = make([]rc, len(rows))
		copy(m[t][R], rows)
		for i := range rows {
			rows[i] = rows[i].next(R)
		}

		// gen col
		if t < H {
			m[t][C] = make([]rc, len(cols))
			copy(m[t][C], cols)

			for i := range cols {
				cols[i] = cols[i].next(C)
			}
		}
	}

	return &m
}

var code = []byte{
	'.': 0, '<': 1, '>': 2, '^': 3, 'v': 4, '#': 5,
}

func (w *waze) append(r []byte) *waze {
	row := make(rc, len(r))
	for i, b := range r {
		row[i] = code[b]
	}
	w[R] = append(w[R], row)

	return w
}

func (w *waze) split() {
	w[C] = make([]rc, len(w[R][0]))
	for i := range w[C] {
		w[C][i] = make(rc, len(w[R]))
	}
	for j := range w[R] {
		for i, b := range w[R][j] {
			switch b {
			case 1, 2:
				w[C][i][j] = 0
			case 3, 4:
				w[R][j][i] = 0
				fallthrough
			default:
				w[C][i][j] = b
			}
		}
	}
}

// row/col indices
const (
	R = iota
	C
)

// row or column
type rc []byte

// rc multi-values are encoded base 8
// here are base 10 for the principle
//   |  . |  < |  > |  ^ |  v
// . |    | 01 | 02 | 03 | 04
// < | 10 |    | 12 | 13 | 14
// > | 20 | 21 |    | 23 | 24
// ^ | 30 | 31 | 32 |    | 34
// v | 40 | 41 | 42 | 43 |

func (x rc) next(d int) rc {
	w := len(x)
	nxt := make(rc, w)

	l, r := code['<'], code['>']
	if d == C {
		l, r = code['^'], code['v']
	}

	for i, c := range x {
		for ; c > 0; c >>= 3 {
			switch c & 0x7 {
			case l:
				i := mod(i-1, w)
				nxt[i] = nxt[i]<<3 + c&0x7
			case r:
				i := mod(i+1, w)
				nxt[i] = nxt[i]<<3 + c&0x7
			default:
				nxt[i] = nxt[i]<<3 + c
			}
		}
	}
	return nxt
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

var DEBUG = true

func debug(a ...any) (int, error) {
	if DEBUG {
		return fmt.Println(a...)
	}
	return 0, nil
}
