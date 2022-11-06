-- name: FindLink :one
SELECT url
FROM links
WHERE base = $1
    AND deleted_at IS NULL;

-- name: AllLinks :many
SELECT base,
    url
FROM links
WHERE deleted_at IS NULL;

-- name: UpsertLink :exec
INSERT INTO links (base, url)
VALUES ($1, $2) ON conflict (base) DO
UPDATE
SET url = $2,
    updated_at = now(),
    deleted_at = NULL;

-- name: DeleteLink :exec
UPDATE links
SET deleted_at = now()
WHERE base = $1;
