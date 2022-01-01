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

func NewCard() *card {
	c := new(card)
	c.val = make(map[int]coo)
	return c
}

func (ca *card) add(n, r, c int) {
	ca.row[r] += n
	ca.col[c] += n
	ca.val[n] = coo{r, c}
}

func (ca *card) biff(n int) bool {
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

func (c *card) sum() int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += c.row[i]
	}
	return sum
}

func main() {
	draw := make([]int, 0, 128)
	deck := make([]*card, 0, 128)
	cur, row := NewCard(), 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		args := strings.Fields(line)
		switch len(args) {
		case 0: // sep
			if row > 0 {
				deck = append(deck, cur)
				cur, row = NewCard(), 0
			}
		case 5: // cardboard row
			for col, s := range args {
				n, _ := strconv.Atoi(s)
				cur.add(n, row, col)
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

	// stack!! heavy but easy & reliable
	type item struct {
		n int
		c *card
	}

	stack := make([]item, 0, 128)

	push := func(n int, c *card) {
		stack = append(stack, item{n, c})
	}

	pop := func() (int, *card) {
		if len(stack) == 0 {
			return 0, nil
		}
		i := len(stack) - 1

		pop := stack[i]
		stack, stack[i] = stack[:i], item{} // remove
		return pop.n, pop.c
	}

	for _, n := range draw {
		i := 0
		for _, card := range deck {
			if card.biff(n) { // win
				push(n, card) // save
				continue
			}
			deck[i] = card
			i++
		}
		deck = deck[:i]
	}
	n, c := pop() // last to win
	fmt.Println(n * c.sum())
}
