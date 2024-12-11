package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	part1Example1FrequencyMapInput = "" +
		"............\n" +
		"........0...\n" +
		".....0......\n" +
		".......0....\n" +
		"....0.......\n" +
		"......A.....\n" +
		"............\n" +
		"............\n" +
		"........A...\n" +
		".........A..\n" +
		"............\n" +
		"............\n"
	part1Example1FrequencyMap = &FrequencyMap{
		antinodes: make(map[Point]int),
		frequencies: map[Frequency][]Point{
			'0': {
				NewPoint(8, 1),
				NewPoint(5, 2),
				NewPoint(7, 3),
				NewPoint(4, 4),
			},
			'A': {
				NewPoint(6, 5),
				NewPoint(8, 8),
				NewPoint(9, 9),
			},
		},
		maxX: 12,
		maxY: 12,
	}
)

func TestDay8_FrequencyMap_NewFrequencyMap(t *testing.T) {
	fm := NewFrequencyMap()
	assert.NotNil(t, fm)
}

func TestDay8_FrequencyMap_AddPoint(t *testing.T) {
	tests := []struct {
		name      string
		fm        *FrequencyMap
		point     Point
		frequency Frequency
		expected  *FrequencyMap
	}{
		{
			name:      "add point to empty map",
			fm:        NewFrequencyMap(),
			point:     NewPoint(0, 0),
			frequency: '0',
			expected: &FrequencyMap{
				antinodes: make(map[Point]int),
				frequencies: map[Frequency][]Point{
					'0': {NewPoint(0, 0)},
				},
			},
		},
		{
			name: "add point to existing map",
			fm: &FrequencyMap{
				antinodes: make(map[Point]int),
				frequencies: map[Frequency][]Point{
					'0': {NewPoint(0, 0)},
				},
			},
			point:     NewPoint(1, 1),
			frequency: '0',
			expected: &FrequencyMap{
				antinodes: make(map[Point]int),
				frequencies: map[Frequency][]Point{
					'0': {NewPoint(0, 0), NewPoint(1, 1)},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fm.AddPoint(test.frequency, test.point)
			assert.Equal(t, test.expected, test.fm)
		})
	}
}

func TestDay8_FrequencyMap_FindAllAntinodes(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrequencyMap
		expected map[Point]int
	}{
		{
			name: "part 1 example 1 antinodes",
			fm:   part1Example1FrequencyMap,
			expected: map[Point]int{
				NewPoint(6, 0):   1,
				NewPoint(11, 0):  1,
				NewPoint(3, 1):   2,
				NewPoint(4, 2):   1,
				NewPoint(10, 2):  1,
				NewPoint(2, 3):   1,
				NewPoint(9, 4):   1,
				NewPoint(1, 5):   1,
				NewPoint(6, 5):   1,
				NewPoint(3, 6):   1,
				NewPoint(0, 7):   1,
				NewPoint(7, 7):   1,
				NewPoint(10, 10): 1,
				NewPoint(10, 11): 1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fm.FindAllAntinodes()
			assert.True(t, compareMaps(test.fm.antinodes, test.expected))
		})
	}

}

func TestDay8_FrequencyMap_InBounds(t *testing.T) {
	tests := []struct {
		name     string
		fm       *FrequencyMap
		point    Point
		expected bool
	}{
		{
			name: "point in bounds",
			fm: &FrequencyMap{
				antinodes: make(map[Point]int),
				maxX:      10,
				maxY:      10,
			},
			point:    NewPoint(5, 5),
			expected: true,
		},
		{
			name: "point out of bounds",
			fm: &FrequencyMap{
				antinodes: make(map[Point]int),
				maxX:      10,
				maxY:      10,
			},
			point:    NewPoint(11, 11),
			expected: false,
		},
		{
			name: "point out of bounds negative",
			fm: &FrequencyMap{
				antinodes: make(map[Point]int),
				maxX:      10,
				maxY:      10,
			},
			point:    NewPoint(-1, -1),
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.fm.InBounds(test.point)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay8_Point_NewPoint(t *testing.T) {
	p := NewPoint(0, 0)
	assert.NotNil(t, p)
}

func TestDay8_ParseFrequencyMap(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected *FrequencyMap
	}{
		{
			name:     "parse part 1 example 1",
			input:    strings.NewReader(part1Example1FrequencyMapInput),
			expected: part1Example1FrequencyMap,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := ParseFrequencyMap(test.input)
			for f, points := range test.expected.frequencies {
				for i, p := range points {
					assert.Equal(t, p, actual.frequencies[f][i])
				}
			}
		})
	}
}

func TestDay8_FindAntinode(t *testing.T) {
	tests := []struct {
		name        string
		pointA      Point
		pointB      Point
		expected    Point
		expectedErr error
	}{
		{
			name:        "antinode heading southwest",
			pointA:      NewPoint(8, 1),
			pointB:      NewPoint(5, 2),
			expected:    NewPoint(2, 3),
			expectedErr: nil,
		},
		{
			name:        "antinode heading northwest",
			pointA:      NewPoint(7, 3),
			pointB:      NewPoint(5, 2),
			expected:    NewPoint(3, 1),
			expectedErr: nil,
		},
		{
			name:        "antinode heading northeast",
			pointA:      NewPoint(4, 4),
			pointB:      NewPoint(7, 3),
			expected:    NewPoint(10, 2),
			expectedErr: nil,
		},
		{
			name:        "antinode heading southeast",
			pointA:      NewPoint(5, 2),
			pointB:      NewPoint(7, 3),
			expected:    NewPoint(9, 4),
			expectedErr: nil,
		},
		{
			name:        "antinode heading west out of bounds",
			pointA:      NewPoint(1, 1),
			pointB:      NewPoint(0, 1),
			expected:    Point{},
			expectedErr: ErrPointOutOfBounds,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := FindAntinode(test.pointA, test.pointB)
			assert.Equal(t, test.expected, actual)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func compareMaps(m1, m2 map[Point]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}
