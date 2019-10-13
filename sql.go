package scim

import (
	"fmt"
	"strings"

	filter "github.com/di-wu/scim-filter-parser"
)

const (
	placeholder = "?"
)

var comparatorMap = map[filter.Token]string{
	filter.EQ:  "=",
	filter.NE:  "!=",
	filter.CO:  "LIKE",
	filter.SW:  "LIKE",
	filter.EW:  "LIKE",
	filter.PR:  "IS NOT NULL",
	filter.GT:  ">",
	filter.GE:  ">=",
	filter.LT:  "<",
	filter.LE:  "<=",
	filter.AND: "AND",
	filter.OR:  "OR",
	filter.NOT: "NOT",
}

var likeExpressionMap = map[filter.Token]likeExpression{
	filter.CO: likeExpression{Prefix: "%", Suffix: "%"},
	filter.SW: likeExpression{Prefix: "", Suffix: "%"},
	filter.EW: likeExpression{Prefix: "%", Suffix: ""},
}

type (
	likeExpression struct {
		Prefix string
		Suffix string
	}

	// Parser provides an interface for transpiling SCIM filter expression.
	Parser struct {
		AttributeMap map[string]string
		params       []interface{}
	}
)

// NewParser builds a new filter parser from a filter expression.
func NewParser(attributeMap map[string]string) *Parser {
	return &Parser{
		AttributeMap: attributeMap,
	}
}

// ToSqlFromString parses raw SCIM filter and then transpiles it to a SQL query.
// You can optionally customize the select query, it will by default select "*".
func (p *Parser) ToSqlFromString(rawFilter string) (query string, params []interface{}, err error) {
	var (
		exp filter.Expression
	)

	rawFilter = strings.TrimSpace(rawFilter)

	if rawFilter != "" {
		parser := filter.NewParser(strings.NewReader(rawFilter))
		exp, err = parser.Parse()

		if err != nil {
			return
		}
	}

	query, params = p.ToSql(exp)
	return
}

// ToSql transpiles parsed filter to a SQL query. You can optionally customize the select query, it will by default
// select "*".
func (p *Parser) ToSql(expression filter.Expression) (query string, params []interface{}) {
	whereClause := p.process(expression, p.AttributeMap)

	return whereClause, p.params
}

func (p *Parser) process(exp filter.Expression, attrMap map[string]string) string {
	if attrExp, ok := exp.(filter.AttributeExpression); ok {
		return p.processAttributeStatement(attrExp, attrMap)
	}

	if biExp, ok := exp.(filter.BinaryExpression); ok {
		return fmt.Sprintf(
			"(%s %s %s)",
			p.process(biExp.X, attrMap),
			getComparator(biExp.CompareOperator),
			p.process(biExp.Y, attrMap),
		)
	}

	if uniExp, ok := exp.(filter.UnaryExpression); ok {
		return fmt.Sprintf("%s (%s)", getComparator(uniExp.CompareOperator), p.process(uniExp.X, attrMap))
	}

	// Should never happen but handled nonetheless
	panic("unsupported expression type")
}

func (p *Parser) processAttributeStatement(exp filter.AttributeExpression, attrMap map[string]string) string {
	path := p.processAttributePath(exp.AttributePath, attrMap)
	comparator := getComparator(exp.CompareOperator)
	value := p.processAttributeValue(exp.CompareValue, exp.CompareOperator)

	return fmt.Sprintf(`%s %s %s`, path, comparator, value)
}

func (p *Parser) processAttributePath(path string, attrMap map[string]string) string {
	for k, v := range attrMap {
		if strings.ToLower(k) == strings.ToLower(path) {
			return v
		}
	}

	return path
}

func (p *Parser) processAttributeValue(value string, op filter.Token) string {
	tokens, ok := likeExpressionMap[op]

	if !ok {
		tokens = likeExpression{}
	}

	p.params = append(p.params, fmt.Sprintf("%s%s%s", tokens.Prefix, value, tokens.Suffix))

	return placeholder
}

func getComparator(comparator filter.Token) string {
	if v, ok := comparatorMap[comparator]; ok {
		return v
	}

	// Should never happen but handled nonetheless
	panic("unsupported expression comparator")
}
