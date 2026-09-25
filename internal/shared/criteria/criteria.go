package criteria

type CriteriaInput struct {
	OrderBy   *string `in:"query=order_by"`
	OrderType *string `in:"query=order_type"`
	Page      int     `in:"query=page"`
	PageSize  int     `in:"query=page_size"`
}

type CriteriaPrimitive struct {
	Filters   []*FilterPrimitive
	OrderBy   *string
	OrderType *string
	Page      int
	PageSize  int
}

type Criteria struct {
	Filters  []*Filter
	Order    *Order
	Page     Page
	PageSize PageSize
}

func FromPrimitive(primitive *CriteriaPrimitive) (*Criteria, error) {
	if primitive == nil {
		return &Criteria{}, nil
	}

	filters := make([]*Filter, 0, len(primitive.Filters))

	for _, filter := range primitive.Filters {
		f, err := NewFilter(filter.Field, filter.Operator, filter.Value)
		if err != nil {
			return nil, err
		}

		filters = append(filters, f)
	}

	var orderCriteria *Order

	if primitive.OrderBy != nil {
		if primitive.OrderType != nil {
			order, err := NewOrder(*primitive.OrderBy, *primitive.OrderType)
			if err != nil {
				return nil, err
			}

			orderCriteria = order
		}
	}

	return NewCriteria(filters, orderCriteria, primitive.Page, primitive.PageSize)
}

func NewCriteria(filters []*Filter, order *Order, page, pageSize int) (*Criteria, error) {
	criteria := &Criteria{
		Filters: filters,
		Order:   order,
	}

	if page == 0 || pageSize == 0 {
		return criteria, nil
	}

	pageVO, err := NewPage(page)
	if err != nil {
		return nil, err
	}

	criteria.Page = pageVO

	pageSizeVO, err := NewPageSize(pageSize)
	if err != nil {
		return nil, err
	}

	criteria.PageSize = pageSizeVO

	return criteria, nil
}

func (c *Criteria) HasOrder() bool {
	return c.Order != nil
}

func (c *Criteria) HasFilter() bool {
	return len(c.Filters) != 0
}
