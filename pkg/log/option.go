package log

import (
	"time"

	"go.uber.org/zap"
)

type Option func(*option)

type Field = zap.Field

type option struct {
	level      Level
	encoding   string
	callerSkip int
	fields     []Field
}

func (o *option) apply(opts ...Option) *option {
	cloned := o.clone()

	for _, opt := range opts {
		opt(cloned)
	}

	return cloned
}

func (o *option) clone() *option {
	cloned := *o
	return &cloned
}

func WithLevel(level Level) Option {
	return func(o *option) {
		o.level = level
	}
}

func WithEncoding(encoding string) Option {
	return func(o *option) {
		o.encoding = encoding
	}
}

func WithCallerSkip(skip int) Option {
	return func(o *option) {
		o.callerSkip = skip
	}
}

func WithBool(key string, value bool) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Bool(key, value))
	}
}

func WithBinary(key string, value []byte) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Binary(key, value))
	}
}

func WithDuration(key string, value time.Duration) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Duration(key, value))
	}
}

func WithError(err error) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Error(err))
	}
}

func WithFloat32(key string, value float32) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Float32(key, value))
	}
}

func WithFloat64(key string, value float64) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Float64(key, value))
	}
}

func WithInt(key string, value int) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Int(key, value))
	}
}

func WithInt32(key string, value int32) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Int32(key, value))
	}
}

func WithInt64(key string, value int64) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Int64(key, value))
	}
}

func WithStack(key string) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Stack(key))
	}
}

func WithString(key string, value string) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.String(key, value))
	}
}

func WithTime(key string, value time.Time) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Time(key, value))
	}
}

func WithObject(key string, value interface{}) Option {
	return func(o *option) {
		o.fields = append(o.fields, zap.Any(key, value))
	}
}
