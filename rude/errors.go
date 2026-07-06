package rude

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
)

type Error struct {
	Err      error                 `json:"-"`
	Type     ErrorType             `json:"type,omitempty"`
	Code     int                   `json:"code,omitempty"`
	Message  string                `json:"message,omitempty"`
	MetaData map[string]any        `json:"metadata,omitempty"`
	builder  func() ProblemDetails `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s (code=%d)", e.Type, e.Message, e.Code)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) WithMetadata(k string, v any) *Error {
	if e.MetaData == nil {
		e.MetaData = make(map[string]any)
	}
	e.MetaData[k] = v
	return e
}

func (e *Error) WithProblemBuilder(builder func() ProblemDetails) *Error {
	e.builder = builder
	return e
}

func NewError(errType ErrorType, code int, message string) *Error {
	return &Error{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

func WrapError(e *Error, err error) *Error {
	if err == nil {
		return e
	}

	message := e.Message
	if message == "" {
		message = err.Error()
	} else {
		message = fmt.Sprintf("%s: %v", e.Message, err)
	}

	if ee, ok := errors.AsType[*Error](err); ok {
		metadata := maps.Clone(e.MetaData)
		if len(metadata) == 0 {
			metadata = nil
		}

		maps.Copy(metadata, ee.MetaData)

		return &Error{
			Err:      ee.Err,
			Type:     ee.Type,
			Code:     ee.Code,
			Message:  message,
			MetaData: metadata,
			builder:  e.builder,
		}
	}

	return &Error{
		Err:      err,
		Type:     e.Type,
		Code:     e.Code,
		Message:  message,
		MetaData: maps.Clone(e.MetaData),
		builder:  e.builder,
	}
}

func (e *Error) Write(w http.ResponseWriter, r *http.Request) {
	if e.Type == "" {
		e.Type = "about:blank"
	}

	if e.Code == 0 {
		e.Code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(e.Code)

	_ = json.NewEncoder(w).Encode(e)
}

func (e *Error) ToProblemDetails() ProblemDetails {
	if e.builder != nil {
		return e.builder()
	}
	detail := ""
	if e.Err != nil {
		detail = e.Err.Error()
	} else if e.Message != "" {
		detail = e.Message
	}

	return ProblemDetails{
		Type:       string(e.Type),
		Title:      e.Message,
		Status:     e.Code,
		Detail:     detail,
		Extensions: maps.Clone(e.MetaData),
	}
}
