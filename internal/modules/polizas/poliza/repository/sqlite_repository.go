// Package repository holds the SQL adapter for the poliza port.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/poliza/repository/db"
	"joao-api-comercio/internal/shared/types"
)

// Constraint messages SQLite gives when a row is rejected.
const (
	correlativoConstraint = "UNIQUE constraint failed: polizas.correlativo"
	referenceConstraint   = "FOREIGN KEY constraint failed"
)

type SQLiteRepository struct {
	q *db.Queries
}

func NewSQLiteRepository(q db.DBTX) *SQLiteRepository {
	return &SQLiteRepository{q: db.New(q)}
}

func (r *SQLiteRepository) Insert(ctx context.Context, p *poliza.Poliza) error {
	id, err := r.q.InsertPoliza(ctx, db.InsertPolizaParams{
		Correlativo:      p.Correlativo,
		FechaDeclaracion: p.FechaDeclaracion,
		TipoCambioDolar:  p.TipoCambioDolar,
		AduanaID:         p.AduanaId,
		TipoRegimenID:    p.TipoRegimenId,
	})
	if err != nil {
		return mapError(err, p.Correlativo)
	}

	for _, item := range p.Items {
		err := r.q.InsertPolizaItem(ctx, db.InsertPolizaItemParams{
			PolizaID:                id,
			SacID:                   item.SacId,
			PaisID:                  item.PaisId,
			TipoUnidadMedidaID:      item.TipoUnidadMedidaId,
			CantidadFraccion:        item.CantidadFraccion,
			TasaDai:                 item.TasaDai,
			ValorDai:                item.ValorDai,
			ValorCifUds:             item.ValorCifUds,
			TasaCifCantidadFraccion: item.TasaCifCantidadFraccion,
		})
		if err != nil {
			return mapError(err, p.Correlativo)
		}
	}

	p.Id = id
	return nil
}

func (r *SQLiteRepository) GetById(ctx context.Context, id int64) (types.Maybe[poliza.Poliza], error) {
	header, err := r.q.GetPolizaById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return types.None[poliza.Poliza](), nil
	}
	if err != nil {
		return types.None[poliza.Poliza](), fmt.Errorf("getting poliza: %w", err)
	}

	rows, err := r.q.ListItemsByPolizaId(ctx, id)
	if err != nil {
		return types.None[poliza.Poliza](), fmt.Errorf("listing poliza items: %w", err)
	}
	items := make([]poliza.Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, poliza.Item{
			SacId:                   row.SacID,
			PaisId:                  row.PaisID,
			TipoUnidadMedidaId:      row.TipoUnidadMedidaID,
			CantidadFraccion:        row.CantidadFraccion,
			TasaDai:                 row.TasaDai,
			ValorDai:                row.ValorDai,
			ValorCifUds:             row.ValorCifUds,
			TasaCifCantidadFraccion: row.TasaCifCantidadFraccion,
		})
	}

	return types.Some(poliza.Poliza{
		Id:               header.ID,
		Correlativo:      header.Correlativo,
		FechaDeclaracion: header.FechaDeclaracion,
		TipoCambioDolar:  header.TipoCambioDolar,
		AduanaId:         header.AduanaID,
		TipoRegimenId:    header.TipoRegimenID,
		Items:            items,
	}), nil
}

func mapError(err error, correlativo int64) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, correlativoConstraint):
		return &poliza.CorrelativoTakenError{Correlativo: correlativo}
	case strings.Contains(msg, referenceConstraint):
		return &poliza.InvalidReferenceError{}
	default:
		return fmt.Errorf("inserting poliza: %w", err)
	}
}
