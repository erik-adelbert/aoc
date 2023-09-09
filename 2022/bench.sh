#!/bin/bash

hyperfine --warmup 5 --min-runs 1000 "cat input.txt > /dev/null" "./${1} < input.txt" --export-markdown time.md
