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
	DirectionUp = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

const (
	SpaceFree = iota
	SpaceFailedSuitPrototypes
	SpaceSpoolOfVeryLongPolymers
	SpaceCrates
	SpaceTankOfUniversalSolvent
)

var (
	ErrInputCannotBeNegative = errors.New("input cannot be negative")
	ErrInvalidGuardInput     = errors.New("invalid guard input")
	ErrInvalidPatrolMapInput = errors.New("invalid patrol map input")
)

type Direction int

// Guard represents a guard that can patrol a map
type Guard struct {
	locationsVisited []Location
	location         Location
	direction        Direction
}

// NewGuard creates a new guard with a starting location and direction
func NewGuard(location Location, direction Direction) *Guard {
	return &Guard{[]Location{location}, location, direction}
}

// changeDirection changes the direction of the guard
func (g *Guard) changeDirection(direction Direction) {
	g.direction = direction
}

// changeLocation changes the location of the guard
func (g *Guard) changeLocation(location Location) {
	g.locationsVisited = append(g.locationsVisited, location)
	g.location = location
}

// LocationsVisited returns the locations visited by the guard
func (g *Guard) LocationsVisited() []Location {
	return g.locationsVisited
}

// Move moves the guard in the direction it is facing
//
// The guard will move in the direction it is facing until it hits an obstacle
// or leaves the map. If the guard hits an obstacle, it will turn right and try
// to move again. If the guard leaves the map, it will stop moving.
func (g *Guard) Move(patrolMap PatrolMap) {
	// if we are not on the map, dont move
	if !patrolMap.OnMap(g.location) {
		return
	}

	var newLocation Location
	switch g.direction {
	case DirectionUp:
		newLocation = NewLocation(g.location.x, g.location.y-1)
	case DirectionRight:
		newLocation = NewLocation(g.location.x+1, g.location.y)
	case DirectionDown:
		newLocation = NewLocation(g.location.x, g.location.y+1)
	case DirectionLeft:
		newLocation = NewLocation(g.location.x-1, g.location.y)
	}

	// if the new location is free, move there
	if patrolMap.Free(newLocation) {
		g.changeLocation(newLocation)
		g.changeDirection(g.direction)
		return
	}
	// otherwise turn right and try again
	g.turnRight()
	g.Move(patrolMap)
}

// turnRight turns the guard to the right
func (g *Guard) turnRight() {
	switch g.direction {
	case DirectionUp:
		g.changeDirection(DirectionRight)
	case DirectionRight:
		g.changeDirection(DirectionDown)
	case DirectionDown:
		g.changeDirection(DirectionLeft)
	case DirectionLeft:
		g.changeDirection(DirectionUp)
	}
}

// Location represents a location on the map
type Location struct {
	x int
	y int
}

// NewLocation creates a new location
func NewLocation(x, y int) Location {
	return Location{x, y}
}

// X returns the x coordinate of the location
func (l Location) X() int {
	return l.x
}

// Y returns the y coordinate of the location
func (l Location) Y() int {
	return l.y
}

// PatrolMap represents a map that a guard can patrol
type PatrolMap [][]Space

// ParseInput parses the input and returns a patrol map and a guard
//
// A free space is represented by a "."
// An obstacle is represented by a "#"
// The guard is represented by "^", ">", "v", or "<" which denotes the direction
func ParseInput(input io.Reader) (PatrolMap, *Guard, error) {
	var patrolMap PatrolMap
	var guard *Guard

	scanner := bufio.NewScanner(input)
	for y := 0; scanner.Scan(); y++ {
		var patrolRow []Space
		spaces := strings.Split(scanner.Text(), "")
		for x, space := range spaces {
			var spaceType Space
			switch space {
			case ".":
				spaceType = SpaceFree
			case "#":
				spaceType = SpaceCrates
			case "^":
				spaceType = SpaceFree
				guard = NewGuard(NewLocation(x, y), DirectionUp)
			case ">":
				spaceType = SpaceFree
				guard = NewGuard(NewLocation(x, y), DirectionRight)
			case "v":
				spaceType = SpaceFree
				guard = NewGuard(NewLocation(x, y), DirectionDown)
			case "<":
				spaceType = SpaceFree
				guard = NewGuard(NewLocation(x, y), DirectionLeft)
			default:
				return PatrolMap{}, &Guard{}, ErrInvalidPatrolMapInput
			}
			patrolRow = append(patrolRow, spaceType)
		}
		patrolMap = append(patrolMap, patrolRow)
	}

	return patrolMap, guard, nil
}

// Free returns true if the location is a free space
func (pm PatrolMap) Free(l Location) bool {
	if !pm.OnMap(l) {
		return true
	}
	return pm[l.Y()][l.X()] == SpaceFree
}

// OnMap returns true if the location is on the map
func (pm PatrolMap) OnMap(l Location) bool {
	if l.X() < 0 || l.Y() < 0 {
		return false
	}
	if l.X() >= len(pm) || l.Y() >= len(pm[0]) {
		return false
	}
	return true
}

// Space represents a space on the map
type Space int

func main() {
	// Open the input file
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer input.Close()

	// parse the input
	patrolMap, guard, err := ParseInput(input)
	if err != nil {
		log.Fatalf("could not parse input: %v", err)
	}

	// move the guard until it leaves the map
	for patrolMap.OnMap(guard.location) {
		guard.Move(patrolMap)
	}

	// Get a distinct list of locations visited
	locationsVisited := guard.LocationsVisited()
	distinctLocations := make(map[Location]struct{})
	for _, location := range locationsVisited {
		distinctLocations[location] = struct{}{}
	}

	// print the number of DISTINCT locations visited
	// minus one location for the last location which is off the map
	log.Printf("(PART 1) The guard visited %d distinct locations", len(distinctLocations)-1)

}
