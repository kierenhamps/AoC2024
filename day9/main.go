package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Blocks map[int]int

// FileSystem represents a file system
type FileSystem struct {
	files      FilePointer
	freeBlocks BlockPointer
	dataBlocks BlockPointer
	blocks     Blocks
}

// NewFileSystem returns a new FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{
		files:      *NewFilePointer(),
		freeBlocks: *NewBlockPointer(),
		dataBlocks: *NewBlockPointer(),
		blocks:     make(Blocks, 0),
	}
}

// Checksum returns the checksum of the file system
// according to part1 rules
func (fs *FileSystem) Checksum() int {
	var checksum int
	for i, v := range fs.blocks {
		checksum += i * v
	}
	return checksum
}

// Compact compacts the file system
// While there are blocks of free space in between blocks of data,
// data at the end of the FileSystem will be relocated to the free space
func (fs *FileSystem) Compact() {
	// get the last free block and the last data block IDs
	lastFreeBlockID, ok := fs.freeBlocks.PeekFirst()
	if !ok {
		log.Println("no free blocks found")
		return
	}
	lastDataBlockID, ok := fs.dataBlocks.PeekLast()
	if !ok {
		log.Println("no data blocks found")
		return
	}

	// while there are free blocks before the last data block
	// we need to compact the data
	for lastDataBlockID > lastFreeBlockID {
		// get the last data block
		dataBlockID, _ := fs.dataBlocks.PopLast()
		dataBlock := fs.blocks[dataBlockID]

		// get the first free block
		freeBlockID, _ := fs.freeBlocks.PopFirst()

		// move the data block to the free block and add to data stack
		fs.blocks[freeBlockID] = dataBlock
		fs.dataBlocks.Push(freeBlockID)
		// update the file map
		fs.files.MoveBlock(dataBlock, dataBlockID, freeBlockID)

		// remove the old data block and add to free stack
		delete(fs.blocks, dataBlockID)
		fs.freeBlocks.Push(dataBlockID)

		// update the loop variables
		lastFreeBlockID, ok = fs.freeBlocks.PeekFirst()
		if !ok {
			break
		}
		lastDataBlockID, ok = fs.dataBlocks.PeekLast()
		if !ok {
			break
		}
	}
}

// CompactByFile compacts the file system moving one file worth of blocks at a time
// Files are moved from biggest to smallest id, and will only move to the lowest
// available free block if the entire file can be moved.
// If there is not enough free space in one contiguous block, the file will not be moved.
func (fs *FileSystem) CompactByFile() {
	// foreach file in our file map
	sortedIDs := []int{}
	for k, _ := range fs.files.files {
		sortedIDs = append(sortedIDs, k)
	}
	sort.Ints(sortedIDs)

	for i := len(sortedIDs) - 1; i >= 0; i-- {
		// get the last file by id
		fileID := sortedIDs[i]
		blocks, ok := fs.files.files[fileID]
		if !ok {
			break
		}

		// file size
		fileSize := len(blocks)

		// find where the file can fit in free space from the start
		canBeMoved := false
		freeBlocks := fs.freeBlocks.ids
		var i int
		for i = 0; i < len(freeBlocks); i++ {
			// get the first free block
			freeBlockID := freeBlocks[i]
			// only move if free space is before current location
			if freeBlockID > blocks[0] {
				break
			}
			// does the file fit?
			if fs.freeBlocks.IsContiguous(freeBlockID, freeBlockID+fileSize-1) {
				canBeMoved = true
				break
			}
		}

		// if the file can fit in free space, move the file
		if canBeMoved {
			// move the file
			for j := 0; j < fileSize; j++ {
				// get the block to move
				blockID := blocks[j]
				// move the block
				fs.blocks[freeBlocks[i]+j] = fileID
				// update the file map
				fs.dataBlocks.Push(freeBlocks[i] + j)
				// update the file map
				fs.files.MoveBlock(fileID, blockID, freeBlocks[i]+j)
				// remove the old free blocks
				fs.freeBlocks.ids = append(fs.freeBlocks.ids[i+1:], fs.freeBlocks.ids[:i]...)
				// remove the old data block and add to free stack
				delete(fs.blocks, blockID)
				fs.freeBlocks.Push(blockID)
			}
		}
	}

}

// BlockPointer is a stack of block IDs
type BlockPointer struct {
	ids []int
}

// NewBlockPointer returns a new BlockPointer
func NewBlockPointer() *BlockPointer {
	return &BlockPointer{
		ids: make([]int, 0),
	}
}

// Find returns the block ID and true if it exists
func (bp *BlockPointer) Find(id int) (int, bool) {
	for _, v := range bp.ids {
		if v == id {
			return v, true
		}
	}
	return 0, false
}

// IsContiguous returns true if the block IDs are contiguous
// startID and endID are inclusive
func (bp *BlockPointer) IsContiguous(startID, endID int) bool {
	for i := startID; i <= endID; i++ {
		if _, ok := bp.Find(i); !ok {
			return false
		}
	}
	return true
}

