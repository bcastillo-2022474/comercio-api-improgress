package sqlite_test

import (
	"database/sql"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/identity/user/repository"
	"joao-api-comercio/internal/shared/contracts"
	outboxrepo "joao-api-comercio/internal/shared/outbox/repository"
	"joao-api-comercio/internal/shared/sqlite"
	"joao-api-comercio/internal/shared/sqlite/testdb"
)

func TestSQLiteUnitOfWorkContract(t *testing.T) {
	contracts.RunUnitOfWork(t, func(t *testing.T) contracts.Harness {
		db := testdb.Open(t)
		uow := sqlite.NewUnitOfWork(db, func(tx *sql.Tx) user.Repository {
			return repository.NewSQLiteRepository(tx)
		}, time.Now)
		return contracts.Harness{
			UoW:    uow,
			Users:  repository.NewSQLiteRepository(db),
			Outbox: outboxrepo.NewSQLiteRepository(db),
		}
	})
}
