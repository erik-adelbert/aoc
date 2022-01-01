package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"

	bs "github.com/bearmini/bitstream-go"
)

type seg struct {
	ver uint8
	typ uint8
	val int
	sub []seg
}

var nread uint // global bit count

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

func op(r *bs.Reader) []seg {
	count, _ := r.ReadBool() // len id
	nread++

	usize := uint8(15)
	if count {
		usize = 11
	}

	n, _ := r.ReadNBitsAsUint16BE(usize)
	nread += uint(usize)
	last, nsub := nread+uint(n), int(n) // only one in use according to count (bool)

	subs := make([]seg, 0, 16)
	for {
		sub := load(r) // recurse
		subs = append(subs, sub...)
		nsub--

		if (!count && last <= nread) || (count && nsub <= 0) {
			break
		}
	}
	return subs
}

const ( // cmd map
	ADD = iota
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
	case LIT: // data from segment
		val := lit(r)
		return append(segs, seg{ver, typ, val, nil})
	default:
		subs := op(r)
		return append(segs, seg{ver, typ, 0, subs})
	}
}

const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -int(^uint(0)>>1) - 1
)

func eval(cmd seg) int { // command from segment
	acc, args := 0, cmd.sub // args from segments
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

func main() {
	n, input := new(big.Int), bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n, _ = n.SetString(input.Text(), 16)
	}
	bits := bs.NewReader(bytes.NewReader(n.Bytes()), nil) // go pipeline
	data := load(bits)[0]                                 // datagram from bit stream (singleton)

	fmt.Println(eval(data))
}
