package log

var defaultLogger CloneableLogger

func InitDefaultLogger(opts ...Option) error {
	logger, err := NewZapLogger(opts...)
	if err != nil {
		return err
	}

	defaultLogger = logger

	return nil
}
