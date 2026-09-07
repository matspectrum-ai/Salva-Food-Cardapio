-- +goose Up

CREATE TABLE app_users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    cpf TEXT NOT NULL CHECK (char_length(cpf) = 11 AND cpf ~ '^[0-9]{11}$'),
    email TEXT NOT NULL CHECK (btrim(email) <> '' AND email = lower(email)),
    phone TEXT NOT NULL CHECK (btrim(phone) <> ''),
    image_url TEXT,
    password_hash TEXT NOT NULL CHECK (btrim(password_hash) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (cpf)
);

CREATE UNIQUE INDEX app_users_email_lower_uidx
    ON app_users (lower(email));

CREATE TABLE tenant_memberships (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    title TEXT NOT NULL CHECK (btrim(title) <> ''),
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    permissions TEXT[] NOT NULL CHECK (cardinality(permissions) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, user_id)
);
CREATE INDEX tenant_memberships_tenant_status_idx
    ON tenant_memberships (tenant_id, status, user_id);

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    active_tenant_id UUID NOT NULL,
    token_hash TEXT NOT NULL CHECK (char_length(token_hash) = 64),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT auth_sessions_membership_fk
        FOREIGN KEY (active_tenant_id, user_id)
        REFERENCES tenant_memberships (tenant_id, user_id)
        ON DELETE CASCADE,
    UNIQUE (token_hash)
);

CREATE INDEX auth_sessions_user_active_idx
    ON auth_sessions (user_id, active_tenant_id, expires_at DESC)
    WHERE revoked_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS tenant_memberships;
DROP TABLE IF EXISTS app_users;
