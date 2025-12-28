import sys

spans = []
queries = []
reading_spans = True
for line in sys.stdin:
    line = line.strip()
    if not line:
        reading_spans = False
        continue
    if reading_spans:
        if "-" in line:
            start, end = map(int, line.split("-"))
            spans.append((start, end))
    else:
        queries.append(int(line))

if not spans or not queries:
    print("No spans or queries found.")
    sys.exit(1)


# Sort and merge spans
def merge_spans(spans):
    spans = sorted(spans)
    merged = []
    for s in spans:
        if not merged or s[0] > merged[-1][1] + 1:
            merged.append(list(s))
        else:
            merged[-1][1] = max(merged[-1][1], s[1])
    return merged


merged = merge_spans(spans)

# Span stats
span_lengths = [end - start + 1 for start, end in merged]
span_gaps = [merged[i][0] - merged[i - 1][1] - 1 for i in range(1, len(merged))]
covered_points = sum(span_lengths)
span_min = min(span_lengths)
span_max = max(span_lengths)
span_mean = covered_points / len(merged)

gap_min = min(span_gaps) if span_gaps else 0
gap_max = max(span_gaps) if span_gaps else 0
gap_mean = sum(span_gaps) / len(span_gaps) if span_gaps else 0

# Range stats
all_starts = [s[0] for s in merged]
all_ends = [s[1] for s in merged]
overall_min = min(all_starts)
overall_max = max(all_ends)
overall_range = overall_max - overall_min + 1

density = covered_points / overall_range
scarcity = 1 - density

# Query stats
queries = sorted(queries)
q_min = min(queries)
q_max = max(queries)
q_mean = sum(queries) / len(queries)

# Coverage ratio
covered_queries = 0
j = 0
for q in queries:
    while j < len(merged) and q > merged[j][1]:
        j += 1
    if j < len(merged) and merged[j][0] <= q <= merged[j][1]:
        covered_queries += 1
coverage_ratio = covered_queries / len(queries)

print(f"Spans: {len(spans)} (merged: {len(merged)})")
print(f"Span length: min={span_min}, max={span_max}, mean={span_mean:.2f}")
print(f"Span gaps: min={gap_min}, max={gap_max}, mean={gap_mean:.2f}")
print(f"Total covered points: {covered_points}")
print(f"Overall range: {overall_min}-{overall_max} (size {overall_range})")
print(f"Density: {density:.6f}, Scarcity: {scarcity:.6f}")
print(f"Queries: {len(queries)} (min={q_min}, max={q_max}, mean={q_mean:.2f})")
print(f"Queries covered: {covered_queries} ({coverage_ratio:.2%})")
