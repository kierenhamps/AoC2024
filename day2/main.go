package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DirectionUnknown = iota
	DirectionIncreasing
	DirectionDecreasing
)

var (
	ErrInputCannotBeZero     = errors.New("input cannot be zero")
	ErrInputCannotBeNegative = errors.New("input cannot be negative")
	ErrReportIsEmpty         = errors.New("report is empty")
	ErrInvalidFormat         = errors.New("invalid format, each line should be 5 integers separated by spaces")

	// Unsafe conditions
	ErrLevelsIncreasedByMoreThanThree          = errors.New("levels increased by more than 3")
	ErrLevelsDecreasedByMoreThanThree          = errors.New("levels decreased by more than 3")
	ErrLevelsAreIncreasingAndDecreasing        = errors.New("levels are increasing and decreasing")
	ErrLevelsAreNeitherIncreasingNorDecreasing = errors.New("levels are neither increasing nor decreasing")
)

type Level int

func NewLevel(l int) (Level, error) {
	if l == 0 {
		return 0, ErrInputCannotBeZero
	}
	if l < 0 {
		return 0, ErrInputCannotBeNegative
	}
	return Level(l), nil
}

func (l Level) Equal(r Level) bool {
	return l == r
}

type Report struct {
	levels []Level
}

func NewReport() *Report {
	return &Report{}
}

func (r *Report) AddLevel(l Level) {
	r.levels = append(r.levels, l)
}

func (r *Report) IsSafe() (bool, error) {
	if r.Size() == 0 {
		return false, ErrReportIsEmpty
	}

	safe := true
	direction := DirectionUnknown
	var err error
	previousLevel := r.levels[0]
	for i, l := range r.levels {
		if i == 0 {
			continue
		}
		// No adjacent levels can be the same
		if l.Equal(previousLevel) {
			safe = false
			err = ErrLevelsAreNeitherIncreasingNorDecreasing
			break
		}
		// Check direction
		currentDirection := DirectionIncreasing
		difference := l - previousLevel
		if l < previousLevel {
			currentDirection = DirectionDecreasing
			difference = previousLevel - l
		}
		if direction == DirectionUnknown {
			direction = currentDirection
		}
		// Direction cannot change
		if direction != currentDirection {
			safe = false
			err = ErrLevelsAreIncreasingAndDecreasing
			break
		}
		// Variation cannot be more than 3
		if difference > 3 {
			safe = false
			switch currentDirection {
			case DirectionIncreasing:
				err = ErrLevelsIncreasedByMoreThanThree
			case DirectionDecreasing:
				err = ErrLevelsDecreasedByMoreThanThree
			}
			break
		}

		previousLevel = l
	}

	return safe, err
}

func (r *Report) IsSafeWithProblemDampner() (bool, error) {
	safe, err := r.IsSafe()
	if safe {
		return true, nil
	}

	for i := range r.Size() {
		// Check if removing the level makes the report safe
		levels := make([]Level, r.Size())
		_ = copy(levels, r.levels)
		levels = append(levels[:i], levels[i+1:]...)
		report := NewReport()
		for _, level := range levels {
			report.AddLevel(level)
		}
		safe, _ := report.IsSafe()
		if safe {
			return true, nil
		}
	}

	return false, err
}

func (r *Report) Size() int {
	return len(r.levels)
}

func main() {
	// Read inputs
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Scan the inputs and convert into the domain
	reports := make([]*Report, 0)
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), " ")

		// create the report and populate
		report := NewReport()
		for _, v := range values {
			levelInt, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			level, err := NewLevel(levelInt)
			if err != nil {
				log.Fatal(err)
			}
			report.AddLevel(level)
		}

		// add report to the list
		reports = append(reports, report)
	}

	// Check if the reports are safe and sum
	var sum int
	for _, r := range reports {
		safe, _ := r.IsSafe()
		if safe {
			sum++
		}
	}
	log.Printf("Safe reports: %d", sum)

	// Check if the reports are safe with problem dampner and sum
	var pbSum int
	for _, r := range reports {
		safe, _ := r.IsSafeWithProblemDampner()
		if safe {
			pbSum++
		}
	}
	log.Printf("Safe reports with Problem Dampener: %d", pbSum)
}
