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
const DistanceRegex = `^distance\(([+-]?(?:[0-9]+[.])?[0-9]+),([+-]?(?:[0-9]+[.])?[0-9]+)\)(?:(<|>|=|>=|<=)([+-]?(?:[0-9]+[.])?[0-9]+))?$`
const Separator = ";"

type filterExpr struct {
	Operator string
	Value    string
}

type DistanceFilterData struct {
	X   string
	Y   string
	Sql string
}

type Translator struct {
	Translations []string
	Err          error
}

func (translator *Translator) Init() {
	translator.Translations = make([]string, 0)
	translator.Err = nil
}

func (translator *Translator) Translate(field string, expr string, exprType ExprType) {
	if translator.Err != nil || field == "" || expr == "" {
		return
	}

	var t string
	t, translator.Err = TranslateToSql(field, expr, exprType)
	if translator.Err != nil {
		return
	}
	translator.Translations = append(translator.Translations, t)
}

func (translator *Translator) GetSqlTranslation() string {
	return strings.Join(translator.Translations, " and ")
}

func (translator *Translator) TranslateDistanceExpr(field string, expr string) DistanceFilterData {
	if translator.Err != nil || field == "" || expr == "" {
		return DistanceFilterData{}
	}

	regex := regexp.MustCompile(DistanceRegex)
	match := regex.FindStringSubmatch(expr)
	fmt.Println(len(match))
	// Should match 3 or 5 parts
	if match == nil || (len(match) != 3 && len(match) != 5) {
		translator.Err = fmt.Errorf(`distance expression "%s" is not valid`, expr)
		return DistanceFilterData{}
	}

	var data DistanceFilterData
	data.X = match[1]
	data.Y = match[2]
	// Check if additional operator and value has been provided besides the distance()
	if len(match) == 5 && match[3] != "" && match[4] != "" {
		data.Sql = fmt.Sprintf("%s %s %s", field, match[3], match[4])
		translator.Translations = append(translator.Translations, data.Sql)
	}

	return data
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
