package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	REGEX_EQUATION = `^(\d+):\s*(\d+(?:\s+\d+)*)+$`
)

// Operator represents an operator that can perform operations on numbers
type Operator interface {
	Evaluate(left, right Number) Number
	EvaluateMultiple(left []Number, right Number) []Number
}

// AdditionOperator represents an operator that can perform additions
type AdditionOperator struct{}

// NewAdditionOperator creates a new AdditionOperator
func NewAdditionOperator() *AdditionOperator {
	return &AdditionOperator{}
}

// Evaluate performs an addition operation on two numbers
func (a AdditionOperator) Evaluate(left, right Number) Number {
	return NewNumber(left.value + right.value)
}

// EvaluateMultiple performs all addition operations on a list of inputs and returns the multiple results
func (a AdditionOperator) EvaluateMultiple(left []Number, right Number) []Number {
	var results []Number
	for _, number := range left {
		results = append(results, a.Evaluate(number, right))
	}
	return results
}

// ConcatenationOperator represents an operator that can perform concatenations
type ConcatenationOperator struct{}

// NewConcatenationOperator creates a new ConcatenationOperator
func NewConcatenationOperator() *ConcatenationOperator {
	return &ConcatenationOperator{}
}

// Evaluate performs a concatenation operation on two numbers
func (c ConcatenationOperator) Evaluate(left, right Number) Number {
	result := left.String() + right.String()
	i, _ := strconv.Atoi(result)
	return NewNumber(i)
}

// EvaluateMultiple performs all concatenation operations on a list of inputs and returns the multiple results
func (c ConcatenationOperator) EvaluateMultiple(left []Number, right Number) []Number {
	var results []Number
	for _, number := range left {
		results = append(results, c.Evaluate(number, right))
	}
	return results
}

// MultiplcationOperator represents an operator that can perform multiplications
type MultiplcationOperator struct{}

// NewMultiplicationOperator creates a new MultiplicationOperator
func NewMultiplicationOperator() *MultiplcationOperator {
	return &MultiplcationOperator{}
}

// Evaluate performs a multiplication operation on two numbers
func (m MultiplcationOperator) Evaluate(left, right Number) Number {
	return NewNumber(left.value * right.value)
}

// EvaluateMultiple performs all multiplication operations on a list of inputs and returns the multiple results
func (m MultiplcationOperator) EvaluateMultiple(left []Number, right Number) []Number {
	var results []Number
	for _, number := range left {
		results = append(results, m.Evaluate(number, right))
	}
	return results
}

// Equation represents a calibration equation
type Equation struct {
	testValue Number
	numbers   map[int]Number
}

// NewEquation creates a new Equation
func NewEquation(testValue Number, numbers map[int]Number) *Equation {
	return &Equation{testValue, numbers}
}

// EvaluateTrue evaluates the equation and returns true if the test value is possible
// to be made with the given operators
func (e Equation) EvaluateTrue(operators []Operator) bool {
	var combined []Number
	for i := 0; i < len(e.numbers); i++ {
		if combined == nil {
			combined = []Number{e.numbers[i]}
			continue
		}
		var newCombined []Number
		for _, op := range operators {
			opResult := op.EvaluateMultiple(combined, e.numbers[i])
			newCombined = append(newCombined, opResult...)
		}
		combined = newCombined
	}
	for _, result := range combined {
		if e.testValue.value == result.value {
			return true
		}
	}
	return false
}

// TestValue returns the test value of the equation
func (e Equation) TestValue() Number {
	return e.testValue
}

// Number represents any number used in an equation
type Number struct {
	value int
}

// NewNumber creates a new Number Value Object
//
// Add validation rules here for 0 and negative numbers
func NewNumber(i int) Number {
	return Number{i}
}

// Int returns the integer value of the Number
func (n Number) Int() int {
	return n.value
}

// String returns the string value of the Number
func (n Number) String() string {
	return strconv.Itoa(n.value)
}

func main() {
	// Load test data
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("could not open input file: %v", err)
	}
	defer input.Close()

	// parse the input
	equations, err := ParseEquations(input)
	if err != nil {
		log.Fatalf("could not parse input: %v", err)
	}

	// Create the operators
	var operators []Operator
	operators = append(operators, NewAdditionOperator())
	operators = append(operators, NewMultiplicationOperator())

	// Evaluate the equations
	calibrationTotal := 0
	for _, equation := range equations {
		if equation.EvaluateTrue(operators) {
			calibrationTotal += equation.TestValue().Int()
		}
	}

	log.Println("(Part 1) Total Calibration Result: ", calibrationTotal)

	// Add the Concatenation Operator
	operators = append(operators, NewConcatenationOperator())

	// Re-evaluate the equations
	calibrationTotal = 0
	for _, equation := range equations {
		if equation.EvaluateTrue(operators) {
			calibrationTotal += equation.TestValue().Int()
		}
	}

	log.Println("(Part 2) Total Calibration Result: ", calibrationTotal)
}

func ParseEquations(input io.Reader) ([]*Equation, error) {
	equations := make([]*Equation, 0)

	// regex
	regexEquation := regexp.MustCompile(REGEX_EQUATION)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		values := regexEquation.FindStringSubmatch(scanner.Text())
		testValueRaw, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		numbers := make(map[int]Number, len(values[2]))
		numbersRaw := strings.Split(values[2], " ")
		for i := range numbersRaw {
			n, err := strconv.Atoi(numbersRaw[i])
			if err != nil {
				return nil, err
			}
			numbers[i] = NewNumber(n)
		}
		equations = append(equations, NewEquation(NewNumber(testValueRaw), numbers))
	}

	return equations, nil
}
