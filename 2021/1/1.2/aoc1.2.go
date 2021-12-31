package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MaxInt = int(^uint(0) >> 1)

func main() {
	old3, old2, old1 := MaxInt, MaxInt, MaxInt // 3 last depths window

	n, input := 0, bufio.NewScanner(os.Stdin)
	for input.Scan() {
		cur, _ := strconv.Atoi(input.Text())
		if old3 < cur { // increase detected!
			n++
		}
		old3, old2, old1 = old2, old1, cur // shift/update window
	}
	fmt.Println(n)
}
