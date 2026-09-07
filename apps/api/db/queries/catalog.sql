-- name: CreateCategory :one
INSERT INTO catalog_categories (id, tenant_id, name, sort_order, sold_out)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCategory :one
SELECT *
FROM catalog_categories
WHERE tenant_id = $1 AND id = $2;

-- name: ListCategories :many
SELECT *
FROM catalog_categories
WHERE tenant_id = $1
ORDER BY sort_order, id;

-- name: CreateItem :one
INSERT INTO catalog_items (
    id, tenant_id, category_id, name, price_cents, sort_order, sold_out, sold_by_weight
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetItem :one
SELECT *
FROM catalog_items
WHERE tenant_id = $1 AND id = $2;

-- name: ListItems :many
SELECT *
FROM catalog_items
WHERE tenant_id = $1
ORDER BY sort_order, id;

-- name: ListItemsByCategory :many
SELECT *
FROM catalog_items
WHERE tenant_id = $1 AND category_id = $2
ORDER BY sort_order, id;

-- name: SetCategorySoldOut :one
UPDATE catalog_categories
SET sold_out = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2
RETURNING *;

-- name: SetItemSoldOut :one
UPDATE catalog_items
SET sold_out = $3, updated_at = now()
WHERE tenant_id = $1 AND id = $2
RETURNING *;
