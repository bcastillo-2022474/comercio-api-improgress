-- name: UpsertUser :exec
INSERT INTO users (id, email, password_hash, created_at)
VALUES (@id, @email, @password_hash, @created_at)
ON CONFLICT (id) DO UPDATE SET
    email = excluded.email,
    password_hash = excluded.password_hash;

-- name: FindUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE email = @email
LIMIT 1;
