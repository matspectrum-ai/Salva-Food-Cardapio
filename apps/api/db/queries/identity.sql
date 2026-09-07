-- name: CreateAppUser :one
INSERT INTO app_users (
    id, name, cpf, email, phone, image_url, password_hash
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, cpf, email, phone, image_url, password_hash, created_at, updated_at;

-- name: GetAppUserByEmail :one
SELECT id, name, cpf, email, phone, image_url, password_hash, created_at, updated_at
FROM app_users
WHERE lower(email) = lower($1)
LIMIT 1;

-- name: GetAppUserByID :one
SELECT id, name, cpf, email, phone, image_url, password_hash, created_at, updated_at
FROM app_users
WHERE id = $1;

-- name: CreateTenantMembership :one
INSERT INTO tenant_memberships (
    tenant_id, user_id, title, status, permissions
)
VALUES ($1, $2, $3, $4, $5)
RETURNING tenant_id, user_id, title, status, permissions, created_at, updated_at;
-- name: GetTenantMembership :one
SELECT tenant_id, user_id, title, status, permissions, created_at, updated_at
FROM tenant_memberships
WHERE tenant_id = $1 AND user_id = $2;

-- name: SetTenantMembershipStatus :one
UPDATE tenant_memberships
SET status = $3, updated_at = now()
WHERE tenant_id = $1 AND user_id = $2
RETURNING tenant_id, user_id, title, status, permissions, created_at, updated_at;

-- name: ListCollaborators :many
SELECT
    u.id, u.name, u.cpf, u.email, u.phone, u.image_url, u.password_hash,
    m.tenant_id, m.title, m.status, m.permissions
FROM tenant_memberships m
JOIN app_users u ON u.id = m.user_id
WHERE m.tenant_id = sqlc.arg(tenant_id)
  AND (
    sqlc.arg(search)::text = '' OR
    lower(u.name) LIKE '%' || lower(sqlc.arg(search)::text) || '%' OR
    lower(u.email) LIKE '%' || lower(sqlc.arg(search)::text) || '%' OR
    lower(u.phone) LIKE '%' || lower(sqlc.arg(search)::text) || '%' OR
    lower(m.title) LIKE '%' || lower(sqlc.arg(search)::text) || '%'
  )
ORDER BY u.name, u.id;
-- name: CreateAuthSession :one
INSERT INTO auth_sessions (
    id, user_id, active_tenant_id, token_hash, expires_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, active_tenant_id, token_hash, expires_at, revoked_at, created_at;

-- name: GetAuthSessionByTokenHash :one
SELECT id, user_id, active_tenant_id, token_hash, expires_at, revoked_at, created_at
FROM auth_sessions
WHERE token_hash = $1;

-- name: RevokeAuthSession :one
UPDATE auth_sessions
SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL
RETURNING id, user_id, active_tenant_id, token_hash, expires_at, revoked_at, created_at;

-- name: ListTenantMembershipsByUser :many
SELECT tenant_id, user_id, title, status, permissions, created_at, updated_at
FROM tenant_memberships
WHERE user_id = $1
ORDER BY tenant_id;
