CREATE TABLE "users" (
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY
);

CREATE TABLE "products" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "stock_quantity" INT NOT NULL
);
CREATE TABLE "cart" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "user_id" UUID REFERENCES "users"("id") NOT NULL,
    "product_id" UUID REFERENCES "products"("id") NOT NULL,
    "quantity" INT NOT NULL,
    "modified_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE("user_id", "product_id")
);
CREATE TABLE "orders" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "customer_id" UUID REFERENCES "users"("id") NOT NULL,
    "order_date" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "total_amount" DECIMAL(10, 2) NOT NULL,
    "status" VARCHAR(50) DEFAULT 'Pending' NOT NULL
);
CREATE TABLE "order_items" (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "order_id" UUID REFERENCES "orders"("id") NOT NULL,
    "product_id" UUID REFERENCES "products"("id") NOT NULL,
    "quantity" INT NOT NULL,
    "subtotal" DECIMAL(10, 2) NOT NULL
);
