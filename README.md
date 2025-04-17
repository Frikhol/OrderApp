# MyOrder

MyOrder is a web application for managing orders, built with Go and PostgreSQL.

## Features

- User Authentication
  - Secure registration with email validation
  - Login with session management
  - Protected routes
  - Logout functionality

- User Management
  - Email-based registration
  - Password hashing
  - Unique email validation

## Prerequisites

- Docker
- Docker Compose
- Go 1.24 or later (for local development)

## Installation

Start the application using Docker Compose:
```bash
docker-compose up --build
```

The application will be available at http://localhost:8080

## Project Structure

```
myorder/
├── cmd/
│   └── web/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── user.go
│   └── handler/
│       └── handler.go
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── scripts/
│   └── init-db.sh
├── templates/
│   ├── index.html
│   ├── landing.html
│   └── register.html
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Development

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL:
```bash
docker-compose up -d postgres
```

3. Run the application:
```bash
go run cmd/web/main.go
```

### Database Migrations

Database migrations are automatically applied when starting the application with Docker Compose. The initialization script checks if the users table exists and creates it if necessary.

## API Endpoints

- `GET /` - Landing page with login form
- `POST /login` - Handle login
- `GET /register` - Show registration form
- `POST /register` - Handle registration
- `GET /logout` - Handle logout
- `GET /index.html` - Protected dashboard page

## Security Features

- Password hashing using bcrypt
- Session management with secure cookies
- Protected routes
- Input validation
- SQL injection prevention

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request