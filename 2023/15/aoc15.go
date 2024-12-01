package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	boxes := newBoxes()
	sum := 0

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		ops := strings.Split(input.Text(), ",")
		for i, s := range ops {
			sum += hash(ops[i])

			op := mkop(i, s)
			boxes[hash(op.name)].enqueue(op)
		}
	}

	pwr := 0
	for i := range boxes {
		if len(boxes[i].add) > 0 {
			pwr += boxes[i].process()
		}
	}

	fmt.Println(sum, pwr)
}

func hash(s string) (h int) {
	for i := range s {
		h = ((h + int(s[i])) * 17) & 0xff
	}
	return
}

func newBoxes() []*queue {
	boxes := make([]*queue, 256)
	for i := range boxes {
		q := new(queue)

		q.idx = i
		q.del = make(map[string]int, 16) // max == 9
		q.add = make([]ops, 0, 64)       // max == 48

		boxes[i] = q
	}
	return boxes
}

type ops struct {
	name     string
	idx, val int
}

func mkop(i int, s string) ops {
	op := ops{idx: i}
	switch s[len(s)-1] {
	case '-':
		op.name = s[:len(s)-1]
	default:
		args := strings.Split(s, "=")

		op.name, op.val = args[0], atoi(args[1])
	}
	return op
}

type queue struct {
	del map[string]int
	add []ops
	idx int
}

func (q *queue) enqueue(op ops) {
	if op.val == 0 {
		q.del[op.name] = op.idx
		return
	}

	if i, ok := q.del[op.name]; !ok || (ok && op.idx > i) {
		q.add = append(q.add, op)
	}

	return
}

func (q *queue) process() int {
	slots := make([]ops, 0, len(q.add))

	for _, op := range q.add {
		if i, ok := q.del[op.name]; !ok || (ok && op.idx > i) {
			if i := index(slots, op); i >= 0 {
				slots[i].val = op.val
			} else {
				slots = append(slots, op)
			}

		}
	}

	sum := 0
	for i := range slots {
		if slots[i].name != "" {
			sum += (q.idx + 1) * (i + 1) * slots[i].val
		}
	}
	return sum
}

func index(list []ops, o ops) int {
	return slices.IndexFunc(list, func(x ops) bool { return x.name == o.name })
}

// strconv.Atoi simplified core loop
// s is ^\d+$
func atoi(s string) (n int) {
	for i := range s {
		n = 10*n + int(s[i]-'0')
	}
	return
}
