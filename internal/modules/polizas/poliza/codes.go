package poliza

import "joao-api-comercio/internal/shared/apperr"

const (
	CodeInvalidCorrelativo apperr.Code = "invalid_correlativo"
	CodeInvalidFecha       apperr.Code = "invalid_fecha_declaracion"
	CodeInvalidTipoCambio  apperr.Code = "invalid_tipo_cambio"
	CodeNoItems            apperr.Code = "poliza_without_items"
	CodeInvalidItem        apperr.Code = "invalid_poliza_item"
	CodeCorrelativoTaken   apperr.Code = "correlativo_taken"
	CodeInvalidReference   apperr.Code = "invalid_reference"
	CodeNotFound           apperr.Code = "poliza_not_found"
)
