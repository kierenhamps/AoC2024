package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	REGEX_MUL = `mul\((\d+),(\d+)\)`
)

type Instruction interface {
	Result() int
}

type Mul struct {
	left  int
	right int
}

func NewMul(left, right int) *Mul {
	return &Mul{
		left:  left,
		right: right,
	}
}

func (m *Mul) Result() int {
	return m.left * m.right
}

type Scanner struct {
	scanner *bufio.Scanner
}

func NewScanner(input io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(input),
	}
}

func (s *Scanner) Scan() []Instruction {
	var instructions []Instruction

	mulRegex := regexp.MustCompile(REGEX_MUL)

	for s.scanner.Scan() {
		// find all mul instructions
		scannedMuls := mulRegex.FindAllStringSubmatch(s.scanner.Text(), -1)

		// create mul instructions
		for _, mul := range scannedMuls {
			left, _ := strconv.Atoi(mul[1])
			right, _ := strconv.Atoi(mul[2])
			instructions = append(instructions, NewMul(left, right))
		}
	}

	return instructions
}

func main() {
	// Open input file
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Create scanner
	scanner := NewScanner(inputFile)

	// Scan the input file for instructions
	instructions := scanner.Scan()

	// Calculate the result for each instruction
	var mulResult int
	for _, instruction := range instructions {
		switch instruction := instruction.(type) {
		case *Mul:
			mulResult += instruction.Result()
		}
	}

	log.Printf("Mul Result: %d\n", mulResult)
}
