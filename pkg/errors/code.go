package errors

import (
	"net/http"
	"strings"
)

var (
	ErrInvalidCode = Define("code_error.invalid_code")
)

type Code string

const (
	InternalCode   Code = "INTERNAL_ERROR"
	BadRequestCode Code = "BAD_REQUEST"
	NotFoundCode   Code = "NOT_FOUND"
	Unauthorized   Code = "UNAUTHORIZED"
)

var allowedCode = map[string]Code{
	InternalCode.String():   InternalCode,
	BadRequestCode.String(): BadRequestCode,
	NotFoundCode.String():   NotFoundCode,
	Unauthorized.String():   Unauthorized,
}

var codeToHttpStatusCode = map[Code]int{
	InternalCode:   http.StatusInternalServerError,
	BadRequestCode: http.StatusBadRequest,
	NotFoundCode:   http.StatusNotFound,
	Unauthorized:   http.StatusUnauthorized,
}

func NewCode(code string) (Code, error) {
	if errorCode, ok := allowedCode[strings.ToUpper(code)]; ok {
		return errorCode, nil
	}

	return "", New(ErrInvalidCode, "invalid code of error", WithMetadata("code", code))
}

func (c Code) String() string {
	return string(c)
}

func (c Code) HttpStatus() int {
	return codeToHttpStatusCode[c]
}
