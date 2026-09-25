package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*ZapLogger)(nil)
var _ Cloner = (*ZapLogger)(nil)

type ZapLogger struct {
	logger *zap.Logger
	option *option
}

func NewZapLogger(opts ...Option) (*ZapLogger, error) {
	options := &option{
		level:    InfoLevel,
		encoding: "json",
		fields: []Field{
			zap.String("env", os.Getenv("APP_ENV")),
			zap.String("service", os.Getenv("APP_NAME")),
		},
		callerSkip: 1,
	}

	options = options.apply(opts...)

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger.WithOptions(zap.AddCallerSkip(options.callerSkip)),
		option: options,
	}, nil
}

func (z *ZapLogger) Debug(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != DebugLevel {
		return
	}

	z.logger.Debug(msg, option.fields...)
}

func (z *ZapLogger) Error(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != ErrorLevel {
		return
	}

	z.logger.Error(msg, option.fields...)
}

func (z *ZapLogger) Fatal(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != FatalLevel {
		return
	}

	z.logger.Fatal(msg, option.fields...)
}

func (z *ZapLogger) Info(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != InfoLevel {
		return
	}

	z.logger.Info(msg, option.fields...)
}

func (z *ZapLogger) Panic(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != PanicLevel {
		return
	}

	z.logger.Panic(msg, option.fields...)
}

func (z *ZapLogger) Warn(msg string, opts ...Option) {
	option := z.option.apply(opts...)

	if option.level != WarnLevel {
		return
	}

	z.logger.Warn(msg, option.fields...)
}

func (z *ZapLogger) CloneWithOptions(opts ...Option) Logger {
	logger := z.logger
	options := z.option.apply(opts...)

	if callerSkipDiff := options.callerSkip - z.option.callerSkip; callerSkipDiff != 0 {
		logger = logger.WithOptions(zap.AddCallerSkip(callerSkipDiff))
	}

	return &ZapLogger{
		logger: logger,
		option: options,
	}
}
