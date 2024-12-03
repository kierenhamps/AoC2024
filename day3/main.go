package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const (
	REGEX_DO         = `(do\(\))`
	REGEX_DONT       = `(don\'t\(\))`
	REGEX_MUL        = `(mul\(\d+,\d+\))`
	REGEX_MUL_VALUES = `mul\((\d+),(\d+)\)`
)

type Instruction interface {
	Result() int
}

type Do struct {
}

func NewDo() *Do {
	return &Do{}
}

func (d *Do) Result() int {
	return 0
}

type Dont struct {
}

func (d *Dont) Result() int {
	return 0
}

func NewDont() *Dont {
	return &Dont{}
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

	doRegex := regexp.MustCompile(REGEX_DO)
	dontRegex := regexp.MustCompile(REGEX_DONT)
	mulRegex := regexp.MustCompile(REGEX_MUL)
	mulValuesRegex := regexp.MustCompile(REGEX_MUL_VALUES)

	for s.scanner.Scan() {
		line := s.scanner.Text()

		// find all indexes of all instructions in this line
		scannedDos := doRegex.FindAllStringSubmatchIndex(line, -1)
		scannedDonts := dontRegex.FindAllStringSubmatchIndex(line, -1)
		scannedMuls := mulRegex.FindAllStringSubmatchIndex(line, -1)

		// add all indexes to a list
		allIndexes := make([][]int, 0)
		allIndexes = append(allIndexes, scannedDos...)
		allIndexes = append(allIndexes, scannedDonts...)
		allIndexes = append(allIndexes, scannedMuls...)
		// and sort them to ensure they are in order
		sort.Slice(allIndexes, func(i, j int) bool {
			return allIndexes[i][0] < allIndexes[j][0]
		})

		// Add all instructions in order of index
		for _, i := range allIndexes {
			instruction := line[i[0]:i[1]]
			if instruction == "do()" {
				instructions = append(instructions, NewDo())
			}
			if instruction == "don't()" {
				instructions = append(instructions, NewDont())
			}
			if instruction[0:3] == "mul" {
				mulValues := mulValuesRegex.FindStringSubmatch(instruction)
				left, _ := strconv.Atoi(mulValues[1])
				right, _ := strconv.Atoi(mulValues[2])
				instructions = append(instructions, NewMul(left, right))
			}
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

	// Calculate the result for each instruction
	var mulResultWithOthers int
	var recording bool = true
	for _, instruction := range instructions {
		switch instruction := instruction.(type) {
		case *Do:
			recording = true
		case *Dont:
			recording = false
		case *Mul:
			if recording {
				mulResultWithOthers += instruction.Result()
			}
		}
	}

	log.Printf("Mul Result with Do and Don't: %d\n", mulResultWithOthers)
}
