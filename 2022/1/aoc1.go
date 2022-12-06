package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var m1, m2, m3, sum int

	max3 := func() {
		switch {
		case sum > m1:
			m1, m2, m3 = sum, m1, m2
		case sum > m2:
			m2, m3 = sum, m2
		case sum > m3:
			m3 = sum
		}
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if n, err := strconv.Atoi(input.Text()); err == nil {
			sum += n
		} else {
			max3()
			sum = 0
		}
	}
	max3()

	fmt.Println(m1, m1+m2+m3) // part 1 & 2
}
