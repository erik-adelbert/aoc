# pylint: disable=C0103,R0914 we don't care about naming or function length here
# 3D Graph Visualization for Advent of Code input.txt
#
# Requirements:
#   pip install networkx plotly scipy
#
# Usage:
#   python3 webgraph.py
#
# This script reads 'input.txt' in the current directory, parses the graph,
# and displays a 3D interactive plot.


import sys

import networkx as nx
import plotly.graph_objs as go


def parse(filename):
    """Parse the input file and return a list of (src, dst) edges."""
    edges = []
    with open(filename, encoding="utf-8") as f:
        for line in f:
            if not line.strip():
                continue
            src, dsts = line.split(":")
            src = src.strip()
            for dst in dsts.strip().split():
                edges.append((src, dst))
    return edges


def build(edges):
    """Build and return a NetworkX graph from a list of edges."""
    G = nx.Graph()
    G.add_edges_from(edges)
    return G


def plot(G):
    """Visualize the graph G in 3D, highlighting special nodes and their neighbors."""
    pos = nx.spring_layout(G, dim=3, seed=42)
    specials = {"you", "svr", "fft", "dac", "out"}
    # Find neighbors of special nodes
    neighbors = set()
    for n in specials:
        if n in G:
            neighbors.update(G.neighbors(n))
    # Extract node positions and colors
    Xn, Yn, Zn, colors, texts = [], [], [], [], []
    for k in G.nodes():
        Xn.append(pos[k][0])
        Yn.append(pos[k][1])
        Zn.append(pos[k][2])
        if k in specials:
            colors.append("red")
            texts.append(k)
        elif k in neighbors:
            colors.append("orange")
            texts.append("")
        else:
            colors.append("skyblue")
            texts.append("")

    # Edges
    Xe, Ye, Ze = [], [], []
    for e in G.edges():
        Xe += [pos[e[0]][0], pos[e[1]][0], None]
        Ye += [pos[e[0]][1], pos[e[1]][1], None]
        Ze += [pos[e[0]][2], pos[e[1]][2], None]

    # Create traces
    edge_trace = go.Scatter3d(
        x=Xe,
        y=Ye,
        z=Ze,
        mode="lines",
        line={"color": "gray", "width": 2},
        hoverinfo="none",
    )

    node_trace = go.Scatter3d(
        x=Xn,
        y=Yn,
        z=Zn,
        mode="markers+text",
        marker={"symbol": "circle", "size": 4, "color": colors},
        text=texts,
        textposition="top center",
        hoverinfo="text",
    )

    layout = go.Layout(
        title={"text": "AoC 2025 Graph Day 11", "font": {"color": "white"}},
        paper_bgcolor="black",
        showlegend=False,
        scene={
            "xaxis": {
                "showgrid": True,
                "gridcolor": "#6BACDD",
                "zeroline": True,
                "zerolinecolor": "black",
                "showline": True,
                # "linecolor": "#6BACDD",
                "backgroundcolor": "#f0f0f0",
                "tickfont": {"color": "white"},
                "title": {"font": {"color": "white"}},
            },
            "yaxis": {
                "showgrid": True,
                "gridcolor": "#6BACDD",
                "zeroline": True,
                "zerolinecolor": "black",
                "showline": True,
                # "linecolor": "#6BACDD",
                "backgroundcolor": "#f0f0f0",
                "tickfont": {"color": "white"},
                "title": {"font": {"color": "white"}},
            },
            "zaxis": {
                "showgrid": True,
                "gridcolor": "#6BACDD",
                "zeroline": True,
                "zerolinecolor": "black",
                "showline": True,
                # "linecolor": "#6BACDD",
                "backgroundcolor": "#f0f0f0",
                "tickfont": {"color": "white"},
                "title": {"font": {"color": "white"}},
            },
            "camera": {
                "center": {"x": 0, "y": 0, "z": 0},
                "eye": {"x": 1.8, "y": 1.8, "z": 1.0},
            },
            "bgcolor": "black",
        },
        margin={"l": 0, "r": 0, "b": 0, "t": 40},
    )

    fig = go.Figure(data=[edge_trace, node_trace], layout=layout)
    fig.show()


def main():
    """Main entry point: parse input, build graph, and plot in 3D."""
    if len(sys.argv) > 1:
        input_file = sys.argv[1]
    else:
        input_file = "input.txt"

    edges = parse(input_file)

    G = build(edges)
    plot(G)


if __name__ == "__main__":
    main()
