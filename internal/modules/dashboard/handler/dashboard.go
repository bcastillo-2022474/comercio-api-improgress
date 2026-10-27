// Package handler exposes the dashboard KPIs over HTTP with Huma.
package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/dashboard/readmodel"
	kpiquery "joao-api-comercio/internal/modules/dashboard/use_cases/query"
	"joao-api-comercio/internal/shared/auth"
	"joao-api-comercio/internal/shared/httperr"
	"joao-api-comercio/internal/shared/httpx"
)

type Dashboard struct {
	totals      *kpiquery.KpiTotalsUseCase
	topFive     *kpiquery.KpiTopFiveUseCase
	valorPorMes *kpiquery.KpiValorPorMesUseCase
	log         *slog.Logger
}

func NewDashboard(totals *kpiquery.KpiTotalsUseCase, topFive *kpiquery.KpiTopFiveUseCase, valorPorMes *kpiquery.KpiValorPorMesUseCase, log *slog.Logger) *Dashboard {
	return &Dashboard{totals: totals, topFive: topFive, valorPorMes: valorPorMes, log: log}
}

// Mount adds the KPI endpoints to api. Every one needs a signed-in caller, so protect is required.
func (h *Dashboard) Mount(api huma.API, protect huma.Middlewares) {
	huma.Register(api, huma.Operation{
		OperationID: "kpi-totals",
		Method:      http.MethodGet,
		Path:        "/dashboard/kpis/totals",
		Summary:     "Headline totals: value, declarations and items",
		Tags:        []string{"dashboard"},
		Middlewares: protect,
		Security:    auth.BearerSecurity,
	}, h.getTotals)

	huma.Register(api, huma.Operation{
		OperationID: "kpi-top-five",
		Method:      http.MethodGet,
		Path:        "/dashboard/kpis/top-5",
		Summary:     "Top 5 countries and top 5 products by import value",
		Tags:        []string{"dashboard"},
		Middlewares: protect,
		Security:    auth.BearerSecurity,
	}, h.getTopFive)

	huma.Register(api, huma.Operation{
		OperationID: "kpi-valor-por-mes",
		Method:      http.MethodGet,
		Path:        "/dashboard/kpis/valor-por-mes",
		Summary:     "Import value per month for a year",
		Tags:        []string{"dashboard"},
		Middlewares: protect,
		Security:    auth.BearerSecurity,
	}, h.getValorPorMes)
}

type totalsOutput struct {
	Body readmodel.KpiTotals
}

func (h *Dashboard) getTotals(ctx context.Context, _ *struct{}) (*totalsOutput, error) {
	totals, err := h.totals.Execute(ctx)
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "kpi_totals"), err)
	}
	return &totalsOutput{Body: totals}, nil
}

type topFiveOutput struct {
	Body kpiquery.TopFive
}

func (h *Dashboard) getTopFive(ctx context.Context, _ *struct{}) (*topFiveOutput, error) {
	top, err := h.topFive.Execute(ctx)
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "kpi_top_five"), err)
	}
	return &topFiveOutput{Body: top}, nil
}

type valorPorMesInput struct {
	Year httpx.OptionalParam[int64] `query:"year" doc:"Year to read. Omit it for the latest year in the data."`
}

type valorPorMesOutput struct {
	Body kpiquery.ValorPorMes
}

func (h *Dashboard) getValorPorMes(ctx context.Context, in *valorPorMesInput) (*valorPorMesOutput, error) {
	out, err := h.valorPorMes.Execute(ctx, in.Year.Value)
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", "kpi_valor_por_mes"), err)
	}
	return &valorPorMesOutput{Body: out}, nil
}
