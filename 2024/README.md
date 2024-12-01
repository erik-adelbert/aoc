# Timings

| day | time |
|-----|-----:|
| 1 | 0.8 |
| total | 0.8 |

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

On this first day of AoC 2024, the challenge seems reasonable. For today’s solution, I’m using `sort` instead of `slices` because the problem invites presorted integer data. This choice allows the code to perform [binary searches](https://en.wikipedia.org/wiki/Binary_search) on the right dataset using the left dataset. When factoring in the presorting, the overall [runtime complexity](https://en.wikipedia.org/wiki/Time_complexity) is [O(n log n)](https://go.dev/src/sort/sort.go). Another optimization comes from the data itself:

```bash
❯ zsh stats.sh
Duplicates in column 1: 0
Duplicates in column 2: 20
```

This reveals that the left column is [monotonic](https://en.wikipedia.org/wiki/Monotonic_function). As a result, instead of searching the entire right list, the code narrows the search range by slicing out previously searched parts:`popcnt()` is implemented as a [closure](https://en.wikipedia.org/wiki/Closure_(computer_programming)) that maintains a (monotonic) growing base index.
