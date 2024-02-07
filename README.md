# E-Commerce-Server

## Overview

A robust and secure backend service designed for managing user authentication, product catalog, shopping cart, and order processing in an e-commerce application. This API is built using the Go programming language and the Gin framework, providing a scalable and performant solution for basic e-commerce backend needs.

## Features

- **User Authentication:** Secure user authentication with JWT token-based authorization.
- **Product Management:** CRUD operations for managing product information.
- **Shopping Cart:** Ability to add, remove, view, and clear items in the user's shopping cart.
- **Order Processing:** Place orders, view order history, update order status, and view order items.

## Getting Started

### Prerequisites

- Go programming language
- PostgreSQL database
- [Gin](https://github.com/gin-gonic/gin) framework

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/adarsh4500/E-Commerce-Server.git

2. Navigate to the project directory:

    ```bash
    cd E-Commerce-Server

3. Create a .env file in the root directory and update it with your configuration using the template below.

    ```.env
    JWT_SECRET_KEY=<your-secret-key-here>
    DB_USERNAME=<username>
    DB_PASSWORD=<password>
    DB_NAME=<database>

4. Install dependencies:

    ```bash
    go mod tidy

5. Run the main.go file.

   ```bash
    go run main.go

6. The Ecom API will be accessible at http://localhost:8080 by default. Explore the API documentation using Swagger at http://localhost:8080/docs/index.html after running the application.
