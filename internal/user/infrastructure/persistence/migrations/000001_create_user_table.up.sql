BEGIN;

CREATE TABLE IF NOT EXISTS my_schema.users
(
    id uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

COMMIT;
