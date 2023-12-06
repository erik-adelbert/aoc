# Timings

| day | time |
|-----|-----:|
| 6 | 0.6 |
| 2 | 0.7 |
| 5 | 0.7 |
| 4 | 0.8 |
| 1 | 0.9 |
| 3 | 1.1 |
| total | 4.8 |

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

Challenge is related to [Aoc2022/day23](https://adventofcode.com/2022/day/23), I'm using multiple bit arrays supported by a custom `u192` type. This solution is amazingly fast mainly because it is cache and CPU friendly. Almost all 2D ops are separated in 1D vector ops (think [`SIMD`](https://en.wikipedia.org/wiki/Single_instruction,_multiple_data)) and occur in a rolling window of 3 input lines.

Anyway, the challenge is akin to a [*static analysis*](https://en.wikipedia.org/wiki/Static_program_analysis) of a [*multi-valued game of life*](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) that's why it shares some of the techniques I used [last year](https://github.com/erik-adelbert/aoc/blob/2576e62f51f3bf653bf95084bca1815c534bf6e2/2022/23/aoc23.go).

I prize E. Wastl continuous effort on delivering such neat subjects every year. Here the challenge story is about a complex machinery with many cogwheels precisely timed and sized for the task: this is a fair description of what today's program feels like, a complex yet efficient machine with many simple parts that intricately but gracefully fall in place.

PS. To the young coders that might read this: don't be afraid! It's *not* a common day3 solution and certainly not the easiest way to solve it (but one of the fastest). The thing is, last november when warming up, I happened to refactor/improve AoC22/23.
So the bitpacking technique is still vivid in my memory. If it wasn't I may not have succeded in conjuring, factorizing and finally getting right all the corner cases and details of this solution in a fair amount of time.

## Day4

Finally, day1 has come! Today's challenge is about typing speed with a few pauses here and there to actually think through the needed ops. As standard Go package `strings` has already proven usefull to tokenize inputs, I'm once again using it here.

The solution is totally linear, that is it follows closely the challenge tale and its runtime complexity is bounded by `O(n)` with `n` the input (deck) size. Given the small size of today's input (~200 lines), it is very fast.

PS. Isn't this awesome that at the heart of today's score calculation lies this beautiful gem:

```C
score += 1 << nmatch >> 1
```

`<EDIT>` following [`u/masklinn`](https://www.reddit.com/r/adventofcode/comments/18actmy/comment/kbzqx3e/?utm_source=share&utm_medium=web2x&context=3) advice, I went the extra mile consisting of replacing the winning number map by a [bitmap](https://en.wikipedia.org/wiki/Bitmap). I've also replaced the static 200+ deck by a [ring buffer](https://en.wikipedia.org/wiki/Circular_buffer). The resulting improvement is not measurable with hyperfine though.

## Day5

[Intervals!](https://en.wikipedia.org/wiki/Interval_(mathematics)) 

Given that it is only day5 (and the input size), I'm *not* going to talk about [`Interval Trees`](https://en.wikipedia.org/wiki/Interval_tree). I am going to [`brute-force`](https://en.wikipedia.org/wiki/Brute-force_search) the thing!
`<Spoiler>` Well brute-forcing it, doesn't mean testing billions of points! It rather means to brute-force the interval *boundaries* checking and by the way to turn `part1` into a peculiar `part2` problem in order to use the *same* code for the two. `</Spoiler>`.

PS. `Interval Trees` can also be found in the [big book](https://en.wikipedia.org/wiki/Introduction_to_Algorithms) 3rd ed. from pp. 348-353.

PS2. Look how fast the solution is \o/


## Day6

Today is a direct application of solving this [`quadratic formula`](https://en.wikipedia.org/wiki/Quadratic_formula):

```C
    (x - t)*t - d = 0
```

For a very long time [`FPU`](https://en.wikipedia.org/wiki/Floating-point_unit) was slow but at the turn of y2k, `OS` and users alike were putting so much pressure on `CPU` that actually, `FPU` pipeline was usually free (and faster than before anyway) making it usable for a variety of computing (the lore of General Purpose FPU was born) that were previously carried on by `CPU`. I remember the astonishment around me when one day I decided to benchmark the `FPU` against the `CPU` and showed that it won hands down in almost all situations.

But, for the sake of remembering those old days, I still don't want to switch to `FPU` when computing a square-root in an otherwise integer problem. I usually use a (fast) [`integer square root`](https://en.wikipedia.org/wiki/Integer_square_root) computation. In this very case, there's no reasonnable way to see the difference.
