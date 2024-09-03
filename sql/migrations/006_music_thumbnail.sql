-- +goose Up
ALTER TABLE music
ADD COLUMN thumbnail_location TEXT NOT NULL DEFAULT 'music/thumbnails/default.png';

-- +goose Down
ALTER TABLE music
DROP COLUMN thumbnail_location;