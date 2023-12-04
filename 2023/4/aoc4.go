package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	stock := make([]int, 256)
	score, ncard := 0, 0 // part 1 & 2

	input := bufio.NewScanner(os.Stdin)
	var i int
	for i = 0; input.Scan(); i++ {
		input := input.Text()
		row := Split(input[Index(input, ":")+1:], " | ")
		w, deck := Fields(row[0]), Fields(row[1])

		wins := make([]bool, 100)
		for i := range w {
			wins[atoi(w[i])] = true
		}

		nmatch := 0
		for i := range deck {
			if wins[atoi(deck[i])] {
				nmatch++
			}
		}
		score += (1 << nmatch) >> 1

		stock[i] += 1
		for ii := i + 1; ii <= i+nmatch; ii++ {
			stock[ii] += stock[i]
		}
		ncard += stock[i]
	}
	fmt.Println(score, ncard)
}

var Fields, Index, Split = strings.Fields, strings.Index, strings.Split

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
