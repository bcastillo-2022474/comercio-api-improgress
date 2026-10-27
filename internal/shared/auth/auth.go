// Package auth protects HTTP operations with the bearer access token.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/apperr"
	"joao-api-comercio/internal/shared/httperr"
	"joao-api-comercio/internal/shared/types"
)

// Codes for the token middleware. They travel in the body like every other domain code.
const (
	CodeMissingToken apperr.Code = "missing_token"
	CodeInvalidToken apperr.Code = "invalid_token"
)

// Verifier checks an access token and returns the user it was issued for.
type Verifier interface {
	Verify(token string) (types.UserId, user.Email, error)
}

// Principal is the caller of a protected operation.
type Principal struct {
	UserId types.UserId
	Email  user.Email
}

type principalKey struct{}

// PrincipalFrom returns the caller that Require put on the context.
func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

// BearerSecurity marks an operation as needing the bearer token in the OpenAPI spec.
var BearerSecurity = []map[string][]string{{"bearer": {}}}

// Require rejects a request without a valid bearer token with a 401 and a problem body that carries
// the code. Otherwise it puts the caller on the context. Rejections are client errors, so they are
// not logged.
func Require(verifier Verifier) huma.Middlewares {
	return huma.Middlewares{func(ctx huma.Context, next func(huma.Context)) {
		token, found := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !found || token == "" {
			reject(ctx, CodeMissingToken, "missing bearer token")
			return
		}

		id, email, err := verifier.Verify(token)
		if err != nil {
			reject(ctx, CodeInvalidToken, "invalid token")
			return
		}

		p := Principal{UserId: id, Email: email}
		next(huma.WithValue(ctx, principalKey{}, p))
	}}
}

// reject writes a 401 problem body. huma.WriteErr cannot set the problem type, so the body is written here.
func reject(ctx huma.Context, code apperr.Code, detail string) {
	ctx.SetHeader("Content-Type", "application/problem+json")
	ctx.SetStatus(http.StatusUnauthorized)
	_ = json.NewEncoder(ctx.BodyWriter()).Encode(&huma.ErrorModel{
		Type:   httperr.ErrorTypeURN(code),
		Title:  http.StatusText(http.StatusUnauthorized),
		Status: http.StatusUnauthorized,
		Detail: detail,
	})
}
