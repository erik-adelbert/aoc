# Timings

| day | time |
|-----|-----:|
| 5 | 0.5 |
| 1 | 0.6 |
| 6 | 0.6 |
| 18 | 0.7 |
| 2 | 0.7 |
| 11 | 0.8 |
| 4 | 0.8 |
| 9 | 0.8 |
| 10 | 0.9 |
| 13 | 0.9 |
| 3 | 1.0 |
| 15 | 1.1 |
| 7 | 1.1 |
| 8 | 1.1 |
| 12 | 5.7 |
| 14 | 15.2 |
| total | 32.5 |

fastest end-to-end timing minus `cat` time of 100+ runs for part1&2 in ms - mbair M1/16GB - darwin 23.2.0 - go version go1.21.4 darwin/arm64 - hyperfine 1.18.0 - 2023-12

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

## Day 3

Challenge is related to [Aoc2022/day23](https://adventofcode.com/2022/day/23), I'm using multiple bit arrays supported by a custom `u192` type. This solution is amazingly fast mainly because it is cache and CPU friendly. Almost all 2D ops are separated in 1D vector ops (think [`SIMD`](https://en.wikipedia.org/wiki/Single_instruction,_multiple_data)) and occur in a rolling window of 3 input lines.

Anyway, the challenge is akin to a [*static analysis*](https://en.wikipedia.org/wiki/Static_program_analysis) of a [*multi-valued game of life*](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) that's why it shares some of the techniques I used [last year](https://github.com/erik-adelbert/aoc/blob/2576e62f51f3bf653bf95084bca1815c534bf6e2/2022/23/aoc23.go).

I prize E. Wastl continuous effort on delivering such neat subjects every year. Here the challenge story is about a complex machinery with many cogwheels precisely timed and sized for the task: this is a fair description of what today's program feels like, a complex yet efficient machine with many simple parts that intricately but gracefully fall into place.

PS. To the young coders that might read this: don't be afraid! It's *not* a common day3 solution and certainly not the easiest way to solve it (but one of the fastest). The thing is, last november when warming up, I happened to refactor/improve AoC22/23.
So the bitpacking technique is still vivid in my memory. If it wasn't I may not have succeded in conjuring, factorizing and finally getting right all the corner cases and details of this solution in a fair amount of time.

## Day 4

Finally, day1 has come! Today's challenge is about typing speed with a few pauses here and there to actually think through the needed ops. As standard Go package `strings` has already proven usefull to tokenize inputs, I'm once again using it here.

The solution is totally linear, that is it follows closely the challenge tale and its runtime complexity is bounded by `O(n)` with `n` the input (deck) size. Given the small size of today's input (~200 lines), it is very fast.

PS. Isn't this awesome that at the heart of today's score calculation lies this beautiful gem:

```C
score += 1 << nmatch >> 1
```

`<EDIT>` following [`u/masklinn`](https://www.reddit.com/r/adventofcode/comments/18actmy/comment/kbzqx3e/?utm_source=share&utm_medium=web2x&context=3) advice, I went the extra mile consisting of replacing the winning number map by a [bitmap](https://en.wikipedia.org/wiki/Bitmap). I've also replaced the static 200+ deck by a [ring buffer](https://en.wikipedia.org/wiki/Circular_buffer). The resulting improvement is not measurable with hyperfine though.

## Day 5

[Intervals!](https://en.wikipedia.org/wiki/Interval_(mathematics))

Given that it is only day5 (and the input size), I'm *not* going to talk about [`Interval Trees`](https://en.wikipedia.org/wiki/Interval_tree). I am going to [`brute-force`](https://en.wikipedia.org/wiki/Brute-force_search) the thing!
`<Spoiler>` Well brute-forcing it, doesn't mean testing millions of points! It rather means to brute-force the interval *boundaries* checking.`</Spoiler>`.
I don't keep the challenge data format: it's not a so classical interval representation and therefore *will* get in the way. Today
program representation is: `{[α, ω), α'}` where `α'` is the destination category from the challenge.
By the way, the *same* code can be used for part 1 & 2, the trick is to output `part1` as `intervals` of length 0 instead of say, points.

PS. `Interval Trees` can also be found in the [big book](https://en.wikipedia.org/wiki/Introduction_to_Algorithms) 3rd ed. from pp. 348-353.

PS2. Look how fast the solution is \o/

## Day 6

Today is a direct application of solving this [`quadratic formula`](https://en.wikipedia.org/wiki/Quadratic_formula):

```C
    V.x = d <=> (x - t)*x - d = 0
```

And then adjusting for speed and time to exceed d.

For a very long time [`FPU`](https://en.wikipedia.org/wiki/Floating-point_unit) was slow but at the turn of y2k, `OS` and users alike were putting so much pressure on `CPU` that actually, `FPU` pipeline was usually free (and faster than before anyway) making it usable for a variety of computing that were previously carried on by `CPU` (the lore of General Purpose FPU was born). I remember the astonishment around me when one day I decided to benchmark the `FPU` against the `CPU` and showed that it won hands down in almost all situations.

But, for the sake of remembering those old days, I still don't want to switch to `FPU` when computing a square-root in an otherwise integer problem. I usually use a (fast) [`integer square root`](https://en.wikipedia.org/wiki/Integer_square_root) computation. In this very case, there's no reasonnable way to see the difference.

## Day 7

Now we have to [rank some hands](https://en.wikipedia.org/wiki/List_of_poker_hands#Hand-ranking_categories) akin to Poker but with slight variations. To this end and for maximum speed, we'd like to build the smallest representation of a hand that would:

1. fit a machine word
2. store the initial input as it's needed for comparison
3. lexographically sortable with standard integer sorting methods

Here we go, first of all there are 13 cards, namely `123456789TJQKA` so we need `4bits` per card and we have `5 cards` in a hand. That is `20bits` which is good because we could put the remainings `12bits` to good work. We know from our `wish #3` that we will compare all hands lexicographically, here it means that the we have to *reverse* the endianess of our hands. We can already bitmap the thing for more clarity:

<pre>
|0123456789abcdef|0123456789abcdef|
|_C5__C4__C3__C2_|_C1_............|
|1101010110010001|1110............| QJT98 actual hand: 89TJQ

_Cn_ are the 4 bits representing the nth card from 0 to 13
</pre>

One good thing with this mapping is that cards are aligned on word boundaries!

Now, that we have fullfilled almost all of our wishes, all we have left to do is to rank a hand in a [*monotonic*](https://en.wikipedia.org/wiki/Monotonic_function) way. That is we want a rank function `r` with `r(High card) < r(One pair) < ... r(Four of a kind) < r(Five of a kind)`. The naïve way would be to map from 1 to 6 all hand types:

```bash
{
    1: High Card, 2: One Pair, 3: Two Pairs, 4: Three of a kind, 5: Full house, 6: Four of a kind, 7: Five of kind
}
```

But we can already see that it would add some difficulties because it has no connection whatsoever with the *structure* of the hand.
I mean what is it to rank a hand? First, it is grouping and counting the card as in `A23A3 -> AA233` to make the rank obvious: `{A: 2}, {2: 1}, {3: 2}`. This last representation is easy to build and we know that it's the hand [*histogram*](https://en.wikipedia.org/wiki/Histogram). Building it will instantly provides a *base* value for our hand which is the highest card frequency:

```bash
{
    1: High, 2: One, 2: Two, 3: Three, 3: Full, 4: Four, 4: Five,
}
```

As you can see I've written `{4: Five of kind}` but why? The main reason is that we're still trying to build the smallest representation  and while `4` states is `3bits`, `5` is `4bits`. The second reason is that with this left over bit we can do something
better: we could consider that a `Full` is a special (stronger) case of a `Three` and the same goes for `Five` and `Four` or `Two` and `One`. This as the desirable effect to naturally insert the special hands in between the regular ones. So now if we use our bit as `X` flag for thoses special cases, it comes:

```bash
{
    (0,1): High, (0,1): One, (1,2): Two, (0,3): Three, (1,3): Full, (0,4): Four, (1,4): Five,
}
```

and the resulting bit mapping:

<pre>
|0123456789abcdef|0123456789abcdef|
|_C5__C4__C3__C2_|_C1_XRRR........| X: special bit, R: rank bit
|1101010110010001|11100100........| QJT98 actual hand: 89TJQ High

r(JKKK2) = 6999233
r(QQQQ2) = 9157553
</pre>

\o/ and that's it! Upon enconding a hand its unique `24bits` code is in lexicographic order.

```C
func cmp(a, b game) int {
    return int(a.hand - b.hand)
}
```

`part2` has it's own pitfalls and the details including a change in card scale are fun to study. If you look at it you will find a very good reason to use the `X` flag in regard to what is a `wildcard` and what it does to a special hand. But this write-up is already too long.

## Day 8

The crux of this challenge is to correctly encode the input. I mean obviously we could go for `type node struct{name, left, right str}` arranged in a `map[string]node` and it would fit. But do we really want to hash a million or so 3-letters strings?
Can we spare the hashing time? I tend to see programming as balancing time vs. space, so if we want to gain time we have to give space, but how much at most? A letter is `5bits` long: `A == 0, Z == 26 == 11010b`. Thus a node name is `15bits` long. There are `26x26x26 = 26^3 = 17576` unique node names with repetition. If we build a full storage of say 1 word (left, right) indexed by `15bits` node names it would be `17576 * 4bytes ≈ 68KB`, what a great deal!

So the idea here is to *encode the nodes on 3x5bits* and everything becomes natural. The command iterator is a light weight closure that returns 2 functions: one to get the current command and one to step the iterator. And last but not least:

```C
    hash&0x1f ==  0 <=> hash is ??A last letter is A <=> hash is root
    hash%0x1f == 25 <=> hash is ??Z last letter is Z <=> hash is goal
```

PS. There's a port of [binary gcd](https://en.wikipedia.org/wiki/Binary_GCD_algorithm) in this solution.

## Day 9

I don't understand the (lack of) difficulty today. The solution is totally described in the challenge, it is [recursive](https://en.wikipedia.org/wiki/Recursion) by design. This is the [method of differences](https://en.wikipedia.org/wiki/Telescoping_series). Let's call it a day!

## Day 10

This one is tricky because we can easily commit ourselves to tasks that were never needed. Look at this, the challenge occurs in a manhattan space and all the speeches about counting dots and tiles are just obfuscating the crux of this problem: we just want to compute the loop [area](https://en.wikipedia.org/wiki/Shoelace_formula)!

PS. I've included a nice `utf8` vizualisation in this one!

## Day 11

One thing I like so much with AoC is that comes december and suddenly I'm like a witch invited to a month long sabbath. I mean, every day I have to go for the challenge and then I can browse the subreddit to see what others have been up to. But from time to time one can learn much there: This morning I coded a boring `O(N²)` solution and it was ok fast (`~1ms`) but it felt like off to me. Long after, golfing and benchmarking, my mind was still wandering looking for something else to think about today's *universe*.

And finally, I found `u/edoannunziata`! Look at his [amazing solution](https://github.com/edoannunziata/jardin/blob/master/aoc23/AdventOfCode23.ipynb) (and everything else) which I understand but wouldn't have been able to come out with.

So I decided to study it by porting it, and the result is rewarding: It is so neat and fast! The point was to separate the 2D problem into independant 2x1D problems. When this technique exists (and around), the problem is dimmed [*separable*](https://en.wikipedia.org/wiki/Separable_filter).

PS. there's a minimal Python3 `Counter` port that was fun to compose in my solution.

## Day 12

Today's challenge is about one liner [nonograms](https://en.wikipedia.org/wiki/Nonogram)!

I believe there's no other way than [dynamic programming](https://en.wikipedia.org/wiki/Dynamic_programming) to solve today's challenge. This is why we may well find very fast `py3` cached recursive functions in the subreddit. My point was to make the fastest iterative *arranger* possible and to store directly the result instead of say caching the function call.
And I did!

Later on, I saw [this](https://www.reddit.com/r/adventofcode/comments/18ge41g/comment/kd0ohrj/?utm_source=share&utm_medium=web2x&context=3) idea about trimming and rolling the dp table and implemented it, resulting in the same result (but mine published later).

Every year, around day 12 to 15, there's a steep increase in AoC challenge runtimes, this could be it and surely `Dijkstra` is coming soon now.

## Day 13

Due to family commitments, I'm lagging a bit behind right now.

Nevermind, today's challenge is about finding the longuest [palindrome](https://en.wikipedia.org/wiki/Palindrome) and I've found a very fast (linear) way to find one but the program is not ready yet for release. I expect a runtime in the low ~~`~4ms`~~.

Here it is, running in __less than `1ms`__.

I knew finding palindromes was a quadratic task at best except it's not!! 

Look at [this awesome work](https://www.akalin.com/longest-palindrome-linear-time) by Fred Akalin. Actually, it is akin to build a `trie` but tailored to the task at hand. Instead of running in `O(n³)` for the naïve solution or in `O(n²)` for an improved technique, this one runs in __`Θ(n)`__ which is even faster than usually expected. The difference is really sensible even with small inputs like the ones from today. It comes with many good properties like:

- being bound by __`Θ(n)`__
- being bound by __`O(1)` space__
- finding __all palindromes in one pass__
- locating them on and __in between__ items

Solving today's challenge with this technology is like being blessed with a really cool super power \o/ (at least for `rio`, my mbair).

But I was in the mood for more, so I've build the solution upon a custom `bitarray32` which uses the transpose algorithm summerized in Warren's [hacker's delight](https://doc.lagout.org/security/Hackers%20Delight.pdf).

1) First, each input line is obviously summerized in a integer and stored in a `bitarray32`. 
2) Then we pipelin this `bitarray32` __directly__ to `flp` the fast palindrome finding bit.
3) We fast transpose the `bitarray32` and
4) Pipeline again the transposed `bitarray32` to `flp`

There's a lot of pitfalls along the way but all of them are manageable. Step #2 scores the columns while step #4 scores the rows.

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `cat input.txt > /dev/null` | 1.2 ± 0.2 | 1.0 | 3.4 | 1.00 |
| `./aoc13 < input.txt` | 2.2 ± 0.1 | 1.9 | 2.9 | 1.73 ± 0.28 |

PS. This may be one of the fastest if not the fastest solution for this challenge.
I've also included 6 additional sample cases some simple, some tricky.

## Day 14

This challenge is a beast! I went for the simulation but needed to be faster than fast: I built a custom `bitarray128` with *fast transpose*, *fast rotate* and *fast hash* built upon a custom `uint128` type with almost all the bells and whistles. And here it is, solving this challenge in `15ms`! The hashing speed enables a standard [`hash map`](https://en.wikipedia.org/wiki/Hash_table#:~:text=In%20computing%2C%20a%20hash%20table,that%20maps%20keys%20to%20values.) to support the cycle detection.

A fast `bitarray128` transpose and hashing is not easily available and the performances of this bitarray implementation could be worth publishing separately.

## Day 15

I don't like to fiddle with arrays and slices: Inserting/deleting from arrays is inefficient by nature and usually the sign of a poor design. Creating/Updating and then removing `lens` from `slots` is the perfect example:

- First of all, the problem is almost a *pure* byte one and the (hopefully small) `lens` names can be hashed efficiently. 

- The challenge really describes some kind of [buckets](https://en.wikipedia.org/wiki/Bucket_sort) that support *ops* in *time* this is exactly what [`queues`](https://en.wikipedia.org/wiki/Queue_(abstract_data_type)) are made for, not [arrays](https://en.wikipedia.org/wiki/Array_(data_structure))! We could say `arrays` are linked to `space` while `queues` are more linked to `time`.

- Actually [this](https://en.wikipedia.org/wiki/Bucket_queue) is the serious version of today's story.

So my idea from the start was to shuffle input commands to the various `queues` (boxes in the challenge) without doing anything more. A queue is built from a `hashmap` that record when a `lens` was removed for the last time, and an array of all the other bucketted commands. Once done, it was easy to built the slots without removing any `lens`: to this end it suffices to see if the candidate `lens` is in the delete list and if yes at what time it did enter there. If the current `lens` was allowed or arrived after the last `delete` command it was ok to add it, and voilà!

## Day 16

Is not ready yet.

## Day 17

Is not ready yet.

## Day 18

There's not to much to say for today challenge, except that we were compelled to use the [*shoelace formula*](https://en.wikipedia.org/wiki/Shoelace_formula) in 2023, in case some of us didn't use it on last Day 10.
