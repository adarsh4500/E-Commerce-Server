CREATE TABLE "users" (
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY
);

CREATE TABLE "product" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "stock_quantity" INT NOT NULL
);
