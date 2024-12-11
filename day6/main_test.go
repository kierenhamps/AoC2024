package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// Small test map that guards can patrol.
	//
	// 0 X 0 0
	// 0 0 0 X
	// 0 0 0 0
	// X 0 X 0
	smallPatrolMap = PatrolMap{
		Location{0, 0}: SpaceFree, Location{1, 0}: SpaceFailedSuitPrototypes, Location{2, 0}: SpaceFree, Location{3, 0}: SpaceFree,
		Location{0, 1}: SpaceFree, Location{1, 1}: SpaceFree, Location{2, 1}: SpaceFree, Location{3, 1}: SpaceSpoolOfVeryLongPolymers,
		Location{0, 2}: SpaceFree, Location{1, 2}: SpaceFree, Location{2, 2}: SpaceFree, Location{3, 2}: SpaceFree,
		Location{0, 3}: SpaceCrates, Location{1, 3}: SpaceFree, Location{2, 3}: SpaceTankOfUniversalSolvent, Location{3, 3}: SpaceFree,
	}

	// Part 2 test maps.
	//
	partTwoPatrolMap = "" +
		"....#.....\n" +
		".........#\n" +
		"..........\n" +
		"..#.......\n" +
		".......#..\n" +
		"..........\n" +
		".#..^.....\n" +
		"........#.\n" +
		"#.........\n" +
		"......#...\n"
	partTwoPatrolMapTest = "" +
		"..........\n" +
		".#........\n" +
		"..........\n" +
		"....#.....\n" +
		".......#..\n" +
		"..........\n" +
		"....^.....\n" +
		"#.........\n" +
		"......#...\n" +
		"..........\n"
	partTwoPatrolMapCloseLoop = "" +
		"..........\n" +
		".#........\n" +
		"..........\n" +
		"....#.....\n" +
		".....#....\n" +
		"..........\n" +
		"....^.....\n" +
		"#.........\n" +
		"....#.....\n" +
		"..........\n"
	partTwoPatrolMapCloseLoop2 = "" +
		"..........\n" +
		"......#...\n" +
		"..........\n" +
		"....#.....\n" +
		".......#..\n" +
		"..........\n" +
		"....^.....\n" +
		".....#....\n" +
		"......#...\n" +
		"..........\n"
	// What is an obstruction would bounce us into a loop that we have already visited?
	partTwoPatrolMapCloseLoop3 = "" +
		"..........\n" +
		"....#.....\n" +
		"........#.\n" +
		".....#....\n" +
		"..........\n" +
		"..........\n" +
		"....^.....\n" +
		"..#.......\n" +
		"....#..#..\n" +
		"..........\n"
	// You cant have an loop formed in a place you need to go through to get to the loop, otherwise you would never get to the loop.
	partTwoPatrolMapNoEntry = "" +
		"..........\n" +
		"....##....\n" +
		".......#..\n" +
		"..........\n" +
		"..........\n" +
		"......#...\n" +
		"....^.....\n" +
		"..........\n" +
		"..........\n" +
		"..........\n"
	partTwoPatrolMapNoEntry2 = "" +
		"##....#...\n" +
		"....#....#\n" +
		"..........\n" +
		"...#.#.#..\n" +
		".........#\n" +
		".....#....\n" +
		".........#\n" +
		"...^......\n" +
		"#.........\n" +
		"........#.\n"
	partTwoPatrolMapNoEntry3 = "" +
		"...#......\n" +
		"........#.\n" +
		"..#..#....\n" +
		"..........\n" +
		"..........\n" +
		".#.^......\n" +
		".......#..\n" +
		".#........\n" +
		"..#.......\n" +
		"..........\n"
		// loop not duplicated as its on existing path
	partTwoPatrolMapDuplicateLoop = "" +
		"..........\n" +
		"..........\n" +
		"..........\n" +
		"....#.....\n" +
		".......#..\n" +
		"..........\n" +
		"..........\n" +
		".#..^.....\n" +
		"#.........\n" +
		"..#...#...\n"
)

func TestDay6_Guard_NewGuard(t *testing.T) {
	g := NewGuard(Location{0, 0}, DirectionUp)
	assert.NotNil(t, g)
}

func TestDay6_Guard_AddVisitedLocation(t *testing.T) {
	tests := []struct {
		name     string
		guard    *Guard
		expected PositionMap
	}{
		{
			name: "one location visited",
			guard: &Guard{
				positionsVisited: PositionMap{},
				location:         Location{0, 2},
				direction:        DirectionUp,
			},
			expected: PositionMap{Location{0, 2}: []Direction{DirectionUp}},
		},
		{
			name: "two locations visited",
			guard: &Guard{
				positionsVisited: PositionMap{Location{0, 1}: []Direction{DirectionUp}},
				location:         Location{0, 2},
				direction:        DirectionUp,
			},
			expected: PositionMap{Location{0, 1}: []Direction{DirectionUp}, Location{0, 2}: []Direction{DirectionUp}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.guard.AddVisitedPosition()
			locations := test.guard.PositionsVisited()
			assert.Equal(t, test.expected, locations)
		})
	}
}

