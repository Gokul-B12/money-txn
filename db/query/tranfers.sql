-- name: CreateTransfer :one
INSERT INTO Transfers (
  from_accound_id,
  to_accound_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM Transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM Transfers WHERE from_accound_id = $1 OR to_accound_id = $2
ORDER BY id
LIMIT $3
OFFSET $4; 



