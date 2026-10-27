package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/identity/user/repository"
	"joao-api-comercio/internal/shared"
	"joao-api-comercio/internal/shared/sqlite"
	"joao-api-comercio/internal/shared/sqlite/testdb"
)

var created = time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)

func newUnitOfWork(db *sql.DB) *sqlite.UnitOfWork[user.Repository] {
	return sqlite.NewUnitOfWork(db, func(tx *sql.Tx) user.Repository {
		return repository.NewSQLiteRepository(tx)
	}, func() time.Time { return created })
}

// registerUser saves u and publishes its registration event inside the unit of work.
func registerUser(u user.User) func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
	return func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
		if err := users.Save(ctx, u); err != nil {
			return err
		}
		pub.Publish(user.NewUserRegisteredEvent(u))
		return nil
	}
}

func TestUnitOfWork_CommitWritesUserAndOutboxRecord(t *testing.T) {
	db := testdb.Open(t)
	u := user.NewUser("ana@example.com", "hash", created)

	if err := newUnitOfWork(db).Execute(context.Background(), registerUser(u)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCount(t, db, "users", 1)
	assertCount(t, db, "outbox", 1)
	var eventName string
	if err := db.QueryRow(`SELECT event_name FROM outbox`).Scan(&eventName); err != nil {
		t.Fatalf("reading outbox: %v", err)
	}
	if eventName != "user.registered" {
		t.Errorf("event_name = %q, want user.registered", eventName)
	}
}

func TestUnitOfWork_FailureRollsBackUserAndOutbox(t *testing.T) {
	db := testdb.Open(t)
	boom := errors.New("boom")

	err := newUnitOfWork(db).Execute(context.Background(), func(ctx context.Context, users user.Repository, pub shared.Publisher) error {
		if err := registerUser(user.NewUser("ana@example.com", "hash", created))(ctx, users, pub); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}

	assertCount(t, db, "users", 0)
	assertCount(t, db, "outbox", 0)
}

func TestUnitOfWork_RejectedDuplicateWritesNoOutboxRecord(t *testing.T) {
	db := testdb.Open(t)
	uow := newUnitOfWork(db)
	if err := uow.Execute(context.Background(), registerUser(user.NewUser("ana@example.com", "hash-1", created))); err != nil {
		t.Fatalf("seeding user: %v", err)
	}

	err := uow.Execute(context.Background(), registerUser(user.NewUser("ana@example.com", "hash-2", created)))
	var taken *user.EmailTakenError
	if !errors.As(err, &taken) {
		t.Fatalf("err = %v, want *EmailTakenError", err)
	}

	assertCount(t, db, "users", 1)
	assertCount(t, db, "outbox", 1)
}

func assertCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&got); err != nil {
		t.Fatalf("counting %s: %v", table, err)
	}
	if got != want {
		t.Errorf("%s rows = %d, want %d", table, got, want)
	}
}
