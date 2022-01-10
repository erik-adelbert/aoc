#!/usr/bin/env python
import sys

from statistics import mean, median

def load(fd):
    line = fd.readline()
    return [int(x) for x in line.strip().split(",")]

def part1(data) :
    m = median(sorted(data))
    return [abs(x - m) for x in data]

def part2(data):
    def g(x):
        return (x * (x+1))/2

    m = int(mean(data))  # round doesn't work on my input
    return [g(abs(x - m)) for x in data]

def main():
    data = load(sys.stdin)
    for fun in [part1, part2]:
        print(int(sum(fun(data))))

if __name__ == "__main__":
    main()