-- name: GetMusicByID :one
SELECT * FROM music WHERE id = $1;

-- name: GetMusicList :many
SELECT *
FROM music
LIMIT $1
OFFSET $2;

-- name: CreateMusic :one
INSERT INTO music (title, artist, album, location, year)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;