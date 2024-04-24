# E-Commerce Server API

A production-ready, secure backend service for e-commerce applications built with **Go**, **Gin**, and **PostgreSQL**. This API provides comprehensive user authentication, product management, shopping cart functionality, and order processing with proper validation, error handling, and security measures.

## 🚀 Features

### 🔐 Authentication & Security
- **JWT-based authentication** with secure cookie handling
- **Password hashing** using bcrypt
- **Input validation** for email format and password strength
- **Thread-safe user context** (no global state)
- **Secure cookie flags** (HttpOnly, Secure)
- **Role-Based Access Control (RBAC)** with user and admin roles

### 🛍️ E-Commerce Functionality
- **User Management**: Registration, login, logout
- **Product Catalog**: CRUD operations with validation
- **Shopping Cart**: Add, remove, view, clear items with stock validation
- **Order Processing**: Place orders with atomic transactions, view history, update status
- **Stock Management**: Real-time stock checking and updates

### 🛡️ Production Features
- **Database transactions** for data consistency
- **Comprehensive error handling** with proper HTTP status codes
- **Input validation** and sanitization
- **Authorization checks** (users can only access their own data)
- **Structured logging** and response formatting
- **Role-based route protection** with middleware

## 🔐 Role-Based Access Control (RBAC)

### User Roles
- **User** (default): Can view products, manage cart, place orders, view their own orders
- **Admin**: Can perform all user actions plus product management and order status updates

### Role Assignment
- New users are automatically assigned the "user" role
- Admin users must be manually set in the database
- Roles are included in JWT tokens for authorization

### Protected Endpoints

#### Admin-Only Endpoints
- `POST /products/new` - Create new products
- `POST /products/update/:id` - Update existing products  
- `POST /products/delete/:id` - Delete products
- `POST /orders/updatestatus` - Update order status

#### User-Accessible Endpoints
- `GET /products` - View all products
- `GET /products/:id` - View specific product
- All cart endpoints (`/cart/*`)
- `GET /orders` - View user's orders
- `GET /orders/:id` - View specific order details

### Setting Up Admin Users
To create an admin user, manually update the database:
```sql
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
```

## 📋 Prerequisites

- **Go 1.21+**
- **PostgreSQL 12+**
- **Docker** (optional, for sqlc code generation)

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/adarsh4500/E-Commerce-Server.git
cd E-Commerce-Server
```

### 2. Environment Configuration
Create a `.env` file in the root directory:
```env
JWT_SECRET_KEY=your-super-secret-jwt-key-here
DB_USERNAME=your_db_username
DB_PASSWORD=your_db_password
DB_NAME=your_database_name
```

### 3. Database Setup
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Run the schema.sql file to create tables
-- (Located in sql/schema.sql)
```

### 4. Generate Database Code
```bash
# Using Docker (recommended for Windows)
docker run --rm -v /path/to/your/project:/src -w /src kjconroy/sqlc generate

# Or using WSL/Linux
sqlc generate
```

### 5. Install Dependencies
```bash
go mod tidy
```

### 6. Run the Application
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /signup
Content-Type: application/json

{
  "fullname": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 200,
  "message": "Signed Up Successfully"
}
```

**Note:** New users are automatically assigned the "user" role.

#### Login
```http
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 200,
  "message": "Successfully Logged in"
}
```

**Note:** JWT token includes the user's role for authorization.

#### Logout
```http
POST /logout
```

### Product Endpoints

#### Get All Products (User/Admin)
```http
GET /products
Authorization: Bearer <jwt-token>
```

#### Get Product by ID (User/Admin)
```http
GET /products/{id}
Authorization: Bearer <jwt-token>
```

#### Create Product (Admin Only)
```http
POST /products/new
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "Product Name",
  "price": 29.99,
  "description": "Product description",
  "stock_quantity": 100
}
```

**Error Response (Non-Admin):**
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 403,
  "message": "admin access required"
}
```

#### Update Product (Admin Only)
```http
POST /products/update/{id}
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "Updated Product Name",
  "price": 34.99
}
```

#### Delete Product (Admin Only)
```http
POST /products/delete/{id}
Authorization: Bearer <jwt-token>
```

### Cart Endpoints (User/Admin)

#### Add to Cart
```http
POST /cart/new
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "product_id": "uuid",
  "quantity": 2
}
```

#### View Cart
```http
GET /cart
Authorization: Bearer <jwt-token>
```

