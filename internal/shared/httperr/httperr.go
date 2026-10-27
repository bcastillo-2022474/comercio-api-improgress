// Package httperr is the transport edge for errors: it maps domain error codes to HTTP statuses,
// and logs the errors it cannot map, once, at the top of the request.
package httperr

import (
	"context"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/dashboard/kpi"
	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/shared/apperr"
)

// codeToStatus has a row per code the API reports. A code with no row is internal.
var codeToStatus = map[apperr.Code]func(msg string, errs ...error) huma.StatusError{
	user.CodeInvalidEmail:       huma.Error422UnprocessableEntity,
	user.CodeWeakPassword:       huma.Error422UnprocessableEntity,
	user.CodeEmailTaken:         huma.Error409Conflict,
	user.CodeInvalidCredentials: huma.Error401Unauthorized,

	poliza.CodeInvalidCorrelativo: huma.Error422UnprocessableEntity,
	poliza.CodeInvalidFecha:       huma.Error422UnprocessableEntity,
	poliza.CodeInvalidTipoCambio:  huma.Error422UnprocessableEntity,
	poliza.CodeNoItems:            huma.Error422UnprocessableEntity,
	poliza.CodeInvalidItem:        huma.Error422UnprocessableEntity,
	poliza.CodeInvalidReference:   huma.Error422UnprocessableEntity,
	poliza.CodeCorrelativoTaken:   huma.Error409Conflict,
	poliza.CodeNotFound:           huma.Error404NotFound,

	kpi.CodeInvalidYear: huma.Error422UnprocessableEntity,
}

// ErrorTypeURN is the RFC 9457 "type" of a domain error: a stable identifier in the body that a
// client branches on, instead of matching the message text. It is a URN, so it needs no web host.
func ErrorTypeURN(code apperr.Code) string {
	return "urn:joao-api-comercio:error:" + string(code)
}

// ToHumaError turns err into the HTTP error the client sees. A mapped error carries its domain code
// as the problem "type" in the body. An unmapped error is logged in full with the caller's log, and
// the client gets a plain 500 with no internal detail.
func ToHumaError(ctx context.Context, log *slog.Logger, err error) error {
	code := apperr.CodeOf(err)
	status, known := codeToStatus[code]
	if !known {
		log.ErrorContext(ctx, "request failed", "err", err)
		return huma.Error500InternalServerError("internal error")
	}

	mapped := status(err.Error())
	if model, ok := mapped.(*huma.ErrorModel); ok {
		model.Type = ErrorTypeURN(code)
	}
	return mapped
}
