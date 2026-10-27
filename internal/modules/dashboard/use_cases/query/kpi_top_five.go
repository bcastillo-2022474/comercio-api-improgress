package query

import (
	"context"
	"log/slog"
	"time"

	"joao-api-comercio/internal/modules/dashboard/readmodel"
	"joao-api-comercio/internal/shared/cache"
)

type TopFive struct {
	Paises    []readmodel.CountryValue `json:"paises"`
	Productos []readmodel.ProductValue `json:"productos"`
}

type KpiTopFiveUseCase struct {
	countries readmodel.TopCountriesQuery
	products  readmodel.TopProductsQuery
	log       *slog.Logger
	cache     *cache.TTL[struct{}, TopFive]
}

func NewKpiTopFiveUseCase(countries readmodel.TopCountriesQuery, products readmodel.TopProductsQuery, log *slog.Logger, ttl time.Duration) *KpiTopFiveUseCase {
	return &KpiTopFiveUseCase{countries: countries, products: products, log: log.With("use_case", "kpi_top_five"), cache: cache.NewTTL[struct{}, TopFive](ttl, nil)}
}

func (uc *KpiTopFiveUseCase) Execute(ctx context.Context) (TopFive, error) {
	uc.log.InfoContext(ctx, "starting kpi top five")

	if cached, ok := uc.cache.Get(struct{}{}); ok {
		uc.log.InfoContext(ctx, "kpi served from cache")
		return cached, nil
	}

	countries, err := uc.countries.TopCountries(ctx)
	if err != nil {
		return TopFive{}, err
	}
	products, err := uc.products.TopProducts(ctx)
	if err != nil {
		return TopFive{}, err
	}

	result := TopFive{Paises: countries, Productos: products}
	uc.cache.Set(struct{}{}, result)
	uc.log.InfoContext(ctx, "kpi top five finished", "countries", len(countries), "products", len(products))
	return result, nil
}
