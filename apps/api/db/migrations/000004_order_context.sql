-- +goose Up

ALTER TABLE orders
    ADD COLUMN customer_id UUID,
    ADD COLUMN customer_name_snapshot TEXT,
    ADD COLUMN customer_phone_snapshot TEXT,
    ADD COLUMN customer_email_snapshot TEXT,
    ADD COLUMN address_id UUID,
    ADD COLUMN address_snapshot JSONB,
    ADD COLUMN fulfillment_type TEXT NOT NULL DEFAULT 'PICKUP',
    ADD COLUMN scheduled_at TIMESTAMPTZ,
    ADD COLUMN notes TEXT NOT NULL DEFAULT '';

ALTER TABLE orders
    ADD CONSTRAINT orders_customer_fk
    FOREIGN KEY (tenant_id, customer_id)
    REFERENCES customers (tenant_id, id)
    ON DELETE SET NULL,
    ADD CONSTRAINT orders_address_fk
    FOREIGN KEY (tenant_id, address_id)
    REFERENCES customer_addresses (tenant_id, id)
    ON DELETE SET NULL;

ALTER TABLE orders
    ADD CONSTRAINT orders_fulfillment_type_ck
    CHECK (fulfillment_type IN ('DELIVERY', 'PICKUP', 'DINE_IN'));

CREATE INDEX orders_tenant_customer_idx
    ON orders (tenant_id, customer_id, created_at DESC, id);

CREATE INDEX orders_tenant_scheduled_idx
    ON orders (tenant_id, scheduled_at, id)
    WHERE scheduled_at IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS orders_tenant_scheduled_idx;
DROP INDEX IF EXISTS orders_tenant_customer_idx;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_fulfillment_type_ck;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_address_fk;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_customer_fk;
ALTER TABLE orders
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS scheduled_at,
    DROP COLUMN IF EXISTS fulfillment_type,
    DROP COLUMN IF EXISTS address_snapshot,
    DROP COLUMN IF EXISTS address_id,
    DROP COLUMN IF EXISTS customer_email_snapshot,
    DROP COLUMN IF EXISTS customer_phone_snapshot,
    DROP COLUMN IF EXISTS customer_name_snapshot,
    DROP COLUMN IF EXISTS customer_id;
