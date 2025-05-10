-- name: CreateUser :one
INSERT INTO users (
  name, bio, profile_image, location, website,
  followers, following, email, stories_count, is_verified, password
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET name = $2, bio = $3, profile_image = $4, location = $5, website = $6,
followers = $7, following = $8, stories_count = $9, is_verified = $10
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users;


-- name: ListStories :many
SELECT * FROM stories;

-- name: CreateStory :one
INSERT INTO stories (
  id, title, description, cover_image, user_id, likes, views,
  published_date, last_edited, story_type, status, genres
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetStoryByID :one
SELECT * FROM stories WHERE id = $1;

-- name: UpdateStory :one
UPDATE stories SET
  title = $2,
  description = $3,
  cover_image = $4,
  user_id = $5,
  likes = $6,
  views = $7,
  published_date = $8,
  last_edited = $9,
  story_type = $10,
  status = $11,
  genres = $12
WHERE id = $1
RETURNING *;

-- name: DeleteStory :exec
DELETE FROM stories WHERE id = $1;

-- name: GetStoriesByUser :many
SELECT * FROM stories WHERE user_id = $1;



-- name: CreateChapter :one
INSERT INTO chapters (id, story_id, title, content, is_complete)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetChapterByID :one
SELECT * FROM chapters WHERE id = $1;

-- name: UpdateChapter :one
UPDATE chapters SET story_id = $2, title = $3, content = $4, is_complete = $5
WHERE id = $1
RETURNING *;

-- name: DeleteChapter :exec
DELETE FROM chapters WHERE id = $1;



-- name: AddCollaborator :exec
INSERT INTO story_collaborators (story_id, user_id)
VALUES ($1, $2);

-- name: GetCollaborator :one
SELECT * FROM story_collaborators WHERE story_id = $1 AND user_id = $2;

-- name: RemoveCollaborator :exec
DELETE FROM story_collaborators WHERE story_id = $1 AND user_id = $2;
