package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay3_Mul_NewMul(t *testing.T) {
	mul := NewMul(2, 3)
	assert.NotNil(t, mul)
}

func TestDay3_Mul_Result(t *testing.T) {
	tests := []struct {
		name     string
		left     int
		right    int
		expected int
	}{
		{"2 * 3", 2, 3, 6},
		{"3 * 3", 3, 3, 9},
		{"4 * 3", 4, 3, 12},
		{"5 * 5", 5, 5, 25},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mul := NewMul(test.left, test.right)
			assert.Equal(t, test.expected, mul.Result())
		})
	}
}

func TestDay3_Scanner_NewScanner(t *testing.T) {
	inputFile := strings.NewReader("mul(2,3)")
	scanner := NewScanner(inputFile)
	assert.NotNil(t, scanner)
}

func TestDay3_Scanner_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected []Instruction
	}{
		{"single valid input", strings.NewReader("mul(2,3)"), []Instruction{&Mul{2, 3}}},
		{"single valid input messy", strings.NewReader("m()from()*mul(810,344)?what()m"), []Instruction{&Mul{810, 344}}},
		{"multipl valid input", strings.NewReader("mul(2,3)\nmul(3,4)"), []Instruction{&Mul{2, 3}, &Mul{3, 4}}},
		{"multiple valid input messy", strings.NewReader("m()from()*mul(810,344)?what()mmul(223,22)asf22"), []Instruction{&Mul{810, 344}, &Mul{223, 22}}},
		{"invalid input", strings.NewReader("mul(2*3)"), nil},
		{"invalid input", strings.NewReader("mul[2,3]"), nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := NewScanner(test.input)
			instructions := scanner.Scan()
			assert.Equal(t, test.expected, instructions)
		})
	}
}