// PeekFirst returns the first ID from the stack without removing it
func (bp *BlockPointer) PeekFirst() (int, bool) {
	if len(bp.ids) == 0 {
		return 0, false
	}
	return bp.ids[0], true
}

// PeekLast returns the last ID from the stack without removing it
func (bp *BlockPointer) PeekLast() (int, bool) {
	if len(bp.ids) == 0 {
		return 0, false
	}
	return bp.ids[len(bp.ids)-1], true
}

// PopFirst removes the first ID from the stack and returns it
func (bp *BlockPointer) PopFirst() (int, bool) {
	if len(bp.ids) == 0 {
		return 0, false
	}
	// get the first ID
	first := bp.ids[0]

	// remove the first ID
	bp.ids = bp.ids[1:]

	return first, true
}

// PopLast removes the last ID from the stack and returns it
func (bp *BlockPointer) PopLast() (int, bool) {
	if len(bp.ids) == 0 {
		return 0, false
	}
	// get the last ID
	last := bp.ids[len(bp.ids)-1]

	// remove the last ID
	bp.ids = bp.ids[:len(bp.ids)-1]

	return last, true
}

// Push adds an ID to the end of the stack
func (bp *BlockPointer) Push(id int) {
	bp.ids = append(bp.ids, id)
	bp.Sort()
}

// Sort sorts the stack of IDs
func (bp *BlockPointer) Sort() {
	sort.Ints(bp.ids)
}

// FilePointer is a stack of files with their block IDs
type FilePointer struct {
	files map[int][]int
}

// NewFilePointer returns a new FilePointer
func NewFilePointer() *FilePointer {
	return &FilePointer{
		// initialise the files map
		files: make(map[int][]int, 0),
	}
}

// AppendToFile adds a blocks to a file
func (fp *FilePointer) AppendToFile(fileID int, blockIDs []int) {
	if _, exists := fp.files[fileID]; !exists {
		fp.files[fileID] = []int{}
	}
	fp.files[fileID] = append(fp.files[fileID], blockIDs...)
}

// MoveBlock moves a block from one location to another for a given file
func (fp *FilePointer) MoveBlock(file, from, to int) {
	// get the file blocks
	blocks, ok := fp.files[file]
	if !ok {
		return
	}
	// find the block to move
	var blockIndex int
	for _, b := range blocks {
		if b == from {
			// replace the block
			fp.files[file][blockIndex] = to
			return
		}
		blockIndex++
	}
}

// PeekLast returns the last file from the stack without removing it
func (fp *FilePointer) PeekLast() (int, []int, bool) {
	if len(fp.files) == 0 {
		return 0, nil, false
	}
	// get last key in map
	var lastID int
	var lastBlocks []int
	for k, v := range fp.files {
		if k > lastID {
			lastID = k
			lastBlocks = v
		}
	}
	return lastID, lastBlocks, true
}

// Push adds a file to the stack
func (fp *FilePointer) Push(file map[int][]int) {
	for k, v := range file {
		fp.files[k] = v
	}
}

// ParseDiskMap reads the input and returns a new FileSystem
// with populated Blocks
func ParseDiskMap(input io.Reader) *FileSystem {
	scanner := bufio.NewScanner(input)
	fs := NewFileSystem()
	for scanner.Scan() {
		line := scanner.Text()
		blocks := strings.Split(line, "")
		blockID := 0
		fileID := -1
		// get position and data for each block
		for p, d := range blocks {
			// determin	if the block is a file or a free block
			// using the position (p) of the block
			isFree := true
			if p%2 == 0 {
				isFree = false
				fileID++
			}

			// number is the number of times the block is repeated
			number, err := strconv.Atoi(d)
			if err != nil {
				log.Fatalf("failed to convert %s to int: %v", d, err)
			}

			// for each time the block is repeated
			for k := 0; k < number; k++ {
				if isFree {
					// record the free block
					fs.freeBlocks.Push(blockID)
					blockID++
				} else {
					// record the data block
					fs.blocks[blockID] = fileID
					fs.dataBlocks.Push(blockID)
					// Add file to FilePointer
					fs.files.AppendToFile(fileID, []int{blockID})
					blockID++
				}
			}
		}
	}
	return fs
}

func main() {
	// read the input file
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}
	defer input.Close()

	// parse the disk map
	fs := ParseDiskMap(input)

	// compact the file system
	fs.Compact()

	// print the checksum
	log.Printf("(Part 1) checksum: %d\n", fs.Checksum())

	// compact the file system by file for part 2
	// first we need to reset the file system
	input2, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}
	defer input2.Close()
	fs = ParseDiskMap(input2)
	fs.CompactByFile()

	// print the checksum
	log.Printf("(Part 2) checksum: %d\n", fs.Checksum())

}
