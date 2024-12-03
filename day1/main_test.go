package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay1_NewLocation(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		expected    location
		expectedErr error
	}{
		{
			name:        "valid input",
			input:       5,
			expected:    5,
			expectedErr: nil,
		},
		{
			name:        "input is zero",
			input:       0,
			expected:    0,
			expectedErr: ErrInputCannotBeZero,
		},
		{
			name:        "negative input",
			input:       -5,
			expected:    0,
			expectedErr: ErrInputCannotBeNegative,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewLocation(test.input)
			assert.Equal(t, test.expected, result)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay1_Location_Distance(t *testing.T) {
	tests := []struct {
		name     string
		left     location
		right    location
		expected int
	}{
		{"left is smaller than right", 1, 2, 1},
		{"right is smaller than left", 2, 1, 1},
		{"both sides match", 1, 1, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.left.Distance(test.right))
		})
	}
}

func TestDay1_LocationList_NewLocationList(t *testing.T) {
	ll := NewLocationList()
	assert.NotNil(t, ll)
}

func TestDay1_LocationList_Next(t *testing.T) {
	tests := []struct {
		name           string
		list           []location
		expectedFirst  location
		expectedSize   int
		expectedSecond location
	}{
		{
			name:           "normal jumbled list",
			list:           []location{90, 10, 80, 20, 70, 30, 60, 40, 50},
			expectedFirst:  10,
			expectedSize:   8,
			expectedSecond: 20,
		},
		{
			name:           "list with duplicates",
			list:           []location{77, 99, 33, 22, 77, 55, 33, 22},
			expectedFirst:  22,
			expectedSize:   7,
			expectedSecond: 22,
		},
		{
			name:           "errors without enough numbers in list",
			list:           []location{20},
			expectedFirst:  20,
			expectedSize:   0,
			expectedSecond: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ll := NewLocationList()
			for _, i := range test.list {
				ll.AddLocation(i)
			}

			next := ll.Next()
			assert.Equal(t, test.expectedFirst, next)
			assert.Equal(t, test.expectedSize, ll.Size())

			next = ll.Next()
			assert.Equal(t, test.expectedSecond, next)
		})
	}
}

func TestDay1_LocationList_Size(t *testing.T) {
	tests := []struct {
		name     string
		list     []location
		expected int
	}{
		{"normal list", []location{90, 10, 80, 20, 70, 30, 60, 40, 50}, 9},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ll := NewLocationList()
			for _, i := range test.list {
				ll.AddLocation(i)
			}

			assert.Equal(t, test.expected, ll.Size())
		})
	}
}

func TestDay1_CreateLists(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "valid input",
			input:       "10   20\n30   40\n50   60\n",
			expectedErr: nil,
		},
		{
			name:        "invalid input format",
			input:       "10   20\n30,  40\n50   60\n",
			expectedErr: ErrInvalidInputFormat,
		},
		{
			name:        "non-numeric input",
			input:       "10 20\n30 abc\n50 60\n",
			expectedErr: ErrInvalidInputFormat,
		},
		{
			name:        "a zero value is encountered on the right",
			input:       "10 20\n30 40\n50 0\n",
			expectedErr: ErrInputCannotBeZero,
		},
		{
			name:        "a zero value is encountered on the left",
			input:       "10 20\n0 40\n50 60\n",
			expectedErr: ErrInputCannotBeZero,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputFile := createTempFile(t, test.input)
			defer os.Remove(inputFile.Name())

			leftList, rightList, err := createLists(inputFile)
			assert.ErrorIs(t, err, test.expectedErr)
			assert.NotNil(t, leftList)
			assert.NotNil(t, rightList)
		})
	}
}

func TestDay1_SumDistances(t *testing.T) {
	tests := []struct {
		name      string
		leftList  []location
		rightList []location
		expected  int
	}{
		{
			name:      "normal lists",
			leftList:  []location{10, 30, 50},
			rightList: []location{20, 40, 60},
			expected:  30,
		},
		{
			name:      "lists with same values",
			leftList:  []location{10, 20, 30},
			rightList: []location{10, 20, 30},
			expected:  0,
		},
		{
			name:      "lists with one element",
			leftList:  []location{10},
			rightList: []location{20},
			expected:  10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			leftList := NewLocationList()
			rightList := NewLocationList()

			for _, i := range test.leftList {
				leftList.AddLocation(i)
			}
			for _, i := range test.rightList {
				rightList.AddLocation(i)
			}

			result := sumDistances(leftList, rightList)
			assert.Equal(t, test.expected, result)
		})
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	return tmpfile
}

func TestDay1_LocationList_CountMatches(t *testing.T) {
	tests := []struct {
		name          string
		list          LocationList
		value         location
		expectedCount int
	}{
		{"1 match", LocationList{[]location{9, 2, 5, 6, 5, 6, 5}}, 9, 1},
		{"2 matches", LocationList{[]location{9, 2, 5, 6, 5, 6, 5}}, 6, 2},
		{"3 matches", LocationList{[]location{9, 2, 5, 6, 5, 6, 5}}, 5, 3},
		{"no matches", LocationList{[]location{9, 2, 5, 6, 5, 6, 5}}, 7, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedCount, test.list.CountMatches(test.value))
		})
	}

}

func TestDay1_SumSimilarities(t *testing.T) {
	tests := []struct {
		name      string
		listLeft  *LocationList
		listRight *LocationList
		expected  int
	}{
		{
			name:      "test with two normal lists with matches",
			listLeft:  &LocationList{[]location{1, 2, 3, 4, 5}},
			listRight: &LocationList{[]location{1, 2, 3, 4, 5}},
			expected:  15,
		},
		{
			name:      "test with a non matching number",
			listLeft:  &LocationList{[]location{1, 2, 3, 4, 5}},
			listRight: &LocationList{[]location{1, 2, 3, 4, 5, 6}},
			expected:  15,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, sumSimilarities(test.listLeft, test.listRight))
		})
	}
}

func TestDay1_Location_Int(t *testing.T) {
	assert.Equal(t, 5, location(5).Int())
	assert.Equal(t, 2, location(2).Int())
}
