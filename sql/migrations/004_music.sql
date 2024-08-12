-- +goose Up
CREATE TABLE
    music (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        artist TEXT NOT NULL,
        album TEXT NOT NULL,
        location TEXT NOT NULL,
        year INT NOT NULL
    );

-- +goose Down
DROP TABLE music;