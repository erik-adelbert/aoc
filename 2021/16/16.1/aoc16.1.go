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

var nread uint

func chunk(r *bs.Reader) (uint8, bool) {
	n, _ := r.ReadNBitsAsUint8(5)
	nread += 5
	return n & 0xf, (n & 0x10) > 0
}

func lit(r *bs.Reader) int {
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

	usize := uint(15)
	if count {
		usize = 11
	}

	n, _ := r.ReadNBitsAsUint16BE(uint8(usize))
	nread += usize
	end, nsub := nread+uint(n), int(n)

	subs := make([]seg, 0, 16)
	for {
		sub := load(r)
		subs = append(subs, sub...)
		nsub--

		if (!count && end <= nread) || (count && nsub <= 0) {
			break
		}
	}

	return subs
}

func load(r *bs.Reader) []seg {
	var segs []seg

	ver, _ := r.ReadNBitsAsUint8(3)
	nread += 3
	typ, _ := r.ReadNBitsAsUint8(3)
	nread += 3

	switch typ {
	case 4:
		val := lit(r)
		return append(segs, seg{ver, typ, val, []seg(nil)})
	default:
		subs := op(r)
		return append(segs, seg{ver, typ, 0, subs})
	}
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
	n := new(big.Int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n, _ = n.SetString(input.Text(), 16)
	}
	bits := bs.NewReader(bytes.NewReader(n.Bytes()), nil) // go pipeline
	dgram := load(bits)                                   // datagram from bitstream

	fmt.Println(sum(dgram))
}
