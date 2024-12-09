// aoc9.go --
// advent of code 2024 day 9
//
// https://adventofcode.com/2024/day/9
// https://github.com/erik-adelbert/aoc
//
// (ɔ) Erik Adelbert - erik_AT_adelbert_DOT_fr
// -------------------------------------------
// 2024-12-9: initial commit

package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
	MAXFILE = 20000
)

func main() {

	fs1 := φInit()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		var start int

		line := input.Bytes()
		for par, c := range line {
			if n := btoi(c); n > 0 {
				switch {
				case par%2 == 0:
					fs1.Store(File{{start, n}})
				default:
					fs1.Free(Block{start, n})
				}
				start += n
			}
		}
	}

	fs2 := fs1.Clone()
	fs2.Defrag()
	fs1.Compact()

	fmt.Println(fs1.Checksum(), fs2.Checksum()) // part 1 & 2
}

type Block struct {
	start int
	size  int
}

type File []Block

func (f File) Size() int {
	size := 0
	for _, b := range f {
		size += b.size
	}
	return size
}

type φFS struct {
	fat  []File
	free []Block
}

func φInit() *φFS {
	return &φFS{
		fat:  make([]File, 0, MAXFILE/2),
		free: make([]Block, 0, MAXFILE/2),
	}
}

func (fs *φFS) Clone() *φFS {
	clone := φInit()
	clone.fat = slices.Clone(fs.fat)
	clone.free = slices.Clone(fs.free)
	return clone
}

func (fs *φFS) Store(f File) {
	fs.fat = append(fs.fat, f)
}

func (fs *φFS) Free(b Block) {
	i, _ := slices.BinarySearchFunc(fs.free, b.start, func(x Block, start int) int {
		return x.start - start
	})

	fs.free = slices.Insert(fs.free, i, b)

	// merge with previous block if contiguous
	if i > 0 && fs.free[i-1].start+fs.free[i-1].size == fs.free[i].start {
		fs.free[i-1].size += fs.free[i].size
		fs.free = append(fs.free[:i], fs.free[i+1:]...)
		i-- // adjust index after merging
	}

	// merge with next block if contiguous
	if i+1 < len(fs.free) && fs.free[i].start+fs.free[i].size == fs.free[i+1].start {
		fs.free[i].size += fs.free[i+1].size
		fs.free = append(fs.free[:i+1], fs.free[i+2:]...)
	}
}

func (fs *φFS) Move(fid int) {
	size := fs.fat[fid].Size()

	allocated := 0

	blocks := make([]Block, 0, size)
	used := make([]int, 0, size)

	for i, block := range fs.free {
		if block.start > fs.fat[fid][0].start {
			return
		}

		if block.size < size {
			continue
		}

		blocks = append(blocks, block)
		used = append(used, i)
		if allocated += block.size; allocated > size {
			// keep the current block in the free list
			used = used[:len(used)-1]

			// split block
			free := allocated - size
			used := block.size - free

			blocks[len(blocks)-1].size = used
			fs.free[i] = Block{block.start + used, free}

			// unlink file from old blocks
			for i := range fs.fat[fid] {
				fs.Free(fs.fat[fid][i])
			}
		}
		break
	}
	if len(used) > 0 {
		// update free list
		fs.free = slices.Delete(fs.free, slices.Min(used), slices.Max(used)+1)
	}

	// link file to new blocks
	fs.fat[fid] = blocks
	return
}

func (fs *φFS) Realloc(fid int) {
	size := fs.fat[fid].Size()
	allocated := 0

	blocks := make([]Block, 0, size)
	used := make([]int, 0, size)

	for i := range fs.fat[fid] {
		fs.Free(fs.fat[fid][i])
	}
ALLOC:
	for i, block := range fs.free {
		allocated += block.size

		blocks = append(blocks, block)
		used = append(used, i)

		switch {
		case allocated < size:
			continue
		case allocated > size:
			// keep the current block in the free list
			used = used[:len(used)-1]

			// split block
			free := allocated - size
			used := block.size - free

			blocks[len(blocks)-1].size = used
			fs.free[i] = Block{block.start + used, free}

			fallthrough
		case allocated == size:
			// done allocating
			break ALLOC
		}
	}
	if len(used) > 0 {
		// update free list
		fs.free = slices.Delete(fs.free, slices.Min(used), slices.Max(used)+1)
	}

	// link file to new blocks
	fs.fat[fid] = blocks
	return
}

func (fs *φFS) Checksum() int {
	checksum := 0
	for i, file := range fs.fat {
		for _, block := range file {
			for j := 0; j < block.size; j++ {
				checksum += (block.start + j) * i
			}
		}
	}
	return checksum
}

func (fs *φFS) Compact() {
	for i := len(fs.fat) - 1; i >= 0 && len(fs.free) > 0; i-- {
		if fs.fat[i][0].start < fs.free[0].start {
			return
		}
		fs.Realloc(i)
	}
}

func (fs *φFS) Defrag() {
	for i := len(fs.fat) - 1; i >= 0; i-- {
		if fs.fat[i][0].start < fs.free[0].start {
			return
		}
		fs.Move(i)
	}
}

func btoi(b byte) int {
	return int(b - '0')
}
