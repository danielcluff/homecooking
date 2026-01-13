-- name: CreateUserInvite :one
INSERT INTO user_invites (code, email, role, created_by, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInviteByCode :one
SELECT * FROM user_invites
WHERE code = $1
  AND (expires_at IS NULL OR expires_at > NOW())
  AND used_at IS NULL
LIMIT 1;

-- name: UseInvite :one
UPDATE user_invites
SET
    used_at = NOW(),
    used_by = $2
WHERE id = $1
RETURNING *;

-- name: ListInvites :many
SELECT * FROM user_invites
ORDER BY created_at DESC;

-- name: DeleteInvite :exec
DELETE FROM user_invites WHERE id = $1;
