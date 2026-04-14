package rude

type ErrorType string

const (
	TypeNotFound     ErrorType = "NOT_FOUND"
	TypeUnauthorized ErrorType = "UNAUTHORIZED"
	TypeValidation   ErrorType = "VALIDATION"
	TypeInternal     ErrorType = "INTERNAL"
)
