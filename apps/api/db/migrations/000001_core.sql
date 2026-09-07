-- +goose Up

CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE establishments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    timezone TEXT NOT NULL DEFAULT 'America/Santarem',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id)
);

CREATE TABLE catalog_categories (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    sort_order INTEGER NOT NULL CHECK (sort_order >= 0),
    sold_out BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id)
);

CREATE INDEX catalog_categories_tenant_sort_idx
    ON catalog_categories (tenant_id, sort_order, id);

CREATE TABLE catalog_items (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id UUID NOT NULL,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    sort_order INTEGER NOT NULL CHECK (sort_order >= 0),
    sold_out BOOLEAN NOT NULL DEFAULT FALSE,
    sold_by_weight BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id),
    CONSTRAINT catalog_items_category_fk
        FOREIGN KEY (tenant_id, category_id)
        REFERENCES catalog_categories (tenant_id, id)
        ON DELETE RESTRICT
);

CREATE INDEX catalog_items_tenant_category_sort_idx
    ON catalog_items (tenant_id, category_id, sort_order, id);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    source TEXT NOT NULL CHECK (btrim(source) <> ''),
    status TEXT NOT NULL CHECK (status IN ('ANALYSIS', 'PRODUCTION', 'READY', 'FINALIZED', 'CANCELLED')),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id)
);

CREATE INDEX orders_tenant_created_idx
    ON orders (tenant_id, created_at DESC, id);
CREATE INDEX orders_tenant_status_created_idx
    ON orders (tenant_id, status, created_at DESC, id);

CREATE TABLE order_items (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id UUID NOT NULL,
    line_no INTEGER NOT NULL CHECK (line_no > 0),
    item_id UUID NOT NULL,
    name_snapshot TEXT NOT NULL CHECK (btrim(name_snapshot) <> ''),
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    PRIMARY KEY (tenant_id, order_id, line_no),
    CONSTRAINT order_items_order_fk
        FOREIGN KEY (tenant_id, order_id)
        REFERENCES orders (tenant_id, id)
        ON DELETE CASCADE
);

CREATE TABLE order_idempotency (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    idempotency_key TEXT NOT NULL CHECK (btrim(idempotency_key) <> ''),
    request_fingerprint CHAR(64) NOT NULL,
    order_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, idempotency_key),
    CONSTRAINT order_idempotency_order_fk
        FOREIGN KEY (tenant_id, order_id)
        REFERENCES orders (tenant_id, id)
        ON DELETE CASCADE
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    aggregate_type TEXT NOT NULL CHECK (btrim(aggregate_type) <> ''),
    aggregate_id UUID NOT NULL,
    event_type TEXT NOT NULL CHECK (btrim(event_type) <> ''),
    payload JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ,
    attempts INTEGER NOT NULL DEFAULT 0 CHECK (attempts >= 0)
);

CREATE INDEX outbox_events_unpublished_idx
    ON outbox_events (occurred_at, id)
    WHERE published_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS order_idempotency;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS catalog_items;
DROP TABLE IF EXISTS catalog_categories;
DROP TABLE IF EXISTS establishments;
DROP TABLE IF EXISTS tenants;
