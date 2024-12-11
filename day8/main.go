package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
)

const (
	// The rune that represents a space in the grid
	FrequencyMapEmptySpace = '.'
)

var (
	ErrPointOutOfBounds = errors.New("point is out of bounds")
)

// Frequency is represented by a rune
type Frequency rune

// FrequencyMap reperesents a grid of frequencies at 2D coordinates (points)
type FrequencyMap struct {
	antinodes   map[Point]int
	frequencies map[Frequency][]Point
	maxX        int
	maxY        int
}

// NewFrequencyMap returns a new FrequencyMap
func NewFrequencyMap() *FrequencyMap {
	fm := &FrequencyMap{}
	fm.antinodes = make(map[Point]int)
	fm.frequencies = make(map[Frequency][]Point)
	return fm
}

// AddPoint adds a Point to the FrequencyMap for the given Frequency
func (fm *FrequencyMap) AddPoint(f Frequency, p Point) {
	fm.frequencies[f] = append(fm.frequencies[f], p)
}

// FindAllAntinodes searches and updates the FrequencyMap for all antinodes found
func (fm *FrequencyMap) FindAllAntinodes() {
	for _, points := range fm.frequencies {
		for _, p1 := range points {
			for _, p2 := range points {
				if p1 == p2 {
					continue
				}
				antinode, err := FindAntinode(p1, p2)
				if err == nil && fm.InBounds(antinode) {
					fm.antinodes[antinode]++
				}
			}
		}
	}
}

// InBounds returns true if the given Point is within the bounds of the FrequencyMap
func (fm *FrequencyMap) InBounds(p Point) bool {
	return p.x >= 0 && p.x < fm.maxX && p.y >= 0 && p.y < fm.maxY
}

// Point is a 2D coordinate
type Point struct {
	x int
	y int
}

// NewPoint returns a new Point ValueObject
func NewPoint(x, y int) Point {
	return Point{x, y}
}

// ParseFrequencyMap reads a 2D grid from an io.Reader and returns a FrequencyMap
func ParseFrequencyMap(r io.Reader) *FrequencyMap {
	fm := NewFrequencyMap()
	scanner := bufio.NewScanner(r)
	var x int
	var y int
	var c rune
	for y = 0; scanner.Scan(); y++ {
		line := scanner.Text()
		for x, c = range line {
			if c != FrequencyMapEmptySpace {
				// Add the point to the FrequencyMap
				fm.AddPoint(Frequency(c), NewPoint(x, y))
			}
		}
	}
	fm.maxX = x + 1
	fm.maxY = y
	return fm
}

// FindAntinode returns a point that is the antinode of the given point projected forward
func FindAntinode(p1, p2 Point) (Point, error) {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	newX := p2.x + dx
	newY := p2.y + dy
	if newX < 0 || newY < 0 {
		return Point{}, ErrPointOutOfBounds
	}
	return NewPoint(p2.x+dx, p2.y+dy), nil
}

func main() {
	// open the test data
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer input.Close()

	// parse the frequency map
	fm := ParseFrequencyMap(input)

	// find all antinodes
	fm.FindAllAntinodes()

	// Part1 count the antinodes found
	log.Println("Antinodes found:", len(fm.antinodes))

}
