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

// PositionMap represents a map of locations to directions
// Locations are keys and must be unique allowing for multiple directions at each location
type PositionMap map[Location][]Direction

// Guard represents a guard that can patrol a map
type Guard struct {
	positionsVisited PositionMap
	location         Location
	direction        Direction
}

// NewGuard creates a new guard with a starting location and direction
func NewGuard(l Location, d Direction) *Guard {
	return &Guard{make(PositionMap), l, d}
}

// Move moves the guard to the new location
func (g *Guard) Move(l Location) {
	g.location = l
}

// AddVisitedPosition adds the current position to the list of visited positions
func (g *Guard) AddVisitedPosition() {
	g.positionsVisited[g.location] = append(g.positionsVisited[g.location], g.direction)
}

// CurrentDirection returns the guards current direction
func (g *Guard) CurrentDirection() Direction {
	return g.direction
}

// CurrentLocation returns the guards current location
func (g *Guard) CurrentLocation() Location {
	return g.location
}

// PositionsVisited returns the positions visited by the guard
func (g *Guard) PositionsVisited() PositionMap {
	return g.positionsVisited
}

// DejaVu returns true if the guard has visited its current position before (in the same direction)
// This is used to detect loops
func (g *Guard) DejaVu() bool {
	if _, ok := g.positionsVisited[g.location]; ok {
		for _, direction := range g.positionsVisited[g.location] {
			if direction == g.direction {
				return true
			}
		}
	}
	return false
}

// NextPosition returns the next position in front of the guard
func (g *Guard) NextLocation(patrolMap PatrolMap) Location {
	var newLocation Location
	switch g.CurrentDirection() {
	case DirectionUp:
		newLocation = Location{g.location.X(), g.location.Y() - 1}
	case DirectionRight:
		newLocation = Location{g.location.X() + 1, g.location.Y()}
	case DirectionDown:
		newLocation = Location{g.location.X(), g.location.Y() + 1}
	case DirectionLeft:
		newLocation = Location{g.location.X() - 1, g.location.Y()}
	}
	return newLocation
}

// turnRight turns the guard to the right
func (g *Guard) turnRight() {
	var newDirection Direction
	switch g.direction {
	case DirectionUp:
		newDirection = DirectionRight
	case DirectionRight:
		newDirection = DirectionDown
	case DirectionDown:
		newDirection = DirectionLeft
	case DirectionLeft:
		newDirection = DirectionUp
	}

	g.direction = newDirection
}

// Location represents a location on the map
type Location struct {
	x int
	y int
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
type PatrolMap map[Location]Space

// ParseInput parses the input and returns a patrol map and a guard
//
// A free space is represented by a "."
// An obstacle is represented by a "#"
// The guard is represented by "^" and is facing up
func ParseInput(input io.Reader) (PatrolMap, *Guard, error) {
	patrolMap := make(PatrolMap)
	var guard *Guard

	scanner := bufio.NewScanner(input)
	for y := 0; scanner.Scan(); y++ {
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
				guard = NewGuard(Location{x, y}, DirectionUp)
			}
			patrolMap[Location{x, y}] = spaceType
		}
	}

	return patrolMap, guard, nil
}

// Free returns true if the location is a free space
func (pm PatrolMap) Free(l Location) bool {
	if _, ok := pm[l]; ok {
		return pm[l] == SpaceFree
	}
	return false
}

// OnMap returns true if the guard is on the map
func (pm PatrolMap) OnMap(l Location) bool {
	_, ok := pm[l]
	return ok
}

// Patrol returns the number of distinct locations visited by the guard
// before it leaves the map
//
// nil will be returned if the guard gets stuck in a loop (has seen this road before!)
func (pm *PatrolMap) Patrol(l Location, d Direction, end Location) PositionMap {
	// setup the guard
	guard := NewGuard(l, d)
	for {
		// if off map return list of found positions
		if !pm.OnMap(guard.CurrentLocation()) {
			return guard.PositionsVisited()
		}
		// if dejavu return nil (for loop detection)
		if guard.DejaVu() {
			return nil
		}
		// Track visited position
		guard.AddVisitedPosition()
		// get next location
		newLocation := guard.NextLocation(*pm)
		// if obstacle in front or the given location (the given pretend loop obstruction) then change direction
		if !pm.Free(newLocation) || newLocation == end {
			if pm.OnMap(newLocation) {
				guard.turnRight()
			} else {
				guard.Move(newLocation)
			}
		} else {
			// otherwise move forward
			guard.Move(newLocation)
		}
	}
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

	// Save start location and direction for part 2
	startLocation := guard.CurrentLocation()
	startDirection := guard.CurrentDirection()

	// Patrol till we leave the map
	visitedPositions := patrolMap.Patrol(startLocation, startDirection, Location{-1, -1})

	// count the number of locations visited for part 1
	log.Printf("(PART 1) The guard visited %d distinct locations", len(visitedPositions))

	// Loop through every step we took to get through the map and see if we can add an obstruction
	// to create a loop
	part2 := 0
	for visitedLocation := range visitedPositions {
		if patrolMap.Patrol(startLocation, startDirection, visitedLocation) == nil {
			part2++
		}
	}
	log.Printf("(PART 2) We could add %d obstructions to create a loop", part2)
}
