package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

const (
	DirectionEast Direction = iota
	DirectionSouth
	DirectionSouthEast
	DirectionWest
	DirectionSouthWest
	DirectionNorth
	DirectionNorthEast
	DirectionNorthWest
)

var (
	ErrInvalidGrid = errors.New("invalid grid")
)

// Direction indicates the direction of the pattern
type Direction int

// Grid is a 2D array of strings that represents a word search puzzle
type Grid [][]string

// IsPatternAt checks if the pattern is at the location in the grid
func (g Grid) IsPatternAt(word string, pattern Pattern, atLocation Coordinate) bool {
	// Check if the pattern is at the location
	for i, c := range pattern.coordinates {
		if atLocation.x+c.x < 0 || atLocation.x+c.x >= len(g) || atLocation.y+c.y < 0 || atLocation.y+c.y >= len(g[0]) {
			return false
		}
		if g[atLocation.x+c.x][atLocation.y+c.y] != string(word[i]) {
			return false
		}
	}
	return true
}

// Coordinate represents a location in 2D space
type Coordinate struct {
	x int
	y int
}

// Match represents a match found in the grid
type Match struct {
	direction Direction
	location  Coordinate
}

// Pattern represents a place that a word can take up in a 2D space
type Pattern struct {
	direction   Direction
	coordinates []Coordinate
}

// Word represents a word that can be found in a word search puzzle
type Word struct {
	word    string
	pattern []Pattern
}

// NewWord creates a new Word and all possible patterns for the word
func NewWord(w string) *Word {
	word := &Word{word: w}

	// Calculate all possible Positions
	word.createPattern(DirectionEast)
	word.createPattern(DirectionSouthEast)
	word.createPattern(DirectionSouth)
	word.createPattern(DirectionSouthWest)
	word.createPattern(DirectionWest)
	word.createPattern(DirectionNorthWest)
	word.createPattern(DirectionNorth)
	word.createPattern(DirectionNorthEast)

	return word
}

// createPattern creates a pattern for the word based on the direction given
func (w *Word) createPattern(direction Direction) {
	// Create a Pattern based on the direction and start location
	p := Pattern{direction: direction}
	c := []Coordinate{}
	for i := 0; i < len(w.word); i++ {
		switch direction {
		case DirectionEast:
			c = append(c, Coordinate{0, 0 + i})
		case DirectionSouthEast:
			c = append(c, Coordinate{0 + i, 0 + i})
		case DirectionSouth:
			c = append(c, Coordinate{0 + i, 0})
		case DirectionSouthWest:
			c = append(c, Coordinate{0 + i, 0 - i})
		case DirectionWest:
			c = append(c, Coordinate{0, 0 - i})
		case DirectionNorthWest:
			c = append(c, Coordinate{0 - i, 0 - i})
		case DirectionNorth:
			c = append(c, Coordinate{0 - i, 0})
		case DirectionNorthEast:
			c = append(c, Coordinate{0 - i, 0 + i})
		}
	}
	p.coordinates = c
	w.pattern = append(w.pattern, p)
}

// WordSearch represents a word search puzzle
type WordSearch struct {
	grid Grid
}

// NewWordSearch creates a new WordSearch
func NewWordSearch(g Grid) (*WordSearch, error) {
	// if length of one dimension is not the same as the other, return an error
	log.Println("Creating wordsearch for grid of size:", len(g), "x", len(g[0]))
	if len(g) != len(g[0]) {
		return &WordSearch{}, ErrInvalidGrid
	}
	return &WordSearch{grid: g}, nil
}

// FindWord finds all occurances of a word in a grid
func (ws *WordSearch) FindWord(w *Word) *[]Match {
	log.Println("Finding word:", w.word)
	matches := &[]Match{}
	for i := 0; i < len(ws.grid); i++ {
		for j := 0; j < len(ws.grid[i]); j++ {
			for _, p := range w.pattern {
				if ws.grid.IsPatternAt(w.word, p, Coordinate{i, j}) {
					match := Match{direction: p.direction, location: Coordinate{i, j}}
					*matches = append(*matches, match)
				}
			}
		}
	}
	log.Println("Matches found for", w.word, ":", len(*matches))
	return matches
}

func main() {
	// read the inputs
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Create the Grid
	grid, err := createGrid(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create the WordSearch
	ws, err := NewWordSearch(grid)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Word
	word := NewWord("XMAS")

	// Find the Word
	ws.FindWord(word)

}

// createGrid creates a grid from an input file
func createGrid(inputFile io.Reader) (Grid, error) {
	scanner := bufio.NewScanner(inputFile)
	grid := Grid{}
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		grid = append(grid, row)
	}
	return grid, nil
}
