package log

import (
	"go.uber.org/zap"
)

type Level int

const (
	DebugLevel = Level(zap.DebugLevel)
	InfoLevel  = Level(zap.InfoLevel)
	WarnLevel  = Level(zap.WarnLevel)
	ErrorLevel = Level(zap.ErrorLevel)
	FatalLevel = Level(zap.FatalLevel)
	PanicLevel = Level(zap.PanicLevel)
)
