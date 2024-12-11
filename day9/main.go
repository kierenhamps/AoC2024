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

// FileSystem represents a file system
type FileSystem struct {
	freeBlocks Pointer
	dataBlocks Pointer
	blocks     map[int]int
}

// NewFileSystem returns a new FileSystem
func NewFileSystem() *FileSystem {
	return &FileSystem{
		blocks: make(map[int]int, 0),
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

// Equals compares two FileSystems
func (fs *FileSystem) Equals(other *FileSystem) bool {
	if len(fs.blocks) != len(other.blocks) {
		return false
	}
	for k, v := range fs.blocks {
		if other.blocks[k] != v {
			return false
		}
	}
	if !fs.freeBlocks.Equals(other.freeBlocks) {
		return false
	}
	if !fs.dataBlocks.Equals(other.dataBlocks) {
		return false
	}
	return true
}

// Pointer represents a stack of IDs to track blocks
type Pointer struct {
	ids []int
}

// NewPointer returns a new Pointer
func NewPointer() *Pointer {
	return &Pointer{
		ids: make([]int, 0),
	}
}

// Equals compares two Pointers
func (p *Pointer) Equals(other Pointer) bool {
	if len(p.ids) != len(other.ids) {
		return false
	}
	for k, v := range p.ids {
		if other.ids[k] != v {
			return false
		}
	}
	return true
}

// PeekFirst returns the first ID from the stack without removing it
func (p *Pointer) PeekFirst() (int, bool) {
	if len(p.ids) == 0 {
		return 0, false
	}
	return p.ids[0], true
}

// PeekLast returns the last ID from the stack without removing it
func (p *Pointer) PeekLast() (int, bool) {
	if len(p.ids) == 0 {
		return 0, false
	}
	return p.ids[len(p.ids)-1], true
}

// PopFirst removes the first ID from the stack and returns it
func (p *Pointer) PopFirst() (int, bool) {
	if len(p.ids) == 0 {
		return 0, false
	}
	// get the first ID
	first := p.ids[0]

	// remove the first ID
	p.ids = p.ids[1:]

	return first, true
}

// PopLast removes the last ID from the stack and returns it
func (p *Pointer) PopLast() (int, bool) {
	if len(p.ids) == 0 {
		return 0, false
	}
	// get the last ID
	last := p.ids[len(p.ids)-1]

	// remove the last ID
	p.ids = p.ids[:len(p.ids)-1]

	return last, true
}

// Push adds an ID to the end of the stack
func (p *Pointer) Push(id int) {
	p.ids = append(p.ids, id)
	p.Sort()
}

// Sort sorts the stack of IDs
func (p *Pointer) Sort() {
	sort.Ints(p.ids)
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
}
