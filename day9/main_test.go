package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay9_FileSystem_NewFileSystem(t *testing.T) {
	fs := NewFileSystem()
	assert.NotNil(t, fs)
}

func TestDay9_FileSystem_Checksum(t *testing.T) {
	tests := []struct {
		name     string
		fs       *FileSystem
		expected int
	}{
		{
			name: "checksum for part 1 simple example",
			fs: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0},
						1: {3, 4, 5},
						2: {1, 2, 6, 7, 8},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{9, 10, 11, 12, 13, 14},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
				blocks: Blocks{
					0: 0,
					1: 2,
					2: 2,
					3: 1,
					4: 1,
					5: 1,
					6: 2,
					7: 2,
					8: 2,
				},
			},
			expected: 60,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.fs.Checksum())
		})
	}
}

func TestDay9_FileSystem_Compact(t *testing.T) {
	tests := []struct {
		name     string
		fs       *FileSystem
		expected *FileSystem
	}{
		{
			name: "compact file system for part 1 simple example",
			fs: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0},
						1: {3, 4, 5},
						2: {10, 11, 12, 13, 14},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{1, 2, 6, 7, 8, 9},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 3, 4, 5, 10, 11, 12, 13, 14},
				},
				blocks: Blocks{
					0:  0,
					3:  1,
					4:  1,
					5:  1,
					10: 2,
					11: 2,
					12: 2,
					13: 2,
					14: 2,
				},
			},
			expected: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0},
						1: {3, 4, 5},
						2: {8, 7, 6, 2, 1},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{9, 10, 11, 12, 13, 14},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
				blocks: Blocks{
					0: 0,
					1: 2,
					2: 2,
					3: 1,
					4: 1,
					5: 1,
					6: 2,
					7: 2,
					8: 2,
				}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fs.Compact()
			assert.Equal(t, test.expected.blocks, test.fs.blocks)
		})
	}
}

func TestDay9_FileSystem_CompactByFile(t *testing.T) {
	tests := []struct {
		name     string
		fs       *FileSystem
		expected *FileSystem
	}{
		{
			name: "compact file system for part 2 by file example",
			fs: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0, 1},
						1: {5, 6, 7},
						2: {11},
						3: {15, 16, 17},
						4: {19, 20},
						5: {22, 23, 24, 25},
						6: {27, 28, 29, 30},
						7: {32, 33, 34},
						8: {36, 37, 38, 39},
						9: {40, 41},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{2, 3, 4, 8, 9, 10, 12, 13, 14, 18, 21, 26, 31, 35},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 1, 5, 6, 7, 11, 15, 16, 17, 19, 20, 22, 23, 24, 25, 27, 28, 29, 30, 32, 33, 34, 36, 37, 38, 39, 40, 41},
				},
				blocks: Blocks{
					0:  0,
					1:  0,
					5:  1,
					6:  1,
					7:  1,
					11: 2,
					15: 3,
					16: 3,
					17: 3,
					19: 4,
					20: 4,
					22: 5,
					23: 5,
					24: 5,
					25: 5,
					27: 6,
					28: 6,
					29: 6,
					30: 6,
					32: 7,
					33: 7,
					34: 7,
					36: 8,
					37: 8,
					38: 8,
					39: 8,
					40: 9,
					41: 9,
				},
			},
			expected: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0, 1},
						1: {5, 6, 7},
						2: {4},
						3: {15, 16, 17},
						4: {12, 13},
						5: {22, 23, 24, 25},
						6: {27, 28, 29, 30},
						7: {8, 9, 10},
						8: {36, 37, 38, 39},
						9: {2, 3},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{11, 14, 18, 19, 20, 21, 26, 31, 32, 33, 34, 35, 40, 41},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 15, 16, 17, 22, 23, 24, 25, 27, 28, 29, 30, 36, 37, 38, 39},
				},
				blocks: Blocks{
					0:  0,
					1:  0,
					2:  9,
					3:  9,
					4:  2,
					5:  1,
					6:  1,
					7:  1,
					8:  7,
					9:  7,
					10: 7,
					12: 4,
					13: 4,
					15: 3,
					16: 3,
					17: 3,
					22: 5,
					23: 5,
					24: 5,
					25: 5,
					27: 6,
					28: 6,
					29: 6,
					30: 6,
					36: 8,
					37: 8,
					38: 8,
					39: 8,
				}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fs.CompactByFile()

			assert.Equal(t, test.expected.blocks, test.fs.blocks)
		})
	}
}

func TestDay9_BlockPointer_NewBlockPointer(t *testing.T) {
	p := NewBlockPointer()
	assert.NotNil(t, p)
}

func TestDay9_BlockPointer_Find(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	found, ok := p.Find(2)
	assert.True(t, ok)
	assert.Equal(t, 2, found)
}

func TestDay9_BlockPointer_IsContiguous(t *testing.T) {
	tests := []struct {
		name     string
		p        *BlockPointer
		startID  int
		endID    int
		expected bool
	}{
		{
			name: "contiguous blocks",
			p: &BlockPointer{
				ids: []int{1, 2, 3},
			},
			startID:  1,
			endID:    3,
			expected: true,
		},
		{
			name: "non-contiguous blocks",
			p: &BlockPointer{
				ids: []int{1, 2, 4},
			},
			startID:  2,
			endID:    4,
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.p.IsContiguous(test.startID, test.endID))
		})
	}
}

