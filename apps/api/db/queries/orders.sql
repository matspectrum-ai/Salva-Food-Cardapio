-- name: CreateOrder :one
INSERT INTO orders (id, tenant_id, source, status, total_cents)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    tenant_id, order_id, line_no, item_id, name_snapshot, quantity, unit_price_cents
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CreateOrderIdempotency :one
INSERT INTO order_idempotency (
    tenant_id, idempotency_key, request_fingerprint, order_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderIdempotency :one
SELECT *
FROM order_idempotency
WHERE tenant_id = $1 AND idempotency_key = $2;

-- name: GetOrder :one
SELECT *
FROM orders
WHERE tenant_id = $1 AND id = $2;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE tenant_id = $1
ORDER BY created_at DESC, id;

-- name: ListOrderItems :many
SELECT *
FROM order_items
WHERE tenant_id = $1 AND order_id = $2
ORDER BY line_no;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2
RETURNING *;

-- name: LockOrderForUpdate :one
SELECT *
FROM orders
WHERE tenant_id = $1 AND id = $2
FOR UPDATE;

-- name: DeleteOrderIdempotency :exec
DELETE FROM order_idempotency
WHERE tenant_id = $1 AND idempotency_key = $2;
