# pylint: disable=C0103,R0914 # we don't care about naming or function length here

"""
3D Graph Visualization for Advent of Code input.txt

Requirements:
    pip3 install networkx plotly scipy

Usage:
    python3 webgraph.py

This script reads 'input.txt' in the current directory, parses the graph,
and displays a 3D interactive plot.
"""

import sys

import networkx as ntx
import plotly.graph_objs as gro


def parse(filename):
    """Parse the input file and return a list of (src, dst) edges."""

    edges = []
    with open(filename, encoding="utf-8") as f:
        for raw in f:
            line = raw.strip()
            if not line:
                continue

            src, sep, dsts = line.partition(":")
            if not sep:
                continue  # skip lines without a colon

            src = src.strip()
            edges.extend((src, dst) for dst in dsts.split())

    return edges


def graph(edges):
    """Build and return a NetworkX graph from a list of edges."""
    G = ntx.Graph()
    G.add_edges_from(edges)
    return G


def plot(G):
    """Visualize the graph G in 3D, highlighting special nodes and their neighbors."""
    pos = ntx.spring_layout(G, dim=3, seed=42)
    specials = {"you", "svr", "fft", "dac", "out"}

    # Find neighbors of special nodes
    neighbors = set()
    for n in specials:
        if n in G:
            neighbors.update(G.neighbors(n))

    # Extract node positions and colors
    nodeX, nodeY, nodeZ, colors, texts = [], [], [], [], []
    append_x, append_y, append_z = nodeX.append, nodeY.append, nodeZ.append
    append_c, append_t = colors.append, texts.append
    for k in G.nodes():
        x, y, z = pos[k]
        append_x(x)
        append_y(y)
        append_z(z)

        color, text = "skyblue", ""
        if k in neighbors:
            color = "orange"
        if k in specials:  # specials override neighbors
            color = "red"
            text = k

        append_c(color)
        append_t(text)

    # Edges
    edgeX, edgeY, edgeZ = [], [], []
    extend_x, extend_y, extend_z = edgeX.extend, edgeY.extend, edgeZ.extend

    for u, v in G.edges():
        x0, y0, z0 = pos[u]
        x1, y1, z1 = pos[v]

        extend_x([x0, x1, None])
        extend_y([y0, y1, None])
        extend_z([z0, z1, None])

    # Create traces
    edge_trace = gro.Scatter3d(
        x=edgeX,
        y=edgeY,
        z=edgeZ,
        mode="lines",
        line={"color": "gray", "width": 2},
        hoverinfo="none",
    )

    node_trace = gro.Scatter3d(
        x=nodeX,
        y=nodeY,
        z=nodeZ,
        mode="markers+text",
        marker={"symbol": "circle", "size": 4, "color": colors},
        text=texts,
        textposition="top center",
        hoverinfo="text",
    )

    # Layout configuration for the 3D plot
    blue_axis = {
        "showgrid": True,
        "gridcolor": "#6BACDD",
        "zeroline": True,
        "zerolinecolor": "black",
        "showline": True,
        # "linecolor": "#6BACDD",
        "backgroundcolor": "#f0f0f0",
        "tickfont": {"color": "white"},
        "title": {"font": {"color": "white"}},
    }

    layout = gro.Layout(
        title={"text": "AoC 2025 Day 11 Graph", "font": {"color": "white"}},
        paper_bgcolor="black",
        showlegend=False,
        scene={
            "xaxis": blue_axis,
            "yaxis": blue_axis,
            "zaxis": blue_axis,
            "camera": {
                "center": {"x": 0, "y": 0, "z": 0},
                "eye": {"x": 1.8, "y": 1.8, "z": 1.0},
            },
            "bgcolor": "black",
        },
        margin={"l": 0, "r": 0, "b": 0, "t": 40},
    )

    # Create the figure and display it
    fig = gro.Figure(data=[edge_trace, node_trace], layout=layout)
    fig.show()


def main():
    """Main entry point: parse input, build graph, and plot in 3D."""
    if len(sys.argv) > 1:
        input_file = sys.argv[1]
    else:
        input_file = "input.txt"

    edges = parse(input_file)

    plot(graph(edges))


if __name__ == "__main__":
    main()
