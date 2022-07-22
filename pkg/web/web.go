// Package web defines useful structs to wrap http responses and errors.
package web

import (
	"github.com/aws/aws-lambda-go/events"
)

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request through the
// application with web specific context.
type RequestError struct {
	Err    error
	Status int
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (re *RequestError) Error() string {
	return re.Err.Error()
}

// A Middleware is a type that wraps a Handler function that returns another Handler function.
type Middleware func(Handler) Handler

// A Handler is a type that handles a http request within our project.
type Handler func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
