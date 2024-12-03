# Timings

| day | time |
|-----|-----:|
| 1 | 0.8 |
| 2 | 0.9 |
| 3 | 1.4 |
| total | 3.1 |


fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 23.6.0 - go version go1.23.3 darwin/arm64 - hyperfine 1.18.0 - 2024-12

## Installation and benchmark

0. optionnally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2024`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals

## Day 1

*This year, I’m freelancing and available to take on projects—preferably in Go or Python. Please help spread the word!*

On this first day of AoC 2024, the challenge seems reasonable. For today’s solution, I’m using `sort` instead of `slices` because the problem invites presorted integer data. This choice allows the code to perform [binary searches](https://en.wikipedia.org/wiki/Binary_search) on the right dataset using the left dataset. When factoring in the presorting, the overall [runtime complexity](https://en.wikipedia.org/wiki/Time_complexity) is [O(n log n)](https://go.dev/src/sort/sort.go).

## Day 2

My solution to today’s problem is fairly straightforward. It involves a left-to-right scan to ensure the safety constraints, evolving into a [tail-recursive call](https://en.wikipedia.org/wiki/Tail_call) to tolerate exactly one misplaced element. The only tricky part is when the first item is the misplaced one. In that case, we can simply check whether the report starting from the second element is safe. This approach applies anytime there’s a misplaced element, significantly simplifying the flow of control.

It’s interesting to see that the vast majority of other coders opted to generate all possible reports and brute-force the solution. In contrast, my approach requires at most two generated reports to validate or invalidate any given report and only one of them needs extra memory allocation and data copy.

## Day 3

Go should be renowned for its blazing-fast [regular expression](https://en.wikipedia.org/wiki/Regular_expression) matching [engine](https://swtch.com/~rsc/regexp/). I thoroughly enjoyed today’s problem as it reminded me of [Exercice 12](https://clc-wiki.net/wiki/K%26R2_solutions:Chapter_1:Exercise_12) from *[The C Programming Language](https://en.wikipedia.org/wiki/The_C_Programming_Language)*. Combining the two concepts was a delight—and a satisfying way to wrap up the day.

Here's a fun fact collected while analysing my input:

```sh
❯ grep -oE '\w+\([0-9]+,[0-9]+\)' input.txt | sed -E 's/\(([0-9]+,[0-9]+)\)//' | sort | uniq | tr '\n' ' '
354who 953why from how mul select what when where who why
❯ grep -oE '[a-z]+' input.txt | sed -E 's/\(([0-9]+,[0-9]+)\)//' | sort | uniq | tr '\n' ' '
bin do don from how mul mulselect mulwhen mulwho mulwhy perl select t usr what when where who why
❯ grep -oE '.{0,10}perl.{0,10}' input.txt
!/usr/bin/perl@~mul?what
```

Today's problem seems to be a tribute to Perl by [u/topaz](https://x.com/ericwastl/status/1465082878073753600?lang=en), the creator of AoC.

Note: I may revisit this solution at some point.
