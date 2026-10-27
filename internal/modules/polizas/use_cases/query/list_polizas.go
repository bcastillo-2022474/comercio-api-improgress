package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/readmodel"
)

type ListPolizasUseCase struct {
	polizas readmodel.ListPolizasQuery
	log     *slog.Logger
}

func NewListPolizasUseCase(polizas readmodel.ListPolizasQuery, log *slog.Logger) *ListPolizasUseCase {
	return &ListPolizasUseCase{polizas: polizas, log: log.With("use_case", "list_polizas")}
}

func (uc *ListPolizasUseCase) Execute(ctx context.Context, filter readmodel.PolizaFilter) (readmodel.PolizaPage, error) {
	uc.log.InfoContext(ctx, "starting list polizas", "page", filter.Page, "page_size", filter.PageSize)

	if filter.FechaDesde != nil {
		if err := poliza.ValidFecha(*filter.FechaDesde); err != nil {
			return readmodel.PolizaPage{}, err
		}
	}
	if filter.FechaHasta != nil {
		if err := poliza.ValidFecha(*filter.FechaHasta); err != nil {
			return readmodel.PolizaPage{}, err
		}
	}

	page, err := uc.polizas.ListPolizas(ctx, filter)
	if err != nil {
		return readmodel.PolizaPage{}, err
	}

	uc.log.InfoContext(ctx, "list polizas finished", "total", page.Total)
	return page, nil
}
