package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	u := make(universe, 0, 128)

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		u = u.readline(j, input.Text())
	}

	fmt.Println(u.pairdist(2), u.pairdist(1_000_000))
}

type galaxy struct {
	y, x int
}

func cmp(a, b galaxy) int {
	if a.x == b.x {
		return a.y - b.y
	}
	return a.x - b.x
}

type universe []galaxy

func (u universe) readline(j int, line string) universe {
	for i := range line {
		if line[i] == '#' {
			u = append(u, galaxy{j, i})
		}
	}
	return u
}

func (u universe) split2D() [2][]int {
	X, Y := make([]int, len(u)), make([]int, len(u))
	for i := range u {
		Y[i], X[i] = u[i].y, u[i].x
	}
	return [2][]int{Y, X}
}

func expand(C counter, k int) (CC counter) {
	type sd struct {
		src, dst int
	}

	remap := func(L []int) (T []sd) {
		T = make([]sd, 0, len(L))
		CC = make(counter, len(L))

		old := sd{-1, -1}
		for _, x := range L {
			n := old.dst + k*(x-old.src-1) + 1
			T = append(T, sd{x, n})
			old = sd{x, n}
		}
		return
	}

	for _, x := range remap(C.list()) {
		CC[x.dst] = C[x.src]
	}

	return CC
}

// k factor, d pairwise distance sum
func (u universe) pairdist(k int) (d int) {
	d = 0
	for _, dim := range u.split2D() {
		n := pairdist1D(expand(count(dim), k))
		d += n
	}
	return
}

func pairdist1D(C counter) int {
	off, K, X := C.sum(), 0, 0
	for k, v := range C {
		K += k * v
	}

	acc := 0
	for _, x := range C.list() {
		k := K - off*(x-X)
		acc += k * C[x]
		off -= C[x]
		K, X = k, x
	}

	return acc
}

type counter map[int]int

func count(L []int) (C counter) {
	C = make(counter, len(L))
	for i := range L {
		C[L[i]]++
	}
	return
}

func (C counter) list() (L []int) {
	L = make([]int, 0, len(C))
	for i := range C {
		L = append(L, i)
	}
	slices.Sort(L)

	return
}

func (C counter) sum() (S int) {
	S = 0
	for _, x := range C {
		S += x
	}
	return
}
