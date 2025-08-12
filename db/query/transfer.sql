-- name: CreateTransaction :one
INSERT INTO transfer (
    from_account_id,
    to_account_id,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transfer
WHERE transfer_id = $1;

-- name: ListTransaction :many
SELECT * FROM transfer
ORDER BY transfer_id
LIMIT $1
OFFSET $2;

-- name: UpdateTransaction :one
UPDATE transfer
SET amount = $2
WHERE transfer_id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transfer 
WHERE transfer_id = $1;

