package main

import (
	"io"
	"log"
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
				freeBlocks: Pointer{
					ids: []int{9, 10, 11, 12, 13, 14},
				},
				dataBlocks: Pointer{
					ids: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
				blocks: map[int]int{
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
				freeBlocks: Pointer{
					ids: []int{1, 2, 6, 7, 8, 9},
				},
				dataBlocks: Pointer{
					ids: []int{0, 3, 4, 5, 10, 11, 12, 13, 14},
				},
				blocks: map[int]int{
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
				freeBlocks: Pointer{
					ids: []int{9, 10, 11, 12, 13, 14},
				},
				dataBlocks: Pointer{
					ids: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
				blocks: map[int]int{
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
			log.Println(test.expected.Equals(
				test.fs))
			assert.True(t, test.expected.Equals(test.fs))
		})
	}
}

func TestDay9_Pointer_NewPointer(t *testing.T) {
	p := NewPointer()
	assert.NotNil(t, p)
}

func TestDay9_Pointer_PeekFirst(t *testing.T) {
	p := NewPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	peek, ok := p.PeekFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, peek)
}

func TestDay9_Pointer_PeekLast(t *testing.T) {
	p := NewPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)

	peek, ok := p.PeekLast()
	assert.True(t, ok)
	assert.Equal(t, 3, peek)
}

func TestDay9_Pointer_Push(t *testing.T) {
	p := NewPointer()
	p.Push(1)
	p.Push(2)
	p.Push(3)
	assert.Equal(t, []int{1, 2, 3}, p.ids)
}

func TestDay9_Pointer_PopFirst(t *testing.T) {
	p := NewPointer()
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

func TestDay9_Pointer_PopLast(t *testing.T) {
	p := NewPointer()
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

func TestDay9_Pointer_Sort(t *testing.T) {
	p := NewPointer()
	p.Push(3)
	p.Push(1)
	p.Push(2)

	p.Sort()

	assert.Equal(t, []int{1, 2, 3}, p.ids)
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
				freeBlocks: Pointer{
					ids: []int{1, 2, 6, 7, 8, 9},
				},
				dataBlocks: Pointer{
					ids: []int{0, 3, 4, 5, 10, 11, 12, 13, 14},
				},
				blocks: map[int]int{
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
			assert.True(t, test.expected.Equals(fs))
		})
	}
}
