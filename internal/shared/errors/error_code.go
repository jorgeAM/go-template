package errors

type ErrorCode struct {
	code string
}

func (c *ErrorCode) Error() string {
	return c.code
}

func Define(code string) *ErrorCode {
	if len(code) == 0 {
		panic("empty error code")
	}

	return &ErrorCode{code}
}
