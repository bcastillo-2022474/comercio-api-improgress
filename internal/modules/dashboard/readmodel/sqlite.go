package readmodel

import (
	"context"
	"fmt"

	"joao-api-comercio/internal/modules/dashboard/readmodel/db"
)

// SQLite implements every dashboard query over the import tables.
type SQLite struct {
	q *db.Queries
}

func NewSQLite(q db.DBTX) *SQLite {
	return &SQLite{q: db.New(q)}
}

func (s *SQLite) KpiTotals(ctx context.Context) (KpiTotals, error) {
	row, err := s.q.KpiTotals(ctx)
	if err != nil {
		return KpiTotals{}, fmt.Errorf("reading kpi totals: %w", err)
	}
	return KpiTotals{TotalValorCifUsd: row.TotalValorCifUds, PolizasCount: row.PolizasCount, ItemsCount: row.ItemsCount}, nil
}

func (s *SQLite) TopCountries(ctx context.Context) ([]CountryValue, error) {
	rows, err := s.q.TopCountries(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading top countries: %w", err)
	}
	out := make([]CountryValue, 0, len(rows))
	for _, r := range rows {
		out = append(out, CountryValue{PaisId: *r.PaisID, PaisName: r.PaisName, ValorCifUsd: r.ValorCifUds})
	}
	return out, nil
}

func (s *SQLite) TopProducts(ctx context.Context) ([]ProductValue, error) {
	rows, err := s.q.TopProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading top products: %w", err)
	}
	out := make([]ProductValue, 0, len(rows))
	for _, r := range rows {
		out = append(out, ProductValue{SacId: *r.SacID, Fraccion: r.Fraccion, Descripcion: r.Descripcion, ValorCifUsd: r.ValorCifUds})
	}
	return out, nil
}

func (s *SQLite) ValorPorMes(ctx context.Context, year int64) ([]MonthValue, error) {
	rows, err := s.q.ValorPorMes(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("reading monthly value: %w", err)
	}
	out := make([]MonthValue, 0, len(rows))
	for _, r := range rows {
		out = append(out, MonthValue{Mes: r.Mes, ValorCifUsd: r.ValorCifUds})
	}
	return out, nil
}

func (s *SQLite) LatestYear(ctx context.Context) (int64, error) {
	year, err := s.q.LatestYear(ctx)
	if err != nil {
		return 0, fmt.Errorf("reading latest year: %w", err)
	}
	return year, nil
}
