// aoc22.go --
// advent of code 2022 day 22
//
// https://adventofcode.com/2022/day/22
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-22: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// fieldalignment auto-layout
type cubemap struct {
	faces   map[UV]Euv
	edges   map[EU]UV
	tex     bmp
	path    []string
	ij      []lim
	ji      []lim
	i0, j0  int
	h, w, N int
}

func main() {
	input := bufio.NewScanner(os.Stdin)

	cube := mkcubemap(150)
	cube.load(input)

	fmt.Println(cube.walk(1))
	fmt.Println(cube.walk(2))
}

// axis
const (
	X = iota
	Y
	Z
)

// bounds
const (
	Min = iota
	Max
)

type (
	lim [2]int
	UV  [2]int
	XYZ [3]int
	EU  struct{ E, U XYZ }
	Euv struct{ E, u, v XYZ }
)

func mkcubemap(N int) *cubemap {
	c := new(cubemap)
	c.tex = make([][]byte, 0, N)
	return c
}

func (c *cubemap) load(input *bufio.Scanner) {
	// cubemap
	c.h, c.w = 0, 0
	for input.Scan() {
		var line []byte
		if line = input.Bytes(); len(line) == 0 {
			break
		}
		c.h++
		c.w = max(c.w, len(line))
		r := make([]byte, len(line))
		copy(r, line)
		c.tex = append(c.tex, r)
	}

	minmax := func(n int) []lim {
		a := make([]lim, n)
		for i := range a {
			a[i][Min] = 1000
			a[i][Max] = -1
		}
		return a
	}

	c.ij = minmax(c.h)
	c.ji = minmax(c.w)

	for i, r := range c.tex {
		for j, x := range r {
			if x != ' ' {
				c.ij[i][Min] = min(c.ij[i][Min], j)
				c.ij[i][Max] = max(c.ij[i][Max], j)
				c.ji[j][Min] = min(c.ji[j][Min], i)
				c.ji[j][Max] = max(c.ji[j][Max], i)
				c.N++
			}
		}
	}
	c.N = isqrt(c.N / 6)

	// cmds
	input.Scan() // one line
	r := strings.NewReplacer(
		"L", " L ",
		"R", " R ",
	)
	c.path = strings.Fields(r.Replace(input.Text()))

	i0, j0 := 0, 1000
	for j, c := range c.tex[0] {
		if c == '.' {
			j0 = min(j0, j)
		}
	}
	c.i0, c.j0 = i0, j0

	c.edges = make(map[EU]UV)
	c.faces = make(map[UV]Euv)
	c.fold(i0, j0, XYZ{0, 0, 0}, XYZ{1, 0, 0}, XYZ{0, 1, 0})

	return
}

func (c *cubemap) out(i, j int) bool {
	return i < 0 || i >= c.h || j < 0 || j >= len(c.tex[i]) || c.tex[i][j] == ' '
}

func (c *cubemap) fold(i, j int, E, u, v XYZ) {
	_, seen := c.faces[UV{i, j}]
	if c.out(i, j) || seen {
		return
	}

	c.faces[UV{i, j}] = Euv{E, u, v}

	N := c.N
	U, D := u.cross(v), v.cross(u) // Up/Down
	for r := 0; r < N; r++ {
		// E + u*r, U
		eu := EU{E.add(u.mul(r)), U}
		c.edges[eu] = UV{i + r, j}

		// E + u*r+v*(N-1), U
		eu = EU{E.add(u.mul(r), v.mul(N-1)), U}
		c.edges[eu] = UV{i + r, j + N - 1}

		// E + v*r, U
		eu = EU{E.add(v.mul(r)), U}
		c.edges[eu] = UV{i, j + r}

		// E + v*r+u*(N-1), U
		eu = EU{E.add(v.mul(r), u.mul(N-1)), U}
		c.edges[eu] = UV{i + N - 1, j + r}
	}

	c.fold(i+N, j, E.add(u.mul(N-1)), U, v)
	c.fold(i-N, j, E.add(U.mul(N-1)), U, v)
	c.fold(i, j+N, E.add(v.mul(N-1)), u, U)
	c.fold(i, j-N, E.add(U.mul(N-1)), u, D)
}

