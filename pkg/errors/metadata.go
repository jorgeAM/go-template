package errors

type Metadata map[string]interface{}

func NewMetadata() Metadata {
	return make(Metadata)
}

func WithMetadata(key string, value interface{}) Metadata {
	return NewMetadata().And(key, value)
}

func (m Metadata) And(key string, value interface{}) Metadata {
	m[key] = value
	return m
}

func (m Metadata) Merge(other Metadata) Metadata {
	for k, v := range other {
		m[k] = v
	}

	return m
}
