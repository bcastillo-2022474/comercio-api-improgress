package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dimensions/readmodel"
)

type ListCountriesUseCase struct {
	countries readmodel.CountriesQuery
	log       *slog.Logger
}

func NewListCountriesUseCase(countries readmodel.CountriesQuery, log *slog.Logger) *ListCountriesUseCase {
	return &ListCountriesUseCase{countries: countries, log: log.With("use_case", "list_countries")}
}

func (uc *ListCountriesUseCase) Execute(ctx context.Context) ([]readmodel.Option, error) {
	uc.log.InfoContext(ctx, "starting countries query")
	options, err := uc.countries.ListCountries(ctx)
	if err != nil {
		return nil, err
	}
	uc.log.InfoContext(ctx, "countries query finished", "count", len(options))
	return options, nil
}
