-- name: CreateEntry :one
INSERT INTO entry (
    account_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entry
WHERE entry_id = $1;

-- name: ListEntry :many
SELECT * FROM entry
ORDER BY entry_id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE entry
SET amount = $2
WHERE entry_id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entry 
WHERE entry_id = $1;

