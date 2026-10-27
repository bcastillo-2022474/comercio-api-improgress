// Package apperr is the shared error-code contract. Every domain error that the transport layer
// can report carries a Code. The transport layer maps codes to statuses, and nothing else.
package apperr

import "errors"

type Code string

// Coded is implemented by every domain error with a code.
type Coded interface {
	error
	Code() Code
}

// CodeOf returns the code of the first Coded error in err's chain, or "" if there is none.
func CodeOf(err error) Code {
	var coded Coded
	if errors.As(err, &coded) {
		return coded.Code()
	}
	return ""
}
