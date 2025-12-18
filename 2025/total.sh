#!/bin/bash
# Usage: total.sh
# Runs 'make run' in each day subdirectory (1..25), collects and sums timings from their outputs.

for d in {1..25}; do
  if [ -d "$d" ]; then
    (cd "$d" && out=$(make bench 2>&1 | grep -Eo '[0-9]+(\.[0-9]+)? ?(µs|us|ms)' | head -1)
      val=$(echo "$out" | grep -Eo '^[0-9]+(\.[0-9]+)?')
      unit=$(echo "$out" | grep -Eo '(ms|[µu]s)$')
      if [ -z "$unit" ]; then unit="us"; fi
      if [ "$unit" = "ms" ]; then
        val_us=$(echo "$val * 1000" | bc -l)
      else
        val_us=$val
      fi
      printf "%s: %s %s\n" "$d" "$val" "$unit"
      echo "$val_us"
    )
  fi
done | awk '/^[0-9]/ { sum += $1 } /^[0-9]+: / { next } END { printf "Total: %.3f µs (%.3f ms)\n", sum, sum/1000 }'
