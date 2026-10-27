// Package repository holds the SQL adapter for the user port.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"joao-api-comercio/internal/modules/identity/user"
	"joao-api-comercio/internal/modules/identity/user/repository/db"
	"joao-api-comercio/internal/shared/types"
)

// emailConstraint is the message SQLite gives when the unique email index rejects a row.
const emailConstraint = "UNIQUE constraint failed: users.email"

type SQLiteRepository struct {
	q *db.Queries
}

func NewSQLiteRepository(q db.DBTX) *SQLiteRepository {
	return &SQLiteRepository{q: db.New(q)}
}

func (r *SQLiteRepository) Save(ctx context.Context, u user.User) error {
	err := r.q.UpsertUser(ctx, db.UpsertUserParams{
		ID:           u.Id.String(),
		Email:        string(u.Email),
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt.UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		if strings.Contains(err.Error(), emailConstraint) {
			return &user.EmailTakenError{Email: u.Email}
		}
		return fmt.Errorf("upserting user: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) FindByEmail(ctx context.Context, email user.Email) (types.Maybe[user.User], error) {
	row, err := r.q.FindUserByEmail(ctx, string(email))
	if errors.Is(err, sql.ErrNoRows) {
		return types.None[user.User](), nil
	}
	if err != nil {
		return types.None[user.User](), fmt.Errorf("finding user by email: %w", err)
	}

	id, err := types.ParseUserId(row.ID)
	if err != nil {
		return types.None[user.User](), fmt.Errorf("parsing stored user id: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, row.CreatedAt)
	if err != nil {
		return types.None[user.User](), fmt.Errorf("parsing stored created_at: %w", err)
	}

	return types.Some(user.User{
		Id:           id,
		Email:        user.Email(row.Email),
		PasswordHash: row.PasswordHash,
		CreatedAt:    createdAt.UTC(),
	}), nil
}
