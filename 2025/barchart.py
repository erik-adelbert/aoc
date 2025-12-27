import matplotlib.pyplot as plt

# Data
days = ["2", "7", "5", "9", "12", "1", "6", "11", "3", "10", "4", "8"]
times = [8, 30, 95, 106, 119, 124, 150, 157, 189, 248, 695, 863]

# Colors and text color
bar_color = "#808080"  # medium gray bars
text_color = "#808080"  # medium gray text

# Create figure
fig, ax = plt.subplots(figsize=(10, 6))

# Horizontal bar chart
bars = ax.barh(days, times, color=bar_color)

# Add labels on bars
for bar in bars:
    width = bar.get_width()
    ax.text(
        width + 10,
        bar.get_y() + bar.get_height() / 2,
        f"{width}",
        va="center",
        ha="left",
        color=text_color,
    )

# Style axes
ax.set_xlabel("Time (μs)", color=text_color)
ax.set_ylabel("Day", color=text_color)
ax.tick_params(colors=text_color)  # axis tick labels
ax.spines["top"].set_visible(False)
ax.spines["right"].set_visible(False)
ax.spines["left"].set_color(text_color)
ax.spines["bottom"].set_color(text_color)

# Transparent background
fig.patch.set_alpha(0.0)
ax.patch.set_alpha(0.0)

# Invert y-axis for descending order like a typical bar chart
ax.invert_yaxis()

plt.tight_layout()
plt.savefig("barchart.png", transparent=True)
plt.show()
