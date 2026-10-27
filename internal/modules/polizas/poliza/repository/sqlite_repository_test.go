package repository_test

import (
	"testing"

	"joao-api-comercio/internal/modules/polizas/poliza"
	"joao-api-comercio/internal/modules/polizas/poliza/contracts"
	"joao-api-comercio/internal/modules/polizas/poliza/repository"
	"joao-api-comercio/internal/shared/sqlite/testdb"
)

func TestSQLiteRepositoryContract(t *testing.T) {
	contracts.RunRepository(t, func(t *testing.T) poliza.Repository {
		db := testdb.Open(t)
		// The contract's items point at lookup rows, which the real database requires to exist.
		for _, stmt := range []string{
			`INSERT INTO sac (id, fraccion) VALUES (1, 1)`,
			`INSERT INTO pais (id, name) VALUES (1, 'PAIS')`,
			`INSERT INTO tipo_unidad_medida (id, name) VALUES (1, 'UNIDAD')`,
			`INSERT INTO aduanas (id, name) VALUES (1, 'ADUANA')`,
		} {
			if _, err := db.Exec(stmt); err != nil {
				t.Fatalf("seeding lookup row: %v", err)
			}
		}
		return repository.NewSQLiteRepository(db)
	})
}
