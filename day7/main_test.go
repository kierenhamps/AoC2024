package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay7_Equation_NewEquation(t *testing.T) {
	tests := []struct {
		name      string
		testValue Number
		numbers   map[int]Number
	}{
		{"equation with two numbers", Number{190}, map[int]Number{0: {10}, 1: {19}}},
		{"equation with three numbers", Number{3267}, map[int]Number{0: {81}, 1: {40}, 2: {27}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NewEquation(test.testValue, test.numbers)
			assert.NotNil(t, actual)
		})
	}
}

func TestDay7_Equation_EvaluateTrue(t *testing.T) {
	tests := []struct {
		name      string
		equation  *Equation
		operators []Operator
		expected  bool
	}{
		{
			name: "two numbers with just the multiplication operator",
			equation: NewEquation(
				Number{190},
				map[int]Number{0: {10}, 1: {19}},
			),
			operators: []Operator{NewMultiplicationOperator()},
			expected:  true,
		},
		{
			name: "two numbers with just the addition operator",
			equation: NewEquation(
				Number{190},
				map[int]Number{0: {10}, 1: {19}},
			),
			operators: []Operator{NewAdditionOperator()},
			expected:  false,
		},
		{
			name: "two numbers with both operators",
			equation: NewEquation(
				Number{190},
				map[int]Number{0: {10}, 1: {19}},
			),
			operators: []Operator{NewAdditionOperator(), NewMultiplicationOperator()},
			expected:  true,
		},
		{
			name: "three numbers with all operator",
			equation: NewEquation(
				Number{3267},
				map[int]Number{0: {81}, 1: {40}, 2: {27}},
			),
			operators: []Operator{NewAdditionOperator(), NewMultiplicationOperator()},
			expected:  true,
		},
		{
			name: "two numbers with all operator, but not possible",
			equation: NewEquation(
				Number{83},
				map[int]Number{0: {17}, 1: {5}},
			),
			operators: []Operator{NewAdditionOperator(), NewMultiplicationOperator()},
			expected:  false,
		},
		{
			name: "two numbers with all operator, but not possible",
			equation: NewEquation(
				Number{156},
				map[int]Number{0: {15}, 1: {6}},
			),
			operators: []Operator{NewAdditionOperator(), NewMultiplicationOperator()},
			expected:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.equation.EvaluateTrue(test.operators)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_Equation_TestValue(t *testing.T) {
	tests := []struct {
		name     string
		equation *Equation
		expected Number
	}{
		{
			name: "single digit",
			equation: NewEquation(
				Number{190},
				map[int]Number{0: {10}, 1: {19}},
			),
			expected: Number{190},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.equation.TestValue()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_Number_NewNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected Number
	}{
		{"single digit", 1, Number{1}},
		{"double digit", 10, Number{10}},
		{"triple digit", 100, Number{100}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NewNumber(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_Number_Int(t *testing.T) {
	tests := []struct {
		name     string
		number   Number
		expected int
	}{
		{"single digit", Number{1}, 1},
		{"double digit", Number{10}, 10},
		{"triple digit", Number{100}, 100},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.number.Int()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_AdditionOperator_NewAdditionOperator(t *testing.T) {
	operator := NewAdditionOperator()
	assert.NotNil(t, operator)
}

func TestDay7_AdditionOperator_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     Number
		right    Number
		expected Number
	}{
		{"single digit", Number{1}, Number{2}, Number{3}},
		{"double digit", Number{12}, Number{34}, Number{46}},
		{"triple digit", Number{123}, Number{345}, Number{468}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := NewAdditionOperator()
			actual := op.Evaluate(test.left, test.right)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_AdditionOperator_EvaluateMultiple(t *testing.T) {
	tests := []struct {
		name     string
		left     []Number
		right    Number
		expected []Number
	}{
		{"single digit", []Number{{1}, {2}}, Number{3}, []Number{{4}, {5}}},
		{"double digit", []Number{{12}, {34}}, Number{46}, []Number{{58}, {80}}},
		{"triple digit", []Number{{123}, {345}}, Number{468}, []Number{{591}, {813}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := NewAdditionOperator()
			actual := op.EvaluateMultiple(test.left, test.right)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_MultiplicationOperator_NewMultiplicationOperator(t *testing.T) {
	operator := NewMultiplicationOperator()
	assert.NotNil(t, operator)
}

func TestDay7_MultiplicationOperator_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     Number
		right    Number
		expected Number
	}{
		{"single digit", Number{1}, Number{2}, Number{2}},
		{"double digit", Number{12}, Number{34}, Number{408}},
		{"triple digit", Number{123}, Number{345}, Number{42435}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := NewMultiplicationOperator()
			actual := op.Evaluate(test.left, test.right)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_MultiplicationOperator_EvaluateMultiple(t *testing.T) {
	tests := []struct {
		name     string
		left     []Number
		right    Number
		expected []Number
	}{
		{"single digit", []Number{{1}, {2}}, Number{3}, []Number{{3}, {6}}},
		{"double digit", []Number{{12}, {34}}, Number{46}, []Number{{552}, {1564}}},
		{"triple digit", []Number{{123}, {345}}, Number{468}, []Number{{57564}, {161460}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := NewMultiplicationOperator()
			actual := op.EvaluateMultiple(test.left, test.right)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDay7_ParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []*Equation
	}{
		{
			name: "single equation",
			input: "" +
				"190: 10 19\n",
			expected: []*Equation{
				NewEquation(
					Number{190},
					map[int]Number{0: {10}, 1: {19}},
				),
			},
		},
		{
			name: "multiple equations",
			input: "" +
				"190: 10 19\n" +
				"3267: 81 40 27\n" +
				"83: 17 5\n",
			expected: []*Equation{
				NewEquation(
					Number{190},
					map[int]Number{0: {10}, 1: {19}},
				),
				NewEquation(
					Number{3267},
					map[int]Number{0: {81}, 1: {40}, 2: {27}},
				),
				NewEquation(
					Number{83},
					map[int]Number{0: {17}, 1: {5}},
				),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ParseEquations(strings.NewReader(test.input))
			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
