-- name: AddUser :exec
INSERT INTO users (fullname, email, password)
VALUES ($1, $2, $3);

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
SELECT * FROM cart WHERE user_id = $1;

-- name: AddToCart :one
INSERT INTO cart (user_id, product_id, quantity)
VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveFromCart :one
DELETE FROM cart
WHERE user_id = $1 AND product_id = $2 RETURNING *;

-- name: ClearCart :exec
DELETE FROM cart
WHERE user_id = $1;

-- name: AddOrder :one
INSERT INTO orders (customer_id, total_amount)
VALUES ($1, $2) RETURNING id;

-- name: UpdateOrderTotal :one
UPDATE orders SET total_amount = $2
WHERE id = $1 RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $1
WHERE id = $2 RETURNING *;

-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, subtotal)
VALUES ($1, $2, $3, (SELECT price * CAST($5 AS DECIMAL(10, 2)) FROM products WHERE products.id = $4))
RETURNING subtotal;

-- name: ViewOrders :many
SELECT id, order_date, total_amount, status
FROM orders
WHERE customer_id = $1;

-- name: ViewOrderItems :many
SELECT id, order_id, product_id, quantity, subtotal
FROM order_items
WHERE order_id = $1;