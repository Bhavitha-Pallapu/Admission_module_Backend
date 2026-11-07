package errors

import (
	"encoding/json"
	stderrors "errors"
)

// Error is a standard application error.
type Error struct {
	Kind       Kind   `json:"kind"`
	Message    string `json:"message"`
	WrappedErr error  `json:"wrapped_err,omitempty"`
}

// Error returns the JSON representation of the application error.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	data, err := json.Marshal(e)
	if err != nil {
		return e.Message
	}
	return string(data)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.WrappedErr
}

// NewError returns a standard Go error with the given message.
func NewError(message string) error {
	return stderrors.New(message)
}

// Kind defines the kind or class of an error.
type Kind uint8

// Transport-agnostic error kinds.
const (
	Other        Kind = iota // Unclassified error.
	Internal                 // Internal server error.
	Conflict                 // Conflict with the current state.
	Invalid                  // Invalid input.
	NotFound                 // Requested entity was not found.
	Unauthorized             // Authentication is required or failed.
	Forbidden                // The operation is not permitted.
)

// String returns the human-readable description of the error kind.
func (k Kind) String() string {
	switch k {
	case Other:
		return "unclassified error"
	case Internal:
		return "internal error"
	case Conflict:
		return "conflict"
	case Invalid:
		return "invalid input"
	case NotFound:
		return "entity not found"
	case Unauthorized:
		return "unauthorized"
	case Forbidden:
		return "forbidden"
	default:
		return "unknown error kind"
	}
}

// MarshalJSON encodes the human-readable error kind.
func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

// E constructs an application error from a kind, message, and/or wrapped error.
// Arguments of other types are ignored; when an argument type is repeated, the
// last value takes precedence.
func E(args ...interface{}) error {
	err := &Error{}
	for _, arg := range args {
		switch value := arg.(type) {
		case Kind:
			err.Kind = value
		case error:
			err.WrappedErr = value
		case string:
			err.Message = value
		}
	}
	return err
}

// NewInternalServerError creates a new internal server error.
func NewInternalServerError(message string) error {
	return E(Internal, message)
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) error {
	return E(NotFound, message)
}

// NewInvalidParamsError creates a new invalid-parameters error.
func NewInvalidParamsError(message string) error {
	return E(Invalid, message)
}

// NewUnauthorizedError creates a new unauthorized error.
func NewUnauthorizedError(message string) error {
	return E(Unauthorized, message)
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(message string) error {
	return E(Forbidden, message)
}

// NewConflictError creates a new conflict error.
func NewConflictError(message string) error {
	return E(Conflict, message)
}

var (
	As = stderrors.As
	Is = stderrors.Is
)
