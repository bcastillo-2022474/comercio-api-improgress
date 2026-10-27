package user

import "fmt"

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return fmt.Sprintf("invalid email %q", e.Email)
}

type WeakPasswordError struct {
	MinLength int
	MaxBytes  int
}

func (e *WeakPasswordError) Error() string {
	return fmt.Sprintf("password must be between %d and %d bytes", e.MinLength, e.MaxBytes)
}

type EmailTakenError struct {
	Email Email
}

func (e *EmailTakenError) Error() string {
	return fmt.Sprintf("email %q is already registered", e.Email)
}

// InvalidCredentialsError is returned for an unknown email and for a wrong
// password alike, so a caller cannot tell which accounts exist.
type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials"
}
