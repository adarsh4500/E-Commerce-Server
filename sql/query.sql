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