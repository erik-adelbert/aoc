import matplotlib.pyplot as plt

# Read coordinates from input.txt
xs, ys = [], []
with open("input.txt") as f:
    for line in f:
        x, y = map(int, line.strip().split(","))
        xs.append(x)
        ys.append(y)

# Close the path if it's a polygon
if xs[0] != xs[-1] or ys[0] != ys[-1]:
    xs.append(xs[0])
    ys.append(ys[0])

plt.figure(figsize=(8, 8))

# Fill the polygon (interior surface) in green
plt.fill(xs, ys, color="limegreen", alpha=0.3, zorder=1)

# Draw the edges in green
plt.plot(xs, ys, color="green", linewidth=1, zorder=2)

# Draw waypoints as small, thin red dots
plt.scatter(xs, ys, color="red", s=2, zorder=3)

plt.title("Path Plot")
plt.xlabel("X")
plt.ylabel("Y")
plt.axis("equal")
plt.grid(True)
plt.show()
