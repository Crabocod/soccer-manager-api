# Soccer Manager API

RESTful API for managing soccer teams and player transfers.

## Technologies

- Go 1.24.6
- PostgreSQL 16
- Redis 7
- Gin Web Framework
- JWT Authentication
- Uber FX
- Circuit Breaker (gobreaker)

## Features

- User registration and authentication
- Automatic team creation with 20 players upon registration
- Team and player management
- Transfer market (buying/selling players)
- Dynamic player value changes
- Localization (EN/KA)
- Redis caching
- Rate limiting for logins

## Localization

The API supports English and Georgian languages via the `Accept-Language` header:

```bash
curl -H "Accept-Language: en" http://localhost:8080/api/v1/...
curl -H "Accept-Language: ka" http://localhost:8080/api/v1/...
```

## Quick Start

### Requirements

- Docker and Docker Compose

### Running

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Start the project:

```bash
docker-compose up --build
```

API is available at `http://localhost:8080`

## Make Commands

```bash
make help   # Show available commands
make test   # Run tests with race detector and coverage
make up     # Start Docker containers
make clean  # Clean build artifacts
```

## Project Structure

```
.
├── cmd/
│   └── soccer_manager_service/     # Entry point
├── internal/
│   ├── api/rest/                   # REST handlers & middleware
│   ├── bootstrap/                  # Application initialization (DI, DB, Redis, JWT)
│   ├── config/                     # Configuration
│   ├── dto/                        # Data Transfer Objects for API
│   ├── entity/                     # Domain models
│   ├── ports/                      # Repository interfaces
│   ├── repository/                 # PostgreSQL and Redis repositories
│   └── usecase/                    # Business logic
├── migrations/                     # SQL migrations
├── pkg/
│   ├── errors/                     # Custom errors
│   ├── i18n/                       # Localization (en.json, ka.json)
│   └── jwt/                        # JWT utilities
├── postman/                        # Postman collection
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## API Documentation

Postman collection is available at `postman/Soccer_Manager_API.postman_collection.json`

Main endpoints:
- `POST /api/v1/auth/register` - Registration
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/team` - Get your team
- `PATCH /api/v1/team` - Update team
- `PATCH /api/v1/players/:id` - Update player
- `POST /api/v1/players/:id/transfer` - List for transfer
- `GET /api/v1/transfers` - List transfers
- `POST /api/v1/transfers/:id/buy` - Buy player

## Testing

```bash
make test
```

Running a specific test:

```bash
go test -v ./internal/usecase -run TestAuthService
```