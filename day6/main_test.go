package main

import (
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
		{SpaceFree, SpaceFailedSuitPrototypes, SpaceFree, SpaceFree},
		{SpaceFree, SpaceFree, SpaceFree, SpaceSpoolOfVeryLongPolymers},
		{SpaceFree, SpaceFree, SpaceFree, SpaceFree},
		{SpaceCrates, SpaceFree, SpaceTankOfUniversalSolvent, SpaceFree},
	}
)

func TestDay6_Guard_NewGuard(t *testing.T) {
	location := NewLocation(0, 0)
	g := NewGuard(location, DirectionUp)
	assert.NotNil(t, g)
}

func TestDay6_Guard_LocationsVisited(t *testing.T) {
	tests := []struct {
		name     string
		guard    *Guard
		expected []Location
	}{
		{"one location visited", &Guard{[]Location{{0, 3}}, Location{0, 2}, DirectionUp}, []Location{{0, 3}}},
		{"two locations visited", &Guard{[]Location{{0, 1}, {0, 2}}, Location{0, 2}, DirectionDown}, []Location{{0, 1}, {0, 2}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			locations := test.guard.LocationsVisited()
			assert.Equal(t, test.expected, locations)
		})
	}
}

func TestDay6_Guard_Move(t *testing.T) {
	tests := []struct {
		name      string
		guard     *Guard
		patrolMap PatrolMap
		expected  *Guard
	}{
		{"move without obstacle and staying on the map", &Guard{[]Location{{0, 2}}, Location{0, 2}, DirectionUp}, smallPatrolMap, &Guard{[]Location{{0, 2}, {0, 1}}, Location{0, 1}, DirectionUp}},
		{"move with obstacle and staying on map", &Guard{[]Location{{2, 1}}, Location{2, 1}, DirectionRight}, smallPatrolMap, &Guard{[]Location{{2, 1}, {2, 2}}, Location{2, 2}, DirectionDown}},
		{"move without obstacle leaving the map", &Guard{[]Location{{0, 2}}, Location{0, 2}, DirectionLeft}, smallPatrolMap, &Guard{[]Location{{0, 2}, {-1, 2}}, Location{-1, 2}, DirectionLeft}},
		{"move with obstacle leaving the map", &Guard{[]Location{{0, 2}}, Location{0, 2}, DirectionDown}, smallPatrolMap, &Guard{[]Location{{0, 2}, {-1, 2}}, Location{-1, 2}, DirectionLeft}},
		{"dont move if off the map", &Guard{[]Location{{-1, 0}}, Location{-1, 0}, DirectionUp}, smallPatrolMap, &Guard{[]Location{{-1, 0}}, Location{-1, 0}, DirectionUp}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.guard.Move(test.patrolMap)
			assert.Equal(t, test.expected, test.guard)
		})
	}
}

func TestDay6_Guard_turnRight(t *testing.T) {
	tests := []struct {
		name     string
		guard    *Guard
		expected *Guard
	}{
		{"up to right", &Guard{[]Location{}, Location{0, 0}, DirectionUp}, &Guard{[]Location{}, Location{0, 0}, DirectionRight}},
		{"right to down", &Guard{[]Location{}, Location{0, 0}, DirectionRight}, &Guard{[]Location{}, Location{0, 0}, DirectionDown}},
		{"down to left", &Guard{[]Location{}, Location{0, 0}, DirectionDown}, &Guard{[]Location{}, Location{0, 0}, DirectionLeft}},
		{"left to up", &Guard{[]Location{}, Location{0, 0}, DirectionLeft}, &Guard{[]Location{}, Location{0, 0}, DirectionUp}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.guard.turnRight()
			assert.Equal(t, test.expected, test.guard)
		})
	}
}

func TestDay6_Location_NewLocation(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected Location
	}{
		{"valid input", 1, 2, Location{1, 2}},
		{"x is negative", -1, 2, Location{-1, 2}},
		{"y is negative", 2, -1, Location{2, -1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NewLocation(test.x, test.y)
			assert.Equal(t, test.expected, result)
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
		{"free space", NewLocation(0, 0), smallPatrolMap, true},
		{"obstacle", NewLocation(1, 0), smallPatrolMap, false},
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
			name:  "valid input guard right",
			input: ".#>.\n...#\n....\n#.#.",
			expectedMap: PatrolMap{
				{SpaceFree, SpaceCrates, SpaceFree, SpaceFree},
				{SpaceFree, SpaceFree, SpaceFree, SpaceCrates},
				{SpaceFree, SpaceFree, SpaceFree, SpaceFree},
				{SpaceCrates, SpaceFree, SpaceCrates, SpaceFree},
			},
			expectedGuard: &Guard{[]Location{{2, 0}}, Location{2, 0}, DirectionRight},
			expectedErr:   nil,
		},
		{
			name:  "valid input guard up",
			input: ".^\n.#",
			expectedMap: PatrolMap{
				{SpaceFree, SpaceFree},
				{SpaceFree, SpaceCrates},
			},
			expectedGuard: &Guard{[]Location{{1, 0}}, Location{1, 0}, DirectionUp},
			expectedErr:   nil,
		},
		{
			name:  "valid input guard down",
			input: ".v\n.#",
			expectedMap: PatrolMap{
				{SpaceFree, SpaceFree},
				{SpaceFree, SpaceCrates},
			},
			expectedGuard: &Guard{[]Location{{1, 0}}, Location{1, 0}, DirectionDown},
			expectedErr:   nil,
		},
		{
			name:  "valid input guard left",
			input: ".<\n.#",
			expectedMap: PatrolMap{
				{SpaceFree, SpaceFree},
				{SpaceFree, SpaceCrates},
			},
			expectedGuard: &Guard{[]Location{{1, 0}}, Location{1, 0}, DirectionLeft},
			expectedErr:   nil,
		},
		{
			name:          "invalid input",
			input:         "^#00\n000#\n0000\n#0#",
			expectedMap:   PatrolMap{},
			expectedGuard: &Guard{},
			expectedErr:   ErrInvalidPatrolMapInput,
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
