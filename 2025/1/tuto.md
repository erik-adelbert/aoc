
# How can one solve today's challenge efficiently?

Efficiency here means our goal is to solve the problem as fast as possible while avoiding any complication. Here is a technique that is generally useful: we are going to solve two simpler problems (A & B) and then decompose the challenge into these in order to solve it.

## Simplified problem

First let's think of a simplified problem where the dial lies at 0 before we start moving it and every turn is any number of **wraps** and we have to solve **n** such moves:

- `L100` or `R100`
- `L2000` or `R2000`
- `Lk` or `Rk` where `0 < k and k is a multiple of 100` (general case)

From there, the answer for the challenge parts is obvious (very good):

- For part 1: each move (starts at 0 and) ends at 0, so the answer is $part1 = n$
- For part 2: each move *i* crosses 0 $k_i$ times, so the answer is $\sum_{i=0}^{n} k_i$

| move               | part1 | part2   |
|--------------------|------:|--------:|
| `L100` or `R100`   |     1 |       1 |
| `L2000` or `R2000` |     1 |      20 |
| `Lk` or `Rk`       |     1 |       k |
| **Total**          | **3** |**21+k** |

And in code:

     part1 := n

     part2 := 0
     for i range n {
        part2 += move[i]/100 // k
     }

## Less simplified problem A

Let's just say that instead of laying at 0 initially, now the dial lies at 1.

I leave the reasoning to you, but now we have:

     part1 := 0

     part2 := 0
     for i range n {
        part2 += move[i]/100 // k, /!\ int division: 1/2 = 0, 5/2 = 2
     }

It's time to say that this is true not only for a dial starting at 1 but for any starting positions from 1 to 99. So the conclusion, for now is:

- **part1 depends on the original position**
- *whatever* the move or the starting position, **part2 always depends on the number of wraps**.

## Less simplified problem B

Now consider that we are starting from any position in `[0, 99]`, but the move *m* is always in `[1, 99]` (it is never a full wrap). What happens then?

With *s* as the starting position and *e* as the ending one, and going **left**:

- when *s* is 0 it is impossible to land on 0, *e* != 0
- if *e* > *s* we have crossed 0

going **right** now:

- when *s* is 0 it is impossible to land on 0, *e* != 0
- if *e* < *s* we have crossed 0

## Putting everything together

Consider move L101, it is a full wrap and then a single left step. If we generalize for any move of distance *d*:

     wraps := d/100
     steps := d%100

But wait a minute! **wraps** is **problemA** and **steps** is **problemB**!

For any move with direction *dir*, distance *d* from starting position *s* to ending position *e*:

    wraps := d/100 // problem A
    steps := d%100 // problem B

    part2 += wraps // problem A: part2 always depends on the number of full wraps

    d = steps // problem B: steps in [1..99]

    // compute ending position
    // s in [0, 99], d in [1, 99] => e in [-99, 198]
    if dir == Left {
        e = (s - d)
    } else {
        e = (s + d)
    }

    // handle circular dial
    e = e%100 // e in [-99,99]
    if e < 0 {
        e += 100 // e always in [0, 99]
    }

    // problem B: handle steps
    switch {
    case s == 0:
        // can't reach or cross 0 from 0
        // count nothing
    case e == 0:
        // landed on 0
        part1++
        part2++
    case (e > s) == (dir == Left):
        // 0-crossings
        // position increased/decreased when turning left/right
        part2++
    }

This has to be wrapped in a loop and *s* advanced to *e*... And that's it!
