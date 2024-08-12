-- name: GetMusicByID :one
SELECT * FROM music WHERE id = $1;