func TestDay9_BlockPointer_PeekFirst(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	peek, ok := p.PeekFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, peek)
}

func TestDay9_BlockPointer_PeekLast(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	peek, ok := p.PeekLast()
	assert.True(t, ok)

	assert.Equal(t, 3, peek)
}

func TestDay9_BlockPointer_Push(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)
	assert.Equal(t, []int{1, 2, 3}, p.ids)
}

func TestDay9_BlockPointer_PopFirst(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	pop, ok := p.PopFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, pop)

	pop, ok = p.PopFirst()
	assert.True(t, ok)
	assert.Equal(t, 2, pop)

	assert.Equal(t, p.ids, []int{3})
}

func TestDay9_BlockPointer_PopLast(t *testing.T) {
	p := NewBlockPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	pop, ok := p.PopLast()
	assert.True(t, ok)
	assert.Equal(t, 3, pop)

	pop, ok = p.PopLast()
	assert.True(t, ok)
	assert.Equal(t, 2, pop)

	assert.Equal(t, p.ids, []int{1})
}

func TestDay9_BlockPointer_Sort(t *testing.T) {
	p := NewBlockPointer()
	p.Push(3)
	p.Push(1)
	p.Push(2)

	p.Sort()

	assert.Equal(t, []int{1, 2, 3}, p.ids)
}

func TestDay9_FilePointer_NewFilePointer(t *testing.T) {
	p := NewFilePointer()
	assert.NotNil(t, p)
}

func TestDay9_FilePointer_AppendToFile(t *testing.T) {
	tests := []struct {
		name     string
		p        *FilePointer
		file     int
		blockIDs []int
		expected *FilePointer
	}{
		{
			name: "append to file that already exists",
			p: &FilePointer{
				files: map[int][]int{
					0: {0, 1},
				},
			},
			file: 0,
			blockIDs: []int{
				2,
			},
			expected: &FilePointer{
				files: map[int][]int{
					0: {0, 1, 2},
				},
			},
		},
		{
			name: "append to a file that doesnt yet exist",
			p:    NewFilePointer(),
			file: 1,
			blockIDs: []int{
				3, 4, 5,
			},
			expected: &FilePointer{
				files: map[int][]int{
					1: {3, 4, 5},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.p.AppendToFile(test.file, test.blockIDs)
			assert.Equal(t, test.expected.files, test.p.files)
		})
	}
}

func TestDay9_FilePointer_Equals(t *testing.T) {
	p1 := NewFilePointer()
	p1.Push(map[int][]int{0: {0}})
	p1.Push(map[int][]int{1: {1, 2}})
	p1.Push(map[int][]int{2: {3, 4, 5}})

	p2 := NewFilePointer()
	p2.Push(map[int][]int{0: {0}})
	p2.Push(map[int][]int{1: {1, 2}})
	p2.Push(map[int][]int{2: {3, 4, 5}})

	assert.Equal(t, p1, p2)
}
func TestDay9_FilePointer_MoveBlock(t *testing.T) {
	tests := []struct {
		name        string
		p           *FilePointer
		fileID      int
		fromBlockID int
		toBlockID   int
		expected    *FilePointer
	}{
		{
			name: "move block in one file from one location to another",
			p: &FilePointer{
				files: map[int][]int{
					0: {0, 1, 2},
				},
			},
			fileID:      0,
			fromBlockID: 1,
			toBlockID:   3,
			expected: &FilePointer{
				files: map[int][]int{
					0: {0, 3, 2},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.p.MoveBlock(test.fileID, test.fromBlockID, test.toBlockID)
			assert.Equal(t, test.expected.files, test.p.files)
		})
	}
}

func TestDay9_FilePointer_PeekLast(t *testing.T) {
	p := NewFilePointer()
	p.Push(map[int][]int{0: {0}})
	p.Push(map[int][]int{1: {1, 2}})
	p.Push(map[int][]int{2: {3, 4, 5}})

	peekID, peekBlocks, ok := p.PeekLast()
	assert.True(t, ok)
	assert.Equal(t, 2, peekID)
	assert.Equal(t, []int{3, 4, 5}, peekBlocks)
}

func TestDay9_FilePointer_Push(t *testing.T) {
	p := NewFilePointer()
	p.Push(map[int][]int{0: {0}})
	p.Push(map[int][]int{1: {1, 2}})
	p.Push(map[int][]int{2: {3, 4, 5}})

	assert.Equal(t, map[int][]int{0: {0}, 1: {1, 2}, 2: {3, 4, 5}}, p.files)
}

func TestDay9_ParseDiskMap(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected *FileSystem
	}{
		{
			name:  "parse disk map for part 1 simple example",
			input: strings.NewReader("12345\n"),
			expected: &FileSystem{
				files: FilePointer{
					files: map[int][]int{
						0: {0},
						1: {3, 4, 5},
						2: {10, 11, 12, 13, 14},
					},
				},
				freeBlocks: BlockPointer{
					ids: []int{1, 2, 6, 7, 8, 9},
				},
				dataBlocks: BlockPointer{
					ids: []int{0, 3, 4, 5, 10, 11, 12, 13, 14},
				},
				blocks: Blocks{
					0:  0,
					3:  1,
					4:  1,
					5:  1,
					10: 2,
					11: 2,
					12: 2,
					13: 2,
					14: 2,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := ParseDiskMap(test.input)
			assert.Equal(t, test.expected.blocks, fs.blocks)
		})
	}
}
