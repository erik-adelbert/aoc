package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type coo struct {
	r, c int
}

type card struct {
	row [5]int
	col [5]int
	val map[int]coo
}

func newCard() *card {
	c := new(card)
	c.val = make(map[int]coo)
	return c
}

func (ca *card) add(r, c, n int) {
	ca.row[r] += n
	ca.col[c] += n
	ca.val[n] = coo{r, c}
}

func (ca *card) biff(n int) bool { // biff a number, return true if it's a win
	if rc, ok := ca.val[n]; ok {
		r, c := rc.r, rc.c
		ca.row[r] -= n
		ca.col[c] -= n

		if ca.row[r] == 0 || ca.col[c] == 0 { // win
			return true
		}
	}
	return false
}

func (c *card) sum() int { // sum remaining numbers (kinda)
	sum := 0
	for i := 0; i < 5; i++ {
		sum += c.row[i]
	}
	return sum
}

func main() {
	draw := make([]int, 0, 128)
	deck := make([]*card, 0, 128)
	cur, row := newCard(), 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args := strings.Fields(line)
		switch len(args) {
		case 0: // sep
			if row > 0 {
				deck = append(deck, cur)
				cur, row = newCard(), 0
			}
		case 5: // cardboard row
			for col, s := range args {
				n, _ := strconv.Atoi(s)
				cur.add(row, col, n)
			}
			row++
		default: // first line
			args = strings.Split(args[0], ",")
			for _, s := range args {
				n, _ := strconv.Atoi(s)
				draw = append(draw, n)
			}
		}
	}
	if row > 0 { // no newline at eof
		deck = append(deck, cur)
	}

	for _, n := range draw {
		for _, card := range deck {
			if card.biff(n) { // first win
				fmt.Println(n * card.sum())
				return
			}
		}
	}
}
