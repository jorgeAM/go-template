package middleware

type loggerOptions struct {
	ignoredPaths map[string]struct{}
}

type Option func(*loggerOptions)

func (o *loggerOptions) apply(opts ...Option) *loggerOptions {
	cloned := o.clone()

	for _, opt := range opts {
		opt(cloned)
	}

	return cloned
}

func (o *loggerOptions) clone() *loggerOptions {
	cloned := *o
	return &cloned
}

func WithIgnoreRoutes(paths ...string) Option {
	return func(o *loggerOptions) {
		for _, path := range paths {
			o.ignoredPaths[path] = struct{}{}
		}
	}
}
