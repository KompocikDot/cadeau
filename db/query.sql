-- name: CreateUser :exec
INSERT INTO users(username, password) values(?, ?);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: CreateOccasion :exec
INSERT INTO occasions(name, gift_receiver) values(?, ?);

-- name: GetOccasionByUserId :many
SELECT * FROM occasions WHERE gift_receiver = ?;

-- name: DeleteOccasion :exec
DELETE FROM occasions WHERE id = ?;
