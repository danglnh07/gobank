-- name: CreateAccount :one
INSERT INTO account (
    owner, 
    balance, 
    currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM account
WHERE account_id = $1;

-- name: GetAccountForUpdate :one
SELECT * FROM account
WHERE account_id = $1
FOR NO KEY UPDATE;

-- name: ListAccount :many
SELECT * FROM account
ORDER BY account_id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
SET balance = $2
WHERE account_id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE account_id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE account_id = $1;