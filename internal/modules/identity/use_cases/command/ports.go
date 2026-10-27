package command

import (
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/types"
)

// TokenIssuer signs the access token for a logged-in user.
type TokenIssuer interface {
	Issue(userId types.UserId, email user.Email) (token string, expiresAt time.Time, err error)
}
