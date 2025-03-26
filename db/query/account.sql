-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
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
  set name = $2,
  bio = $3
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;