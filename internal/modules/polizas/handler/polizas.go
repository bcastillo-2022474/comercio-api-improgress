// Package handler exposes the polizas module over HTTP with Huma.
package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/readmodel"
	"joao-api-comercio/internal/modules/polizas/use_cases/command"
	polizaquery "joao-api-comercio/internal/modules/polizas/use_cases/query"
	"joao-api-comercio/internal/shared/auth"
	"joao-api-comercio/internal/shared/httperr"
	"joao-api-comercio/internal/shared/httpx"
)

type Polizas struct {
	create *command.CreatePolizaUseCase
	list   *polizaquery.ListPolizasUseCase
	get    *polizaquery.GetPolizaUseCase
	log    *slog.Logger
}

func NewPolizas(create *command.CreatePolizaUseCase, list *polizaquery.ListPolizasUseCase, get *polizaquery.GetPolizaUseCase, log *slog.Logger) *Polizas {
	return &Polizas{create: create, list: list, get: get, log: log}
}

// Mount adds the declaration endpoints to api. Every one needs a signed-in caller, so protect is required.
func (h *Polizas) Mount(api huma.API, protect huma.Middlewares) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-poliza",
		Method:        http.MethodPost,
		Path:          "/polizas",
		Summary:       "Create a declaration with its items",
		Tags:          []string{"polizas"},
		DefaultStatus: http.StatusCreated,
		Middlewares:   protect,
		Security:      auth.BearerSecurity,
	}, h.createPoliza)

	huma.Register(api, huma.Operation{
		OperationID: "list-polizas",
		Method:      http.MethodGet,
		Path:        "/polizas",
		Summary:     "List declarations, filtered and paginated",
		Tags:        []string{"polizas"},
		Middlewares: protect,
		Security:    auth.BearerSecurity,
	}, h.listPolizas)

	huma.Register(api, huma.Operation{
		OperationID: "get-poliza",
		Method:      http.MethodGet,
		Path:        "/polizas/{id}",
		Summary:     "Get a declaration with its items",
		Tags:        []string{"polizas"},
		Middlewares: protect,
		Security:    auth.BearerSecurity,
	}, h.getPoliza)
}

type itemInput struct {
	SacId                   int64   `json:"sac_id"`
	PaisId                  int64   `json:"pais_id"`
	TipoUnidadMedidaId      int64   `json:"tipo_unidad_medida_id"`
	CantidadFraccion        float64 `json:"cantidad_fraccion"`
	TasaDai                 float64 `json:"tasa_dai"`
	ValorDai                float64 `json:"valor_dai"`
	ValorCifUds             float64 `json:"valor_cif_uds"`
	TasaCifCantidadFraccion float64 `json:"tasa_cif_cantidad_fraccion"`
}

type createPolizaInput struct {
	Body struct {
		Correlativo      int64       `json:"correlativo"`
		FechaDeclaracion string      `json:"fecha_declaracion" doc:"YYYY-MM-DD"`
		TipoCambioDolar  float64     `json:"tipo_cambio_dolar"`
		AduanaId         *int64      `json:"aduana_id,omitempty"`
		TipoRegimenId    *int64      `json:"tipo_regimen_id,omitempty"`
		Items            []itemInput `json:"items"`
	}
}

type createPolizaOutput struct {
	Body struct {
		Id          int64 `json:"id"`
		Correlativo int64 `json:"correlativo"`
	}
}

func (h *Polizas) createPoliza(ctx context.Context, in *createPolizaInput) (*createPolizaOutput, error) {
	items := make([]poliza.Item, 0, len(in.Body.Items))
	for _, it := range in.Body.Items {
		items = append(items, poliza.Item{
			SacId:                   it.SacId,
			PaisId:                  it.PaisId,
			TipoUnidadMedidaId:      it.TipoUnidadMedidaId,
			CantidadFraccion:        it.CantidadFraccion,
			TasaDai:                 it.TasaDai,
			ValorDai:                it.ValorDai,
			ValorCifUds:             it.ValorCifUds,
			TasaCifCantidadFraccion: it.TasaCifCantidadFraccion,
		})
	}

	p, err := h.create.Execute(ctx, command.CreatePolizaInput{
		Correlativo:      in.Body.Correlativo,
		FechaDeclaracion: in.Body.FechaDeclaracion,
		TipoCambioDolar:  in.Body.TipoCambioDolar,
		AduanaId:         in.Body.AduanaId,
		TipoRegimenId:    in.Body.TipoRegimenId,
		Items:            items,
	})
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "create_poliza"), err)
	}

	out := &createPolizaOutput{}
	out.Body.Id = p.Id
	out.Body.Correlativo = p.Correlativo
	return out, nil
}

type listPolizasInput struct {
	Page          int64                       `query:"page" default:"1" minimum:"1"`
	PageSize      int64                       `query:"page_size" default:"20" minimum:"1" maximum:"100"`
	AduanaId      httpx.OptionalParam[int64]  `query:"aduana_id"`
	TipoRegimenId httpx.OptionalParam[int64]  `query:"tipo_regimen_id"`
	FechaDesde    httpx.OptionalParam[string] `query:"fecha_desde" doc:"YYYY-MM-DD"`
	FechaHasta    httpx.OptionalParam[string] `query:"fecha_hasta" doc:"YYYY-MM-DD"`
}

type listPolizasOutput struct {
	Body readmodel.PolizaPage
}

func (h *Polizas) listPolizas(ctx context.Context, in *listPolizasInput) (*listPolizasOutput, error) {
	page, err := h.list.Execute(ctx, readmodel.PolizaFilter{
		AduanaId:      in.AduanaId.Ptr(),
		TipoRegimenId: in.TipoRegimenId.Ptr(),
		FechaDesde:    in.FechaDesde.Ptr(),
		FechaHasta:    in.FechaHasta.Ptr(),
		Page:          in.Page,
		PageSize:      in.PageSize,
	})
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "list_polizas"), err)
	}
	return &listPolizasOutput{Body: page}, nil
}

type getPolizaInput struct {
	Id int64 `path:"id"`
}

type getPolizaOutput struct {
	Body readmodel.PolizaDetail
}

func (h *Polizas) getPoliza(ctx context.Context, in *getPolizaInput) (*getPolizaOutput, error) {
	detail, err := h.get.Execute(ctx, in.Id)
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "get_poliza"), err)
	}
	return &getPolizaOutput{Body: detail}, nil
}
