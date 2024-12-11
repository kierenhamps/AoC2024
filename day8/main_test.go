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
	part1Example1Antinodes = map[Point]int{
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
	}
	part2Example1FrequencyMapInput = "" +
		"T.........\n" +
		"...T......\n" +
		".T........\n" +
		"..........\n" +
		"..........\n" +
		"..........\n" +
		"..........\n" +
		"..........\n" +
		"..........\n" +
		"..........\n"
	part2Example1FrequencyMap = &FrequencyMap{
		antinodes: make(map[Point]int),
		frequencies: map[Frequency][]Point{
			'T': {
				NewPoint(0, 0),
				NewPoint(3, 1),
				NewPoint(1, 2),
			},
		},
		maxX: 10,
		maxY: 10,
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
		af       AntinodeFinder
		expected map[Point]int
	}{
		{
			name:     "part 1 example 1 antinodes",
			fm:       part1Example1FrequencyMap,
			af:       SimpleAntinodeFinder{},
			expected: part1Example1Antinodes,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fm.FindAllAntinodes(test.af)
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
		{
			name:     "parse part 2 example 1",
			input:    strings.NewReader(part2Example1FrequencyMapInput),
			expected: part2Example1FrequencyMap,
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

func TestDay8_SimpleAntinodeFinder_FindAntinodes(t *testing.T) {
	tests := []struct {
		name     string
		pointA   Point
		pointB   Point
		maxX     int
		maxY     int
		expected []Point
	}{
		{
			name:     "antinode heading southwest",
			pointA:   NewPoint(8, 1),
			pointB:   NewPoint(5, 2),
			maxX:     12,
			maxY:     12,
			expected: []Point{NewPoint(2, 3)},
		},
		{
			name:     "antinode heading northwest",
			pointA:   NewPoint(7, 3),
			pointB:   NewPoint(5, 2),
			maxX:     12,
			maxY:     12,
			expected: []Point{NewPoint(3, 1)},
		},
		{
			name:     "antinode heading northeast",
			pointA:   NewPoint(4, 4),
			pointB:   NewPoint(7, 3),
			maxX:     12,
			maxY:     12,
			expected: []Point{NewPoint(10, 2)},
		},
		{
			name:     "antinode heading southeast",
			pointA:   NewPoint(5, 2),
			pointB:   NewPoint(7, 3),
			maxX:     12,
			maxY:     12,
			expected: []Point{NewPoint(9, 4)},
		},
		{
			name:     "antinode heading west out of bounds",
			pointA:   NewPoint(1, 1),
			pointB:   NewPoint(0, 1),
			maxX:     12,
			maxY:     12,
			expected: []Point{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			af := SimpleAntinodeFinder{}
			actual := af.FindAntinodes(test.pointA, test.pointB, test.maxX, test.maxY)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay8_HarmonicAntinodeFinder_FindAntinodes(t *testing.T) {
	tests := []struct {
		name     string
		pointA   Point
		pointB   Point
		maxX     int
		maxY     int
		expected []Point
	}{
		{
			name:   "antinodes heading south-southeast",
			pointA: NewPoint(0, 0),
			pointB: NewPoint(1, 2),
			maxX:   10,
			maxY:   10,
			expected: []Point{
				NewPoint(0, 0),
				NewPoint(1, 2),
				NewPoint(2, 4),
				NewPoint(3, 6),
				NewPoint(4, 8),
			},
		},
		{
			name:   "antinodes heading southeast",
			pointA: NewPoint(0, 0),
			pointB: NewPoint(3, 1),
			maxX:   10,
			maxY:   10,
			expected: []Point{
				NewPoint(0, 0),
				NewPoint(3, 1),
				NewPoint(6, 2),
				NewPoint(9, 3),
			},
		},
		{
			name:   "antinodes heading northeast",
			pointA: NewPoint(1, 2),
			pointB: NewPoint(3, 1),
			maxX:   10,
			maxY:   10,
			expected: []Point{
				NewPoint(1, 2),
				NewPoint(3, 1),
				NewPoint(5, 0),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			af := HarmonicAntinodeFinder{}
			actual := af.FindAntinodes(test.pointA, test.pointB, test.maxX, test.maxY)
			assert.Equal(t, test.expected, actual)
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
