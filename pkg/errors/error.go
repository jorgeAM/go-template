package errors

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	code     *ErrorCode
	message  string
	cause    error
	metadata Metadata
}

func New(errorCode *ErrorCode, message string, metadata ...Metadata) *Error {
	m := NewMetadata()
	for _, metadata := range metadata {
		m = m.Merge(metadata)
	}

	return &Error{
		code:     errorCode,
		message:  message,
		metadata: m,
	}
}

func Wrap(errorCode *ErrorCode, cause error, message string, metadata ...Metadata) *Error {
	err := New(errorCode, message, metadata...)
	err.cause = cause

	return err
}

func Is(err error, target error) bool {
	if err == nil || target == nil {
		return false
	}

	if err, ok := err.(*Error); ok {
		return err.is(target)
	}

	return errors.Is(err, target)
}

func (err *Error) Code() *ErrorCode {
	return err.code
}

func (err *Error) Message() string {
	return err.message
}

func (err *Error) Cause() error {
	return err.cause
}

func (err *Error) Metadata() Metadata {
	return err.metadata
}

func (err *Error) Unwrap() error {
	return err.cause
}

func (err *Error) Error() string {
	if err.code == nil || err.cause == nil {
		return err.message
	}

	var sb strings.Builder

	sb.WriteString(err.code.code)
	sb.WriteString(": ")
	sb.WriteString(err.message)

	if err.cause != nil {
		sb.WriteString(fmt.Sprintf(" (%s)", err.cause.Error()))
	}

	if len(err.metadata) > 0 {
		metadataStr := make([]string, 0, len(err.metadata))
		for k, v := range err.metadata {
			metadataStr = append(metadataStr, fmt.Sprintf("[%s = %v]", k, v))
		}

		sb.WriteString(fmt.Sprintf(" %s", strings.Join(metadataStr, ", ")))
	}

	return sb.String()
}

func (err *Error) is(target error) bool {
	if target == nil {
		return false
	}

	if t, ok := target.(*Error); ok {
		target = t.code
	}

	if t, ok := target.(*ErrorCode); ok {
		if err.code == t {
			return true
		}

		if err.cause != nil {
			if cause, ok := err.cause.(*Error); ok {
				return cause.is(t)
			}
		}
	}

	return false
}
