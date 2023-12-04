# Timings

| day | time |
|-----|-----:|
| 2 | 0.7 |
| 4 | 0.8 |
| 1 | 0.9 |
| 3 | 1.2 |
| total | 3.6 |

fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 23.0.0 - go version go1.21.4 darwin/arm64 - hyperfine 1.18.0 - 2023-12

## Installation and benchmark

0. optionnally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2022`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals

## Day 1

For this 2023 edition first day, I find the problem unusually sophisticated. The crux of today's challenge is to efficiently [find strings](https://en.wikipedia.org/wiki/String-searching_algorithm) that could *overlap*.

And that's the point, we can easily "unoverlap" those!

Actually, if we scan from left to right for the first digit and from right to left for the second, it doesn't matter if patterns overlap as the two scans are separated.

I'm using a [`trie`](https://web.stanford.edu/class/archive/cs/cs166/cs166.1146/lectures/09/Small09.pdf) abstraction to match the numbers but in this case it is overly simplified. Namely, while LR scanning we have to search for a prefix, conversely while RL scanning we are searching for a suffix. Fortunately, Go provides [`strings.HasPrefix`](https://pkg.go.dev/strings#HasPrefix) and [`strings.HasSuffix`](https://pkg.go.dev/strings#HasSuffix) to just do that. It's easy to
sync LR and RL (ie. in a single core loop).

There's one last trick, LR and RL scans use the same trie at the cost of ~2 extraneous failed comparisons on each number.

PS. Another way is to use [`strings.Replacer`](https://pkg.go.dev/strings#Replacer) and iterating replacement of ie. `two` by `2o` (or `eight` by `8t`) until the string is fixed and then matching only for digits. It is also ok but really slower. Anyway, I guess two successive replacements could do the trick for a vast majority of inputs (3 should kill them all) and the code is pretty neat.

## Day 2

Today, it is all about parsing. There's not much to say except that today's challenge is more like a day1 challenge than actual day1 was.
My program uses [`strings`](https://pkg.go.dev/strings) functions, allmost all variable names are 5 letters long and last but not least, allocations are kept on the low side because input memoization is not needed.

## Day3

Challenge is like [Aoc2022/day23](https://github.com/erik-adelbert/aoc/blob/2576e62f51f3bf653bf95084bca1815c534bf6e2/2022/23/aoc23.go), I'm using multiple bit arrays supported by a custom `u192` type.
This solution is fast.

## Day4

Finally, day1 has come! Today's challenge is about typing speed with a few pauses here and there to actually think through the needed ops.
