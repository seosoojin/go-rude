package rude

import (
	"encoding/json"
	"maps"
	"net/http"
)

type ProblemDetails struct {
	Type       string         `json:"type,omitempty"`
	Status     int            `json:"status,omitempty"`
	Title      string         `json:"title,omitempty"`
	Detail     string         `json:"detail,omitempty"`
	Instance   string         `json:"instance,omitempty"`
	Extensions map[string]any `json:"-"`
}

func (p ProblemDetails) MarshalJSON() ([]byte, error) {
	problemType := p.Type
	if problemType == "" {
		problemType = "about:blank"
	}

	base := map[string]any{"type": problemType}
	if p.Title != "" {
		base["title"] = p.Title
	}
	if p.Status != 0 {
		base["status"] = p.Status
	}
	if p.Detail != "" {
		base["detail"] = p.Detail
	}
	if p.Instance != "" {
		base["instance"] = p.Instance
	}

	extensions := p.Extensions
	for k, v := range extensions {
		base[k] = v
	}

	return json.Marshal(base)
}

func (p ProblemDetails) WithExtension(k string, v any) ProblemDetails {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}
	p.Extensions[k] = v
	return p
}

func (p ProblemDetails) WithExtensions(exts map[string]any) ProblemDetails {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}
	maps.Copy(p.Extensions, exts)
	return p
}

func (p ProblemDetails) Write(w http.ResponseWriter, r *http.Request) {
	if p.Type == "" {
		p.Type = "about:blank"
	}

	if p.Status == 0 {
		p.Status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)

	_ = json.NewEncoder(w).Encode(p)
}

func FromError(err error) ProblemDetails {
	switch e := err.(type) {
	case Error:
		return ProblemDetails{
			Type:       string(e.Type),
			Title:      e.Message,
			Status:     e.Code,
			Detail:     e.Err.Error(),
			Extensions: e.MetaData,
		}
	default:
		return ProblemDetails{
			Type:       "about:blank",
			Title:      err.Error(),
			Status:     http.StatusInternalServerError,
			Detail:     err.Error(),
			Extensions: nil,
		}
	}
}
