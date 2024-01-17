package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	plots := newArea()

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		plots.load(j, input.Bytes())
	}
	plots.eyefill()

	fmt.Println(plots.walk(64), plots.solve(26_501_365))
}

const MAX, OFF = 131, 300

type yx uint32

func YX(j, i int) yx {
	return yx(uint32(i+OFF)&0xffff | uint32(j+OFF)<<16)
}

func Y(p yx) int {
	return int(p>>16) - OFF
}

func X(p yx) int {
	return int(p&0xffff) - OFF
}

func (u yx) add(v yx) yx {
	return u + v - (OFF<<16 | OFF)
	//return YX(Y(u)+Y(v), X(u)+X(v))
}

type area struct {
	m map[yx]bool
	o yx
	w int
}

func newArea() (a *area) {
	a = new(area)
	a.m = make(map[yx]bool, MAX*MAX)
	return
}

func (a *area) load(j int, s []byte) {
	a.w = len(s)
	for i := range s {
		k := YX(j, i)
		switch s[i] {
		case '#':
			a.m[k] = true
		case 'S':
			a.o = k
		}
	}
}

func (a *area) eyefill() {
	Δ := []yx{
		YX(-1, 0), YX(0, -1), YX(0, 1), YX(1, 0),
	}

	for j := 1; j < a.w-1; j++ {
	ROW:
		for i := 1; i < a.w-1; i++ {
			for _, δ := range Δ {
				if a.isplot(YX(j, i).add(δ)) {
					continue ROW
				}
			}
			a.m[YX(j, i)] = true
		}
	}
}

func (a *area) walk(n int) (nreach int) {

	δy := 0
	for x := -n; x <= n; x++ {
		for y := -δy; y <= δy; y += 2 {
			if a.isplot(YX(65+y, 65+x)) {
				nreach++
			}
		}
		δy--
		if x < 0 {
			δy += 2
		}
	}

	return
}

func (a *area) solve(nstep int) int {
	w, h := a.w, a.w/2
	w0h := a.walk(w*0 + h)
	w1h := a.walk(w*1 + h)
	w2h := a.walk(w*2 + h)

	δ1 := w1h - w0h
	δ2 := w2h - w1h

	q := nstep / w

	return w0h + δ1*q + q*(q-1)*(δ2-δ1)/2
}

func (a *area) isrock(k yx) bool {
	θ := func(p yx) yx {
		j, i := Y(p), X(p)
		return YX(mod(j, MAX), mod(i, MAX))
	}

	return a.m[θ(k)]
}

func (a *area) isplot(k yx) bool {
	return !a.isrock(k)
}

func (a *area) String() string {
	var sb strings.Builder

	get := func(j, i int) byte {
		k := YX(j, i)
		switch {
		case k == a.o:
			return 'S'
		case a.m[k]:
			return '#'
		}
		return '.'
	}

	for j := 0; j < a.w; j++ {
		for i := 0; i < a.w; i++ {
			sb.WriteByte(get(j, i))
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}
