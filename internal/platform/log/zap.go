package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(service, env string) (*ZapLogger, error) {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *ZapLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *ZapLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *ZapLogger) Panic(msg string) {
	l.logger.Panic(msg)
}
