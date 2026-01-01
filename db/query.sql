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

-- name: GetOccasionGuests :many
SELECT id, username FROM users AS u JOIN user_occasions AS uo on u.id = uo.user_id WHERE uo.occasion_id = ?;

-- name: GetOccasionById :one
SELECT * FROM occasions WHERE id = ?;

-- name: UpdateOccasion :exec
UPDATE occasions SET name = ? WHERE id = ?;

-- name: DeleteOccasion :exec
DELETE FROM occasions WHERE id = ?;

-- name: CreateGift :execlastid
INSERT INTO gifts(name, url, occasion) values(?, ?, ?);

-- name: DeleteGift :exec
DELETE FROM gifts WHERE id = ?;

-- name: SelectGiftsByOcassionId :many
SELECT g.* FROM gifts AS g
	JOIN occasions AS o ON g.occasion = o.id
	WHERE g.occasion = ? AND o.gift_receiver = ?;

-- name: GetGiftById :one
SELECT g.*, o.gift_receiver FROM gifts AS g
	LEFT JOIN occasions AS o ON g.occasion = o.id
	WHERE g.id = ?;

-- name: UpdateGift :exec
UPDATE gifts SET name = ?, url = ? WHERE id = ?;

-- name: AssignUserToOccasion :exec
INSERT INTO user_occasions(occasion_id, user_id) values(?, ?);

-- name: RemoveUsersFromOccasion :exec
DELETE FROM gifts WHERE id IN sqlc.slice(ids);

-- name: GetUsersList :many
SELECT id, username FROM users WHERE id != ?;
