package readmodel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"joao-api-comercio/internal/modules/polizas/readmodel/db"
	"joao-api-comercio/internal/shared/types"
)

// SQLite implements both query interfaces over the declaration tables.
type SQLite struct {
	q *db.Queries
}

func NewSQLite(q db.DBTX) *SQLite {
	return &SQLite{q: db.New(q)}
}

func (s *SQLite) ListPolizas(ctx context.Context, f PolizaFilter) (PolizaPage, error) {
	offset := (f.Page - 1) * f.PageSize
	rows, err := s.q.ListPolizas(ctx, db.ListPolizasParams{
		AduanaID:      f.AduanaId,
		TipoRegimenID: f.TipoRegimenId,
		FechaDesde:    f.FechaDesde,
		FechaHasta:    f.FechaHasta,
		PageOffset:    offset,
		PageSize:      f.PageSize,
	})
	if err != nil {
		return PolizaPage{}, fmt.Errorf("listing polizas: %w", err)
	}

	total, err := s.q.CountPolizas(ctx, db.CountPolizasParams{
		AduanaID:      f.AduanaId,
		TipoRegimenID: f.TipoRegimenId,
		FechaDesde:    f.FechaDesde,
		FechaHasta:    f.FechaHasta,
	})
	if err != nil {
		return PolizaPage{}, fmt.Errorf("counting polizas: %w", err)
	}

	items := make([]PolizaSummary, 0, len(rows))
	for _, r := range rows {
		items = append(items, PolizaSummary{
			Id:               r.ID,
			Correlativo:      r.Correlativo,
			FechaDeclaracion: r.FechaDeclaracion,
			TipoCambioDolar:  r.TipoCambioDolar,
			AduanaId:         r.AduanaID,
			AduanaName:       r.AduanaName,
			TipoRegimenId:    r.TipoRegimenID,
			TipoRegimenName:  r.TipoRegimenName,
			ItemsCount:       r.ItemsCount,
		})
	}
	return PolizaPage{Items: items, Total: total, Page: f.Page, PageSize: f.PageSize}, nil
}

func (s *SQLite) GetPoliza(ctx context.Context, id int64) (types.Maybe[PolizaDetail], error) {
	header, err := s.q.GetPoliza(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return types.None[PolizaDetail](), nil
	}
	if err != nil {
		return types.None[PolizaDetail](), fmt.Errorf("getting poliza: %w", err)
	}

	rows, err := s.q.ListPolizaItems(ctx, id)
	if err != nil {
		return types.None[PolizaDetail](), fmt.Errorf("listing poliza items: %w", err)
	}
	items := make([]ItemDetail, 0, len(rows))
	for _, r := range rows {
		items = append(items, ItemDetail{
			Id:                      r.ID,
			SacId:                   r.SacID,
			SacFraccion:             r.SacFraccion,
			SacDescripcion:          r.SacDescripcion,
			PaisId:                  r.PaisID,
			PaisName:                r.PaisName,
			TipoUnidadMedidaId:      r.TipoUnidadMedidaID,
			TipoUnidadMedidaName:    r.TipoUnidadMedidaName,
			CantidadFraccion:        r.CantidadFraccion,
			TasaDai:                 r.TasaDai,
			ValorDai:                r.ValorDai,
			ValorCifUds:             r.ValorCifUds,
			TasaCifCantidadFraccion: r.TasaCifCantidadFraccion,
		})
	}

	return types.Some(PolizaDetail{
		PolizaSummary: PolizaSummary{
			Id:               header.ID,
			Correlativo:      header.Correlativo,
			FechaDeclaracion: header.FechaDeclaracion,
			TipoCambioDolar:  header.TipoCambioDolar,
			AduanaId:         header.AduanaID,
			AduanaName:       header.AduanaName,
			TipoRegimenId:    header.TipoRegimenID,
			TipoRegimenName:  header.TipoRegimenName,
			ItemsCount:       header.ItemsCount,
		},
		Items: items,
	}), nil
}
