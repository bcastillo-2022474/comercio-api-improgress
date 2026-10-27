package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dimensions/readmodel"
)

type ListCustomsUseCase struct {
	customs readmodel.CustomsQuery
	log     *slog.Logger
}

func NewListCustomsUseCase(customs readmodel.CustomsQuery, log *slog.Logger) *ListCustomsUseCase {
	return &ListCustomsUseCase{customs: customs, log: log.With("use_case", "list_customs")}
}

func (uc *ListCustomsUseCase) Execute(ctx context.Context) ([]readmodel.Option, error) {
	uc.log.InfoContext(ctx, "starting customs query")
	options, err := uc.customs.ListCustoms(ctx)
	if err != nil {
		return nil, err
	}
	uc.log.InfoContext(ctx, "customs query finished", "count", len(options))
	return options, nil
}
