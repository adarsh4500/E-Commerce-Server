-- name: AddUser :one
INSERT INTO users (fullname, email, password)
VALUES ($1, $2, $3) RETURNING *;

-- name: ValidateCreds :one
SELECT CASE
         WHEN EXISTS (SELECT 1 FROM users u WHERE u.email = $1 AND u.password = $2) THEN 1
         WHEN EXISTS (SELECT 1 FROM users v WHERE v.email = $1) THEN 0
         ELSE 2
       END AS Result;

-- name: AddProduct :one
INSERT INTO product (name, price, description, stock_quantity) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetProducts :many
SELECT * FROM product;

-- name: GetProductById :one
SELECT * FROM product WHERE id = $1;

-- name: DeleteProductById :one
DELETE FROM Product 
WHERE id=$1 RETURNING *;

-- name: UpdateProductById :one
UPDATE Product
SET 
    name = COALESCE($2, name),
    price = COALESCE($3, price),
    description = COALESCE($4, description),
    stock_quantity = COALESCE($5, stock_quantity)
WHERE
    id = $1 RETURNING *;