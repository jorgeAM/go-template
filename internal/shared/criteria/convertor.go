package criteria

import "context"

type Converter interface {
	Convert(ctx context.Context, source string, criteria *Criteria) (string, []interface{}, error)
}
