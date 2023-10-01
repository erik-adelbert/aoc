// benchmark with:
// $ go test -bench=. -benchmem

package main

import (
	"bufio"
	_ "embed"
	"strings"
	"testing"
)

//go:embed input.txt
var input string

var result any

func BenchmarkMkBuros(b *testing.B) {
	var buros []buro

	for i := 0; i < b.N; i++ {
		buros = loadBuro(input)
	}

	result = buros
}

func BenchmarkBuroString(b *testing.B) {
	var r string
	var x buro

	for i := 0; i < b.N; i++ {
		r = x.String()
	}

	result = r
}

func BenchmarkBuroGet(b *testing.B) {
	var r byte
	var x buro

	for i := 0; i < b.N; i++ {
		r = x.get(1, 2)
	}

	result = r
}

func BenchmarkBuroSet(b *testing.B) {
	var x buro

	for i := 0; i < b.N; i++ {
		x.set(1, 2, 'Z')
	}

	result = x
}

func BenchmarkBuroHome(b *testing.B) {
	var r int
	var x buro

	for i := 0; i < b.N; i++ {
		r = x.home('A')
	}

	result = r
}

func loadBuro(input string) []buro {
	return mkburos(bufio.NewScanner(strings.NewReader(input)))
}

func BenchmarkBuroPart1Hcost(b *testing.B) {
	x := loadBuro(input)

	var c cost

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c = x[0].hcost()
	}

	result = c
}

func BenchmarkBuroPart2Hcost(b *testing.B) {
	x := loadBuro(input)

	var c cost

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c = x[1].hcost()
	}

	result = c
}

func BenchmarkBuroPeek(b *testing.B) {
	x := loadBuro(input)

	var c byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c = x[1].peek(3)
	}

	result = c
}

func BenchmarkBuroPop(b *testing.B) {
	x := loadBuro(input)

	var c byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, _ = x[1].pop(3)
	}

	result = c
}

func BenchmarkBuroPush(b *testing.B) {
	x := loadBuro(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x[1].push(1, 'A')
	}

	result = x
}

func BenchmarkBuroIsDead(b *testing.B) {
	x := loadBuro(input)

	dead := false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dead = x[1].isdead()
	}

	result = dead
}

func benchmarkMove(t, s int, inplace bool, b *testing.B) {
	x := loadBuro(input)

	m := newMove(&x[1], 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ = m.move(t, s, inplace)
	}

	result = m
}

func BenchmarkMoveInplaceMoveOk(b *testing.B) {
	benchmarkMove(1, 3, true, b)
}

func BenchmarkMoveAllocMoveOk(b *testing.B) {
	benchmarkMove(1, 3, false, b)
}

func BenchmarkMoveInplaceMoveNotOk(b *testing.B) {
	benchmarkMove(3, 1, true, b)
}

func BenchmarkMoveAllocMoveNotOk(b *testing.B) {
	benchmarkMove(3, 1, false, b)
}

func BenchmarkMoveMoves(b *testing.B) {
	x := loadBuro(input)

	var moves []*move
	m := newMove(&x[1], 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		moves = m.moves()
	}

	result = moves
}

func benchmarkMoveSolve(part int, b *testing.B) {
	x := loadBuro(input)
	m := newMove(&x[part], 0)

	var r cost

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = m.solve()
	}

	result = r
}

func BenchmarkMoveSolvePart1(b *testing.B) {
	benchmarkMoveSolve(0, b)
}

func BenchmarkMoveSolvePart2(b *testing.B) {
	benchmarkMoveSolve(1, b)
}

func BenchmarkStdHashmapRoundTrip(b *testing.B) {
	x := loadBuro(input)

	var r int
	m := make(map[buro]int, 1<<16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[x[1]] = 1
		r = m[x[1]]
	}

	result = r
}

func BenchmarkFNV1HashmapRoundTrip(b *testing.B) {
	x := loadBuro(input)

	var r int
	m := make(map[uint64]int, 1<<16)
	h := x[1].hash()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[h] = 1
		r = m[h]
	}

	result = r
}
