// aoc16.go --
// advent of code 2021 day 16
//
// https://adventofcode.com/2021/day/16
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2021-12-16: initial commit

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"

	bs "github.com/bearmini/bitstream-go"
)

var nread uint // global bit count

type seg struct { // datagram segment
	ver uint8
	typ uint8
	val int
	sub []seg
}

func lit(r *bs.Reader) int {
	chunk := func(r *bs.Reader) (uint8, bool) {
		n, _ := r.ReadNBitsAsUint8(5)
		nread += 5
		return n & 0xf, (n & 0x10) > 0
	}

	val := 0
	for {
		n, more := chunk(r)
		val = (val << 4) | int(n)
		if !more {
			return val
		}
	}
}

func subs(r *bs.Reader) []seg {
	enum, _ := r.ReadBool() // load mode: enumerate or read subsegments
	nread++

	usize := uint(15)
	if enum {
		usize = 11
	}

	n, _ := r.ReadNBitsAsUint16BE(uint8(usize))
	nread += usize
	end, nsub := nread+uint(n), int(n) // only one in use according to (bool) enum

	subs := make([]seg, 0, 16)
	for {
		sub := load(r) // recursive loading
		subs = append(subs, sub...)
		nsub--

		if (!enum && end <= nread) || (enum && nsub <= 0) {
			break
		}
	}
	return subs
}

// Sugars for cmd type
const ( // cmd map
	ADD uint8 = iota
	MUL
	MIN
	MAX
	LIT
	GT
	LT
	EQ
)

func load(r *bs.Reader) []seg {
	var segs []seg

	ver, _ := r.ReadNBitsAsUint8(3)
	nread += 3
	typ, _ := r.ReadNBitsAsUint8(3)
	nread += 3

	switch typ {
	case LIT: // load literal
		return append(segs, seg{ver, typ, lit(r), []seg(nil)})
	default: // load args
		return append(segs, seg{ver, typ, 0, subs(r)})
	}
}

func eval(cmd seg) int { // command from segment
	acc, args := 0, cmd.sub // set accumulator and args from sub segments
	switch cmd.typ {
	case ADD:
		for _, a := range args {
			acc += eval(a)
		}
	case MUL:
		acc = 1
		for _, a := range args {
			acc *= eval(a)
		}
	case MIN:
		acc = MaxInt
		for _, a := range args {
			n := eval(a)
			if n < acc {
				acc = n
			}
		}
	case MAX:
		acc = MinInt
		for _, a := range args {
			n := eval(a)
			if n > acc {
				acc = n
			}
		}
	case LIT:
		acc = cmd.val
	case GT:
		if eval(args[0]) > eval(args[1]) {
			acc = 1
		}
	case LT:
		if eval(args[0]) < eval(args[1]) {
			acc = 1
		}
	case EQ:
		if eval(args[0]) == eval(args[1]) {
			acc = 1
		}
	}
	return acc
}

func sum(dgram []seg) int {
	n := 0
	for _, s := range dgram { // segment
		n += int(s.ver)
		n += sum(s.sub)
	}
	return n
}

func main() {
	n, input := new(big.Int), bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n, _ = n.SetString(input.Text(), 16)
	}
	bits := bs.NewReader(bytes.NewReader(n.Bytes()), nil) // go pipeline
	dgram := load(bits)                                   // datagram from bit stream (singleton)

	fmt.Println(sum(dgram))
	fmt.Println(eval(dgram[0]))
}

// MaxInt and MinInt are defined in the idiomatic way
const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -MaxInt - 1
)
