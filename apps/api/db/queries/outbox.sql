-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (
    id, tenant_id, aggregate_type, aggregate_id, event_type, payload
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListPendingOutboxEvents :many
SELECT *
FROM outbox_events
WHERE published_at IS NULL
ORDER BY occurred_at, id
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events
SET published_at = now()
WHERE id = $1;

-- name: IncrementOutboxEventAttempts :exec
UPDATE outbox_events
SET attempts = attempts + 1
WHERE id = $1;
