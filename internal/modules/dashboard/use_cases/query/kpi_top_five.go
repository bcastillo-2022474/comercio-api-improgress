package query

import (
	"context"
	"log/slog"

	"joao-api-comercio/internal/modules/dashboard/readmodel"
)

type TopFive struct {
	Paises    []readmodel.CountryValue `json:"paises"`
	Productos []readmodel.ProductValue `json:"productos"`
}

type KpiTopFiveUseCase struct {
	countries readmodel.TopCountriesQuery
	products  readmodel.TopProductsQuery
	log       *slog.Logger
}

func NewKpiTopFiveUseCase(countries readmodel.TopCountriesQuery, products readmodel.TopProductsQuery, log *slog.Logger) *KpiTopFiveUseCase {
	return &KpiTopFiveUseCase{countries: countries, products: products, log: log.With("use_case", "kpi_top_five")}
}

func (uc *KpiTopFiveUseCase) Execute(ctx context.Context) (TopFive, error) {
	uc.log.InfoContext(ctx, "starting kpi top five")

	countries, err := uc.countries.TopCountries(ctx)
	if err != nil {
		return TopFive{}, err
	}
	products, err := uc.products.TopProducts(ctx)
	if err != nil {
		return TopFive{}, err
	}

	uc.log.InfoContext(ctx, "kpi top five finished", "countries", len(countries), "products", len(products))
	return TopFive{Paises: countries, Productos: products}, nil
}
