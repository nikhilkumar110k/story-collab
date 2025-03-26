CREATE DATABASE project3;

-- name: GetAuthor :one
SELECT id, name, bio, email, password FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT id, name, bio, email, password FROM authors
ORDER BY name;

-- name: GetAuthorsByEmail :one
SELECT id, password FROM authors
WHERE email = $1
LIMIT 1;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, email, password
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, bio, email, password;

-- name: UpdateAuthor :exec
UPDATE authors
SET name = $2,
    bio = $3,
    email = $4,
    password = $5
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- Ensure required columns exist
SELECT column_name FROM information_schema.columns WHERE table_name = 'authors';

-- Add missing columns if necessary
ALTER TABLE authors ADD COLUMN IF NOT EXISTS email TEXT UNIQUE NOT NULL;
ALTER TABLE authors ADD COLUMN IF NOT EXISTS password TEXT NOT NULL;
