// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (email, username, hashed_password, first_name, last_name, date_of_birth)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING user_id, email, username, hashed_password, first_name, last_name, date_of_birth, created_at, updated_at, verified
`

type CreateUserParams struct {
	Email          string
	Username       string
	HashedPassword string
	FirstName      string
	LastName       string
	DateOfBirth    time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Username,
		arg.HashedPassword,
		arg.FirstName,
		arg.LastName,
		arg.DateOfBirth,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.HashedPassword,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Verified,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE email = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, email)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT user_id, email, username, hashed_password, first_name, last_name, date_of_birth, created_at, updated_at, verified
FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Email,
			&i.Username,
			&i.HashedPassword,
			&i.FirstName,
			&i.LastName,
			&i.DateOfBirth,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Verified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHashedPassword = `-- name: GetHashedPassword :one
SELECT hashed_password
FROM users
WHERE email = $1
`

func (q *Queries) GetHashedPassword(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRowContext(ctx, getHashedPassword, email)
	var hashed_password string
	err := row.Scan(&hashed_password)
	return hashed_password, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, email, username, hashed_password, first_name, last_name, date_of_birth, created_at, updated_at, verified
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.HashedPassword,
		&i.FirstName,
		&i.LastName,
		&i.DateOfBirth,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Verified,
	)
	return i, err
}

const isUserVerified = `-- name: IsUserVerified :one
SELECT verified
FROM users
WHERE email = $1
`

func (q *Queries) IsUserVerified(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRowContext(ctx, isUserVerified, email)
	var verified bool
	err := row.Scan(&verified)
	return verified, err
}

const verifyUser = `-- name: VerifyUser :exec
UPDATE users
SET verified = TRUE
WHERE email = $1
`

func (q *Queries) VerifyUser(ctx context.Context, email string) error {
	_, err := q.db.ExecContext(ctx, verifyUser, email)
	return err
}
