#!/bin/bash

hyperfine --warmup 5 --min-runs 1000 "cat input.txt" "cat input.txt | ./${1}" --export-markdown time.md
