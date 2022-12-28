package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// array based circular doubly linked list
type list struct {
	seq []int
	bck []int
	fwd []int
	o   int
}

func main() {
	part1 := mklist()

	salt := 811_589_153
	part2 := mklist()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		n := atoi(input.Bytes())
		if n == 0 {
			part1.o = len(part1.seq)
			part2.o = len(part2.seq)
		}

		part1.append(n)
		part2.append(n * salt)
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		part2.shuffle(10)
	}()

	go func() {
		defer wg.Done()
		part1.shuffle(1)
	}()

	wg.Wait()

	fmt.Println(part1.key())
	fmt.Println(part2.key())
}

func (l *list) shuffle(nloop int) {
	seq, bck, fwd := l.seq, l.bck, l.fwd
	cnt := len(seq)

	bck = bck[:cnt]
	fwd = fwd[:cnt]

	// loopback lists
	last := cnt - 1
	bck[0] = last
	fwd[last] = 0

	// seek fwd or bck reinsertion point @(index+offset)
	seek := func(i, off int) int {
		// i is unlinked, offset cnt!
		off %= (cnt - 1)

		switch {
		case off > 0:
			// forward compressed scan

			// halve lookup path
			if odd(off) {
				i = fwd[i]
				off--
			}
			// right compressed scan
			for ; off > 0; off -= 2 {
				i = fwd[fwd[i]]
			}
		case off < 0:
			// backward compressed scan

			// later on when reinserting, backward moving keys
			// are inserted *before* j, that is *after* bck[j]
			//
			// offset now to return bck[j] instead of j
			// see move() below
			off--

			// halve lookup path
			if odd(off) {
				i = bck[i]
				off++
			}
			// left compressed scan
			for ; off < 0; off += 2 {
				i = bck[bck[i]]
			}
		}
		return i
	}

	move := func(i, off int) {
		// to get on point, order matters here *first* unlink i,
		// this way the insert point is to be found between the
		// *remaining* items

		if off == 0 {
			// nothing to do
			return
		}

		// unlink i
		// see https://en.wikipedia.org/wiki/Dancing_Links#Main_ideas
		fwd[bck[i]], bck[fwd[i]] = fwd[i], bck[i]

		// seek the new insertion point in the remaining sequence
		// see seek() above
		j := seek(i, off)
		if i == j {
			// no move, relink i and abort
			fwd[bck[i]], bck[fwd[i]] = i, i
			return
		}

		// always relink i *after* j
		bck[i], fwd[i] = j, fwd[j]
		bck[fwd[j]], fwd[j] = i, i
	}

	// no global loop cycle because
	// ppcm(key, nloop) = salt*nloop >>> nloop = 1|10
	for nloop > 0 {
		for i, x := range seq {
			move(i, x)
		}
		nloop--
	}
}

func (l *list) key() int {
	seq, fwd, cnt := l.seq, l.fwd, len(l.seq)

	o := l.o

	// subkeys relative to n
	k1 := 1000 % cnt
	k2 := 2000 % cnt
	k3 := 3000 % cnt

	// synced dual loops
	// i in indices space
	// n is item count
	key := 0
	unk := 3 // unknown subkeys
	for n, i := 0, o; unk > 0; n, i = n+1, fwd[i] {
		switch n {
		case k1, k2, k3:
			// forge key from found subkey
			key += seq[i]
			unk--
		}
	}

	return key
}

func mklist() *list {
	p := new(list)
	p.seq = make([]int, 0, 5000)
	p.bck = make([]int, 5000)
	p.fwd = make([]int, 5000)

	return p
}

func (l *list) append(n int) {
	i := len(l.seq)
	l.bck[i] = i - 1
	l.fwd[i] = i + 1
	l.seq = append(l.seq, n)
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

// parity test
func odd(n int) bool {
	return n&1 == 1
}
