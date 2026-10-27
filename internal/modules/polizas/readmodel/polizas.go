// Package readmodel holds the read side of the polizas module: screen-shaped DTOs and the query
// interfaces that produce them. There are no repositories here, only queries.
package readmodel

import (
	"context"

	"joao-api-comercio/internal/shared/types"
)

// PolizaSummary is one row of the declarations list. It has no items.
type PolizaSummary struct {
	Id               int64    `json:"id"`
	Correlativo      int64    `json:"correlativo"`
	FechaDeclaracion string   `json:"fecha_declaracion"`
	TipoCambioDolar  *float64 `json:"tipo_cambio_dolar"`
	AduanaId         *int64   `json:"aduana_id"`
	AduanaName       *string  `json:"aduana_name"`
	TipoRegimenId    *int64   `json:"tipo_regimen_id"`
	TipoRegimenName  *string  `json:"tipo_regimen_name"`
	ItemsCount       int64    `json:"items_count"`
}

// PolizaFilter narrows the list. Nil filters are not applied. Page is 1-based.
type PolizaFilter struct {
	AduanaId      *int64
	TipoRegimenId *int64
	FechaDesde    *string
	FechaHasta    *string
	Page          int64
	PageSize      int64
}

type PolizaPage struct {
	Items    []PolizaSummary `json:"items"`
	Total    int64           `json:"total"`
	Page     int64           `json:"page"`
	PageSize int64           `json:"page_size"`
}

// ItemDetail is one item of a declaration, with its SAC, country and unit names.
type ItemDetail struct {
	Id                      int64    `json:"id"`
	SacId                   *int64   `json:"sac_id"`
	SacFraccion             int64    `json:"sac_fraccion"`
	SacDescripcion          *string  `json:"sac_descripcion"`
	PaisId                  *int64   `json:"pais_id"`
	PaisName                string   `json:"pais_name"`
	TipoUnidadMedidaId      *int64   `json:"tipo_unidad_medida_id"`
	TipoUnidadMedidaName    string   `json:"tipo_unidad_medida_name"`
	CantidadFraccion        *float64 `json:"cantidad_fraccion"`
	TasaDai                 *float64 `json:"tasa_dai"`
	ValorDai                *float64 `json:"valor_dai"`
	ValorCifUds             *float64 `json:"valor_cif_uds"`
	TasaCifCantidadFraccion *float64 `json:"tasa_cif_cantidad_fraccion"`
}

// PolizaDetail is one declaration with its items.
type PolizaDetail struct {
	PolizaSummary
	Items []ItemDetail `json:"items"`
}

// ListPolizasQuery pages through the declarations that match the filter.
type ListPolizasQuery interface {
	ListPolizas(ctx context.Context, filter PolizaFilter) (PolizaPage, error)
}

// GetPolizaQuery reads one declaration with its items. An unknown id is absent, not an error.
type GetPolizaQuery interface {
	GetPoliza(ctx context.Context, id int64) (types.Maybe[PolizaDetail], error)
}
