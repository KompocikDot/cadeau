-- name: CreateUser :exec
INSERT INTO users(username, password) values(?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE username = ?;

