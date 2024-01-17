package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const MAXN = 300

func main() {

	r := strings.NewReplacer(
		",", "",
		"@", "",
	)

	stones := make([]stone, 0, MAXN)

	input := bufio.NewScanner(os.Stdin)
	for j := 0; input.Scan(); j++ {
		var x stone
		for i, s := range fields(r.Replace(input.Text())) {
			x[i] = atoi(s)
		}
		stones = append(stones, x)
	}

	fmt.Println(intersect(stones), collide(stones))
}

type stone [6]int64

func collide(stones []stone) *big.Int {
	p0, v0 := mkMove(stones[0])
	p1, v1 := mkMove(stones[1])
	p2, v2 := mkMove(stones[2])

	p3, v3 := p1.sub(p0), v1.sub(v0)
	p4, v4 := p2.sub(p0), v2.sub(v0)

	q := v3.cross(p3).gcd()
	r := v4.cross(p4).gcd()
	s := q.cross(r).gcd()

	t := div(
		sub(mul(p3[Y], s[X]), mul(p3[X], s[Y])),
		sub(mul(v3[X], s[Y]), mul(v3[Y], s[X])),
	)

	u := div(
		sub(mul(p4[Y], s[X]), mul(p4[X], s[Y])),
		sub(mul(v4[X], s[Y]), mul(v4[Y], s[X])),
	)

	a := p0.add(p3).trace()
	b := p0.add(p4).trace()
	c := v3.sub(v4).trace()

	return div(
		add(
			sub(
				mul(u, a),
				mul(t, b),
			),
			mul(
				u,
				mul(t, c),
			),
		),
		sub(u, t),
	)
}

func intersect(stones []stone) int {
	ninter := 0

	for i, ii := range stones[1:] {
		a, b, c, d := ii[0], ii[1], ii[3], ii[4]
	INNER:
		for _, jj := range stones[:i+1] {
			e, f, g, h := jj[0], jj[1], jj[3], jj[4]
			Δ := (d * g) - (c * h)
			if Δ == 0 {
				continue INNER
			}

			t := (g*(f-b) - h*(e-a)) / Δ
			u := (c*(f-b) - d*(e-a)) / Δ

			x := a + t*c
			y := b + t*d

			inrange := func(x int64) bool {
				return x >= 200_000_000_000_000 && x <= 400_000_000_000_000
			}

			if t >= 0 && u >= 0 && inrange(x) && inrange(y) {
				ninter++
			}
		}
	}

	return ninter
}

const (
	X = iota
	Y
	Z
)

type vec3 [3]*big.Int

func mkMove(s stone) (p, v vec3) {
	p = vec3{
		big.NewInt(s[0]), big.NewInt(s[1]), big.NewInt(s[2]),
	}
	v = vec3{
		big.NewInt(s[3]), big.NewInt(s[4]), big.NewInt(s[5]),
	}
	return
}

func (u vec3) add(v vec3) (s vec3) {
	for i := range s {
		s[i] = new(big.Int).Add(u[i], v[i])
	}
	return
}

func (u vec3) sub(v vec3) (s vec3) {
	for i := range s {
		s[i] = new(big.Int).Sub(u[i], v[i])
	}
	return
}

func (u vec3) cross(v vec3) (c vec3) {
	for i := range c {
		j, k := (i+1)%3, (i+2)%3
		c[i] = new(big.Int).Sub(
			new(big.Int).Mul(u[j], v[k]),
			new(big.Int).Mul(u[k], v[j]),
		)
	}

	return
}

func (u vec3) gcd() (g vec3) {
	gcd := func(a, b, c *big.Int) (g *big.Int) {
		g = new(big.Int)
		g.GCD(nil, nil, a, b)
		g.GCD(nil, nil, g, c)
		return
	}

	γ := gcd(u[X], u[Y], u[Z])
	for i := range g {
		g[i] = new(big.Int).Div(u[i], γ)
	}
	return
}

func (u vec3) trace() (t *big.Int) {
	t = new(big.Int)
	t.Add(t.Add(u[X], u[Y]), u[Z])
	return
}

func add(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func div(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}

var fields = strings.Fields

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(s string) (n int64) {
	neg := int64(1)
	if s[0] == '-' {
		neg, s = -1, s[1:]
	}

	for i := range s {
		n = 10*n + int64(s[i]-'0')
	}
	return neg * n
}
