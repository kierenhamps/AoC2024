package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay5_PageOrderingRule_NewPageOrderingRule(t *testing.T) {
	left, _ := NewPageNumber(47)
	right, _ := NewPageNumber(53)
	por := NewPageOrderingRule(left, right)
	assert.NotNil(t, por)
}

func TestDay5_PageOrderingRule_Valid(t *testing.T) {
	tests := []struct {
		name        string
		left        PageNumber
		right       PageNumber
		manual      *SafetyManual
		expected    bool
		expectedErr error
	}{
		{
			name:        "rule passes",
			left:        PageNumber{47},
			right:       PageNumber{53},
			manual:      NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}),
			expected:    true,
			expectedErr: nil,
		},
		{
			name:        "rule fails",
			left:        PageNumber{97},
			right:       PageNumber{75},
			manual:      NewSafetyManual([]PageNumber{{75}, {97}, {47}, {61}, {53}}),
			expected:    false,
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			por := NewPageOrderingRule(test.left, test.right)
			result, err := por.Valid(test.manual)
			assert.Equal(t, test.expected, result)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay5_PageNumber_NewPageNumber(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		expectedErr error
	}{
		{"valid input", 47, nil},
		{"input is zero", 0, ErrInputCannotBeZero},
		{"negative input", -5, ErrInputCannotBeNegative},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pn, err := NewPageNumber(test.input)
			assert.NotNil(t, pn)
			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDay5_PageNumber_Int(t *testing.T) {
	tests := []struct {
		name     string
		expected int
	}{
		{"valid input", 47},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pn, _ := NewPageNumber(test.expected)
			assert.Equal(t, test.expected, pn.Int())
		})
	}
}

func TestDay5_SafetyManual_NewSafetyManual(t *testing.T) {
	pages := []PageNumber{{75}, {47}, {61}, {53}, {29}}
	sm := NewSafetyManual(pages)
	assert.NotNil(t, sm)
}

func TestDay5_SafetyManual_GetPageIndex(t *testing.T) {
	tests := []struct {
		name     string
		page     PageNumber
		manual   *SafetyManual
		expected int
	}{
		{"page exists", PageNumber{47}, NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}), 1},
		{"page does not exist", PageNumber{97}, NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}), -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.manual.PageIndex(test.page)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestDay5_SafetyManual_GetMiddlePage(t *testing.T) {
	tests := []struct {
		name     string
		manual   *SafetyManual
		expected PageNumber
	}{
		{"odd number of pages", NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}), PageNumber{61}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.manual.MiddlePage()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestDay5_PageOrderingRuleset_NewPageOrderingRuleset(t *testing.T) {
	pors := NewPageOrderingRuleset()
	assert.NotNil(t, pors)
}

func TestDay5_PageOrderingRuleset_AddRule(t *testing.T) {
	rule := NewPageOrderingRule(PageNumber{47}, PageNumber{53})
	ruleset := NewPageOrderingRuleset()
	ruleset.AddRule(rule)
	assert.NotNil(t, ruleset)
	assert.Equal(t, 1, len(ruleset.rules))
}

func TestDay5_PageOrderingRuleset_Valid(t *testing.T) {
	tests := []struct {
		name     string
		ruleset  *PageOrderingRuleset
		manual   *SafetyManual
		expected bool
	}{
		{
			name: "manual is valid",
			ruleset: &PageOrderingRuleset{
				rules: []PageOrderingRule{
					{PageNumber{47}, PageNumber{53}},
					{PageNumber{97}, PageNumber{75}},
				}},
			manual:   NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}),
			expected: true,
		},
		{
			name: "manual is invalid",
			ruleset: &PageOrderingRuleset{
				rules: []PageOrderingRule{
					{PageNumber{47}, PageNumber{53}},
					{PageNumber{97}, PageNumber{75}},
				}},
			manual:   NewSafetyManual([]PageNumber{{75}, {97}, {47}, {61}, {53}}),
			expected: false,
		},
		{
			name: "manual is valid with no rules",
			ruleset: &PageOrderingRuleset{
				rules: []PageOrderingRule{}},
			manual:   NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}),
			expected: true,
		},
		{
			name: "manual is valid with no matching rules",
			ruleset: &PageOrderingRuleset{
				rules: []PageOrderingRule{
					{PageNumber{75}, PageNumber{97}},
				}},
			manual:   NewSafetyManual([]PageNumber{{75}, {47}, {61}, {53}, {29}}),
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.ruleset.Valid(test.manual)
			assert.Equal(t, test.expected, result)
		})
	}
}
