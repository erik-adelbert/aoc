#!/bin/zsh

# Input file
file="input.txt"

# Extract the first and second columns into separate arrays
col1=($(awk '{print $1}' $file))
col2=($(awk '{print $2}' $file))

# Function to count duplicates in an array
count_duplicates() {
  local arr=("$@")
  local counts=()

  # Count occurrences using associative array
  for val in $arr; do
    ((counts[$val]++))
  done

  # Count and print the number of duplicates
  local duplicate_count=0
  for val count in ${(kv)counts}; do
    if (( count > 1 )); then
      ((duplicate_count++))
    fi
  done

  echo $duplicate_count
}

# Count duplicates for each column
echo "Duplicates in column 1: $(count_duplicates $col1)"
echo "Duplicates in column 2: $(count_duplicates $col2)"

