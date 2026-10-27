package user

import "joao-api-comercio/internal/shared/apperr"

const (
	CodeInvalidEmail       apperr.Code = "invalid_email"
	CodeWeakPassword       apperr.Code = "weak_password"
	CodeEmailTaken         apperr.Code = "email_taken"
	CodeInvalidCredentials apperr.Code = "invalid_credentials"
)
