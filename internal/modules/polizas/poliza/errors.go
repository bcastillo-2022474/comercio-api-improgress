package poliza

import (
	"fmt"

	"joao-api-comercio/internal/shared/apperr"
)

type InvalidCorrelativoError struct {
	Correlativo int64
}

func (e *InvalidCorrelativoError) Error() string {
	return fmt.Sprintf("correlativo must be positive, got %d", e.Correlativo)
}

func (e *InvalidCorrelativoError) Code() apperr.Code { return CodeInvalidCorrelativo }

type InvalidFechaError struct {
	Fecha string
}

func (e *InvalidFechaError) Error() string {
	return fmt.Sprintf("fecha_declaracion must be YYYY-MM-DD, got %q", e.Fecha)
}

func (e *InvalidFechaError) Code() apperr.Code { return CodeInvalidFecha }

type InvalidTipoCambioError struct {
	TipoCambio float64
}

func (e *InvalidTipoCambioError) Error() string {
	return fmt.Sprintf("tipo_cambio_dolar must be positive, got %v", e.TipoCambio)
}

func (e *InvalidTipoCambioError) Code() apperr.Code { return CodeInvalidTipoCambio }

type NoItemsError struct{}

func (e *NoItemsError) Error() string { return "a declaration needs at least one item" }

func (e *NoItemsError) Code() apperr.Code { return CodeNoItems }

type InvalidItemError struct {
	Index int
}

func (e *InvalidItemError) Error() string {
	return fmt.Sprintf("item %d has a missing reference or a negative amount", e.Index)
}

func (e *InvalidItemError) Code() apperr.Code { return CodeInvalidItem }

type CorrelativoTakenError struct {
	Correlativo int64
}

func (e *CorrelativoTakenError) Error() string {
	return fmt.Sprintf("correlativo %d is already registered", e.Correlativo)
}

func (e *CorrelativoTakenError) Code() apperr.Code { return CodeCorrelativoTaken }

// InvalidReferenceError means an aduana, regime, SAC, country or unit does not exist.
type InvalidReferenceError struct{}

func (e *InvalidReferenceError) Error() string {
	return "a referenced aduana, regime, SAC, country or unit does not exist"
}

func (e *InvalidReferenceError) Code() apperr.Code { return CodeInvalidReference }

type NotFoundError struct {
	Id int64
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("no poliza with id %d", e.Id)
}

func (e *NotFoundError) Code() apperr.Code { return CodeNotFound }
