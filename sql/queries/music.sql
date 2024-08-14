-- name: GetMusicByID :one
SELECT * FROM music WHERE id = $1;

-- name: GetMusicList :many
SELECT *
FROM music
LIMIT $1
OFFSET $2;