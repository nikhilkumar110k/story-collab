CREATE DATABASE project3;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT
);
-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id, name, bio;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;

-- name: GetAuthor :one
SELECT id, name, bio 
FROM authors
WHERE id = $1 
LIMIT 1;

-- name: ListAuthors :many
SELECT id, name, bio 
FROM authors
ORDER BY name;
