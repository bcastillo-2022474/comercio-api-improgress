package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/readmodel"
)

type GetPolizaUseCase struct {
	polizas readmodel.GetPolizaQuery
	log     *slog.Logger
}

func NewGetPolizaUseCase(polizas readmodel.GetPolizaQuery, log *slog.Logger) *GetPolizaUseCase {
	return &GetPolizaUseCase{polizas: polizas, log: log.With("use_case", "get_poliza")}
}

func (uc *GetPolizaUseCase) Execute(ctx context.Context, id int64) (readmodel.PolizaDetail, error) {
	uc.log.InfoContext(ctx, "starting get poliza", "poliza_id", id)

	found, err := uc.polizas.GetPoliza(ctx, id)
	if err != nil {
		return readmodel.PolizaDetail{}, err
	}
	if !found.IsPresent {
		return readmodel.PolizaDetail{}, &poliza.NotFoundError{Id: id}
	}

	uc.log.InfoContext(ctx, "get poliza finished", "poliza_id", id, "items", len(found.Value.Items))
	return found.Value, nil
}
