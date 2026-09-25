package criteria

import (
	"context"
	"errors"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

var _ Converter = (*CriteriaToPostgresConverter)(nil)

type CriteriaToPostgresConverter struct{}

func NewCriteriaToPostgresConverter() *CriteriaToPostgresConverter {
	return &CriteriaToPostgresConverter{}
}

func (c *CriteriaToPostgresConverter) Convert(ctx context.Context, source string, criteria *Criteria) (string, []interface{}, error) {
	dialect := goqu.Dialect("postgres")
	ds := dialect.From(source)
	exprs := exp.NewExpressionList(exp.AndType)

	if criteria.HasFilter() {
		for _, filter := range criteria.Filters {
			exp, err := c.generateWhereQuery(filter)
			if err != nil {
				return "", nil, err
			}

			exprs = exprs.Append(exp)
		}
	}

	if criteria.HasOrder() {
		exp, err := c.generateOrderQuery(criteria.Order)
		if err != nil {
			return "", nil, err
		}

		ds = ds.Order(exp)
	}

	if criteria.PageSize > 0 {
		ds = ds.Limit(uint(criteria.PageSize.Raw()))

		if criteria.Page > 0 {
			ds = ds.Offset(uint(criteria.PageSize.Raw() * (criteria.Page.Raw() - 1)))
		}
	}

	return ds.Prepared(true).Where(exprs).ToSQL()
}

var operatorMap = map[Operator]func(field string, value interface{}) exp.Expression{
	EQUAL: func(field string, value interface{}) exp.Expression { return goqu.L(field).Eq(value) },
	GT:    func(field string, value interface{}) exp.Expression { return goqu.L(field).Gt(value) },
	LT:    func(field string, value interface{}) exp.Expression { return goqu.L(field).Lt(value) },
}

func (c *CriteriaToPostgresConverter) generateWhereQuery(filter *Filter) (exp.Expression, error) {
	if fn, ok := operatorMap[filter.Operator]; ok {
		return fn(filter.Field, filter.Value), nil
	}

	return nil, errors.New("unsupported filter operator")
}

var operatorTypeMap = map[OrderType]func(orderBy string) exp.OrderedExpression{
	ASC:  func(orderBy string) exp.OrderedExpression { return goqu.I(orderBy).Asc() },
	DESC: func(orderBy string) exp.OrderedExpression { return goqu.I(orderBy).Desc() },
}

func (c *CriteriaToPostgresConverter) generateOrderQuery(order *Order) (exp.OrderedExpression, error) {
	if fn, ok := operatorTypeMap[order.OrderType]; ok {
		return fn(order.OrderBy), nil
	}

	return nil, errors.New("unsupported order type")
}
