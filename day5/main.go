package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	REGEX_RULE   = `^(\d+)\|(\d+)$`
	REGEX_MANUAL = `^(\d+(?:,(\d+))+)$`
)

var (
	ErrInputCannotBeZero     = errors.New("input cannot be zero")
	ErrInputCannotBeNegative = errors.New("input cannot be negative")
)

// PageNumber is a ValueObject that represents a page number.
type PageNumber struct {
	number int
}

// NewPageNumber creates a new PageNumber at a given index.
//
// A PageNumber cannot be zero or negative.
func NewPageNumber(number int) (PageNumber, error) {
	if number == 0 {
		return PageNumber{}, ErrInputCannotBeZero
	}
	if number < 0 {
		return PageNumber{}, ErrInputCannotBeNegative
	}
	return PageNumber{
		number: number,
	}, nil
}

// Int returns the integer value of the PageNumber.
func (pn *PageNumber) Int() int {
	return pn.number
}

// PageOrderingRule is a ValueObject that represents a rule for ordering pages.
type PageOrderingRule struct {
	left  PageNumber
	right PageNumber
}

// NewPageOrderingRule creates a new PageOrderingRule.
func NewPageOrderingRule(left, right PageNumber) PageOrderingRule {
	return PageOrderingRule{
		left:  left,
		right: right,
	}
}

// Correct changes the order of pages in the given manual based on the rule if possible.
//
// returns bool if the manual was corrected.
func (r *PageOrderingRule) Correct(manual *SafetyManual) bool {
	// skip rule if pages are not in the manual
	if manual.PageIndex(r.left) == -1 || manual.PageIndex(r.right) == -1 {
		return false
	}

	valid, _ := r.Valid(manual)
	if valid {
		return false
	}
	// move the right page to the left of the left page to make rule valid
	leftIndex := manual.PageIndex(r.left)
	manual.MovePage(r.right, leftIndex)
	return true
}

// Valid checks if a rule is valid for a SafetyManual.
func (r *PageOrderingRule) Valid(manual *SafetyManual) (bool, error) {
	// skip rule if pages are not in the manual
	if manual.PageIndex(r.left) == -1 || manual.PageIndex(r.right) == -1 {
		return true, nil
	}

	leftIndex := manual.PageIndex(r.left)
	rightIndex := manual.PageIndex(r.right)

	if leftIndex > rightIndex {
		return false, nil
	}
	return true, nil
}

// PageOrderingRuleset is a ValueObject that represents a set of PageOrderingRules.
type PageOrderingRuleset struct {
	rules []PageOrderingRule
}

// NewPageOrderingRuleset creates a new PageOrderingRuleset for managing a set of rules.
func NewPageOrderingRuleset() *PageOrderingRuleset {
	return &PageOrderingRuleset{}
}

// AddRule adds a new rule to the ruleset.
func (rs *PageOrderingRuleset) AddRule(rule PageOrderingRule) {
	rs.rules = append(rs.rules, rule)
}

// Correct corrects the order of pages in the given manual based on all rules in the ruleset.
//
// This is ugly, but it works.
// It will keep running until no more corrections are made. To ensure no rule breaks another.
// This should problably have a timeout and limit to ensure this doesnt break the program.
func (rs *PageOrderingRuleset) Correct(manual *SafetyManual) {
	running := true
	for running {
		running = false
		for _, rule := range rs.rules {
			if rule.Correct(manual) {
				running = true
			}
		}
	}
}

// Valid checks if all rules that apply to a SafetyManual are valid.
func (rs *PageOrderingRuleset) Valid(manual *SafetyManual) bool {
	for _, rule := range rs.rules {
		valid, _ := rule.Valid(manual)
		if !valid {
			return false
		}
	}
	return true
}

// SafetyManual is an Entity that represents a safety manual.
type SafetyManual struct {
	indexedPages map[int]PageNumber
}

// NewSafetyManual creates a new SafetyManual.
func NewSafetyManual(pages []PageNumber) *SafetyManual {
	// index pages and store them
	indexedPages := make(map[int]PageNumber)
	for i, page := range pages {
		indexedPages[i] = page
	}
	return &SafetyManual{
		indexedPages: indexedPages,
	}
}

