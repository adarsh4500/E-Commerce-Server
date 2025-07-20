-- name: AddUser :exec
INSERT INTO users (fullname, email, password, role)
VALUES ($1, $2, $3, $4);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: AddProduct :one
INSERT INTO products (name, price, description, stock_quantity) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetProducts :many
SELECT * FROM products;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: DeleteProductById :one
DELETE FROM products 
WHERE id=$1 RETURNING *;

-- name: UpdateProductById :one
UPDATE products
SET 
    name = COALESCE($2, name),
    price = COALESCE($3, price),
    description = COALESCE($4, description),
    stock_quantity = COALESCE($5, stock_quantity)
WHERE
    id = $1 RETURNING *;

-- name: ViewCart :many
SELECT cart.id, cart.user_id, cart.product_id, cart.quantity, cart.modified_at,
       products.id as product_id, products.name as product_name, products.price as product_price, products.description as product_description, products.stock_quantity as product_stock_quantity
FROM cart
JOIN products ON cart.product_id = products.id
WHERE cart.user_id = $1;

-- name: AddToCart :one
INSERT INTO cart (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id)
DO UPDATE SET quantity = cart.quantity + EXCLUDED.quantity, modified_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: RemoveFromCart :one
DELETE FROM cart
WHERE user_id = $1 AND product_id = $2 RETURNING *;

-- name: ClearCart :exec
DELETE FROM cart
WHERE user_id = $1;

-- name: AddOrder :one
INSERT INTO orders (customer_id, total_amount, status)
VALUES ($1, $2, 'Pending') RETURNING id;

-- name: UpdateOrderTotal :one
UPDATE orders SET total_amount = $2
WHERE id = $1 RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $1
WHERE id = $2 RETURNING *;

-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, subtotal)
SELECT $1, $2, $3::INTEGER, price * $3::NUMERIC
FROM products 
WHERE products.id = $2
RETURNING subtotal;

-- name: ViewOrders :many
SELECT id, order_date, total_amount, status
FROM orders
WHERE customer_id = $1;

-- name: ViewOrderItems :many
SELECT id, order_id, product_id, quantity, subtotal
FROM order_items
WHERE order_id = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetProductsSorted :many
SELECT * FROM products
WHERE
  LOWER(name) LIKE LOWER('%' || COALESCE($1, '') || '%')
  OR LOWER(description) LIKE LOWER('%' || COALESCE($1, '') || '%')
ORDER BY
  CASE WHEN $2 = 'name-asc' THEN name END ASC,
  CASE WHEN $2 = 'name-desc' THEN name END DESC,
  CASE WHEN $2 = 'price-asc' THEN price::numeric END ASC,
  CASE WHEN $2 = 'price-desc' THEN price::numeric END DESC,
  CASE WHEN $2 = 'stock-asc' THEN stock_quantity END ASC,
  CASE WHEN $2 = 'stock-desc' THEN stock_quantity END DESC,
  id ASC
LIMIT $3 OFFSET $4;

-- name: GetProductsCount :one
SELECT COUNT(*) FROM products
WHERE LOWER(name) LIKE LOWER('%' || COALESCE($1, '') || '%')
  OR LOWER(description) LIKE LOWER('%' || COALESCE($1, '') || '%');

-- name: UpdateCartItem :one
UPDATE cart
SET quantity = $3, modified_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND product_id = $2
RETURNING *;

-- name: GetAllPendingOrders :many
SELECT * FROM orders WHERE LOWER(status) = 'pending';

-- name: GetOrderById :one
SELECT id, customer_id, order_date, total_amount, status FROM orders WHERE id = $1;
