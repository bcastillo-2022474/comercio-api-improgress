package repository_test

import (
	"testing"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/identity/user/contracts"
	"joao-api-comercio/internal/modules/identity/user/repository"
	"joao-api-comercio/internal/shared/sqlite/testdb"
)

func TestSQLiteRepositoryContract(t *testing.T) {
	contracts.RunRepository(t, func(t *testing.T) user.Repository {
		return repository.NewSQLiteRepository(testdb.Open(t))
	})
}
