// Package testdb opens a fresh, fully migrated SQLite database for each test.
package testdb

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"joao-api-comercio/internal/shared/sqlite"
	"joao-api-comercio/migrations"
)

func Open(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	db, err := sqlite.Open(ctx, filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("migrating test database: %v", err)
	}
	return db
}
