package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	ErrNoPossibleLocation = errors.New("no possible location available")
)

// Grid is a 2D map of heights
type Grid map[int]map[int]Height

// Height is a topographical height in the grid
type Height int

// Location is a point on the grid
type Location struct {
	x, y int
}

// Path is a list of locations that make up a trail
type Path []Location

// Trail represents a trail through the TrailMap
//
// Always starting at a trailhead (height 0) and ending at height 9
// Tails can only move N, E, S, W, and never diagonally
// Trails can only travel up 1 height at a time
// Each unique path through the TrailMap from the same origin is a
// extra score
type Trail struct {
	start  Location
	paths  []Path
	score  int
	rating int
}

// NewTrail creates a new Trail starting at the given location
func NewTrail(start Location) Trail {
	return Trail{
		start: start,
		paths: []Path{},
		score: 0,
	}
}

// Walk attempts to walk through the given TrailMap
//
// starting at the given location and ending when it reaches a
// height of 9 or a dead end on all paths
func Walk(l Location, tm TrailMap, p Path) []Path {
	// get height of current location
	h := tm.trailMap[l.y][l.x]

	// Add our current location to the path
	p = append(p, l)

	// if height is 9, we are done
	if h == 9 {
		return []Path{p}
	}

	// Find the next possible locations to move to
	paths := []Path{}
	nextSteps, err := NextStep(l, h, tm.trailMap)
	if err != nil {
		// No possible location to move to
		// return an empty path
		return []Path{}
	}
	for _, nextLocation := range nextSteps {
		// Create a fresh copy of the path
		newPath := make(Path, len(p))
		copy(newPath, p)

		// Walk the path from the next location
		currentPaths := Walk(nextLocation, tm, newPath)

		// Add the path to the list of paths
		// if valid
		if len(currentPaths) > 0 {
			paths = append(paths, currentPaths...)
		}
	}

	return paths
}

// TrailMap is a representation of a map of trails in a grid
// each trail is a path that originates at a trailhead (height 0)
type TrailMap struct {
	trailMap   Grid
	trailheads []Location
}

// DiscoverTrails finds all possible trails through the TrailMap
// returning a list of Trails found
func (tm TrailMap) DiscoverTrails() []Trail {
	trails := []Trail{}
	// order of trails is not important
	for _, th := range tm.trailheads {
		// Create a new trail starting at this trailhead
		trail := NewTrail(th)

		// Walk through the trailMap to discover the paths from this trailhead
		paths := Walk(th, tm, Path{})

		// if no paths were found, continue to the next trailhead
		if len(paths) == 0 {
			continue
		}

		// Add the discovered trail to our list
		trail.paths = paths

		// determine the score of the trail
		// by counting the number of unique ends (using a map)
		uniqueEnds := make(map[Location]struct{})
		for _, path := range paths {
			uniqueEnds[path[len(path)-1]] = struct{}{}
		}
		trail.score = len(uniqueEnds)

		// Add a rating for Part 2
		trail.rating = len(paths)

		trails = append(trails, trail)
	}
	return trails
}

// NextStep provides a list of possible locations to move to from the given location
// based on the current height and the grid
//
// The grid given can be a subset of the full grid, only containing the current
// trail being walked, but must contain the locations of the adjacent heights to
// the current location.
func NextStep(l Location, h Height, g Grid) ([]Location, error) {
	north := Location{l.x, l.y - 1}
	east := Location{l.x + 1, l.y}
	south := Location{l.x, l.y + 1}
	west := Location{l.x - 1, l.y}
	nextHeight := h + 1

	// Find all possible locations to move to
	locations := []Location{}
	if _, ok := g[north.y][north.x]; ok && g[north.y][north.x] == nextHeight {
		locations = append(locations, north)
	}
	if _, ok := g[east.y][east.x]; ok && g[east.y][east.x] == nextHeight {
		locations = append(locations, east)
	}
	if _, ok := g[south.y][south.x]; ok && g[south.y][south.x] == nextHeight {
		locations = append(locations, south)
	}
	if _, ok := g[west.y][west.x]; ok && g[west.y][west.x] == nextHeight {
		locations = append(locations, west)
	}

	if len(locations) == 0 {
		return []Location{}, ErrNoPossibleLocation
	}

	return locations, nil
}

// Parse reads a map from an io.Reader and returns a TrailMap
func Parse(input io.Reader) TrailMap {
	scanner := bufio.NewScanner(input)
	trailMap := TrailMap{
		trailMap:   make(Grid),
		trailheads: []Location{},
	}
	for y := 0; scanner.Scan(); y++ {
		trailMap.trailMap[y] = make(map[int]Height)
		for x, v := range scanner.Text() {
			if v == '.' {
				continue
			}
			height, err := strconv.Atoi(string(v))
			if err != nil {
				log.Fatal(err)
			}
			// Check if this is a trailhead
			if height == 0 {
				trailMap.trailheads = append(trailMap.trailheads, Location{x, y})
			}
			// Add to grid
			trailMap.trailMap[y][x] = Height(height)
		}
	}
	return trailMap
}

func main() {
	// Parse the input
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	trailMap := Parse(input)

	// Discover all trails
	trails := trailMap.DiscoverTrails()

	// Sum all scores for Part 1
	sumOfScores := 0
	for _, trail := range trails {
		sumOfScores += trail.score
	}
	log.Println("Part 1: Sum of scores of all trails on map is", sumOfScores)

	// Sum all ratings for Part 2
	sumOfRatings := 0
	for _, trail := range trails {
		sumOfRatings += trail.rating
	}
	log.Println("Part 2: Sum of ratings of all trails on map is", sumOfRatings)
}