func TestDay6_Guard_turnRight(t *testing.T) {
	tests := []struct {
		name     string
		guard    *Guard
		expected *Guard
	}{
		{
			name:     "up to right",
			guard:    &Guard{PositionMap{}, Location{0, 0}, DirectionUp},
			expected: &Guard{PositionMap{}, Location{0, 0}, DirectionRight},
		},
		{
			name:     "right to down",
			guard:    &Guard{PositionMap{}, Location{0, 0}, DirectionRight},
			expected: &Guard{PositionMap{}, Location{0, 0}, DirectionDown},
		},
		{
			name:     "down to left",
			guard:    &Guard{PositionMap{}, Location{0, 0}, DirectionDown},
			expected: &Guard{PositionMap{}, Location{0, 0}, DirectionLeft},
		},
		{
			name:     "left to up",
			guard:    &Guard{PositionMap{}, Location{0, 0}, DirectionLeft},
			expected: &Guard{PositionMap{}, Location{0, 0}, DirectionUp},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.guard.turnRight()
			assert.Equal(t, test.expected, test.guard)
		})
	}
}

func TestDay6_PatrolMap_IsFree(t *testing.T) {
	tests := []struct {
		name      string
		location  Location
		patrolMap PatrolMap
		expected  bool
	}{
		{"free space", Location{0, 0}, smallPatrolMap, true},
		{"obstacle", Location{1, 0}, smallPatrolMap, false},
		{"not on map", Location{4, 0}, smallPatrolMap, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, smallPatrolMap.Free(test.location))
		})
	}
}

func TestDay6_PatrolMap_ParsePatrolMap(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedMap   PatrolMap
		expectedGuard *Guard
		expectedErr   error
	}{
		{
			name:  "valid input guard up",
			input: ".^\n.#",
			expectedMap: PatrolMap{
				Location{0, 0}: SpaceFree, Location{1, 0}: SpaceFree,
				Location{0, 1}: SpaceFree, Location{1, 1}: SpaceCrates,
			},
			expectedGuard: &Guard{PositionMap{}, Location{1, 0}, DirectionUp},
			expectedErr:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patrolMap, guard, err := ParseInput(strings.NewReader(test.input))
			assert.Equal(t, test.expectedMap, patrolMap)
			assert.Equal(t, test.expectedGuard, guard)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay6_PatrolMap_OnMap(t *testing.T) {
	tests := []struct {
		name      string
		location  Location
		patrolMap PatrolMap
		expected  bool
	}{
		{"on map", Location{0, 0}, smallPatrolMap, true},
		{"on map", Location{3, 3}, smallPatrolMap, true},
		{"not on map", Location{-1, 0}, smallPatrolMap, false},
		{"not on map", Location{0, -1}, smallPatrolMap, false},
		{"not on map", Location{4, 0}, smallPatrolMap, false},
		{"not on map", Location{0, 4}, smallPatrolMap, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, smallPatrolMap.OnMap(test.location))
		})
	}
}

func TestDay6_Guard_DejaVu(t *testing.T) {
	tests := []struct {
		name     string
		guard    *Guard
		expected bool
	}{
		{
			name: "no deja vu",
			guard: &Guard{
				positionsVisited: PositionMap{},
				location:         Location{0, 0},
				direction:        DirectionUp,
			},
			expected: false,
		},
		{
			name: "deja vu",
			guard: &Guard{
				positionsVisited: PositionMap{Location{0, 0}: []Direction{DirectionUp}},
				location:         Location{0, 0},
				direction:        DirectionUp,
			},
			expected: true,
		},
		{
			name: "no deja vu with different direction",
			guard: &Guard{
				positionsVisited: PositionMap{Location{0, 0}: []Direction{DirectionUp}},
				location:         Location{0, 0},
				direction:        DirectionRight,
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.guard.DejaVu())
		})
	}
}
func TestDay6_Patrol(t *testing.T) {
	tests := []struct {
		name            string
		input           io.Reader
		expectedLoops   int
		expectedVisited int
	}{
		{"part two example set", strings.NewReader(partTwoPatrolMap), 6, 41},
		{"part two test set", strings.NewReader(partTwoPatrolMapTest), 3, 26},
		{"part two close loop", strings.NewReader(partTwoPatrolMapCloseLoop), 2, 20},
		{"part two close loop 2", strings.NewReader(partTwoPatrolMapCloseLoop2), 2, 13},
		{"part two close loop 3", strings.NewReader(partTwoPatrolMapCloseLoop3), 3, 24},
		{"part two no entry", strings.NewReader(partTwoPatrolMapNoEntry), 1, 14},
		{"part two no entry 2", strings.NewReader(partTwoPatrolMapNoEntry2), 8, 32},
		{"part two no entry 3", strings.NewReader(partTwoPatrolMapNoEntry3), 2, 24},
		{"part two duplicate loop", strings.NewReader(partTwoPatrolMapDuplicateLoop), 2, 18},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patrolMap, guard, _ := ParseInput(test.input)

			startLocation := guard.location
			startDirection := guard.direction

			// move the guard until it leaves the map
			visitedPositions := patrolMap.Patrol(startLocation, startDirection, Location{-1, -1})

			assert.Equal(t, test.expectedVisited, len(visitedPositions))

			// Part2
			loops := 0
			for visitedPosition := range visitedPositions {
				if patrolMap.Patrol(startLocation, startDirection, visitedPosition) == nil {
					loops++
				}
			}

			// compare loops to correct solution
			assert.Equal(t, test.expectedLoops, loops)
		})
	}
}
