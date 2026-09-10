-- +goose Up

CREATE TABLE customers (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    id UUID NOT NULL,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    phone TEXT NOT NULL CHECK (btrim(phone) <> ''),
    email TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, id)
);
CREATE INDEX customers_tenant_phone_idx ON customers (tenant_id, phone);
CREATE INDEX customers_tenant_name_idx ON customers (tenant_id, name);

CREATE TABLE customer_addresses (
    tenant_id UUID NOT NULL,
    id UUID NOT NULL,
    customer_id UUID NOT NULL,
    label TEXT,
    street TEXT NOT NULL,
    number TEXT NOT NULL,
    complement TEXT,
    neighborhood TEXT NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    reference TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, id),
    CONSTRAINT customer_addresses_customer_fk
        FOREIGN KEY (tenant_id, customer_id)
        REFERENCES customers (tenant_id, id) ON DELETE CASCADE
);
CREATE INDEX customer_addresses_customer_idx ON customer_addresses (tenant_id, customer_id);

-- +goose Down
DROP TABLE IF EXISTS customer_addresses;
DROP TABLE IF EXISTS customers;
