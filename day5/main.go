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

// Valid checks if a rule is valid for a SafetyManual.
func (r *PageOrderingRule) Valid(manual *SafetyManual) (bool, error) {
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

// Valid checks if all rules that apply to a SafetyManual are valid.
func (rs *PageOrderingRuleset) Valid(manual *SafetyManual) bool {
	for _, rule := range rs.rules {
		// skip rule if pages are not in the manual
		if manual.PageIndex(rule.left) == -1 || manual.PageIndex(rule.right) == -1 {
			continue
		}
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
	for _, manual := range manuals {
		if rules.Valid(manual) {
			// get the middle page of the manual
			middlePage := manual.MiddlePage()
			sumPart1 += middlePage.Int()
		}
	}

	// Print part 1 result
	log.Println("(PART 1) Sum of middle pages:", sumPart1)
}