func (c *cubemap) step(p, n, i, j, δi, δj int) (int, int, int, int) {
	N := c.N

	for n > 0 {
		k, l, δk, δl := i+δi, j+δj, δi, δj
		if c.out(k, l) {
			// part
			switch p {
			case 1:
				switch {
				default:
				case δi == 0:
				case k < c.ji[j][Min]:
					k = c.ji[j][Max]
				case k > c.ji[l][Max]:
					k = c.ji[l][Min]
				}
				switch {
				default:
				case δj == 0:
				case l < c.ij[i][Min]:
					l = c.ij[i][Max]
				case l > c.ij[k][Max]:
					l = c.ij[i][Min]
				}
			case 2:
				f := c.faces[UV{i / N * N, j / N * N}]
				E, u, v := f.E, f.u, f.v

				O := E.add(u.mul(i%N), v.mul(j%N))
				uv := u.cross(v)

				IJ := c.edges[EU{O, u.mul(-δi).add(v.mul(-δj))}]
				k, l = IJ[X], IJ[Y]

				f = c.faces[UV{k / N * N, l / N * N}]
				u, v = f.u, f.v

				δk, δl = u.dot(uv), v.dot(uv)
			}
		}
		if c.tex[k][l] == '#' {
			break
		}
		i, j, δi, δj = k, l, δk, δl
		n--
	}
	return i, j, δi, δj
}

func (c *cubemap) walk(p int) int {
	i, j := c.i0, c.j0
	δi, δj := 0, 1
	for _, x := range c.path {
		switch x {
		case "L":
			δi, δj = -δj, δi
		case "R":
			δi, δj = δj, -δi
		default:
			i, j, δi, δj = c.step(p, atoi(x), i, j, δi, δj)
		}
	}

	scale := map[UV]int{
		{0, 1}:  0,
		{1, 0}:  1,
		{0, -1}: 2,
		{-1, 0}: 3,
	}

	return 1000*(i+1) + 4*(j+1) + scale[UV{δi, δj}]
}

type (
	bmp [][]byte
)

func (b bmp) String() string {
	var sb strings.Builder
	for _, r := range b {
		sb.Write(r)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (a XYZ) add(b ...XYZ) XYZ {
	for _, x := range b {
		for i := X; i <= Z; i++ {
			a[i] += x[i]
		}
	}

	return a
}

func (a XYZ) mul(k int) XYZ {
	for i := X; i <= Z; i++ {
		a[i] *= k
	}
	return a
}

func (a XYZ) dot(b XYZ) int {
	for i := X; i <= Z; i++ {
		a[i] *= b[i]
	}
	return a[X] + a[Y] + a[Z]
}

func (a XYZ) cross(b XYZ) XYZ {
	a[X], a[Y], a[Z] = a[Y]*b[Z]-a[Z]*b[Y],
		a[Z]*b[X]-a[X]*b[Z],
		a[X]*b[Y]-a[Y]*b[X]
	return a
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

var tab64 = [64]uint64{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5,
}

func log2(n uint64) uint64 {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return tab64[((n-(n>>1))*0x07EDD5E59A4E28C2)>>58]
}

func isqrt(x int) int {
	x64 := uint64(x)
	var b, r uint64
	for p := uint64(1 << ((uint(log2(x64)) >> 1) << 1)); p != 0; p >>= 2 {
		b = r | p
		r >>= 1
		if x64 >= b {
			x64 -= b
			r |= p
		}
	}
	return int(r)
}

func debug(a ...any) {
	if false {
		fmt.Println(a...)
	}
}
