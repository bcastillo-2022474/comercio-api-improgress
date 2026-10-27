package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dimensions/readmodel"
)

type ListRegimeTypesUseCase struct {
	regimeTypes readmodel.RegimeTypesQuery
	log         *slog.Logger
}

func NewListRegimeTypesUseCase(regimeTypes readmodel.RegimeTypesQuery, log *slog.Logger) *ListRegimeTypesUseCase {
	return &ListRegimeTypesUseCase{regimeTypes: regimeTypes, log: log.With("use_case", "list_regime_types")}
}

func (uc *ListRegimeTypesUseCase) Execute(ctx context.Context) ([]readmodel.Option, error) {
	uc.log.InfoContext(ctx, "starting regime types query")
	options, err := uc.regimeTypes.ListRegimeTypes(ctx)
	if err != nil {
		return nil, err
	}
	uc.log.InfoContext(ctx, "regime types query finished", "count", len(options))
	return options, nil
}
