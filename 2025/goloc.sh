#!/bin/sh
echo '| Directory | Go Lines |'
echo '|:----------|---------:|'
for dir in "$@"; do
  if [ -d "$dir" ]; then
    lines=$(cloc --include-lang=Go --csv "$dir" 2>/dev/null | awk -F, '$2=="Go"{print $5}')
    echo "| $dir | ${lines:-0} |"
  else
    echo "| $dir | 0 |"
  fi
done
