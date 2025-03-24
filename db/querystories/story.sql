-- name: GetStory :one
SELECT * FROM stories
WHERE story_id = $1 LIMIT 1;

-- name: ListStories :many
SELECT * FROM stories
ORDER BY story_id;

-- name: CreateStory :one
INSERT INTO stories (
  originalstory, pulledrequests, updatedstory, author_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateStory :exec
UPDATE stories
SET originalstory = $2,
    pulledrequests = $3,
    updatedstory = $4,
    author_id = $5
WHERE story_id = $1;

-- name: DeleteStory :exec
DELETE FROM stories
WHERE story_id = $1;
