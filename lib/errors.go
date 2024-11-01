package lib

import (
	"fmt"
	"net/http"
)

type ExternalError struct {
	Code    int16
	Message string
	Type    string
}

func (e ExternalError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewExternalError() *ExternalError {
	return &ExternalError{}
}

func (e *ExternalError) BadRequest(message string) ExternalError {
	e.Code = http.StatusBadRequest
	e.Type = "bad_request"

	if message == "" {
		e.Message = "invalid request"
	} else {
		e.Message = message
	}
	return *e
}

func (e *ExternalError) NotFound(message string) ExternalError {
	e.Code = http.StatusNotFound
	e.Type = "not_found"

	if message == "" {
		e.Message = "not found"
	} else {
		e.Message = message
	}
	return *e
}

func (e *ExternalError) Unavailable(message string) ExternalError {
	e.Code = http.StatusServiceUnavailable
	e.Type = "service_unavailable"

	if message == "" {
		e.Message = "not yet available"
	} else {
		e.Message = message
	}
	return *e
}

func (e *ExternalError) InternalServerError() ExternalError {
	e.Code = http.StatusInternalServerError
	e.Message = "internal server error"
	e.Type = "internal_server_error"
	return *e
}