// PageIndex returns the index of a page in the manual.
//
// If the page is not found, it returns -1.
func (m *SafetyManual) PageIndex(page PageNumber) int {
	for i, p := range m.indexedPages {
		if p == page {
			return i
		}
	}
	return -1
}

// MiddlePage returns the middle page of the manual.
func (m *SafetyManual) MiddlePage() PageNumber {
	// get the middle page
	middleIndex := len(m.indexedPages) / 2
	return m.indexedPages[middleIndex]
}

// MovePage moves a page to a new index in the manual, shuffling across the remaining pages.
func (m *SafetyManual) MovePage(page PageNumber, index int) {
	// reorder the indexed pages and skip the page to move
	tempIndexedPages := make(map[int]PageNumber, len(m.indexedPages))
	var j int
	for i := 0; i < len(m.indexedPages); i++ {
		if m.indexedPages[i] == page {
			continue
		}
		tempIndexedPages[j] = m.indexedPages[i]
		j++
	}

	// create a new map to store the indexed pages
	newIndexedPages := make(map[int]PageNumber, len(m.indexedPages))
	for i := 0; i < len(tempIndexedPages)+1; i++ {
		if i == index {
			newIndexedPages[i] = page
		}
		if i > index {
			newIndexedPages[i] = tempIndexedPages[i-1]
		}
		if i < index {
			newIndexedPages[i] = tempIndexedPages[i]
		}
	}

	// update the indexed pages
	m.indexedPages = newIndexedPages
}

func main() {
	inputFile, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	rules := NewPageOrderingRuleset()
	manuals := []*SafetyManual{}

	regexRule := regexp.MustCompile(REGEX_RULE)
	regexManual := regexp.MustCompile(REGEX_MANUAL)

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if regexRule.MatchString(scanner.Text()) {
			// Rule found
			// Extract values
			values := regexRule.FindStringSubmatch(scanner.Text())
			leftValue, err := strconv.Atoi(values[1])
			if err != nil {
				log.Fatal(err)
			}
			rightValue, err := strconv.Atoi(values[2])
			if err != nil {
				log.Fatal(err)
			}

			// Create the rule
			left, err := NewPageNumber(leftValue)
			if err != nil {
				log.Fatal(err)
			}
			right, err := NewPageNumber(rightValue)
			if err != nil {
				log.Fatal(err)
			}
			rule := NewPageOrderingRule(left, right)

			// Add it to our ruleset
			rules.AddRule(rule)
		}
		if regexManual.MatchString(scanner.Text()) {
			// Manual found
			// Extract values
			valuesRaw := regexManual.FindStringSubmatch(scanner.Text())
			values := strings.Split(valuesRaw[0], ",")

			pages := []PageNumber{}
			for _, v := range values {
				pageValue, err := strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
				page, err := NewPageNumber(pageValue)
				if err != nil {
					log.Fatal(err)
				}
				pages = append(pages, page)
			}

			// Create the manual
			manual := NewSafetyManual(pages)
			manuals = append(manuals, manual)
		}
	}

	// Check if the rules are valid for each manual
	// and sum the middle pages for part1 answer
	var sumPart1 int
	badManuals := []*SafetyManual{}
	for _, manual := range manuals {
		if rules.Valid(manual) {
			// get the middle page of the manual
			middlePage := manual.MiddlePage()
			sumPart1 += middlePage.Int()
		} else {
			// All manuals that are not in the correct order
			badManuals = append(badManuals, manual)
		}
	}

	// Print part 1 result
	log.Println("(PART 1) Sum of middle pages:", sumPart1)

	// Part2
	var sumPart2 int
	// Correct the order of the bad manuals
	for _, manual := range badManuals {
		rules.Correct(manual)
		// get the middle page of the manual of any valid
		if rules.Valid(manual) {
			middlePage := manual.MiddlePage()
			sumPart2 += middlePage.Int()
		}
	}

	// Print part 2 result
	log.Println("(PART 2) Sum of middle pages:", sumPart2)
}
