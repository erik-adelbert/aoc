package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	var p1, p2 lagoon

	decode := func(s string) (byte, int) {
		θ := "RDLU"[s[len(s)-1]-'0'] // last char encodes R, D, L or U
		k := htoi(s[:len(s)-1])      // first chars encodes an hex number
		return θ, k
	}

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		args := Fields(input.Text())

		θ, k := args[0][0], atoi(args[1])
		p1 = p1.append(θ, k)

		x := args[2][2 : len(args[2])-1] // slice the hex part off args[2] "^(#\h+)$" with \h hex digit
		p2 = p2.append(decode(x))
	}

	fmt.Println(p1.area(), p2.area())
}

type vec struct {
	y, x int
}

func (p vec) add(u vec) vec {
	return vec{p.y + u.y, p.x + u.x}
}

func (p vec) scale(k int) vec {
	return vec{k * p.y, k * p.x}
}

type lagoon struct {
	cur        vec
	peri, lace int
}

func (p lagoon) append(θ byte, k int) lagoon {
	δ := []vec{
		'R': {0, -1}, 'D': {-1, 0}, 'L': {0, 1}, 'U': {1, 0},
	}

	cur := p.cur
	new := cur.add(δ[θ].scale(k))

	p.peri += k
	p.lace += cur.x*new.y - new.x*cur.y // shoelace formula

	p.cur = new
	return p
}

func (p lagoon) area() int {
	return (p.peri+p.lace)/2 + 1
}

var Fields, Index = strings.Fields, strings.Index

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}

func htoi(s string) (n int) {
	ctoi := func(c byte) int {
		return Index("0123456789abcdef", string(c))
	}

	for i := range s {
		n = 16*n + ctoi(s[i])
	}
	return
}
