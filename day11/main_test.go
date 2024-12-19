package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDay11_Rule_Eval(t *testing.T) {
	tests := []struct {
		name            string
		rule            Rule
		input           Stone
		expected        []Stone
		expectedChanged bool
	}{
		{"Rule zero to one applicable", &RuleZeroToOne{}, Stone(0), []Stone{1}, true},
		{"Rule zero to one not applicable", &RuleZeroToOne{}, Stone(124515), []Stone{124515}, false},
		{"Rule split even digits applicable", &RuleSplitEvenDigits{}, 123321, []Stone{123, 321}, true},
		{"Rule split even digits applicable", &RuleSplitEvenDigits{}, 12332, []Stone{12332}, false},
		{"Rule multiply by 2024 always applicable", &RuleMultiplyBy2024{}, 1, []Stone{2024}, true},
		{"Rule multiply by 2024 always applicable", &RuleMultiplyBy2024{}, 2, []Stone{4048}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, changed := test.rule.Eval(test.input)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedChanged, changed)
		})
	}
}

func TestDay11_Blink(t *testing.T) {
	tests := []struct {
		name     string
		input    []Stone
		rules    []Rule
		expected []Stone
	}{
		{"1 stone 1 rule", []Stone{0}, []Rule{&RuleZeroToOne{}}, []Stone{1}},
		{"1 stone all rules only one applies", []Stone{0}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{1}},
		{"1 stone all rules only one applies", []Stone{123321}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{123, 321}},
		{"1 stone all rules only one applies", []Stone{1}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{2024}},
		{"1 stone all rules two apply", []Stone{0000}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{1}},
		{"example 1 first blink", []Stone{0, 1, 10, 99, 999}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{1, 2024, 1, 0, 9, 9, 2021976}},
		{"example 2 first blink", []Stone{125, 17}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{253000, 1, 7}},
		{"example 2 second blink", []Stone{253000, 1, 7}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{253, 0, 2024, 14168}},
		{"example 2 third blink", []Stone{253, 0, 2024, 14168}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{512072, 1, 20, 24, 28676032}},
		{"example 2 fourth blink", []Stone{512072, 1, 20, 24, 28676032}, []Rule{&RuleZeroToOne{}, &RuleSplitEvenDigits{}, &RuleMultiplyBy2024{}}, []Stone{512, 72, 2024, 2, 0, 2, 4, 2867, 6032}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Blink(test.input, test.rules)
			assert.Equal(t, test.expected, result)
		})
	}
}
