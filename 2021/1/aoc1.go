package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// MaxInt is defined in the idiomatic way
const MaxInt = int(^uint(0) >> 1)

func main() {
	old1, old2, old3 := MaxInt, MaxInt, MaxInt // 3 last depths window

	p1, p2 := 0, 0
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cur, _ := strconv.Atoi(input.Text())
		if old1 < cur { // increase!
			p1++
		}
		if old3 < cur { // increase!
			p2++
		}
		old3, old2, old1 = old2, old1, cur // shift/update window
	}
	fmt.Println(p1)
	fmt.Println(p2)
}
