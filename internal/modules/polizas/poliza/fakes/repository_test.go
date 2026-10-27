package fakes_test

import (
	"testing"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/poliza/contracts"
	"joao-api-comercio/internal/modules/polizas/poliza/fakes"
)

func TestFakeRepositoryContract(t *testing.T) {
	contracts.RunRepository(t, func(*testing.T) poliza.Repository {
		return fakes.NewRepository()
	})
}
