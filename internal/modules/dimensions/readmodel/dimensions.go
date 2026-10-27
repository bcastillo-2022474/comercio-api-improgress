// Package readmodel holds the read side of the dimensions module: screen-shaped DTOs and the
// query interfaces that produce them. There are no repositories here, only queries.
package readmodel

import "context"

// Option is one entry of a filter dropdown. The three lookup lists share it.
type Option struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// CountriesQuery lists every country, for the country filter.
type CountriesQuery interface {
	ListCountries(ctx context.Context) ([]Option, error)
}

// CustomsQuery lists every customs office (aduana), for the customs filter.
type CustomsQuery interface {
	ListCustoms(ctx context.Context) ([]Option, error)
}

// RegimeTypesQuery lists every customs regime type, for the regime filter.
type RegimeTypesQuery interface {
	ListRegimeTypes(ctx context.Context) ([]Option, error)
}
