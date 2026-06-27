"""
Bar chart for AOC 2025 benchmark times
A basic horizontal bar chart with matplotlib
"""

import matplotlib.pyplot as plt

# Data
data = {
    "2": 6,
    "7": 23,
    "5": 64,
    "12": 72,
    "9": 81,
    "3": 83,
    "1": 87,
    "6": 94,
    "11": 106,
    "10": 219,
    "4": 453,
    "8": 542,
}

days, times = zip(
    *sorted(data.items(), key=lambda item: item[1])
)  # sort by time then split

# Colors and text color
BAR_COLOR = "#808080"  # medium gray bars
TEXT_COLOR = "#808080"  # medium gray text

# Create figure
fig, ax = plt.subplots(figsize=(10, 6))

# Horizontal bar chart
bars = ax.barh(days, times, color=BAR_COLOR)

# Add labels on bars
for bar in bars:
    width = bar.get_width()
    ax.text(
        width + 10,
        bar.get_y() + bar.get_height() / 2,
        f"{width}",
        va="center",
        ha="left",
        color=TEXT_COLOR,
    )

# Style axes
ax.set_xlabel("Time (μs)", color=TEXT_COLOR)
ax.set_ylabel("Day", color=TEXT_COLOR)
ax.tick_params(colors=TEXT_COLOR)  # axis tick labels
ax.spines["top"].set_visible(False)
ax.spines["right"].set_visible(False)
ax.spines["left"].set_color(TEXT_COLOR)
ax.spines["bottom"].set_color(TEXT_COLOR)

fig.set_alpha(0.0)  # Transparent figure background
ax.set_alpha(0.0)  # Transparent axes background

# Invert y-axis for descending order
ax.invert_yaxis()

plt.tight_layout()
plt.savefig("barchart.png", transparent=True)
plt.show()
