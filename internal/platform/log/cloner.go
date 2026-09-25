package log

type Cloner interface {
	CloneWithOptions(opts ...Option) Logger
}
