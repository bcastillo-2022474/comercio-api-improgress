package kpi

import (
	"fmt"

	"joao-api-comercio/internal/shared/apperr"
)

type InvalidYearError struct {
	Year int64
}

func (e *InvalidYearError) Error() string {
	return fmt.Sprintf("year must be between 1900 and 2100, got %d", e.Year)
}

func (e *InvalidYearError) Code() apperr.Code { return CodeInvalidYear }
