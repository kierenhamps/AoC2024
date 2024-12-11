package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

const (
	// The rune that represents a space in the grid
	FrequencyMapEmptySpace = '.'
)

// AntinodeFinder is an interface for finding antinodes
type AntinodeFinder interface {
	FindAntinodes(p1, p2 Point, maxX, maxY int) []Point
}

// SimpleAntinodeFinder is a simple implementation of AntinodeFinder
type SimpleAntinodeFinder struct{}

// FindAntinodes returns all points that are antinodes for the given point projected forward
func (saf SimpleAntinodeFinder) FindAntinodes(p1, p2 Point, maxX, maxY int) []Point {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	newX := p2.x + dx
	newY := p2.y + dy
	if newX < 0 || newX >= maxX || newY < 0 || newY >= maxY {
		return []Point{}
	}
	return []Point{NewPoint(newX, newY)}
}

// HarmonicAntinodeFinder is a more complex implementation of AntinodeFinder that takes into
// account the effects of resonant harmonics
type HarmonicAntinodeFinder struct{}

// FindAntinodes returns all points that are antinodes for the given point projected forward using
// resonant harmonics
func (haf HarmonicAntinodeFinder) FindAntinodes(p1, p2 Point, maxX, maxY int) []Point {
	// Figure out the deltas
	dx := p2.x - p1.x
	dy := p2.y - p1.y

	// Add antenna nodes to list of antinodes to start with
	antinodes := []Point{p1, p2}

	// loop until we fall off the grid
	for i := 1; ; i++ {
		newX := p2.x + dx*i
		newY := p2.y + dy*i
		if newX < 0 || newX >= maxX || newY < 0 || newY >= maxY {
			break
		}
		antinodes = append(antinodes, NewPoint(newX, newY))
	}
	return antinodes
}

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
func (fm *FrequencyMap) FindAllAntinodes(af AntinodeFinder) {
	for _, points := range fm.frequencies {
		for _, p1 := range points {
			for _, p2 := range points {
				if p1 == p2 {
					continue
				}
				antinodes := af.FindAntinodes(p1, p2, fm.maxX, fm.maxY)
				for _, antinode := range antinodes {
					if fm.InBounds(antinode) {
						fm.antinodes[antinode]++
					}
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

func main() {
	// open the test data
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer input.Close()

	// parse the frequency map
	fm := ParseFrequencyMap(input)

	// find all antinodes using the simple antinode finder
	saf := SimpleAntinodeFinder{}
	fm.FindAllAntinodes(saf)

	// Part1 count the antinodes found
	log.Println("Antinodes found:", len(fm.antinodes))

	// find all antinodes using the harmonic antinode finder
	haf := HarmonicAntinodeFinder{}
	fm.FindAllAntinodes(haf)

	// Part2 count the antinodes found
	log.Println("Antinodes found with harmonics:", len(fm.antinodes))

}
