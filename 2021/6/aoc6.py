#!/usr/bin/env python
import sys

from io import StringIO
from collections import Counter, deque

def incube(fishes, days):
    popcnt = deque([fishes[i] for i in range(9)])
    for _ in range(days):
        popcnt.rotate(-1)
        popcnt[6] += popcnt[-1]
    return popcnt


def load(fd):
    lines = fd.readline()
    return Counter(int(x) for x in lines.strip().split(","))
    

def main():
    fishes = load(sys.stdin)
    days = int(sys.argv[1])
    print(sum(incube(fishes, days)))

if __name__ == "__main__":
    main()