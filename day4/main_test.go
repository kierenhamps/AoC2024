package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validSmallTestGrid = Grid{
		{"F", "F", "F", "F", "F"},
		{"F", "S", "A", "M", "X"},
		{"F", "F", "F", "F", "F"},
		{"F", "F", "F", "F", "F"},
		{"X", "M", "A", "S", "F"},
	}
	validTestGrid = Grid{
		{"M", "M", "M", "S", "X", "X", "M", "A", "S", "M"},
		{"M", "S", "A", "M", "X", "M", "S", "M", "S", "A"},
		{"A", "M", "X", "S", "X", "M", "A", "A", "M", "M"},
		{"M", "S", "A", "M", "A", "S", "M", "S", "M", "X"},
		{"X", "M", "A", "S", "A", "M", "X", "A", "M", "M"},
		{"X", "X", "A", "M", "M", "X", "X", "A", "M", "A"},
		{"S", "M", "S", "M", "S", "A", "S", "X", "S", "S"},
		{"S", "A", "X", "A", "M", "A", "S", "A", "A", "A"},
		{"M", "A", "M", "M", "M", "X", "M", "M", "M", "M"},
		{"M", "X", "M", "X", "A", "X", "M", "A", "S", "X"},
	}
	invalidTestGrid = Grid{
		{"M", "M", "M", "S", "X", "X", "M", "A", "S", "M"},
		{"M", "S", "A", "M", "X", "M", "S", "M", "S", "A"},
		{"A", "M", "X", "S", "X", "M", "A", "A", "M", "M"},
		{"M", "S", "A", "M", "A", "S", "M", "S", "M", "X"},
	}
)

func TestDay4_WordSearch_NewWordsearch(t *testing.T) {
	tests := []struct {
		name        string
		grid        [][]string
		expected    *WordSearch
		expectedErr error
	}{
		{"valid_grid", validTestGrid, &WordSearch{grid: validTestGrid}, nil},
		{"invalid_grid", invalidTestGrid, &WordSearch{}, ErrInvalidGrid},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewWordSearch(test.grid)
			assert.Equal(t, test.expected, result)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay4_Word_NewWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected *Word
	}{
		{"valid_word", "XMAS", &Word{word: "XMAS", pattern: []Pattern{
			{direction: DirectionEast, coordinates: []Coordinate{{0, 0}, {0, 1}, {0, 2}, {0, 3}}},
			{direction: DirectionSouthEast, coordinates: []Coordinate{{0, 0}, {1, 1}, {2, 2}, {3, 3}}},
			{direction: DirectionSouth, coordinates: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}}},
			{direction: DirectionSouthWest, coordinates: []Coordinate{{0, 0}, {1, -1}, {2, -2}, {3, -3}}},
			{direction: DirectionWest, coordinates: []Coordinate{{0, 0}, {0, -1}, {0, -2}, {0, -3}}},
			{direction: DirectionNorthWest, coordinates: []Coordinate{{0, 0}, {-1, -1}, {-2, -2}, {-3, -3}}},
			{direction: DirectionNorth, coordinates: []Coordinate{{0, 0}, {-1, 0}, {-2, 0}, {-3, 0}}},
			{direction: DirectionNorthEast, coordinates: []Coordinate{{0, 0}, {-1, 1}, {-2, 2}, {-3, 3}}},
		}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewWord(test.word)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestDay4_Grid_FindPatternAtCoodinate(t *testing.T) {
	grid := validTestGrid
	tests := []struct {
		name       string
		word       string
		pattern    Pattern
		atLocation Coordinate
		expected   bool
	}{
		{"match word in position", "XMAS", Pattern{direction: DirectionEast, coordinates: []Coordinate{{0, 0}, {0, 1}, {0, 2}, {0, 3}}}, Coordinate{0, 5}, true},
		{"no match", "TR", Pattern{direction: DirectionEast, coordinates: []Coordinate{{0, 0}, {0, 1}}}, Coordinate{0, 0}, false},
		{"match out of bounds", "MISS", Pattern{direction: DirectionWest, coordinates: []Coordinate{{0, 0}, {0, -1}, {0, -2}, {0, -3}}}, Coordinate{0, 0}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, grid.IsPatternAt(test.word, test.pattern, test.atLocation))
		})
	}
}

func TestDay4_WordSearch_FindWord(t *testing.T) {
	tests := []struct {
		name     string
		word     *Word
		grid     Grid
		expected *[]Match
	}{
		{"match word", NewWord("XMAS"), validSmallTestGrid, &[]Match{{DirectionWest, Coordinate{1, 4}}, {DirectionEast, Coordinate{4, 0}}}},
		{"no match", NewWord("TR"), validSmallTestGrid, &[]Match{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wordSearch, _ := NewWordSearch(test.grid)
			assert.Equal(t, test.expected, wordSearch.FindWord(test.word))
		})
	}
}

func TestDay4_createGrid(t *testing.T) {
	tests := []struct {
		name        string
		input       io.Reader
		expected    Grid
		expectedErr error
	}{
		{"valid_input", strings.NewReader("AAA\nBBB\nCCC"), Grid{
			{"A", "A", "A"},
			{"B", "B", "B"},
			{"C", "C", "C"},
		}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := createGrid(test.input)
			assert.Equal(t, test.expected, result)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}