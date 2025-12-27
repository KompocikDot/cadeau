-- name: CreateUser :exec
INSERT INTO users(username, password) values(?, ?);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: CreateOccasion :execlastid
INSERT INTO occasions(name, gift_receiver) values(?, ?);

-- name: GetUserOccasionsByUserId :many
SELECT * FROM occasions WHERE gift_receiver = ?;

-- name: GetOccasionById :one
SELECT * FROM occasions WHERE id = ? AND gift_receiver = ?;

-- name: DeleteOccasion :exec
DELETE FROM occasions WHERE id = ?;

-- name: CreateGift :exec
INSERT INTO gifts(name, url, occasion) values(?, ?, ?);

-- name: DeleteGift :exec
DELETE FROM gifts WHERE id = ?;

-- name: SelectGiftsByOcassionId :many
SELECT *  FROM gifts AS g JOIN occasions AS o ON g.occasion = o.id WHERE g.occasion = ? AND o.gift_receiver = ?;
