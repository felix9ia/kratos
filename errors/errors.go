package errors

import (
	"fmt"
)

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Details provide more context to an error.
	Details []interface{} `json:"details"`
}

// ErrorInfo is a detailed error code & message from the API frontend.
type ErrorInfo struct {
	// Reason is the typed error code. For example: "some_example".
	Reason string `json:"reason"`
	// Message is the human-readable description of the error.
	Message string `json:"message"`
}

// WithDetails provided details messages appended to the errors.
func (e *Error) WithDetails(details ...interface{}) {
	e.Details = append(e.Details, details...)
}

// Is each error in a chain for a match with a target value.
func (e *Error) Is(target error) bool {
	te, ok := target.(*Error)
	if !ok {
		return false
	}
	if e.Code != te.Code || len(e.Details) != len(te.Details) {
		return false
	}
	if e.Code == te.Code && len(e.Details) == 0 && len(te.Details) == 0 {
		return true
	}
	var (
		n    int
		errs []*ErrorInfo
	)
	for _, d := range e.Details {
		if e, ok := d.(*ErrorInfo); ok {
			errs = append(errs, e)
		}
	}
	for _, d := range te.Details {
		if e, ok := d.(*ErrorInfo); ok {
			for _, err := range errs {
				if e.Reason == err.Reason {
					n++
				}
			}
		}
	}
	return len(errs) == n
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", e.Code, e.Message, e.Details)
}

// New returns a Status representing c and msg.
func New(code int, message string, details ...interface{}) error {
	return &Error{Code: code, Message: message, Details: details}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(code int, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}
