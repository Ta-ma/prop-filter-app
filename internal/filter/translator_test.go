/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateStrExpr(t *testing.T) {
	type testCase struct {
		fieldName string
		e         filterExpr
		expected  string
	}

	testCases := []testCase{
		{fieldName: "description", e: filterExpr{Operator: "=", Value: "alaska"},
			expected: "lower(description)=lower('alaska')"},
		{fieldName: "desc", e: filterExpr{Operator: ":has", Value: "alaska"},
			expected: "lower(desc) like lower('%alaska%')"},
		{fieldName: "descriptiON", e: filterExpr{Operator: ":has", Value: "alASka"},
			expected: "lower(descriptiON) like lower('%alASka%')"},
	}

	for _, test := range testCases {
		actual := translateStrExpr(test.fieldName, test.e)

		assert.Equal(t, test.expected, actual)
	}
}

func TestTranslateNumExpr(t *testing.T) {
	type testCase struct {
		fieldName string
		e         filterExpr
		expected  string
	}

	testCases := []testCase{
		{fieldName: "price", e: filterExpr{Operator: "<", Value: "150458"},
			expected: "price<150458"},
		{fieldName: "price", e: filterExpr{Operator: ">", Value: "100000.63"},
			expected: "price>100000.63"},
		{fieldName: "price", e: filterExpr{Operator: ">=", Value: "0.158"},
			expected: "price>=0.158"},
		{fieldName: "price", e: filterExpr{Operator: "<=", Value: "50000.0"},
			expected: "price<=50000.0"},
	}

	for _, test := range testCases {
		actual := translateNumExpr(test.fieldName, test.e)

		assert.Equal(t, test.expected, actual)
	}
}

func TestSplitExpr(t *testing.T) {
	type testCase struct {
		expr        string
		regExpr     string
		expected    []filterExpr
		errExpected bool
	}

	testCases := []testCase{
		{expr: "<999", regExpr: NumRegex,
			expected: []filterExpr{{Operator: "<", Value: "999"}}, errExpected: false},
		{expr: ">10000;<=20000", regExpr: NumRegex,
			expected: []filterExpr{{Operator: ">", Value: "10000"}, {Operator: "<=", Value: "20000"}}, errExpected: false},
		{expr: "has:yard;has:pool", regExpr: StrRegex,
			expected: []filterExpr{{Operator: "has:", Value: "yard"}, {Operator: "has:", Value: "pool"}}, errExpected: false},
		{expr: "has:yard;=test", regExpr: StrRegex,
			expected: []filterExpr{{Operator: "has:", Value: "yard"}, {Operator: "=", Value: "test"}}, errExpected: false},
		{expr: "has;yard", regExpr: StrRegex, expected: []filterExpr{}, errExpected: true},
		{expr: "has::yard;has:pool", regExpr: StrRegex, expected: []filterExpr{}, errExpected: true},
		{expr: "has:yard;=test;", regExpr: StrRegex, expected: []filterExpr{}, errExpected: true},
	}

	for _, test := range testCases {
		actual, err := splitExpr(test.expr, test.regExpr)

		if test.errExpected {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expected, actual)
	}
}
