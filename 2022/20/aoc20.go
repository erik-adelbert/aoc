// aoc20.go --
// advent of code 2022 day 20
//
// https://adventofcode.com/2022/day/20
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2022-12-20: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
)

func shuffle(input []int, salt, nround int) ([]int, int) {
	const (
		// 256*32 = 2^13 > 5000 items
		maxbin = 256
		// cluster bin capacity
		nbin = 16
		// bin item capacity
		nitem = 32
	)

	type item struct {
		val int
		// global index
		gix int
	}

	// populate bins from input
	bins := make([][]item, 0, maxbin)
	bin := make([]item, 0, nitem)
	for i, x := range input {
		bin = append(bin, item{x * salt, i})
		if len(bin) >= nitem {
			bins = append(bins, bin)
			bin = make([]item, 0, nitem)
		}
	}
	bins = append(bins, bin)

	// compute cluster sizes
	clusters := make([]int, 1+(len(bins)/nbin))
	for i := range bins {
		bin := i / nbin
		clusters[bin] += len(bins[i])
	}

	// map bin/off addresses to gix
	type addr struct {
		bin, off int
	}

	addrs := make([]addr, 0, maxbin*nitem)
	for i := range bins {
		for j := range bins[i] {
			addrs = append(addrs, addr{i, j})
		}
	}

	// shuffle
	for n := nround; n > 0; n-- {
		// shuffle once
		for k, a := range addrs {
			// extract x from cur a
			x := bins[a.bin][a.off]

			// move left
			for i := a.off; i < len(bins[a.bin])-1; i++ {
				bins[a.bin][i] = bins[a.bin][i+1]
				addrs[bins[a.bin][i].gix].off = i
			}
			// shrink
			bins[a.bin] = bins[a.bin][:len(bins[a.bin])-1]

			// update cluster size
			cid := a.bin / nbin
			clusters[cid]--

			// compute cur x global index
			gix := a.off
			for _, n := range clusters[:cid] {
				gix += n
			}
			for i := cid * nbin; i < a.bin; i++ {
				gix += len(bins[i])
			}

			// compute new x global index
			gix = mod(gix+x.val, len(addrs)-1)

			// remap to bin/off
			bin, off := 0, 0
			for off+clusters[bin/nbin] <= gix {
				off += clusters[bin/nbin]
				bin += nbin
			}
			for off+len(bins[bin]) <= gix {
				off += len(bins[bin])
				bin++
			}
			off = gix - off

			// reinsert x

			// update cluster size
			cid = bin / nbin
			clusters[cid]++

			// expand
			bins[bin] = append(bins[bin], item{})
			// move right
			for i := len(bins[bin]) - 1; i > off; i-- {
				bins[bin][i] = bins[bin][i-1]
				addrs[bins[bin][i].gix].off = i
			}
			// insert back
			bins[bin][off] = x

			// write back new a
			addrs[k] = addr{bin, off}
		}

		// rebalance between rounds
		flat := make([]item, 0, len(bins)*nitem)
		for i := range bins {
			flat = append(flat, bins[i]...)
		}

		bins = make([][]item, 0, maxbin)
		bin = make([]item, 0, nitem)
		for _, x := range flat {
			addrs[x.gix] = addr{
				bin: len(bins), off: len(bin),
			}
			bin = append(bin, x)
			if len(bin) >= nitem {
				bins = append(bins, bin)
				bin = make([]item, 0, nitem)
			}
		}
		bins = append(bins, bin)

		for i := range clusters {
			clusters[i] = 0
		}
		for i := range bins {
			bin := i / nbin
			clusters[bin] += len(bins[i])
		}
	}

	// flatten and extract origin
	shuffled, o := make([]int, 0, len(bins)*nitem), -1
	for i := range bins {
		for _, x := range bins[i] {
			if x.val == 0 {
				o = len(shuffled)
			}
			shuffled = append(shuffled, x.val)
		}
	}

	return shuffled, o
}

func key(a []int, o int) int {
	// get key components
	k1 := mod(o+1000, len(a))
	k2 := mod(o+2000, len(a))
	k3 := mod(o+3000, len(a))

	// forge & return key
	return a[k1] + a[k2] + a[k3]
}

func main() {
	seq := make([]int, 0, 8192)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		seq = append(seq, atoi(input.Bytes()))
	}

	// part 1
	fmt.Println(key(shuffle(seq, 1, 1)))

	// part 2
	const salt = 811_589_153
	fmt.Println(key(shuffle(seq, salt, 10)))
}

// strconv.Atoi simplified core loop
// s is ^-?\d+$
func atoi(b []byte) int {
	n, s := 0, 1
	if b[0] == '-' {
		s = -1
		b = b[1:]
	}
	for _, c := range b {
		n = 10*n + int(c-'0')
	}
	return s * n
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

const DEBUG = true

func debug(a ...any) {
	if DEBUG {
		fmt.Println(a...)
	}
}
