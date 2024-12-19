package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Rule is an interface that defines how to run a generic rule
//
// Rules are used to change the state of a Stone
// A Stone can be run against any rule using the Eval method
//
// If the rule is not applicable to the Stone, the Stone is returned as is
// otherwise the new Stone (or Stones) are returned
//
// If the rule is applicable to the Stone, the second return value is true to
// indicate that the Stone was changed
type Rule interface {
	Eval(Stone) ([]Stone, bool)
}

// RuleZeroToOne represents a rule that changes a stone from 0 to 1
type RuleZeroToOne struct{}

// Eval changes a stone from 0 to 1
func (r *RuleZeroToOne) Eval(s Stone) ([]Stone, bool) {
	if s == 0 {
		return []Stone{Stone(1)}, true
	}
	return []Stone{s}, false
}

// RuleSplitEvenDigits represents a rule that splits a stone with even digits into two stones
type RuleSplitEvenDigits struct{}

// Eval splits a stone with even digits into two stones
func (r *RuleSplitEvenDigits) Eval(s Stone) ([]Stone, bool) {
	text := strconv.Itoa(int(s))
	l := len(text)
	if l%2 == 0 {
		firstHalf := text[0 : l/2]
		firstStone, _ := strconv.Atoi(firstHalf)
		secondHalf := text[l/2:]
		secondStone, _ := strconv.Atoi(secondHalf)
		return []Stone{Stone(firstStone), Stone(secondStone)}, true
	}
	return []Stone{s}, false
}

// RuleMultiplyBy2024 represents a rule that changes a stone from 1 to 2024
type RuleMultiplyBy2024 struct{}

// Eval changes a stone from 1 to 2024
func (r *RuleMultiplyBy2024) Eval(s Stone) ([]Stone, bool) {
	return []Stone{Stone(2024 * s)}, true
}

// Stone represents a physics-defying stone
type Stone int

// Blink runs a slice of Stones through the rules provided
//
// Rules are evaluated in order
// A slice of resulting Stones are returned
func Blink(in []Stone, rules []Rule) []Stone {
	out := []Stone{}

	for _, stone := range in {
		for _, rule := range rules {
			result, changed := rule.Eval(stone)
			if changed {
				out = append(out, result...)
				break
			}
		}
	}
	return out
}

func main() {
	input, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	stones := []Stone{}
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), " ")
		for _, value := range values {
			stone, _ := strconv.Atoi(value)
			stones = append(stones, Stone(stone))
		}
	}

	// Part 1 ruleset
	rules := []Rule{
		&RuleZeroToOne{},
		&RuleSplitEvenDigits{},
		&RuleMultiplyBy2024{},
	}

	// Part 1 - Blink 25 times
	for i := 0; i < 25; i++ {
		stones = Blink(stones, rules)
	}
	log.Println("Part 1: After 25 blinks, there are", len(stones), "stones")
}
