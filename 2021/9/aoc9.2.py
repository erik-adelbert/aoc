import numpy as np
from scipy import ndimage

grid = []
with open("input.txt") as f: 
    for args in f: 
        args = list(args.split()[0]) 
        grid.append([int(arg) for arg in args])

basins, n = ndimage.label((np.array(grid) < 9).astype(int))
counts = [(basins == x).sum() for x in range(1, n+1)]
print(np.sort(counts)[-3:].prod())