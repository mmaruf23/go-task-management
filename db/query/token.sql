-- name: CreateToken :exec
INSERT INTO
    tokens (id, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING
    id;

-- name: RevokeToken :execrows
UPDATE tokens
SET
    revoked_at = now()
WHERE
    id = $1
    AND revoked_at IS NULL;

-- name: RevokeAllToken :execrows
UPDATE tokens
SET
    revoked_at = now()
WHERE
    user_id = $1
    AND revoked_at IS NULL;

-- name: ReplaceToken :execrows
UPDATE tokens
SET
    replaced_by = $2
WHERE
    id = $1
    AND replaced_by IS NULL
    AND revoked_at IS NULL
    AND expires_at > now();