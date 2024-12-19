package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay10_Trail_NewTrail(t *testing.T) {
	trail := NewTrail(Location{0, 0})
	assert.NotNil(t, trail)
}

func TestDay10_Trail_NextStep(t *testing.T) {
	tests := []struct {
		name        string
		location    Location
		height      Height
		miniMap     Grid
		expected    []Location
		expectedErr error
	}{
		{
			name:     "No possible location available",
			location: Location{1, 1},
			height:   4,
			// . 1 .
			// 8 4 3
			// . 6 .
			miniMap: Grid{
				1: {1: 1},
				2: {1: 8, 2: 4, 3: 6},
				3: {1: 3},
			},
			expected:    []Location{},
			expectedErr: ErrNoPossibleLocation,
		},
		{
			name:     "Example 2: 1 possible location",
			location: Location{3, 0},
			height:   0,
			miniMap: Grid{
				0: {3: 0},
				1: {3: 1},
			},
			// - - -
			// . 0 .
			// . 1 .
			expected: []Location{
				{3, 1},
			},
			expectedErr: nil,
		},
		{
			name:     "Example 2: 2 possible locations",
			location: Location{3, 3},
			height:   3,
			// . 2 .
			// 4 3 4
			// . . .
			miniMap: Grid{
				2: {3: 2},
				3: {2: 4, 3: 3, 4: 4},
			},
			expected: []Location{
				{4, 3},
				{2, 3},
			},
			expectedErr: nil,
		},
		{
			name:     "3 possible locations of height 9",
			location: Location{7, 6},
			height:   8,
			// . 9 .
			// 9 8 9
			// . . .
			miniMap: Grid{
				5: {7: 9},
				6: {6: 9, 7: 8, 8: 9},
			},
			expected: []Location{
				{7, 5},
				{8, 6},
				{6, 6},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			locations, err := NextStep(test.location, test.height, test.miniMap)
			assert.ErrorIs(t, err, test.expectedErr)
			assert.Equal(t, test.expected, locations)
		})
	}
}

func TestDay10_Trail_Walk(t *testing.T) {
	tests := []struct {
		name     string
		start    Location
		trailMap TrailMap
		path     Path
		expected []Path
	}{
		{
			name:  "simple path",
			start: Location{0, 0},
			trailMap: TrailMap{
				trailMap: Grid{
					0: {0: 8, 1: 9},
					1: {0: 1, 1: 0},
				},
			},
			path: Path{},
			expected: []Path{
				{{0, 0}, {1, 0}},
			},
		},
		{
			name:  "simple path with 2 possible locations",
			start: Location{1, 1},
			trailMap: TrailMap{
				trailMap: Grid{
					0: {1: 9},
					1: {0: 9, 1: 8},
				},
			},
			path: Path{},
			expected: []Path{
				{{1, 1}, {1, 0}},
				{{1, 1}, {0, 1}},
			},
		},
		{
			name:  "No possible location available",
			start: Location{1, 1},
			trailMap: TrailMap{
				trailMap: Grid{
					1: {1: 1},
				},
			},
			path:     Path{},
			expected: []Path{},
		},
		{
			name:  "Example 2: Trail with 2 paths",
			start: Location{3, 0},
			trailMap: TrailMap{
				trailMap: Grid{
					0: {3: 0},
					1: {3: 1},
					2: {3: 2},
					3: {0: 6, 1: 5, 2: 4, 3: 3, 4: 4, 5: 5, 6: 6},
					4: {0: 7, 6: 7},
					5: {0: 8, 6: 8},
					6: {0: 9, 6: 9}},
			},
			path: Path{},
			expected: []Path{
				{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}},
				{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {0, 3}, {0, 4}, {0, 5}, {0, 6}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paths := Walk(test.start, test.trailMap, test.path)
			assert.Equal(t, test.expected, paths)
		})
	}
}

func TestDay10_TrailMap_DiscoverTrails(t *testing.T) {
	tests := []struct {
		name     string
		input    TrailMap
		expected []Trail
	}{
		// {
		// 	name: "Example 2: Map with 1 trailhead and a score of 2",
		// 	input: TrailMap{
		// 		trailMap: Grid{
		// 			0: {3: 0},
		// 			1: {3: 1},
		// 			2: {3: 2},
		// 			3: {0: 6, 1: 5, 2: 4, 3: 3, 4: 4, 5: 5, 6: 6},
		// 			4: {0: 7, 6: 7},
		// 			5: {0: 8, 6: 8},
		// 			6: {0: 9, 6: 9}},
		// 		trailheads: []Location{
		// 			{3, 0},
		// 		},
		// 	},
		// 	expected: []Trail{
		// 		{
		// 			start: Location{3, 0},
		// 			paths: []Path{
		// 				{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {6, 4}, {6, 5}, {6, 6}},
		// 				{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {0, 3}, {0, 4}, {0, 5}, {0, 6}},
		// 			},
		// 			score: 2,
		// 		},
		// 	},
		// },
		{
			name: "Example 3: Map with 1 trailhead and a score of 4",
			input: TrailMap{
				trailMap: Grid{
					0: {2: 9, 3: 0, 6: 9},
					1: {3: 1, 5: 9, 6: 8},
					2: {3: 2, 6: 7},
					3: {0: 6, 1: 5, 2: 4, 3: 3, 4: 4, 5: 5, 6: 6},
					4: {0: 7, 1: 6, 2: 5, 4: 9, 5: 8, 6: 7},
					5: {0: 8, 1: 7, 2: 6},
					6: {0: 9, 1: 8, 2: 7}},
				trailheads: []Location{
					{3, 0},
				},
			},
			expected: []Trail{
				{
					start: Location{3, 0},
					paths: []Path{
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {6, 2}, {6, 1}, {6, 0}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {6, 2}, {6, 1}, {5, 1}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {6, 4}, {5, 4}, {4, 4}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {2, 6}, {1, 6}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {1, 5}, {1, 6}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {2, 5}, {1, 5}, {0, 5}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {1, 4}, {1, 5}, {1, 6}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {1, 4}, {1, 5}, {0, 5}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {2, 4}, {1, 4}, {0, 4}, {0, 5}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {1, 4}, {1, 5}, {1, 6}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {1, 4}, {1, 5}, {0, 5}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {1, 4}, {0, 4}, {0, 5}, {0, 6}},
						{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {0, 3}, {0, 4}, {0, 5}, {0, 6}},
					},
					score: 4,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trails := test.input.DiscoverTrails()
			assert.Equal(t, test.expected, trails)
		})
	}
}