#### Remove from Cart
```http
POST /cart/remove/{product_id}
Authorization: Bearer <jwt-token>
```

#### Clear Cart
```http
POST /cart/clear
Authorization: Bearer <jwt-token>
```

#### Place Order
```http
POST /cart/checkout
Authorization: Bearer <jwt-token>
```

### Order Endpoints

#### View Orders (User/Admin - Own Orders Only)
```http
GET /orders
Authorization: Bearer <jwt-token>
```

#### View Order Items (User/Admin - Own Orders Only)
```http
GET /orders/{order_id}
Authorization: Bearer <jwt-token>
```

#### Update Order Status (Admin Only)
```http
POST /orders/updatestatus
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "id": "order-uuid",
  "status": "Shipped"
}
```

## 🔒 Security Features

### Authentication
- **JWT tokens** stored in secure, HttpOnly cookies
- **Password hashing** with bcrypt (cost factor 14)
- **Token expiration** (1 hour by default)
- **Input validation** for email and password
- **Role-based JWT claims** for authorization

### Authorization
- **Role-based access control** with user and admin roles
- **Middleware-based route protection** for admin-only endpoints
- **User context isolation** - users can only access their own data
- **Order ownership validation** - users can only view/update their own orders
- **Cart isolation** - each user has their own cart

### Data Protection
- **SQL injection prevention** through parameterized queries
- **Input sanitization** and validation
- **Error message sanitization** (no internal details exposed)

## 🏗️ Architecture

### Project Structure
```
E-Commerce-Server/
├── config/          # Configuration management
├── connections/     # Database connection setup
├── controllers/     # HTTP handlers (including RBAC middleware)
├── helpers/         # Utility functions (hashing)
├── models/          # Data models and validation (including role constants)
├── postgres/        # Generated database code
├── routes/          # Route definitions with role-based protection
├── sql/            # SQL schema and queries
├── utils/          # Response utilities
└── main.go         # Application entry point
```

### Database Schema
- **users**: User accounts, authentication, and roles
- **products**: Product catalog with pricing
- **cart**: Shopping cart items
- **orders**: Order headers
- **order_items**: Order line items

### Key Technologies
- **Gin**: HTTP framework
- **PostgreSQL**: Database
- **sqlc**: Type-safe SQL code generation
- **JWT**: Authentication with role claims
- **bcrypt**: Password hashing
- **RBAC**: Role-based access control middleware

## 🚀 Production Deployment

### Environment Variables
```env
JWT_SECRET_KEY=your-production-secret-key
DB_USERNAME=production_db_user
DB_PASSWORD=production_db_password
DB_NAME=production_db_name
```

### Security Checklist
- [ ] Use strong JWT secret key
- [ ] Enable HTTPS in production
- [ ] Configure proper CORS settings
- [ ] Set up database connection pooling
- [ ] Implement rate limiting
- [ ] Add monitoring and logging
- [ ] Set up backup strategy
- [ ] Configure admin user accounts
- [ ] Review role assignments regularly

### Performance Considerations
- **Database indexing** on frequently queried columns
- **Connection pooling** for database connections
- **Caching** for product catalog (Redis recommended)
- **CDN** for static assets

## 🧪 Testing

### Using Postman
1. Import the provided `E-Commerce-API.postman_collection.json`
2. Set the `base_url` variable to your server URL
3. Test the authentication flow first
4. Test both user and admin role scenarios
5. Use the collection to test all endpoints

### Manual Testing Flow
1. **Register** a new user (gets "user" role by default)
2. **Login** to get authentication token
3. **Test user permissions** (view products, manage cart, place orders)
4. **Create admin user** in database
5. **Login as admin** to test admin-only endpoints
6. **Test admin permissions** (create/update/delete products, update order status)

### Role Testing
- **User Role**: Should be able to access all non-admin endpoints
- **Admin Role**: Should be able to access all endpoints including admin-only ones
- **Unauthorized Access**: Should receive 403 Forbidden for admin-only endpoints

## 📝 API Response Format

### Success Response
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 200,
  "data": { ... },        // For data responses
  "message": "Success"     // For message-only responses
}
```

### Error Response
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 400,
  "message": "Error description"
}
```

### Authorization Error Response
```json
{
  "timestamp": "2024-01-01T00:00:00Z",
  "status": 403,
  "message": "admin access required"
}
```

---

**Built with ❤️ using Go and Gin**
