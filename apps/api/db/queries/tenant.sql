-- name: CreateTenant :one
INSERT INTO tenants (id, name)
VALUES ($1, $2)
RETURNING *;

-- name: CreateEstablishment :one
INSERT INTO establishments (id, tenant_id, name, timezone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTenant :one
SELECT * FROM tenants WHERE id = $1;

-- name: GetEstablishment :one
SELECT *
FROM establishments
WHERE tenant_id = $1 AND id = $2;
