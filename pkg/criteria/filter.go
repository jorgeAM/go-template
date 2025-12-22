package criteria

type Filter struct {
	Field    string
	Operator Operator
	Value    any
}

type FilterPrimitive struct {
	Field    string
	Operator string
	Value    any
}

func NewFilter(field, operator string, value any) (*Filter, error) {
	op, err := NewOperator(operator)
	if err != nil {
		return nil, err
	}

	return &Filter{
		Field:    field,
		Operator: op,
		Value:    value,
	}, nil
}
