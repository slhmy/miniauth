# MiniAuth

A simple authentication service with user and organization management.

## Features

- User authentication and authorization
- Organization management
- Role-based access control (Admin/User)
- RESTful API with Swagger documentation
- Modern React frontend

## Quick Start

### Backend

1. Copy the environment configuration:

   ```bash
   cp .env.example .env
   ```

2. Build and run the backend:

   ```bash
   make build
   ./bin/miniauth
   ```

The server will start on `http://localhost:8080`.

### Frontend

1. Navigate to the frontend directory:

   ```bash
   cd website
   ```

2. Install dependencies and start the development server:

   ```bash
   pnpm install
   pnpm dev
   ```

The frontend will be available at `http://localhost:5173`.

## Default Admin User

When the application starts for the first time, it automatically creates a default admin user:

- **Username**: `admin`
- **Email**: `admin@example.com`
- **Password**: `admin123`

### Customizing Default Admin

You can customize the default admin user by setting environment variables:

```bash
DEFAULT_ADMIN_USERNAME=myadmin
DEFAULT_ADMIN_EMAIL=admin@mycompany.com
DEFAULT_ADMIN_PASSWORD=mySecurePassword123
```

### Disabling Default Admin Creation

To disable automatic creation of the default admin user:

```bash
DISABLE_DEFAULT_ADMIN=true
```

**Security Note**: Please change the default admin password immediately after first login!

## API Documentation

The API documentation is available at `http://localhost:8080/swagger/index.html` when the server is running.

## Build Commands

- `make build` - Build the backend binary
- `make swag` - Generate Swagger documentation and frontend API client

## Database

By default, the application uses SQLite with a local `.db` file. You can configure other databases (PostgreSQL, MySQL) via environment variables. See `.env.example` for details.
