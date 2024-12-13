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
	"sync"
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
					fs1.Link(File{{start, n}})
				default:
					fs1.MarkFree(Block{start, n})
				}
				start += n
			}
		}
	}

	fs2 := fs1.Clone()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fs2.Defrag()
	}()
	fs1.Compact()

	wg.Wait()

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

func (fs *φFS) Link(f File) {
	fs.fat = append(fs.fat, f)
}

func (fs *φFS) Unlink(fid int) {
	file := fs.fat[fid]
	for i := range file {
		fs.Free(file[i])
	}
}

func (fs *φFS) MarkFree(b Block) {
	fs.free = append(fs.free, b)
}

func (fs *φFS) Free(b Block) {
	i, _ := slices.BinarySearchFunc(fs.free, b, func(x Block, b Block) int {
		return x.start - b.start
	})

	// merge with previous block if contiguous
	if i > 0 && fs.free[i-1].start+fs.free[i-1].size == b.start {
		fs.free[i-1].size += b.size
	} else {
		fs.free = slices.Insert(fs.free, i, b)
	}
}

func fsmap(fs *φFS, fun func(int)) {
	for fid := len(fs.fat) - 1; fid >= 0; fid-- {
		file := fs.fat[fid]

		if file[0].start < fs.free[0].start {
			return // no more free space at the beginning
		}

		fun(fid)
	}
}

func (fs *φFS) Compact() {
	fsmap(fs, fs.Realloc)
}

func (fs *φFS) Defrag() {
	fsmap(fs, fs.Move)
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

func (fs *φFS) Move(fid int) {
	size := fs.fat[fid].Size()

	nalloc := 0

	blocks := make([]Block, 0)
	reserved := -1

	file := fs.fat[fid]

ALLOC:
	for i, block := range fs.free {
		switch {
		case block.start > file[0].start:
			// can't move file to a lower address
			return
		case block.size < size:
			// not enough space
			continue
		default:
			// move file to new block
			blocks = append(blocks, block)
			reserved = i

			if nalloc += block.size; nalloc > size {
				// keep the current block in the free list
				reserved = -1

				// split block
				free := nalloc - size
				used := block.size - free

				blocks[len(blocks)-1].size = used
				fs.free[i] = Block{block.start + used, free}
			}

			// link file to new blocks
			fs.fat[fid] = blocks
			break ALLOC
		}
	}
	if reserved >= 0 {
		// fs.free = append(fs.free[:reserved], fs.free[reserved+1:]...)
		fs.free = slices.Delete(fs.free, reserved, reserved+1)
	}

	return
}

func (fs *φFS) Realloc(fid int) {
	file := fs.fat[fid]
	size := file.Size()
	nalloc := 0

	blocks := make([]Block, 0, size)
	reserved := make([]int, 0, size)

	if len(fs.free) < 10 {
		fs.Unlink(fid)
	}
ALLOC:
	for i, block := range fs.free {
		nalloc += block.size

		blocks = append(blocks, block)
		reserved = append(reserved, i)

		switch {
		case nalloc < size:
			// not done yet
			continue
		case nalloc > size:
			// keep the current block in the free list
			reserved = reserved[:len(reserved)-1]

			// split block
			free := nalloc - size
			used := block.size - free

			blocks[len(blocks)-1].size = used
			fs.free[i] = Block{block.start + used, free}

			fallthrough
		default:
			// done allocating
			break ALLOC
		}
	}
	if len(reserved) > 0 {
		// update free list
		fs.free = slices.Delete(fs.free, reserved[0], reserved[len(reserved)-1]+1)
	}

	// link file to new blocks
	fs.fat[fid] = blocks
	return
}

func (fs *φFS) Clone() *φFS {
	clone := φInit()
	clone.fat = slices.Clone(fs.fat)
	clone.free = slices.Clone(fs.free)
	return clone
}

func btoi(b byte) int {
	return int(b - '0')
}
