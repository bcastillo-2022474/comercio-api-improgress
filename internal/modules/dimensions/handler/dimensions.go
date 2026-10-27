// Package handler exposes the dimensions read side over HTTP with Huma.
package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"joao-api-comercio/internal/modules/dimensions/readmodel"
	dimensionsquery "joao-api-comercio/internal/modules/dimensions/use_cases/query"
	"joao-api-comercio/internal/shared/httperr"
)

type Dimensions struct {
	listCountries   *dimensionsquery.ListCountriesUseCase
	listCustoms     *dimensionsquery.ListCustomsUseCase
	listRegimeTypes *dimensionsquery.ListRegimeTypesUseCase
	log             *slog.Logger
}

func NewDimensions(listCountries *dimensionsquery.ListCountriesUseCase, listCustoms *dimensionsquery.ListCustomsUseCase, listRegimeTypes *dimensionsquery.ListRegimeTypesUseCase, log *slog.Logger) *Dimensions {
	return &Dimensions{listCountries: listCountries, listCustoms: listCustoms, listRegimeTypes: listRegimeTypes, log: log}
}

// Mount adds the lookup endpoints to api.
func (h *Dimensions) Mount(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-countries",
		Method:      http.MethodGet,
		Path:        "/dimensions/paises",
		Summary:     "List countries",
		Tags:        []string{"dimensions"},
	}, h.getCountries)

	huma.Register(api, huma.Operation{
		OperationID: "list-customs",
		Method:      http.MethodGet,
		Path:        "/dimensions/aduanas",
		Summary:     "List customs offices",
		Tags:        []string{"dimensions"},
	}, h.getCustoms)

	huma.Register(api, huma.Operation{
		OperationID: "list-regime-types",
		Method:      http.MethodGet,
		Path:        "/dimensions/tipos-regimen",
		Summary:     "List customs regime types",
		Tags:        []string{"dimensions"},
	}, h.getRegimeTypes)
}

type optionsOutput struct {
	Body []readmodel.Option
}

func (h *Dimensions) getCountries(ctx context.Context, _ *struct{}) (*optionsOutput, error) {
	return h.respond(ctx, "list_countries", h.listCountries.Execute)
}

func (h *Dimensions) getCustoms(ctx context.Context, _ *struct{}) (*optionsOutput, error) {
	return h.respond(ctx, "list_customs", h.listCustoms.Execute)
}

func (h *Dimensions) getRegimeTypes(ctx context.Context, _ *struct{}) (*optionsOutput, error) {
	return h.respond(ctx, "list_regime_types", h.listRegimeTypes.Execute)
}

func (h *Dimensions) respond(ctx context.Context, operation string, list func(context.Context) ([]readmodel.Option, error)) (*optionsOutput, error) {
	options, err := list(ctx)
	if err != nil {
		return nil, httperr.ToHumaError(ctx, h.log.With("operation", operation), err)
	}
	return &optionsOutput{Body: options}, nil
}
