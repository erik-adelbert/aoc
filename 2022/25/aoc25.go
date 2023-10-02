package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type snafu []byte

func (a snafu) String() string {
	var sb strings.Builder
	for i := len(a) - 1; i >= 0; i-- {
		sb.WriteByte(a[i])
	}
	return sb.String()
}

func main() {
	acc := snafu{'0'}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		acc = add(acc, mksnafu(input.Bytes()))
	}

	fmt.Println(acc)
}

func add(a, b snafu) snafu {
	// sort a, b
	if len(b) > len(a) {
		a, b = b, a
	}

	// add digits x and y return carry and digit
	dadd := func(x, y byte) (byte, byte) {
		base := []int{
			'=': -2,
			'-': -1,
			'0': +0,
			'1': +1,
			'2': +2,
		}

		// adder lookup table
		// tab[x][y] = [carry, x+y]
		tab := [][]struct{ c, d byte }{
			'=': {
				'=': {'-', '1'}, // -2 - 2 = -1*5 + 1
				'-': {'-', '2'}, // -2 - 1 = -1*5 + 1
				'0': {'0', '='}, // -2 + 0 = -2
				'1': {'0', '-'}, // -2 + 1 = -1
				'2': {'0', '0'}, // -2 + 2 =  0
			},
			'-': {
				'-': {'0', '='}, // -1 - 1 = -2
				'0': {'0', '-'}, // -1 + 0 = -1
				'1': {'0', '0'}, // -1 + 1 =  0
				'2': {'0', '1'}, // -1 + 2 =  1
			},
			'0': {
				'0': {'0', '0'}, //  0 + 0 =  0
				'1': {'0', '1'}, //  0 + 1 =  1
				'2': {'0', '2'}, //  0 + 2 =  2
			},
			'1': {
				'1': {'0', '2'}, //  1 + 1 =  2
				'2': {'1', '='}, //  1 + 2 =  1*5 - 2
			},
			'2': {
				'2': {'1', '-'}, //  2 + 2 =  1*5 - 1
			},
		}

		// sort digits
		if base[x] > base[y] {
			x, y = y, x
		}

		return tab[x][y].c, tab[x][y].d
	}

	c := byte('0') // carry
	for i := range a {
		var c1, c2 byte
		// add carry to a[i]
		c1, a[i] = dadd(a[i], c)

		c2, a[i] = dadd(a[i], '0')
		if i < len(b) {
			// add b[i]
			c2, a[i] = dadd(a[i], b[i])
		}
		// propagate carry
		_, c = dadd(c1, c2)

		if i >= len(b) && c == '0' {
			// done
			break
		}
	}

	// overflow
	if c != '0' {
		// expand for carry
		a = append(a, c)
	}

	return a
}

func mksnafu(a []byte) snafu {
	// reverse slice
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return snafu(a)
}
