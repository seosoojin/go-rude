package rude

import (
	"net/http"
	"sync"
)

var registry = map[ErrorType]*ProblemDetails{
	TypeNotFound: {
		Type:   "/not-found",
		Title:  http.StatusText(http.StatusNotFound),
		Status: http.StatusNotFound,
	},
	TypeUnauthorized: {
		Type:   "/unauthorized",
		Title:  http.StatusText(http.StatusUnauthorized),
		Status: http.StatusUnauthorized,
	},
	TypeValidation: {
		Type:   "/validation-error",
		Title:  http.StatusText(http.StatusBadRequest),
		Status: http.StatusBadRequest,
	},
	TypeInternal: {
		Type:   "/internal-error",
		Title:  http.StatusText(http.StatusInternalServerError),
		Status: http.StatusInternalServerError,
	},
}

var registryMu sync.RWMutex

func RegisterErrorType(t ErrorType, p ProblemDetails) {
	if t == "" {
		return
	}
	registryMu.Lock()
	defer registryMu.Unlock()

	registry[t] = &p
}

func getProblemByType(t ErrorType) *ProblemDetails {
	registryMu.RLock()
	p, ok := registry[t]
	registryMu.RUnlock()

	if !ok || p == nil {
		registryMu.RLock()
		p = registry[TypeInternal]
		registryMu.RUnlock()
	}

	if p == nil {
		return &ProblemDetails{Type: "about:blank", Status: http.StatusInternalServerError, Title: http.StatusText(http.StatusInternalServerError)}
	}

	return p
}

func isReservedProblemKey(k string) bool {
	return k == "type" || k == "title" || k == "status" || k == "detail" || k == "instance"
}