func TestDay10_TrailMap_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TrailMap
	}{
		{
			name:  "Empty map",
			input: "",
			expected: TrailMap{
				trailMap:   Grid{},
				trailheads: []Location{},
			},
		},
		{
			name: "Example 1: Map",
			input: "" +
				"0123\n" +
				"1234\n" +
				"8765\n" +
				"9876\n",
			expected: TrailMap{
				trailMap: Grid{
					0: {0: 0, 1: 1, 2: 2, 3: 3},
					1: {0: 1, 1: 2, 2: 3, 3: 4},
					2: {0: 8, 1: 7, 2: 6, 3: 5},
					3: {0: 9, 1: 8, 2: 7, 3: 6}},
				trailheads: []Location{
					{0, 0},
				},
			},
		},
		{
			name: "Example 2: Map",
			input: "" +
				"...0...\n" +
				"...1...\n" +
				"...2...\n" +
				"6543456\n" +
				"7.....7\n" +
				"8.....8\n" +
				"9......9\n",
			expected: TrailMap{
				trailMap: Grid{
					0: {3: 0},
					1: {3: 1},
					2: {3: 2},
					3: {0: 6, 1: 5, 2: 4, 3: 3, 4: 4, 5: 5, 6: 6},
					4: {0: 7, 6: 7},
					5: {0: 8, 6: 8},
					6: {0: 9, 7: 9}},
				trailheads: []Location{
					{3, 0},
				},
			},
		},
		{
			name: "Example 3: Map",
			input: "" +
				"..90..9\n" +
				"...1.98\n" +
				"...2..7\n" +
				"6543456\n" +
				"765.987\n" +
				"876....\n" +
				"987....\n",
			expected: TrailMap{
				trailMap: Grid{
					0: {2: 9, 3: 0, 6: 9},
					1: {3: 1, 5: 9, 6: 8},
					2: {3: 2, 6: 7},
					3: {0: 6, 1: 5, 2: 4, 3: 3, 4: 4, 5: 5, 6: 6},
					4: {0: 7, 1: 6, 2: 5, 4: 9, 5: 8, 6: 7},
					5: {0: 8, 1: 7, 2: 6},
					6: {0: 9, 1: 8, 2: 7}},
				trailheads: []Location{
					{3, 0},
				},
			},
		},
		{
			name: "Example 4: Map",
			input: "" +
				"10..9..\n" +
				"2...8..\n" +
				"3...7..\n" +
				"4567654\n" +
				"...8..3\n" +
				"...9..2\n" +
				".....01\n",
			expected: TrailMap{
				trailMap: Grid{
					0: {0: 1, 1: 0, 4: 9},
					1: {0: 2, 4: 8},
					2: {0: 3, 4: 7},
					3: {0: 4, 1: 5, 2: 6, 3: 7, 4: 6, 5: 5, 6: 4},
					4: {3: 8, 6: 3},
					5: {3: 9, 6: 2},
					6: {5: 0, 6: 1}},
				trailheads: []Location{
					{1, 0},
					{5, 6},
				},
			},
		},
		{
			name: "Example 5: Map",
			input: "" +
				"89010123\n" +
				"78121874\n" +
				"87430965\n" +
				"96549874\n" +
				"45678903\n" +
				"32019012\n" +
				"01329801\n" +
				"10456732\n",
			expected: TrailMap{
				trailMap: Grid{
					0: {0: 8, 1: 9, 2: 0, 3: 1, 4: 0, 5: 1, 6: 2, 7: 3},
					1: {0: 7, 1: 8, 2: 1, 3: 2, 4: 1, 5: 8, 6: 7, 7: 4},
					2: {0: 8, 1: 7, 2: 4, 3: 3, 4: 0, 5: 9, 6: 6, 7: 5},
					3: {0: 9, 1: 6, 2: 5, 3: 4, 4: 9, 5: 8, 6: 7, 7: 4},
					4: {0: 4, 1: 5, 2: 6, 3: 7, 4: 8, 5: 9, 6: 0, 7: 3},
					5: {0: 3, 1: 2, 2: 0, 3: 1, 4: 9, 5: 0, 6: 1, 7: 2},
					6: {0: 0, 1: 1, 2: 3, 3: 2, 4: 9, 5: 8, 6: 0, 7: 1},
					7: {0: 1, 1: 0, 2: 4, 3: 5, 4: 6, 5: 7, 6: 3, 7: 2}},
				trailheads: []Location{
					{2, 0},
					{4, 0},
					{4, 2},
					{6, 4},
					{2, 5},
					{5, 5},
					{0, 6},
					{6, 6},
					{1, 7},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trailMap := Parse(strings.NewReader(test.input))
			assert.Equal(t, test.expected, trailMap)
		})
	}
}
