#!/bin/bash

hyperfine --warmup 100 "cat input.txt > /dev/null" "./${1} < input.txt" --export-markdown time.md
