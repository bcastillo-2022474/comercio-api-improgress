package readmodel

import (
	"context"
	"fmt"

	"joao-api-comercio/internal/modules/dimensions/readmodel/db"
)

// SQLite implements all three query interfaces over the lookup tables.
type SQLite struct {
	q *db.Queries
}

func NewSQLite(q db.DBTX) *SQLite {
	return &SQLite{q: db.New(q)}
}

func (s *SQLite) ListCountries(ctx context.Context) ([]Option, error) {
	rows, err := s.q.ListPaises(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing countries: %w", err)
	}
	out := make([]Option, 0, len(rows))
	for _, r := range rows {
		out = append(out, Option{Id: r.ID, Name: r.Name})
	}
	return out, nil
}

func (s *SQLite) ListCustoms(ctx context.Context) ([]Option, error) {
	rows, err := s.q.ListAduanas(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing customs: %w", err)
	}
	out := make([]Option, 0, len(rows))
	for _, r := range rows {
		out = append(out, Option{Id: r.ID, Name: r.Name})
	}
	return out, nil
}

func (s *SQLite) ListRegimeTypes(ctx context.Context) ([]Option, error) {
	rows, err := s.q.ListTiposRegimen(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing regime types: %w", err)
	}
	out := make([]Option, 0, len(rows))
	for _, r := range rows {
		out = append(out, Option{Id: r.ID, Name: r.Name})
	}
	return out, nil
}
