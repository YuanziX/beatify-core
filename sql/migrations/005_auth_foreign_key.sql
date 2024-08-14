-- +goose Up
ALTER TABLE auth ADD CONSTRAINT fk_user_email FOREIGN KEY (user_email) REFERENCES users (email) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE auth
DROP CONSTRAINT fk_user_email;