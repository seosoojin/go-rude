package rude

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

type Error struct {
	Err      error          `json:"-"`
	Type     ErrorType      `json:"type,omitempty"`
	Code     int            `json:"code,omitempty"`
	Message  string         `json:"message,omitempty"`
	MetaData map[string]any `json:"metadata,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s] %s (code=%d)", e.Type, e.Message, e.Code)
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) WithMetadata(k string, v any) Error {
	if e.MetaData == nil {
		e.MetaData = make(map[string]any)
	}
	e.MetaData[k] = v
	return e
}

func NewError(errType ErrorType, code int, message string) *Error {
	return &Error{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

func WrapError(e Error, err error) Error {
	if err == nil {
		return e
	}

	message := e.Message
	if message == "" {
		message = err.Error()
	} else {
		message = fmt.Sprintf("%s: %v", e.Message, err)
	}

	return Error{
		Err:      err,
		Type:     e.Type,
		Code:     e.Code,
		Message:  message,
		MetaData: maps.Clone(e.MetaData),
	}
}

func (e Error) Write(w http.ResponseWriter, r *http.Request) {
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
