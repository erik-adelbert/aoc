# Summary

This repository contains optimized solutions for Advent of Code 2025, implemented in [Go](https://en.wikipedia.org/wiki/Go_(programming_language)) with a focus on performance and educational value. The whole collection runs all days and all parts in about 4.1 ms. Each day is solved in a single, dependency-free [Go](https://go.dev) file.

## Quick Navigation

### Daily Visual Index

<table>
  <tr>
    <td><a href="#day-1-secret-entrance-"><img src="./images/1606_Mercator_Hondius_Map_of_the_Arctic_(First_Map_of_the_North_Pole)_-_Geographicus_-_NorthPole-mercator-1606.jpg" alt="Day 1" width="80"/></a></td>
    <td><a href="#day-2-gift-shop-"><img src="./images/Serpiente_alquimica.jpg" alt="Day 2" width="80"/></a></td>
    <td><a href="#day-3-lobby-"><img src="./images/Polar_Night_energy.jpg" alt="Day 3" width="80"/></a></td>
    <td><a href="#day-4-printing-department-"><img src="./images/fire.jpg" alt="Day 4" width="80"/></a></td>
    <td><a href="#day-5-cafeteria-"><img src="./images/SpaceVegetables.jpg" alt="Day 5" width="80"/></a></td>
    <td><a href="#day-6-trash-compactor-"><img src="./images/BoxFactory.jpg" alt="Day 6" width="80"/></a></td>
    <td><a href="#day-7-laboratories-"><img src="./images/PrismRoom.jpg" alt="Day 7" width="80"/></a></td>
    <td><a href="#day-8-playground-"><img src="./images/Xmas_Snowball.jpg" alt="Day 8" width="80"/></a></td>
    <td><a href="#day-9-movie-theater-"><img src="./images/industrial_compressor.jpg" alt="Day 9" width="80"/></a></td>
    <td><a href="#day-10-factory-"><img src="./images/hp48.png" alt="Day 10" width="80"/></a></td>
    <td><a href="#day-11-reactor-"><img src="./images/jetcar.jpg" alt="Day 10" width="80"/></a></td>
     <td><a href="#day-12-christmas-tree-farm-"><img src="./images/kfpshifu.jpg" alt="Day 12" width="80"/></a></td>
  </tr>
  <tr>
    <td align="center">Day 1</td>
    <td align="center">Day 2</td>
    <td align="center">Day 3</td>
    <td align="center">Day 4</td>
    <td align="center">Day 5</td>
    <td align="center">Day 6</td>
    <td align="center">Day 7</td>
    <td align="center">Day 8</td>
    <td align="center">Day 9</td>
    <td align="center">Day 10</td>
    <td align="center">Day 11</td>
    <td align="center">Day 12</td>
  </tr>
</table>

- [Timings](#timings-) - Performance metrics and hardware specs
- [Installation and Benchmark](#installation-and-benchmark-) - Setup and testing instructions
- [Day 1: Secret Entrance](#day-1-secret-entrance-) - Modulo operations and dial simulation
- [Day 2: Gift Shop](#day-2-gift-shop-) - Repunit numbers and pattern matching optimization
- [Day 3: Lobby](#day-3-lobby-) - Greedy stack-based string manipulation
- [Day 4: Printing Department](#day-4-printing-department-) - Cellular automata and memory optimization
- [A 5mn crash-introduction to cache and GC friendly solutions](#a-5mn-crash-introduction-to-cache-and-gc-friendly-solutions-) - Slices and memory allocation
- [Day 5: Cafeteria](#day-5-cafeteria-) - Range merging
- [Day 6: Trash Compactor](#day-6-trash-compactor-) - Matrix operations and data organization
- [Day 7: Laboratories](#day-7-laboratories-) - Path propagation and dynamic programming
- [Day 8: Playground](#day-8-playground-) - Modified Kruskal's with distance cutoff
- [Why have I changed the timings?](#why-have-i-changed-the-timings-) - Timings and evaluation
- [Day 9: Movie Theater](#day-9-movie-theater-) - Prefix sums, 2D compress coordinates and cache optimization
- [Day 10: Factory](#day-10-factory-) - BFS and ILP solver
- [Day 11: Reactor](#day-11-reactor-) - Graph DFS and DP, Hashing
- [Day 12: Christmas Tree Farm](#day-12-christmas-tree-farm-) - Parsing and Heuristic
- [How was it?](#how-was-it-)
- [More Metrics](#more-metrics-)

## Timings [↑](#summary)

<div align="center">
  <img src="./images/barchart.png" alt="The barchart of the table below" width="100%" />
</div>

| Day                                | Time (μs) | % of Total  |
| ---------------------------------- | --------: | ----------: |
| [**2**](#day-2-gift-shop-)         |         8 |       0.20% |
| [7](#day-7-laboratories-)          |        30 |       0.75% |
| [5](#day-5-cafeteria-)             |        95 |       2.37% |
| [12](#day-12-christmas-tree-farm-) |       119 |       2.97% |
| [1](#day-1-secret-entrance-)       |       129 |       3.22% |
| [6](#day-6-trash-compactor-)       |       150 |       3.74% |
| [11](#day-11-reactor-)             |       157 |       3.92% |
| [3](#day-3-lobby-)                 |       199 |       4.96% |
| [10](#day-10-factory-)             |       248 |       6.19% |
| [4](#day-4-printing-department-)   |       695 |      17.34% |
| [9](#day-9-movie-theater-)         |     1,037 |      25.87% |
| [8](#day-8-playground-)            |     1,142 |      28.49% |
| **Total**                          | **4,009** | **100.00%** |

fastest of 100 runs for part1&2 in μs - mbair M1/16GB - darwin 24.6.0 - go1.25.3 darwin/arm64 with greentea GC - 2025-12

## Installation and benchmark [↑](#summary)

0. optionally install [gocyclo](https://github.com/fzipp/gocyclo)
1. `git clone` this repository somewhere in your `$GOPATH`
2. `export` envar `$SESSION` with your AoC `session` value (get it from the cookie stored in your browser)
3. `$ cd 2025`
4. `$ make`
5. `$ make run`
6. explore the other `Makefile` goals

## Day 1: [Secret Entrance](https://adventofcode.com/2025/day/1) [↑](#summary)

*This year, I’m freelancing and available to take on projects—preferably in Go or Python. Please help spread the word!*

<div align="center">
  <img src="./images/1606_Mercator_Hondius_Map_of_the_Arctic_(First_Map_of_the_North_Pole)_-_Geographicus_-_NorthPole-mercator-1606.jpg" alt="1606 Mercator Hondius Map of the Arctic (First Map of the North Pole). First day of AoC is always the start of a great journey to the North Pole! We have a map now." width="60%" />
</div>

### Current approach

It is possible to [solve](https://github.com/erik-adelbert/aoc/blob/main/2025/1/aoc1.go) this problem using pure [modular arithmetic](https://en.wikipedia.org/wiki/Modular_arithmetic). Instead of tracking the starting and ending position of the dial we can directly compute the 0-landing click and see if it is reachable. As in the previous approach, we only have to reason on the last single wrap of a move. The challenge 2nd part reduces to solving, with *p* the dial position, *i* the click landing on 0 and *s* the wrap direction as ±1 sign (left is -1):

```math
p + i \cdot s \equiv 0 \pmod{100}
```

If we rearrange:

```math
i \cdot s \equiv −p \pmod{100}
```

As *s* and *M* are coprime, the modular inverse of *s* is *s* and we can multiply both sides by it:

```math
i \equiv −p \cdot s \pmod{100}
```

It makes the solution much simpler. The first zero-landing click is then:

```math
i₀ = (−p \cdot s) \pmod{100}
```

If `i₀ == 0`, the position was already at zero and should be ignored. Otherwise, if
`i₀ ≤ r`, where *r* is the number of remaining steps after completing the full wraps, the zero crossing occurs during the final partial move. The resulting logic is minimal and branchless. It is not only (slightly) faster but also imho much more elegant than before.

### First approach

**For the following discussion please checkout commit [089a922](https://github.com/erik-adelbert/aoc/blob/089a92232b7dd79cb39a4d7d881edf8b9ecce200/2025/1/aoc1.go)**

On this first day of AoC 2025, the challenge is reasonably tricky. It highlights the sign ambiguity of the [modulo](https://en.wikipedia.org/wiki/Modulo) operation when the remainder is negative.

For today’s [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/1/aoc1.go), I’m reimplementing `mod` so that it always returns a positive value, since the problem includes negative integer data (i.e., left turns). Then I validate the inputs and in this case, it allows me to reduce the computation domain to a single wrap of the dial. By doing this consistently, I don’t need to apply any offsets (and neither do you). From there, a switch selects one of the four interesting cases and updates the counts used as passwords for parts 1 and 2.

`<EDIT>` I removed `mod()` because it was called only once.

`<EDIT>` I’ve used  `if` [short statements](https://go.dev/tour/flowcontrol/6) fairly liberally as a stylistic choice.

The code runs with an overall (optimal) [time complexity](https://en.wikipedia.org/wiki/Time_complexity) of `O(n)`, where *n* is the number of moves. What’s interesting here is that ~~I don’t believe it’s possible to accidentally end up with a solution that has a higher complexity~~ it doesn't depend on the distance value of the moves.

`<EDIT>` Actually, naïve solutions might (incorrectly) click through each move — fully simulating the dial — which would increase the total loop count by a factor of the distance value *d* resulting in `o(n * d)` (`d_min` for the best case, `d_avg` for the average case or `d_max` for the worst case). This kind of code would be roughly 50~1000× slower than the showcased solution depending on the input.

<details>
  <summary><strong>SPOILER: Click to reveal</strong></summary>
The password method <span title='CLICK'><code>0x434C49434B</code></span> actually encodes a more sensible name.
</details>

## Day 2: [Gift Shop](https://adventofcode.com/2025/day/2) [↑](#summary)

<div align="center">
  <img src="./images/Serpiente_alquimica.jpg" alt="Ouroboros is a mystical symbol—don't miss its link below. To me, it celebrates the magic that surrounds us. One of my good friends always says that we, as programmers, are closer than ever to magicians. Some of us write code that, if seen from the past, would be indistinguishable from spells.
" width="60%" />
</div>

### Third Approach

As AoC is a gathering, I usually keep a back channel open with my fellow programmer and friend **[hm](https://blog.izissise.net/)**. From the very beginning, he had been insisting on how fast the generation of the numbers we are tasked to find in today’s challenge could be. He was convinced from the start that, given their regular nature, they were natural candidates for efficient generation… and it turns out he was right.

I must admit I wasn’t convinced at first, but once we saw [Tim Visée’s solution](https://github.com/timvisee/advent-of-code-2025/blob/4292ba336b685270886129ee58311943cef5ce63/day02b/src/main.rs)—which does exactly the opposite by sieving the repeating numbers statically—the challenge was on.

And here it is: possibly **the fastest way** to compute the solution to today’s challenge. It runs in under **8 μs**.

This [code](https://github.com/erik-adelbert/aoc/blob/main/2025/2/aoc2.go) runs in `O(k)` time, with *k* being the number of ranges. **It shrinks the original 2M+ search space down to only ~1,800 relevant numbers.** It operates in constant memory and evaluates the repeating-number sums using direct arithmetic formulas.

The flow starts by segmenting the input ranges into sub-ranges aligned on `[10, 1e2, ..., 1e9]` so that the appropriate generating seed values are naturally selected (see the second approach below). Thanks to the structure of the input, this results in only a single split once in a while. Put simply, once the ranges are aligned on successive powers of ten (from 1 to 9 digits), all repeating numbers become multiples of 1 or 2 [Repunit](https://en.wikipedia.org/wiki/Repunit) divisors per range.

From there, given an interval `[a, b]` and its corresponding generating seed `s0`, we simply compute the sum of the multiples of `s0` that fall within the range. The code uses an almost closed-form solution (i.e., an arithmetic formula) to achieve this, along with efficient techniques to handle well-known issues such as eliminating duplicate numbers when different seeds share common multiples within the same range, or merging seeds that become redundant.

### Second Approach

**For the following discussion please checkout commit  [1fd714e](https://github.com/erik-adelbert/aoc/blob/1fd714e7f1a3d37736e4e87a35544bd33a2c852a/2025/2/aoc2.go)**

I have come across a better idea than my original approach and expanded on it. The insight is to exploit the properties of [Repunit](https://en.wikipedia.org/wiki/Repunit) numbers: within each numeric range, the repeating-digit numbers (the ones we need to detect in today’s challenge) are simply multiples of a small set of seed values.

For example, between 10 and 99, it’s easy to see that all repeating numbers are multiples of 11.

The resulting [code](https://github.com/erik-adelbert/aoc/blob/1fd714e7f1a3d37736e4e87a35544bd33a2c852a/2025/2/aoc2.go) stays within the integer domain, the cost effectively disappears — and the routine now runs in **5.8 ms**!

I first saw this idea in [Tim Visée](https://github.com/timvisee/advent-of-code-2025/blob/4292ba336b685270886129ee58311943cef5ce63/day02b/src/main.rs)’s code.
Tim is a performance-oriented programmer of the finest caliber, and I warmly recommend following his work.

As a final note, this solution uses `fallthrough`, which helps improve runtime.
However, it’s not actually necessary: since part 2 includes part 1, all the
`fallthrough` statements can be replaced with a single `part2 += part1` right before producing the final output.

### First Approach

**For the following discussion please checkout commit [a89bc57](https://github.com/erik-adelbert/aoc/blob/a89bc57abece8df39e0ea2acbf5d6d4a9eae6924/2025/2/aoc2.go)**

On this second day, the code speed conundrum begins: the challenge requires us to convert back and forth between integers and ASCII slices, and to check the allocated memory for certain patterns.

For part 1, the second half of the slice should be a copy of the first.

For part 2, a doubled slice should contain the original slice as a subslice — meaning that the slice is a [rotation of itself](https://en.wikipedia.org/wiki/Ouroboros). This idea is demonstrated in this [study](https://www.geeksforgeeks.org/dsa/a-program-to-check-if-strings-are-rotations-of-each-other/) along with various pattern-searching techniques like [Rabin–Karp](https://en.wikipedia.org/wiki/Rabin–Karp_algorithm) and [KMP](https://en.wikipedia.org/wiki/Knuth–Morris–Pratt_algorithm).

As a matter of fact, the Go standard `bytes` package uses a combination of techniques, including an ultimate fallback to [Rabin–Karp](https://cs.opensource.google/go/go/+/refs/tags/go1.25.4:src/bytes/bytes.go;l=1389).

The [code](https://github.com/erik-adelbert/aoc/blob/a89bc57abece8df39e0ea2acbf5d6d4a9eae6924/2025/2/aoc2.go) runs with a time complexity of `k.O(n)` on average, with *n* being the number of digits in the inputs and *k* some big and hard to compute (at least for me) constant. I will get back to this calculation if I don't find a faster idea for this challenge.

It is worth noting that the solution hits the sweet spot where running `part2` *only* if `part1` fails (ie., [predictive branching](https://en.wikipedia.org/wiki/Branch_predictor))— versus *always* running both `part1` and `part2`  — actually hurts the overall runtime.

The solution itself is pretty neat, but the performance, as you can see, isn’t quite there. I’ll call it a day for now.

<`EDIT>` Actually, the performance is interesting to analyze: given my input, there are 2,244,568 candidates (as shown in the `awk` command above), of which 816 are invalid for part 1 and 895 for part 2. This results in a blazing-fast 43.6 ms / 2,244,568 numbers ≈ 19.4 ns per number for parts 1 and 2 combined. This result feels arguably good.

The search space, although it may not seem like it, is actually quite respectable:

```bash
❯ awk -F',' '{for(i=1;i<=NF;i++){split($i,range,"-"); for(j=range[1];j<=range[2];j++){len=length(j); count[len]++}}} END{for(i in count) print i " digits:", count[i] " numbers" | "sort -n"}' input.txt
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

`<EDIT>` In Go, [strings](https://go.dev/blog/strings) are immutable, which means many operations on them require allocations. This is why I prefer [byte slices](https://go.dev/blog/slices-intro) in the solution: they allow me to tightly control memory usage and eliminate all allocations from the hot path.

```bash
❯ go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: github.com/erik-adelbert/aoc/2025/2
cpu: Apple M1
BenchmarkItoa-8         180716007                6.548 ns/op           0 B/op          0 allocs/op
PASS
```

`<EDIT>` This challenge is also interesting because substring matching efficiently is inherently complex. This is one of the rare cases where the standard library’s implementation has a good chance of being the best tool for the job — despite the unserious performance-critical context. **This last idea of using the standard library is almost absolute if you are a beginner**.

The beauty of [`u/topaz2078`](https://www.reddit.com/user/topaz2078/)’s craftsmanship is that, in this very solution, you’ll see me *simultaneously* relying on Go for the heavy lifting *and* deliberately avoiding it for the ASCII translation. I have the room to exercise my judgment to tilt the solution toward the fast side. For that, I am forever in awe.

## Day 3: [Lobby](https://adventofcode.com/2025/day/3) [↑](#summary)

<div align="center">
  <img src="./images/Polar_Night_energy.jpg" alt="A Polar Night Energy monolith that looks like one giant power bank built on the side of a car park. It's straightforward: today I felt like building the largest battery possible out of a stock of batteries. It made me want to see a monolith like in the old days." width="60%" />
</div>

Today's challenge is quite straightforward: the goal is to build the *lexicographically largest string after **k** removals*. I chose a [greedy](https://en.wikipedia.org/wiki/Greedy_algorithm), [stack-based](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) approach to solve it. The [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/3/aoc3.go) is simple, and once again it runs optimally in linear time with respect to the length of the input lines: it is easy  to see that every given digit can only be pushed/popped once.

Having an adhoc `seq` type keeps the main intention obvious while [separating concerns](https://en.wikipedia.org/wiki/Separation_of_concerns). The digit-selection logic becomes a mere implementation detail of the solution. ~~The search space is so small that the Go garbage collector has no time to get in the way, even though the code creates two short-lived small buffers per input line~~.

`<EDIT>` As I wanted to emphasize the `O(1)` space complexity alongside the `O(n)` time complexity of the solution—and to nullify the [Go garbage collector](https://go.dev/doc/gc-guide) pressure altogether—the code now reuses the *same* storage repeatedly.

```bash
❯ wc -lc input.txt # how many lines and cars?
     200   20200 input.txt
```

## Day 4: [Printing Department](https://adventofcode.com/2025/day/4) [↑](#summary)

<div align="center">
  <img src="./images/fire.jpg" alt="People juggling with fire bolas. Bolas is what I felt I was throwing at the problem with my double-buffer idea, which I managed to maintain no matter what." width="60%" />
</div>

### Current Approach

After seeing various solutions online, I decided to move the double buffering off the grid in favor of a [double-buffered](https://wiki.osdev.org/Double_Buffering) queue, where updates are stored and removals are bulk-applied between steps. Moreover, the queue-based approach processes only the cells that might have changed (the neighbors of removed rolls) rather than scanning the entire grid on each iteration. This transforms the algorithm from a grid-scanning problem into a change-propagation problem, with a time complexity linear in the number of total removals.

This also means that the general theory now relates to [cellular automata](https://en.wikipedia.org/wiki/Cellular_automaton).

The [implementation](https://github.com/erik-adelbert/aoc/blob/main/2025/4/aoc4.go) also eliminates the memory overhead of maintaining two grids by collecting removal positions first, then applying them atomically to prevent corruption during neighbor counting. Combined with preallocated queues and direct array indexing instead of hash maps for deduplication, this optimization achieves a **72% performance improvement** over the original double-buffered approach, bringing the runtime down from 6.5ms to 1.8ms.

### First Approach

**For the following discussion please checkout commit [06ede07](https://github.com/erik-adelbert/aoc/blob/06ede07387eb9f7ca4c23409e15c569aa844321f/2025/4/aoc4.go)**

This challenge is the perfect opportunity to go fully old-school with the solution. It’s an AoC [classic](https://adventofcode.com/2021/day/20) that pops up regularly. It has nothing to do with mathematics and everything to do with programming efficiently for our machines when [processing images](https://en.wikipedia.org/wiki/Digital_image_processing). **If you're a beginner, you could benefit from working through this problem and studying its [various solutions](https://www.reddit.com/r/adventofcode/comments/1pdr8x6/2025_day_4_solutions/).**

My [technique](https://github.com/erik-adelbert/aoc/blob/06ede07387eb9f7ca4c23409e15c569aa844321f/2025/4/aoc4.go) of choice here is to [double-buffer](https://wiki.osdev.org/Double_Buffering) the grid. By doing this, the code kills the removal process with a single [double-stone](https://en.wikipedia.org/wiki/Bolas): it becomes natural to go from one step of the roll removals to the next by updating the *next* buffer from the *current* one and then swapping them.

For the 2D grid itself, nothing beats a [1D grid](https://en.wikipedia.org/wiki/Array_(data_structure)). The code uses two preallocated slices and spatially organizes data on the fly. Except for the initial allocations, the solution once again performs **no memory allocation** on the [hot path](https://en.wikipedia.org/wiki/Hot_spot_(computer_programming)).

I didn't add a blank border to the grid because it would interfere with the index computations—and actually, I don't need to. The showcased code features a *branchless* neighborhood scan that is slightly incorrect because it includes the center roll itself. But this turns out to be beneficial: since we only scan *from* the rolls, it is easy to remove the center cell test in favor of thresholding at 4 rolls (3 neighbors + 1 center) during the entire scan.

The time complexity of *one scan* is `O(n)`: it is easy to see that each cell is processed only once per scan. The total runtime depends on the input, its roll count, and the relative positions. For my input, wich contains 64% of rolls, it takes 70 loops to reduce the grid.

<div align="center">
  <img src="./images/aoc20251204.png" alt="The final grid displayed in my terminal (scaled-down)" width="40%" />
</div>

## A 5mn crash-introduction to cache and GC friendly solutions [↑](#summary)

<div align="center">
  <img src="./images/StarWars.jpg" alt="Chewbacca carries a broken C-3PO on his back to save him later—a moment captured on the movie set. This reflects the idea that if you bring broken or inefficient code into good care, it can eventually perform!" width="60%" />
</div>

I solved the Day 4 challenge without losing sight of the [Go GC model](https://go.dev/doc/gc-guide) and, more broadly, how [memory is managed](https://en.wikipedia.org/wiki/Virtual_memory) in our computers (or at least the much simpler real-life version of it—bear with me). I approached it this way because I strongly believe that [mechanical sympathy](https://newsletter.appliedgo.net/archive/2025-11-30-mechanical-sympathy/) improves program efficiency without requiring any energy beyond the effort of thought.

I’m not going to elaborate on what mechanical sympathy is or what it might mean for us to possess it. In everyday life, it’s much simpler than it sounds. Suppose—purely for the sake of demonstration—you need the best possible performance when thinning a cellular automaton through repeated application of the same rule. Everything works: your logic is flawless, and the result is correct. Naturally, Go slices are extremely useful here and well-suited to the task.

Now, let’s talk about two of the seven benchmarks that you can [find](https://github.com/erik-adelbert/aoc/blob/main/2025/4/aoc4_test.go) alongside the Day 4 solution.

First, let’s look at this one:

```Go
// BenchmarkPreallocatedWithCopy shows the efficient approach for copying
func BenchmarkPreallocatedWithCopy(b *testing.B) {
    src := slices.Repeat([]byte{1}, 1024)
    buf := make([]byte, 1024) // GOOD: allocate once

    for b.Loop() {
        // GOOD: just copy to pre-allocated buffer
        copy(buf, src)

        // Prevent optimization
        sink = buf
    }
}
```

The goal is to create a fresh working copy of a source slice at each iteration of a loop. In this example, everything is fine: the buffers exist at the same scope level, and aside from resetting the contents of `buf` each iteration, nothing ever changes. We never need to modify their size, nor do we need to worry about how memory management might behave, because we’re using them consistently. Right?

But the thing is, in real life we often fixate on small details and lose sight of the bigger picture—and that’s when patterns like this can appear:

```Go
// BenchmarkCarelessAllocations demonstrates the performance impact of allocating
// in tight loops - this is what NOT to do in performance-critical code
func BenchmarkCarelessAllocations(b *testing.B) {
    for b.Loop() {
        // BAD: allocating inside the loop
        buf := make([]byte, 1024)

        // Fill buffer with 1s
        for i := range 1024 {
            buf[i] = 1
        }

        // Prevent optimization by assigning to global
        sink = buf
    }
}
```

The point is the same as before and surely the result is correct. But this time, the code applies maximum pressure by claiming short lived memory at a very high pace in its core loop. It is said to create *friction* with the Go runtime.

What the hell does this mean? Actually, It means this:

```bash
❯ go test -bench="BenchmarkCarelessAllocations$|BenchmarkPreallocatedWithCopy$" -benchmem
goos: darwin
goarch: arm64
pkg: github.com/erik-adelbert/aoc/2025/4
cpu: Apple M1
BenchmarkCarelessAllocations-8           2109789               489.6 ns/op          1024 B/op          1 allocs/op
BenchmarkPreallocatedWithCopy-8         84923436                13.77 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/erik-adelbert/aoc/2025/4     2.434s
```

Never having to think about the right time and place to declare a buffer can lead to a ×30 slowdown—wasting at least some amount of computing power for no real benefit in return (ie., sub-optimal efficiency).

If you’re interested in reviewing your own solutions for allocation mishaps, you may find the other [five benchmarks](https://github.com/erik-adelbert/aoc/blob/main/2025/4/aoc4_test.go) useful. They illustrate a variety of good and bad patterns you may have used without realizing it, along with an accompanying analysis summarizing the keypoints.

```bash
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -coverprofile=/var/folders/9y/jfl_qkbs6_9_8xht9qchxhzw0000gn/T/vscode-gobiPhfa/go-code-cover -bench . github.com/erik-adelbert/aoc/2025/4

goos: darwin
goarch: arm64
pkg: github.com/erik-adelbert/aoc/2025/4
cpu: Apple M1
BenchmarkCarelessAllocations-8           	 2826324	       463.9 ns/op	    1024 B/op	       1 allocs/op
BenchmarkPreallocated-8                  	 3608708	       333.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkCarelessAllocationsWithCopy-8   	11459871	       100.8 ns/op	    1024 B/op	       1 allocs/op
BenchmarkPreallocatedWithCopy-8          	88135389	        13.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkWorstCase-8                     	 2897648	       415.0 ns/op	    2560 B/op	       3 allocs/op
BenchmarkRealWorldBad-8                  	 1684596	       711.1 ns/op	    2144 B/op	      11 allocs/op
BenchmarkRealWorldGood-8                 	 2646415	       454.0 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: 0.0% of statements
ok  	github.com/erik-adelbert/aoc/2025/4	8.685s
```

## Day 5: [Cafeteria](https://adventofcode.com/2025/day/5) [↑](#summary)

<div align="center">
  <img src="./images/SpaceVegetables.jpg" alt="Vegetables growing in a space greenhouse aboard an orbital station. Merging food storage units of today's challenge spells SPACE to me." width="60%" />
</div>

I don’t have much to say about today’s challenge. In anticipation of part 2, I used an [interval tree](https://en.wikipedia.org/wiki/Interval_tree). But part 2 ultimately required merging the input ranges and computing the total coverage.

Between the tree querying and the coverage, the time complexity is dominated by `O(m log n)` where *m* is the query count and *n* is the interval count. The storage complexity is, of course, `O(n)`.

The [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/5/aoc5.go) runs in under 1 ms on my inputs, which is perfectly fine. Let’s call it a win!

`<EDIT>` The code now populates the interval tree while tallying coverage from the *merged* intervals. This wasn’t necessary—the speedup is marginal—but it feels more *correct*, and it only required moving a couple of lines around. In the coming days, I’ll remove the tree entirely, since the merged ranges are overlap-free making say a basic binary search perfectly fit for the job.

`<EDIT>` The code now merges the intervals and performs query by bissecting the resulting merged set. It is way lighter now with no visible improvement in runtime (but it was actually ×2 at the μ-level).

```bash
❯ make run
go run ./aoc5.go < input.txt
862 357907198933892
```

### How is it going?

After putting a lot of effort into day 2, I’m quite happy with the total time budget for the first five days: **5.4ms**.

```bash
cloc 1 2 3 4 5
      26 text files.
      23 unique files.
       3 files ignored.

github.com/AlDanial/cloc v 2.06  T=0.03 s (821.8 files/s, 260665.6 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Text                            10              2              0           6339
Go                               6            153            126            653
Markdown                         5              0              0             20
make                             2              0              0              2
-------------------------------------------------------------------------------
SUM:                            23            155            126           7014
-------------------------------------------------------------------------------
```

## Day 6: [Trash Compactor](https://adventofcode.com/2025/day/6) [↑](#summary)

<div align="center">
  <img src="./images/BoxFactory.jpg" alt="A factory’s cardboard conveyor line moving packages. In the top left corner a logo says Process Idea. Yeah! I felt like packaging numbers today, lots of numbers!" width="60%" />
</div>

The challenge presents a problem that’s a perfect opportunity to practice working with Go [slices](https://go.dev/tour/moretypes/7) and understanding how they relate to the [memory management](https://go.dev/doc/gc-guide) provided by the Go runtime.

The [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/6/aoc6.go) is very straightforward and mainly involves retrieving and organizing data from the input considered as a byte matrix. The key insight is to extract the matrix layout from the last line: since the operators are left-aligned within their columns, it’s much easier to determine each column’s fixed width from that line than from any other, avoiding altogether the “what kind of space is this space?” conundrum.

I also went a step further and [transposed](https://en.wikipedia.org/wiki/Transpose) the matrix so the numbers are grouped by column. Transposing the column submatrices (i.e., inducing machines to read from top to bottom) is also essential for part two.

The program isn’t the prettiest, but it gets the job done in 85 lines. I believe the code runs in `O(n)` time, where *n* is the number of digits in the matrix. It executes in under 1 ms.

## Day 7: [Laboratories](https://adventofcode.com/2025/day/7) [↑](#summary)

<div align="center">
  <img src="./images/PrismRoom.jpg" alt="A room with overwhelming, prism-like architecture. If you don't step back from today's challenge to view it fully, your mind could feel like crushed by this design." width="60%" />
</div>

Today's challenge presents a path propagation problem that I solved using [Dynamic Programming](https://en.wikipedia.org/wiki/Dynamic_programming) principles. The algorithm tracks how paths split and multiply as they traverse the grid from top to bottom.

The [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/7/aoc7.go) adheres to the narrative: It *simulates* paths starting from position 'S' and splitting at each '^' character encountered. When a path hits a '^', it disappears and creates two new paths at adjacent positions (left and right). The code strictly does that and then Part 1 counts the total number of splits that occur, while Part 2 sums all active paths remaining at the end.

The algorithm runs with `O(n)` time complexity, where *n* is the number of grid cells. Each row is processed exactly once, and for each row, we iterate through all possible path positions. The space complexity is `O(w)` for the paths array were *w* is the grid width, making it quite memory-efficient.

It runs in under 1ms.

## Day 8: [Playground](https://adventofcode.com/2025/day/8) [↑](#summary)

<div align="center">
  <img src="./images/Xmas_Snowball.jpg" alt="A Christmas snow globe, featuring a decorated Christmas tree inside, rests on the snow outside on a magical, snowy night. It is a mise-en-abîme of today's challenge: surround and then capture!" width="60%" />
</div>

The [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/8/aoc8.go) implements a variant of [Kruskal’s algorithm](https://en.wikipedia.org/wiki/Kruskal%27s_algorithm) for computing a [Minimum Spanning Tree](https://en.wikipedia.org/wiki/Minimum_spanning_tree), but with problem-specific optimizations.

A key observation is that **any edges *after* the one required for Part 2 never affect either answer**.
In other words, the solution has a **distance cutoff**: once we know the maximum edge weight that could possibly matter, every edge longer than that is irrelevant.

By determining this cutoff early, we can **[prune](https://en.wikipedia.org/wiki/Decision_tree_pruning) the edge set from ~500k to ~5k**, dramatically reducing the work.

This greatly improves runtime because Kruskal’s algorithm—along with the ~~heap~~ sorting and the disjoint-set union (DSU)—runs in time proportional to `O(E log E)`, and reducing `E` by two orders of magnitude makes the whole process significantly faster.

`<EDIT>` I replaced the heap by a sort to make the code lighter and actually faster.

The code runs in under `1.3ms`

## Why have I changed the timings? [↑](#summary)

<div align="center">
  <img src="./images/magic-mountains.jpg" alt="A giant roller coaster with crowded cars whizzing past. My point with Hyperfine was to measure whether code performance suffered from the sheer cost of naively copying the input to stdout. Now, the timings are much more comparable with some strong Rust solutions out there, and this illustrates my point: Go is suitable for more than just building web services or infrastructure." width="60%" />
</div>

During AoC I’ve increasingly been comparing my solutions with others written in Rust, and many AoC Rust crates include internal program timers that report raw compute times. On the other hand I have many solutions that are simply too fast for `hyperfine`. Because of this, starting now I will publish **internal timings** instead of external (wall-clock) timings. These internal timings are much more comparable to what Rust and other fast languages report.

For now, my collection of programs solves every day and every part in about **2.7 ms total**.

```bash
❯ make run
go run ./aoc1.go < input.txt
1168 7199 139.792µs
go run ./aoc2.go < input.txt
28146997880 40028128307 11.125µs
go run ./aoc3.go < input.txt
17324 171846613143331 232.791µs
go run ./aoc4.go < input.txt
1474 8910 768.167µs
go run ./aoc5.go < input.txt
862 357907198933892 103.959µs
go run ./aoc6.go < input.txt
5595593539811 10153315705125 160.625µs
go run ./aoc7.go < input.txt
1633 34339203133559 40.292µs
go run ./aoc8.go < input.txt
32103 8133642976 1.257ms
```

## Day 9: [Movie Theater](https://adventofcode.com/2025/day/9) [↑](#summary)

<div align="center">
  <img src="./images/industrial_compressor.jpg" alt="An industrial compressor in full, intricate detail. It is what it is, but this solution surely produces its fair share of gas." width="60%" />
</div>

Today, we’re tackling the problem of finding the largest rectangle inside a rectilinear polyline. Suffice it to say, I’m not a big fan of this type of problem—these puzzles usually end up with a [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/9/aoc9.go) that’s tedious, and today is no exception.

As I see it, there’s no option other than generating all pairs of polyline points. This results in a program that runs in `O(n²)`, where *n* is the number of points—about 500² (250,000) pairs. The good news (yippee!) is that this step dominates everything else we do with the data. The bad news is that we have to keep the core loop as lightweight as possible to avoid a disastrous runtime that would undermine everything we’ve done up to this point. After browsing [various posts](https://www.reddit.com/r/adventofcode/comments/1phywvn/2025_day_9_solutions/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button) today, I've settled to a prefix sum table over compressed coordinates of the polyline, mainly because I am unfamiliar with these techniques.

Namely, today the code must:

- Parse the input `O(n)` where *n* is the point count
- Sort and compact the coordinates `O(nlogn)`
- Build a compress coordinate map `O(m)` where *m* is the number of unique coordinates
- Scanline fill the polyline `O(R·C)` where *R·C* is the number of grid cells
- Build a prefix sum table `O(R·C)`
- Query all pairs of point against this table `O(n²)` (☜ core loop here)

The [prefix sum table](https://www.geeksforgeeks.org/dsa/prefix-sum-array-implementation-applications-competitive-programming/) is built so that each cell at position (i, j) contains the sum of all cells in the rectangle from (0, 0) to (i, j). This is done using the recurrence:

`sums[i][j]=grid[i][j]+sums[i−1][j]+sums[i][j−1]−sums[i−1][j−1]`

It lets us quickly check whether a rectangle is fully inside the polygon (i.e., all its cells are 1), or count how many cells are inside, without scanning every cell in the rectangle. And after doing so much preprocessing, we might as well throw everything we have at the problem —the best 2D matrix is a 1D matrix, ... kind of things— wherever possible, and so on.

At last the resulting code runs in under `1.5ms±0.5` average and the best runs are in under `1ms`.

```bash
❯ best=999999999; for i in {1..100}; do t=$( (make run 2>&1 | grep -oE '[0-9]+\.[0-9]+ms' | head -1 | sed 's/ms//') ); if [ "$(echo "$t < $best" | bc)" -eq 1 ]; then best=$t; fi; done; echo "Best time: $best ms"
Best time: 1.037333 ms
❯ make run
go run ./aoc9.go < input.txt
Parsed input... 27µs
Sorted coordinates... 96.792µs
Built coordinate maps... 158.292µs
Completed scanline fill... 186µs
Built bitmask... 417.333µs
Built prefix sums... 1.200042ms
4294831030 1644094530 1.714583ms
```

### How is it going?

Even if—or maybe because—it’s a 12-day edition, this year’s AoC feels great and keeps motivating me to improve my tooling. The results I’m getting are much easier to compare with other performance-oriented solutions out there, and it’s a ton of fun.

I’ve always seen coding as a way to balance time and space, and to take advantage of that balance. That’s why I love the recurring appearance of intervals throughout this journey—they keep reminding me of *space* every now and then, while every day for a full hour I’m trading everything for *speed*.

This recurrence is also more proof of [`u/topaz2078`](https://www.reddit.com/user/topaz2078/)’s grand wizardry: it’s his way of having us practice *kata (型)*. *Kata* refers to fixed forms or patterns that embody a principle. By practicing many different kata, we gradually internalize the underlying concepts. And then there’s *Shu–Ha–Ri (守破離)*, which I think is the perfect way to describe the stages of mastery:

- *Shu (守)* — obey the forms (learn the kata)
- *Ha (破)* — break the forms (experiment, adapt)
- *Ri (離)* — transcend the forms (create your own)

Today I reached the **1000-LOC** milestone, the collection runs in **3,683 μs**, and we’ve clearly moved past the self-patting comfort zone of `O(n)` into far more adventurous territory. Let the final stretch begin!

<div align="right">
  <img src="./images/golang.png" alt="Flash Gopher from https://wx-chevalier.github.io/books/awesome-lists/01.cs/language/go/gopher-list/" width="10%" />
</div>

```bash
❯ cloc 1 2 3 4 5 6 7 8 9
      46 text files.
      39 unique files.
       7 files ignored.

github.com/AlDanial/cloc v 2.06  T=0.04 s (891.9 files/s, 220230.5 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Text                            18              2              0           8030
Go                              11            307            242           1015
Markdown                         8              0              0             32
make                             2              0              0              2
-------------------------------------------------------------------------------
SUM:                            39            309            242           9079
-------------------------------------------------------------------------------
```

## Day 10: [Factory](https://adventofcode.com/2025/day/10) [↑](#summary)

<div align="center">
  <img src="./images/hp48.png" alt="A close-up of an HP48 calculator screen: it says CHIP8 in a frame everything is line-drawned. Today, it was back-to-school time! Here it is, my fellow pocket computer from that era." width="60%" />
</div>

### Current Approach

You will find the rationale for what is done in the current [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/10/aoc10.go) in this [post](https://www.reddit.com/r/adventofcode/comments/1plzhps/2025_day_10_part_2_pivot_your_way_to_victory/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button), and in particular this [answer](https://www.reddit.com/r/adventofcode/comments/1plzhps/comment/ntx2ggc/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button). The discussion goes back and forth between two of the strongest performance-oriented programmers. [`u/maneatingape`](https://www.reddit.com/user/maneatingape/), for instance, [runs](https://www.reddit.com/r/adventofcode/comments/1hlyocd/500_in_less_than_a_second/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button) all days and all parts ever in under 1 second. The code features a parallelized variation of the aforementioned HNF-based solution to Part 2. It runs in under 248μs with part 1 BFS being the slowest.

### First Approach

**For the following discussion please checkout commit [f04f7e8](https://github.com/erik-adelbert/aoc/blob/f04f7e8e4e39fd28936e1b9a076d1b6d54d40b81/2025/10/aoc10.go)**

Upon reading today's challenge, Part 2, I immediately recognized it as an [ILP](https://en.wikipedia.org/wiki/Integer_programming) problem, but I had no idea how I could compose a solver simply. Relying on [Z3](https://github.com/Z3Prover/z3) or any other general-purpose solver/prover would have been a total disaster for runtime.

So I studied the domain and came up with an approach for what to do and how to do it. Then I browsed solutions hoping to find a skillful reference implementation—and here it was: `u/RussellDash332` had [one](https://github.com/RussellDash332/advent-of-code/blob/main/aoc-2025%2FDay-10%2FPython%2Fmain.py)!

Today's [program](https://github.com/erik-adelbert/aoc/blob/f04f7e8e4e39fd28936e1b9a076d1b6d54d40b81/2025/10/aoc10.go) tackles Part 1 with a neat BFS and includes a fast, compact solver using [simplex](https://en.wikipedia.org/wiki/Simplex_algorithm) and [branch-and-bound](https://en.wikipedia.org/wiki/Branch_and_bound) for Part 2. It runs in under **5 ms**, which is an order of magnitude faster than typical runtimes in the solution megathread.

```bash
❯ make run
go run ./aoc10.go < input.txt
498 17133 5.632792ms
❯ best=999999999
for i in {1..100}; do
  t=$(make run 2>&1 | grep -v part1 | grep -oE '[0-9]+\.[0-9]+ms' | head -1 | sed 's/ms//')
  if [ -n "$t" ] && [ "$(echo "$t < $best" | bc)" -eq 1 ]; then
    best=$t
  fi
done
echo "Best time: $best ms"
Best time: 4.892125 ms
```

## Day 11: [Reactor](https://adventofcode.com/2025/day/10) [↑](#summary)

<div align="center">
  <img src="./images/jetcar.jpg" alt="An old jet car racing by in black and white. Today, I sped through the problem because I was more interested in Day 10. Still, it was fun to be given a reactor by the challenge." width="60%" />
</div>

Today I needed to move quickly through composing the [solution](https://github.com/erik-adelbert/aoc/blob/main/2025/11/aoc11.go), so this reactor comes at a perfect time (TL;DR Yesterday morning I wasn’t able to finish a tiny—and hopefully fast—ILP solver, and I wanted to deliver it today).

To this end, I used [recursive](https://en.wikipedia.org/wiki/Recursion) [DFS](https://en.wikipedia.org/wiki/Depth-first_search) and [DP](https://en.wikipedia.org/wiki/Dynamic_programming) to solve parts 1 and 2, because they are fast to compose and the part 1 result clearly shows that part 2 would otherwise require an intractable search space. I also translated all three-letter tags into fixed indices so I could perform all searches in the integer domain. Believe it or not, there’s still room for additional micro-optimizations like iterate computations, and at some point I’ll come back to implement them.

`<EDIT>` part 1 is iterative now but I think I'll leave part 2 alone because it won't stay as neat otherwise.

`<EDIT>` I’ve upped the ante by speeding up the maps. For the three-letter tags, I eliminated the map entirely thanks to the small hash space. For the DP state, I converted it to a uint32 using an ad-hoc hashing routine. It runs in under 270µs.

```bash
❯ make bench
go build aoc11.go
../bench.sh aoc11
Best time: 280.583 µs
go clean
rm -f aoc11
```

## Day 12: [Christmas Tree Farm](https://adventofcode.com/2025/day/12) [↑](#summary)

<div align="center">
  <img src="./images/kfpshifu.jpg" alt="Kung Fu Panda’s Master Shifu sitting in a lotus pose under a tree. Today, I did not panic because it was not my first AoC edition, but I still laugh thinking about the Piano Guys' faces. What a finish!" width="60%" />
</div>

Oh my gosh! What a roller coaster today. When I first read the challenge, I was feeling overwhelmed and needed a break—so I sat down with a nice breakfast and tried to estimate the computational difficulty on the back of an envelope.

For example, if I wanted to choreograph the solution using Algorithm X or even an MILP, I’d need a matrix with one column per constraint, which for all the polyominoes would mean basically **4 rotations × ~45 occurrences × 6 pieces ≈ 1000 columns**. And we had something like a thousand such grids to solve!

So I already knew that solving it that way was off the table. But that was actually good news: the whole point was to see through the fog of the narrative. It’s the last day—AoC surely doesn’t want to knock us all out right at the end. We’re meant to celebrate, after all.

I needed something simple—something I could search and test on the sample, but that would also scale cleanly to the full input. **I won’t spoil what my solution is today**; I can’t even guarantee that it will work for you, but I can explain why it should. And here’s the good news:

```bash
❯ make bench
go build aoc12.go
../bench.sh aoc12
Best time: 119.375 µs
go clean
rm -f aoc12
```

My collection runs all problems for the entire edition in roughly `4.2ms`.

## How was it? [↑](#summary)

<div align="center">
  <img src="./images/skidoo.jpg" alt="An acrobatic skidoo jump. That's It!" width="60%" />
</div>

It’s fast, but it still has all the ups and downs (yes, Day 10) of a classic past edition nonetheless. I’m quite happy with the result this year, and keeping my coding time limited — and actually coding at an early hour — has been really good for me. I guess this year, the experience from past years kicked in, and 1–2 hours per day was enough to focus on delivering a strong version the first time (except for Day 2).

Day 2 and Day 10 are the highlights of this year’s collection. To finally solve Day 2 the fastest and definitive way possible, I had to keep an open mind while discussing various solutions for the problem and then recognize and seize a good opportunity, even if it proved me wrong the first time. Day 10 was more of a hunch and a desperate attempt to keep the overall runtime low, and I’m really happy it turned out favorably.

I chose the pictures as a way to convey my feelings about the subject; it’s a stream-of-thought–style attempt.

Last, but not least:

<div align="center">
  <img src="./images/leaderboard.png" alt="1st place on my leaderboard" width="80%" />
</div>

All in all, it's been fun.

Happy coding!

PS. I’m trying to compete in the `Red(dit) One` part of AoC, namely for [Day 1](https://www.reddit.com/r/adventofcode/comments/1pb3y8p/comment/ntzpvto/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button). Stay tuned.

PS2. I have been [awarded](https://www.reddit.com/r/adventofcode/comments/1prrvfn/introducing_your_2025_reddit_one_winners_and/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button) an **Arch-Elf** title 🥈 for this work. Many thanks to this awesome community!

<div align="center">
  <img src="./images/Elfo.jpg" alt="A picture of elfo from Disenchantment" width="20%" />
</div>

## More Metrics [↑](#summary)

<div align="center">
  <img src="./images/lumon.png" alt="A Severance Lumon OS screen displaying numbers" width="60%" />
</div>

### Inputs

| Day  | Lines | Chars | Type    |
|:-----|------:|------:|---------|
| 1    |    5K |   20K | numbers |
| 2    |     1 |   500 | numbers |
| 3    |   200 |   20K | strings |
| 4    |   150 |   20K | grid    |
| 5    |    1K |   20K | numbers |
| 6    |    10 |   20K | grid    |
| 7    |   200 |   20K | grid    |
| 8    |    1K |   20K | numbers |
| 9    |   500 |    6K | numbers |
| 10   |   200 |   20K | tokens  |
| 11   |    1K |   10K | strings |
| 12   |    1K |   20K | hybrid  |

### LOC

| Day | Go Lines |
|:----|---------:|
| 7   |       39 |
| 3   |       64 |
| 12  |       66 |
| 5   |       77 |
| 6   |       88 |
| 11  |       94 |
| 1   |      119 |
| 8   |      142 |
| 9   |      121 |
| 2   |      255 |
| 4   |      186 |
| 10  |      555 |

```bash
❯ cloc 1 2 3 4 5 6 7 8 9 10 11 12
      66 text files.
      56 unique files.
      13 files ignored.

github.com/AlDanial/cloc v 2.06  T=0.05 s (1134.5 files/s, 265457.7 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Text                            26             14              0           9903
Go                              16            564            416           1806
Markdown                        10             49              0            148
Python                           2             29             35            139
make                             2              0              0              2
-------------------------------------------------------------------------------
SUM:                            56            656            451          11998
-------------------------------------------------------------------------------
```

### Cyclomatic complexity over 10

```bash
33 main min2D 10/aoc10.go:316:1
20 main main 4/aoc4.go:21:1
19 main main 6/aoc6.go:23:1
18 main main 11/aoc11.go:21:1
17 main main 9/aoc9.go:13:1
16 main fmbounds3D 10/aoc10.go:647:1
14 main main 8/aoc8.go:25:1
13 main hnf 10/aoc10.go:511:1
12 main main 2/aoc2.go:24:1
11 main fmbounds2D 10/aoc10.go:608:1
Average: 5.29
```

Day 10 is a major contributor in this table, and implementing [BLAS](https://en.wikipedia.org/wiki/Basic_Linear_Algebra_Subprograms) algorithms is certainly quite sophisticated. For instance, they are usually as huge (LOC-wise) as they are fast; this means that the code bails out at the first opportunity, leading to a mechanical increase in cyclomatic complexity. It’s no wonder the [simplex algorithm](https://en.wikipedia.org/wiki/Simplex_algorithm) of my first approach was rated 41. Implementing one is not easy and introduces challenges ranging from structural design to numerical instability. It’s a sophisticated technique developed in the mid-20th century that reached peak prominence in the late 1990s. With recent progress in applying [linear programming](https://en.wikipedia.org/wiki/Linear_programming) to AI, it’s seeing a revival in certain hybrid AI [planning](https://arxiv.org/abs/2509.21014) pipelines.

Complexity in [10-20] always ends up in main(). It’s proof that I focused on speed this year, but hopefully the solutions are still easy to read. I am quite happy with an average of ~5.

[Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) is debatable, but as I see it, it is a convenient way to classify the kind of [testing](https://en.wikipedia.org/wiki/Software_testing) needed to prove [correctness](https://en.wikipedia.org/wiki/Correctness_(computer_science)).

| Cyclomatic Complexity | Associated Risk                        |
|-----------------------|----------------------------------------|
| 1–5                   | Low (simple, easy to test)             |
| 6–10                  | Moderate (manageable, some risk)       |
| 11–20                 | High (complex, error-prone)            |
| 21+                   | Very High (difficult to maintain/test) |

<br>
<br>
<br>

<div align="center">
  <img src="./images/calendar.png" alt="AoC 2025 ASCII Art" width="60%" />
</div>

<br>
<br>
<br>

<div align="center">
  <img src="./images/1938-beck-steinway-grand-piano-patent-print-blueprint-greg-edwards.jpg" alt="A Steinway grand piano blueprint, you know!" width="60%" />
</div>
