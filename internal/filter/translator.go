package filter

import (
	"fmt"
	"regexp"
	"strings"
)

type ExprType int

const (
	Str ExprType = iota
	Num
	Lighting
	Ammenity
)

const StrRegex = `^(=|has:)((?:\w|\s|,)+)$`
const NumRegex = `^(<|>|=|>=|<=)([+-]?(?:[0-9]+[.])?[0-9]+)$`
const LightingRegex = `^(=)(low|medium|high)$`
const AmmenityRegex = `^(=|has:)(yard|pool|garage|rooftop|waterfront)$`
const Separator = ";"

type filterExpr struct {
	Operator string
	Value    string
}

func TranslateToSql(field string, expr string, exprType ExprType) (string, error) {
	var sqlCondition string
	var err error

	switch exprType {
	case Str:
		sqlCondition, err = translateFilterExpr(field, expr, StrRegex, translateStrExpr)
	case Num:
		sqlCondition, err = translateFilterExpr(field, expr, NumRegex, translateNumExpr)
	case Lighting:
		sqlCondition, err = translateFilterExpr(field, expr, LightingRegex, translateStrExpr)
	case Ammenity:
		sqlCondition, err = translateFilterExpr(field, expr, AmmenityRegex, translateStrExpr)
	}

	if err != nil {
		return "", err
	}
	return sqlCondition, nil
}

func translateFilterExpr(
	field string, filterExpr string, regex string, translatorFunc func(string, filterExpr) string,
) (string, error) {
	var translation string
	expressions, err := splitExpr(filterExpr, regex)
	if err != nil {
		return "", err
	}

	for _, e := range expressions {
		translation += fmt.Sprintf("%s and ", translatorFunc(field, e))
	}

	translation = strings.TrimSuffix(translation, " and ")
	return translation, nil
}

func splitExpr(expr string, regExpr string) ([]filterExpr, error) {
	var expressions []filterExpr = make([]filterExpr, 0)
	parts := strings.Split(expr, Separator)
	regex := regexp.MustCompile(regExpr)

	for _, p := range parts {
		match := regex.FindStringSubmatch(p)
		// Should match 3 parts, the entire string, the operator and the value
		if match == nil || len(match) != 3 {
			return []filterExpr{}, fmt.Errorf(`filter expression "%s" in "%s" is not valid`, p, expr)
		}

		expressions = append(expressions, filterExpr{Operator: match[1], Value: match[2]})
	}

	return expressions, nil
}

func translateNumExpr(field string, e filterExpr) string {
	return fmt.Sprintf("%s%s%s", field, e.Operator, e.Value)
}

func translateStrExpr(field string, e filterExpr) string {
	if e.Operator == "=" {
		return fmt.Sprintf("lower(%s)=lower('%s')", field, e.Value)
	}

	return fmt.Sprintf("lower(%s) like lower('%%%s%%')", field, e.Value)
}
