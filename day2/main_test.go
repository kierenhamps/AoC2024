package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay2_Level_NewLevel(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		expected    Level
		expectedErr error
	}{
		{"valid input", 50, Level(50), nil},
		{"zero input", 0, Level(0), ErrInputCannotBeZero},
		{"negative input", -5, Level(0), ErrInputCannotBeNegative},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			level, err := NewLevel(test.input)
			assert.Equal(t, test.expected, level)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay2_Report_NewReport(t *testing.T) {
	report := NewReport()
	assert.NotNil(t, report)
}

func TestDay2_Report_AddLevel(t *testing.T) {
	testReport := NewReport()

	tests := []struct {
		name         string
		input        Level
		expectedSize int
	}{
		{"one level", Level(1), 1},
		{"two levels", Level(1), 2},
		{"three levels", Level(1), 3},
		{"four levels", Level(1), 4},
		{"five levels", Level(1), 5},
		{"six levels", Level(1), 6},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testReport.AddLevel(test.input)
			assert.Equal(t, test.expectedSize, testReport.Size())
		})
	}
}

func TestDay2_Report_IsSafe(t *testing.T) {
	tests := []struct {
		name        string
		inputLevels []Level
		expected    bool
		expectedErr error
	}{
		{"empty report", []Level{}, false, ErrReportIsEmpty},
		{"safe: levels are all decreasing at a rate between 1 and 3", []Level{Level(7), Level(6), Level(4), Level(2), Level(1)}, true, nil},
		{"safe: levels are all increasing at a rate between 1 and 3", []Level{Level(1), Level(3), Level(6), Level(7), Level(9)}, true, nil},
		{"unsafe: two adjacent levels increased by more than 3", []Level{Level(1), Level(2), Level(7), Level(8), Level(9)}, false, ErrLevelsIncreasedByMoreThanThree},
		{"unsafe: two adjacent levels decreased by more than 3", []Level{Level(9), Level(7), Level(6), Level(2), Level(1)}, false, ErrLevelsDecreasedByMoreThanThree},
		{"unsafe: levels are increasing and decreasing", []Level{Level(1), Level(3), Level(2), Level(4), Level(5)}, false, ErrLevelsAreIncreasingAndDecreasing},
		{"unsafe: levels are neither increasing nor decreasing", []Level{Level(8), Level(6), Level(4), Level(4), Level(1)}, false, ErrLevelsAreNeitherIncreasingNorDecreasing},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report := NewReport()
			for _, level := range test.inputLevels {
				report.AddLevel(level)
			}
			safe, err := report.IsSafe()
			assert.Equal(t, test.expected, safe)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay2_Report_Size(t *testing.T) {
	testReport := NewReport()

	tests := []struct {
		name     string
		input    []Level
		expected int
	}{
		{"empty report", []Level{}, 0},
		{"one level", []Level{Level(1)}, 1},
		{"two levels", []Level{Level(1)}, 2},
		{"three levels", []Level{Level(1)}, 3},
		{"four levels", []Level{Level(1)}, 4},
		{"five levels", []Level{Level(1)}, 5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, level := range test.input {
				testReport.AddLevel(level)
			}
			assert.Equal(t, test.expected, testReport.Size())
		})
	}
}

func TestDay2_Report_IsSafeWithProblemDampner(t *testing.T) {
	tests := []struct {
		name        string
		inputLevels []Level
		expected    bool
		expectedErr error
	}{
		{"empty report", []Level{}, false, ErrReportIsEmpty},
		{"still safe: levels are all decreasing at a rate between 1 and 3", []Level{Level(7), Level(6), Level(4), Level(2), Level(1)}, true, nil},
		{"still unsafe: two adjacent levels increased by more than 3", []Level{Level(1), Level(2), Level(7), Level(8), Level(9)}, false, ErrLevelsIncreasedByMoreThanThree},
		{"still unsafe: two adjacent levels decreased by more than 3", []Level{Level(9), Level(7), Level(6), Level(2), Level(1)}, false, ErrLevelsDecreasedByMoreThanThree},
		{"now safe: levels are increasing and decreasing", []Level{Level(1), Level(3), Level(2), Level(4), Level(5)}, true, nil},
		{"now safe: levels are neither increasing nor decreasing", []Level{Level(8), Level(6), Level(4), Level(4), Level(1)}, true, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report := NewReport()
			for _, level := range test.inputLevels {
				report.AddLevel(level)
			}
			safe, err := report.IsSafeWithProblemDampner()
			assert.Equal(t, test.expected, safe)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}
