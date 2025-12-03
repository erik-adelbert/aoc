# Timings

| Day  | Time (ms) | % of Total |
|------|----------:|-----------:|
| 1    |       0.8 |      1.75% |
| 3    |       1.2 |      2.62% |
| 2    |      43.6 |     95.63% |
| Total|      45.8 |    100.00% |

fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 24.6.0 - go version go1.25.3 darwin/arm64 - hyperfine 1.20.0 - 2025-12

## Installation and benchmark

0. optionally install [gocyclo](https://github.com/fzipp/gocyclo)
1. install [hyperfine](https://github.com/sharkdp/hyperfine)
2. `git clone` this repository somewhere in your `$GOPATH`
3. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
4. `$ cd 2025`
5. `$ make`
6. `$ make runtime && cat runtime.md`
7. explore the other `Makefile` goals

## Day 1: [Secret Entrance](https://adventofcode.com/2025/day/1)

*This year, I’m freelancing and available to take on projects—preferably in Go or Python. Please help spread the word!*

<div align="center">
  <img src="./images/1606_Mercator_Hondius_Map_of_the_Arctic_(First_Map_of_the_North_Pole)_-_Geographicus_-_NorthPole-mercator-1606.jpg" alt="North Pole Map" width="70%" />
</div>

<!-- ![Secret Entrance](./images/1606_Mercator_Hondius_Map_of_the_Arctic_\(First_Map_of_the_North_Pole\)_-_Geographicus_-_NorthPole-mercator-1606.jpg) -->

On this first day of AoC 2025, the challenge is reasonably tricky. It highlights the sign ambiguity of the [modulo](https://en.wikipedia.org/wiki/Modulo) operation when the remainder is negative.

For today’s solution, I’m reimplementing `mod` so that it always returns a positive value, since the problem includes negative integer data (i.e., left turns). Then, as always — especially when coding for production — I validate the inputs as early as possible. In this case, it allows me to reduce the computation domain to a single wrap of the dial. By doing this consistently, I don’t need to apply any offsets (and neither do you). From there, a switch selects one of the four interesting cases and updates the counts used as passwords for parts 1 and 2.

`<EDIT>` I've removed `mod()` because it was called only once.

`<EDIT>` I’ve used  `if` [short statements](https://go.dev/tour/flowcontrol/6) fairly liberally as a stylistic choice.

The code runs with an overall (optimal) [time complexity](https://en.wikipedia.org/wiki/Time_complexity) of `O(n)`, where *n* is the number of moves. What’s interesting here is that ~~I don’t believe it’s possible to accidentally end up with a solution that has a higher complexity~~ it doesn't depend on the distance value of the moves.

`<EDIT>` Actually, naïve solutions might (incorrectly) click through each move — fully simulating the dial — which would increase the total loop count by a factor of the distance value *d* resulting in `o(n * d)` (`d_min` for the best case, `d_avg` for the average case or `d_max` for the worst case). This kind of code would be roughly 50~1000× slower than the showcased solution depending on the input.

<details>
  <summary><strong>SPOILER: Click to reveal</strong></summary>
The password method <span title='CLICK'><code>0x434C49434B</code></span> actually encodes a more sensible name.
</details>

## Day 2: [Gift Shop](https://adventofcode.com/2025/day/2)

<div align="center">
  <img src="./images/Serpiente_alquimica.jpg" alt="Ouroboros" width="70%" />
</div>

On this second day, the code speed conundrum begins: the challenge requires us to convert back and forth between integers and ASCII slices, and to check the allocated memory for certain patterns.

For part 1, the second half of the slice should be a copy of the first.

For part 2, a doubled slice should contain the original slice as a subslice — meaning that the slice is a [rotation of itself](https://en.wikipedia.org/wiki/Ouroboros). This idea is demonstrated in this [study](https://www.geeksforgeeks.org/dsa/a-program-to-check-if-strings-are-rotations-of-each-other/) along with various pattern-searching techniques like [Rabin–Karp](https://en.wikipedia.org/wiki/Rabin–Karp_algorithm) and [KMP](https://en.wikipedia.org/wiki/Knuth–Morris–Pratt_algorithm).

As a matter of fact, the Go standard `bytes` package uses a combination of techniques, including an ultimate fallback to [Rabin–Karp](https://cs.opensource.google/go/go/+/refs/tags/go1.25.4:src/bytes/bytes.go;l=1389).

The search space, although it may not seem like it, is actually quite respectable:

```bash
❯ cd /Users/erika/go/src/github.com/erik-adelbert/aoc/2025/2 && awk -F',' '{for(i=1;i<=NF;i++){split($i,range,"-"); for(j=range[1];j<=range[2];j++){len=length(j); count[len]++}}} END{for(i in count) print i " digits:", count[i] " numbers" | "sort -n"}' input.txt
1 digits: 8 numbers
2 digits: 81 numbers
3 digits: 758 numbers
4 digits: 8041 numbers
5 digits: 66257 numbers
6 digits: 666270 numbers
7 digits: 413789 numbers
8 digits: 539292 numbers
9 digits: 248595 numbers
10 digits: 301477 numbers
```

The code runs with a time complexity of `k.O(n)` on average, with *n* being the number of digits in the inputs and *k* some big and hard to compute (at least for me) constant. I will get back to this calculation if I don't find a faster idea for this challenge.

It is worth noting that the solution hits the sweet spot where running `part2` *only* if `part1` fails (ie., [predictive branching](https://en.wikipedia.org/wiki/Branch_predictor))— versus *always* running both `part1` and `part2`  — actually hurts the overall runtime.

The solution itself is pretty neat, but the performance, as you can see, isn’t quite there. I’ll call it a day for now.

```bash
cloc .
       5 text files.
       5 unique files.
       0 files ignored.

github.com/AlDanial/cloc v 2.06  T=0.01 s (657.9 files/s, 10131.5 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               1             14              6             50
Markdown                         1              0              0              4
Text                             2              0              0              2
make                             1              0              0              1
-------------------------------------------------------------------------------
SUM:                             5             14              6             57
-------------------------------------------------------------------------------
```

## Day 3: [Lobby](https://adventofcode.com/2025/day/3)

<div align="center">
  <img src="./images/Polar_Night_energy.jpg" alt="Polar Night Energy" width="70%" />
</div>

Today's challenge is quite straightforward: the goal is to build the *lexicographically largest string after **k** removals*. I chose a [greedy](https://en.wikipedia.org/wiki/Greedy_algorithm), [stack-based](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) approach to solve it. The solution is simple, and once again it runs optimally in linear time with respect to the length of the input lines: it is easy  to see that every given digit can only be pushed/popped once.

Having an adhoc `seq` type keeps the main intention obvious while separating concerns. The digit-selection logic becomes a mere implementation detail of the solution. The search space is so small that the Go garbage collector has no time to get in the way, even though the code creates two short-lived buffers per input line.

```bash
❯ wc -lc input.txt # how many lines and cars?
     200   20200 input.txt
```
