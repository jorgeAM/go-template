package log

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	Fatal(msg string)
	Panic(msg string)
}
