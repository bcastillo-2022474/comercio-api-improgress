// Package readmodel holds the read side of the dashboard: the KPI shapes and the queries that produce them.
package readmodel

import "context"

// KpiTotals is the headline numbers. Values are in USD, from valor_cif_uds.
type KpiTotals struct {
	TotalValorCifUsd float64 `json:"total_valor_cif_usd"`
	PolizasCount     int64   `json:"polizas_count"`
	ItemsCount       int64   `json:"items_count"`
}

type CountryValue struct {
	PaisId      int64   `json:"pais_id"`
	PaisName    string  `json:"pais_name"`
	ValorCifUsd float64 `json:"valor_cif_usd"`
}

type ProductValue struct {
	SacId       int64   `json:"sac_id"`
	Fraccion    int64   `json:"fraccion"`
	Descripcion *string `json:"descripcion"`
	ValorCifUsd float64 `json:"valor_cif_usd"`
}

// MonthValue is the import value of one month. Mes is 1-12.
type MonthValue struct {
	Mes         int64   `json:"mes"`
	ValorCifUsd float64 `json:"valor_cif_usd"`
}

// TotalsQuery reads the headline numbers.
type TotalsQuery interface {
	KpiTotals(ctx context.Context) (KpiTotals, error)
}

// TopCountriesQuery reads the five countries with the highest import value.
type TopCountriesQuery interface {
	TopCountries(ctx context.Context) ([]CountryValue, error)
}

// TopProductsQuery reads the five SACs with the highest import value.
type TopProductsQuery interface {
	TopProducts(ctx context.Context) ([]ProductValue, error)
}

// ValorPorMesQuery reads the monthly import value of one year, and the latest year in the data.
type ValorPorMesQuery interface {
	ValorPorMes(ctx context.Context, year int64) ([]MonthValue, error)
	LatestYear(ctx context.Context) (int64, error)
}
