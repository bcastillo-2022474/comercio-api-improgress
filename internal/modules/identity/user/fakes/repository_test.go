package fakes_test

import (
	"testing"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/identity/user/contracts"
	"joao-api-comercio/internal/modules/identity/user/fakes"
)

func TestFakeRepositoryContract(t *testing.T) {
	contracts.RunRepository(t, func(*testing.T) user.Repository {
		return fakes.NewRepository()
	})
}
