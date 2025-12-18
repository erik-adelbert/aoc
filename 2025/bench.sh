#!/bin/bash

# hyperfine --warmup 100 "cat input.txt > /dev/null" "./${1} < input.txt" --export-markdown time.md

if [ -z "$1" ]; then
  echo "Usage: $0 <binary>" >&2
  exit 1
fi

binary="$1"
best=""
for _ in {1..100}; do
  # Extract value and unit (e.g., 1234.56 µs, 1234 µs, 1.234 ms, 1234us)
  out=$("./$binary" < input.txt 2>&1 | grep -oE '[0-9]+(\.[0-9]+)? *([µu]?s|ms)')
  val=$(echo "$out" | grep -oE '^[0-9]+(\.[0-9]+)?')
  unit=$(echo "$out" | grep -oE '(ms|[µu]s)$')

  # Default to µs if unit not found
  if [ -z "$unit" ]; then
    unit="us"
  fi

  # Convert ms to µs if needed
  if [ "$unit" = "ms" ]; then
    val_us=$(echo "$val * 1000" | bc -l)
  else
    val_us=$val
  fi

  if [ -n "$val_us" ]; then
    if [ -z "$best" ] || [ "$(echo "$val_us < $best" | bc -l)" -eq 1 ]; then
      best=$val_us
    fi
  fi
done

# Output in ms if >= 1000 µs, else in µs
if [ -z "$best" ]; then
  echo "No valid timing output found." >&2
  exit 2
fi

if [ "$(echo "$best >= 1000" | bc -l)" -eq 1 ]; then
  ms=$(echo "scale=3; $best/1000" | bc)
  echo "Best time: $ms ms"
else
  echo "Best time: $best µs"
fi
