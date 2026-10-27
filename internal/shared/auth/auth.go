// Package auth protects HTTP operations with the bearer access token.
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/shared/types"
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

// Require rejects a request without a valid bearer token with a 401. Otherwise it puts the caller
// on the context. Rejections are client errors, so they are not logged.
func Require(api huma.API, verifier Verifier) huma.Middlewares {
	return huma.Middlewares{func(ctx huma.Context, next func(huma.Context)) {
		token, found := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !found || token == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing bearer token")
			return
		}

		id, email, err := verifier.Verify(token)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token")
			return
		}

		p := Principal{UserId: id, Email: email}
		next(huma.WithValue(ctx, principalKey{}, p))
	}}
}
