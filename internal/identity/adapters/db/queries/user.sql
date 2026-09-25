-- name: SaveUser :exec
INSERT INTO identity.users (id, name, email, password, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    password = EXCLUDED.password,
    updated_at = EXCLUDED.updated_at,
    deleted_at = EXCLUDED.deleted_at;

-- name: FindUserByID :one
SELECT id, name, email, password, created_at, updated_at, deleted_at
FROM identity.users
WHERE id = $1;
