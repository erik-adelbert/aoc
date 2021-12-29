package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"

	bs "github.com/bearmini/bitstream-go"
)

type segment struct {
	ver uint8
	typ uint8
	val int
	sub []segment
}

var nread int

func chunk(r *bs.Reader) (uint8, bool) {
	n, _ := r.ReadNBitsAsUint8(5)
	return n & 0xf, (n & 0x10) > 0
}

func lit(r *bs.Reader) int {
	val := 0
	for {
		n, more := chunk(r)
		val = (val << 4) | int(n)
		nread += 5
		if !more {
			return val
		}
	}
}

func op(r *bs.Reader) []segment {
	count, _ := r.ReadBool() // len id
	nread++

	usize := uint8(15)
	if count {
		usize = 11
	}

	n, _ := r.ReadNBitsAsUint16BE(usize)
	nread += 15
	last, nsub := nread+int(n), int(n)

	subs := make([]segment, 0, 16)
	for {
		sub := load(r)
		subs = append(subs, sub...)
		nsub--

		if (!count && last <= nread) || (count && nsub <= 0) {
			break
		}
	}

	return subs
}

func load(r *bs.Reader) []segment {
	var segs []segment

	ver, _ := r.ReadNBitsAsUint8(3)
	nread += 3
	typ, _ := r.ReadNBitsAsUint8(3)
	nread += 3

	switch typ {
	case 4:
		val := lit(r)
		return append(segs, segment{ver, typ, val, nil})
	default:
		subs := op(r)
		return append(segs, segment{ver, typ, 0, subs})
	}
}

func sum(datagram []segment) int {
	n := 0
	for _, s := range datagram { // segment
		n += int(s.ver)
		n += sum(s.sub)
	}
	return n
}

func main() {
	n := new(big.Int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n, _ = n.SetString(input.Text(), 16)
	}
	bits := bs.NewReader(bytes.NewReader(n.Bytes()), nil) // go pipeline
	dgram := load(bits)                                   // datagram from bitstream

	fmt.Println(sum(dgram))
}
