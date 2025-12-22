package log

type Logger interface {
	Info(msg string, opts ...Option)
	Warn(msg string, opts ...Option)
	Error(msg string, opts ...Option)
	Debug(msg string, opts ...Option)
	Fatal(msg string, opts ...Option)
	Panic(msg string, opts ...Option)
}

type CloneableLogger interface {
	Logger
	Cloner
}
