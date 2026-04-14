package rude

import (
	"fmt"
	"maps"
)

type Error struct {
	Err     error
	Type    ErrorType
	Code    int
	Message string
	Meta    map[string]any
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s (code=%d)", e.Type, e.Message, e.Code)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func WrapError(e *Error, err error) *Error {
	if err == nil {
		return nil
	}
	if e == nil {
		return &Error{Err: err, Type: TypeInternal, Message: err.Error()}
	}

	message := e.Message
	if message == "" {
		message = err.Error()
	} else {
		message = fmt.Sprintf("%s: %v", e.Message, err)
	}

	return &Error{
		Err:     err,
		Type:    e.Type,
		Code:    e.Code,
		Message: message,
		Meta:    maps.Clone(e.Meta),
	}
}
