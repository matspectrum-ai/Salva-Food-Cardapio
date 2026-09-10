-- name: CreateCustomer :one
INSERT INTO customers (tenant_id, id, name, phone, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE tenant_id = $1 AND id = $2;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE tenant_id = $1
ORDER BY name, id;

-- name: CreateCustomerAddress :one
INSERT INTO customer_addresses (
    tenant_id, id, customer_id, label, street, number,
    complement, neighborhood, city, state, postal_code, reference
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetCustomerAddress :one
SELECT * FROM customer_addresses
WHERE tenant_id = $1 AND id = $2;

-- name: ListCustomerAddresses :many
SELECT * FROM customer_addresses
WHERE tenant_id = $1 AND customer_id = $2
ORDER BY id;
