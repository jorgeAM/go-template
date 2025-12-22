package log

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func fromContext(ctx context.Context) CloneableLogger {
	if logger, ok := ctx.Value(loggerKey{}).(CloneableLogger); ok {
		return logger
	}

	return nil
}

func getLogger(ctx context.Context) CloneableLogger {
	if logger := fromContext(ctx); logger != nil {
		return logger
	}

	if defaultLogger != nil {
		return defaultLogger
	}

	return nil
}

func ContextWithLogger(ctx context.Context, logger CloneableLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger.CloneWithOptions(WithCallerSkip(2)))
}

func AddOptions(ctx context.Context, opts ...Option) context.Context {
	logger := getLogger(ctx)
	if logger == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerKey{}, logger.CloneWithOptions(opts...))
}

func WithFields(fields ...zap.Field) Option {
	return func(o *option) {
		o.fields = append(o.fields, fields...)
	}
}

func Info(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(InfoLevel))

	logger.Info(message, opts...)
}

func Warn(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(WarnLevel))

	logger.Warn(message, opts...)
}

func Error(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(ErrorLevel))

	logger.Error(message, opts...)
}

func Debug(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(DebugLevel))

	logger.Debug(message, opts...)
}

func Fatal(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(FatalLevel))

	logger.Fatal(message, opts...)
}

func Panic(ctx context.Context, message string, opts ...Option) {
	logger := getLogger(ctx)
	if logger == nil {
		return
	}

	opts = append(opts, WithLevel(PanicLevel))

	logger.Panic(message, opts...)
}
