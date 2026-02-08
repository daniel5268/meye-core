# Meye Core

> A comprehensive RPG campaign management system built with Go, Clean Architecture, and Domain-Driven Design.

[![Go Version](https://img.shields.io/badge/Go-1.24.0-00ADD8?style=flat&logo=go)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blueviolet?style=flat)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
[![API](https://img.shields.io/badge/API-OpenAPI%203.0-85EA2D?style=flat&logo=swagger)](./openapi.yml)

## Overview

**Meye Core** is a backend API for managing tabletop RPG campaigns, player characters, and game sessions. It provides a robust system for:

- **User Management**: Admin, Master (GM), and Player roles with JWT authentication
- **Campaign Management**: Create campaigns, invite players, track sessions
- **Character System**: Complex character creation with 40+ stats, talents, and progression
- **XP & Progression**: Sophisticated experience point system with talent-based cost multipliers
- **Session Tracking**: Record game sessions with automatic XP distribution
- **Event-Driven**: RabbitMQ integration for asynchronous processing and event sourcing

### Key Features

✅ **Clean Architecture** - Hexagonal architecture with clear separation of concerns
✅ **Domain-Driven Design** - Rich domain models with business logic encapsulation
✅ **REST API** - Fully documented with OpenAPI 3.0 specification
✅ **JWT Authentication** - Secure token-based authentication with role-based access control
✅ **Event-Driven** - Domain events published to RabbitMQ for async processing
✅ **Event Sourcing** - All domain events persisted for audit trail and potential replay
✅ **PostgreSQL** - Relational database with GORM ORM
✅ **Docker Support** - Docker Compose for easy local development

## Quick Start

### Prerequisites

- **Go 1.24.0+** ([Download](https://golang.org/dl/))
- **Docker & Docker Compose** ([Install](https://docs.docker.com/get-docker/))
- **PostgreSQL 16+** (or use Docker Compose)
- **RabbitMQ 3.12+** (or use Docker Compose)

### Installation

1. **Clone the repository**
   ```bash
   git clone git@github.com:daniel5268/meye-core.git
   cd meye-core
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start infrastructure services**
   ```bash
   docker compose up -d
   ```
   This starts PostgreSQL and RabbitMQ using .env configuration

5. **Run database migrations**
   ```bash
   # Migrations run automatically on server start
   ```

6. **Run the server**
   ```bash
   GO_ENV=development go run ./cmd/server/...
   ```

7. **Run the worker**
   ```bash
   GO_ENV=development go run ./cmd/worker/...
   ```

### Quick Example: Create User and Login

```bash
# 1. Create admin user (requires bootstrapping or existing admin)
curl -X POST http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gamemaster",
    "password": "SecureP@ssw0rd",
    "role": "master"
  }'

# 2. Login
curl -X POST http://localhost:3000/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gamemaster",
    "password": "SecureP@ssw0rd"
  }'

# Response: {"token":"eyJhbGciOiJIUzI1NiIs..."}
```

## Documentation

### 📘 API Documentation

**[OpenAPI 3.0 Specification](./openapi.yml)** - Complete API reference with schemas, examples, and validation rules

**View the API documentation:**
- **Swagger UI**: Paste `openapi.yml` into [Swagger Editor](https://editor.swagger.io/)
- **Postman**: Import `openapi.yml` (File → Import)
- **ReDoc**: Use [ReDoc CLI](https://github.com/Redocly/redoc) for beautiful documentation

**Key endpoints:**
- `POST /api/v1/users/login` - User authentication
- `POST /api/v1/campaigns` - Create campaign (Master role)
- `POST /api/v1/campaigns/{id}/invitations` - Invite players
- `POST /api/v1/campaigns/{id}/pjs` - Create player character
- `POST /api/v1/campaigns/{id}/sessions` - Record game session
- `GET /api/v1/pjs/{id}` - Get character details
- `PUT /api/v1/pjs/{id}/stats` - Update character stats (spend XP)

### 📖 Project Documentation

**[Project Context](./.claude/project_context.md)** - Comprehensive architecture and development guide

This document covers:
- **Architecture**: Hexagonal architecture, DDD patterns, layer responsibilities
- **Domain Model**: Users, Campaigns, Characters (PJs), Sessions, Events
- **Use Cases**: Complete list of business workflows
- **Patterns**: Aggregate roots, repositories, value objects, event sourcing
- **Development Guide**: How to add features, endpoints, use cases
- **Testing**: Unit tests, integration tests, mocks

### 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                   │
│  (HTTP Handlers, PostgreSQL, RabbitMQ, JWT, Bcrypt)    │
└────────────────────┬────────────────────────────────────┘
                     │ Implements Ports
┌────────────────────▼────────────────────────────────────┐
│                   Application Layer                      │
│        (Use Cases: CreateCampaign, CreatePJ, etc.)      │
└────────────────────┬────────────────────────────────────┘
                     │ Orchestrates
┌────────────────────▼────────────────────────────────────┐
│                     Domain Layer                         │
│   (Entities: User, Campaign, PJ, Session + Events)      │
│           (Business Rules & Domain Logic)                │
└──────────────────────────────────────────────────────────┘
```

**Key Principle**: Dependencies point inward. Infrastructure depends on Application and Domain, but never the reverse.

## Project Structure

```
meye-core/
├── cmd/
│   ├── server/              # HTTP Server entry point & dependency injection
│   │   ├── main.go          # Entry point with graceful shutdown
│   │   ├── server.go        # HTTP server setup and migrations
│   │   └── dependencies.go  # Dependency injection container
│   │
│   └── worker/              # Background Worker entry point & dependency injection
│       ├── main.go          # Worker entry point (RabbitMQ consumer)
│       └── dependencies.go  # Worker dependency injection container
│
├── internal/
│   ├── application/         # Use Cases (Business Logic)
│   │   ├── campaign/        # Campaign-related use cases
│   │   ├── session/         # Session management use cases
│   │   └── user/           # User-related use cases
│   │
│   ├── domain/              # Domain Layer (Core Business Rules)
│   │   ├── campaign/        # Campaign aggregate & entities
│   │   ├── event/           # Domain events
│   │   ├── session/         # Session aggregate
│   │   ├── shared/          # Shared domain services
│   │   └── user/           # User aggregate
│   │
│   ├── infrastructure/      # Infrastructure Layer (External Integrations)
│   │   ├── api/             # HTTP API layer (Gin)
│   │   │   ├── handler/     # HTTP request handlers
│   │   │   │   └── dto/     # Data Transfer Objects
│   │   │   └── validator/   # Custom validation rules
│   │   ├── repository/      # Data persistence (GORM)
│   │   ├── hash/           # Password hashing service
│   │   ├── identification/ # ID generation service
│   │   ├── jwt/            # JWT authentication service
│   │   ├── messaging/      # RabbitMQ integration
│   │   └── worker/         # Background job processing
│   │
│   └── config/             # Configuration management
│
├── migrations/             # Database migrations
├── tests/                  # Integration and unit tests
└── openapi.yml            # OpenAPI v3 API Documentation
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate mocks (requires gomock)
go generate ./...
```

### Building

```bash
# Build server
go build -o server cmd/server/main.go

# Build worker
go build -o worker cmd/worker/main.go

# Build both
go build ./...

# Build with optimizations for production
go build -ldflags="-s -w" -o server cmd/server/main.go
```

### Database Migrations

Migrations are located in `migrations/` and run automatically on server start.


### Docker Compose

```bash
# Start all services
docker compose up -d
```

**Services:**
- **PostgreSQL**: Port 5432
- **RabbitMQ**: Port 5672 (AMQP), Port 15672 (Management UI)

**Access RabbitMQ Management**: http://localhost:15672 (guest/guest)

## Configuration

Configuration is managed via environment variables. See `.env.example` for all options.

### Key Configuration

```bash
# API
API_PORT=3000                    # Server port
API_KEY=supersecretapikey        # Internal API key

# JWT
JWT_SECRET=your-secret-key       # JWT signing secret
JWT_ISSUER=meye-core            # Token issuer
JWT_EXPIRATION_TIME=1h          # Token lifetime

# Database
DATABASE_DSN=postgresql://user:pass@localhost:5432/meye_core?sslmode=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EVENTS_QUEUE=domain_events
```

## Technology Stack

| Component | Technology | Version |
|-----------|------------|---------|
| **Language** | Go | 1.24.0 |
| **Web Framework** | Gin | 1.11.0 |
| **Database** | PostgreSQL | 16+ |
| **ORM** | GORM | Latest |
| **Message Queue** | RabbitMQ | 3.12+ |
| **Authentication** | JWT (golang-jwt) | v5 |
| **Password Hashing** | bcrypt | Latest |
| **Validation** | go-playground/validator | v10 |
| **Logging** | Logrus | 1.9.4 |
| **Migrations** | golang-migrate | v4 |
| **Testing** | testify + gomock | Latest |
| **Container** | Docker + Compose | Latest |

## Domain Model

### User Roles

- **admin**: System administrator, can create users
- **master**: Campaign master (Game Master/GM), can create campaigns and sessions
- **player**: Regular player, can create and manage their own characters

### Core Entities

#### Campaign Aggregate
- **Campaign**: Root entity representing an RPG campaign
- **Invitation**: Child entity for player invitations
- **PJ (Player Character)**: Child entity with complex stat system

#### Session Aggregate
- **Session**: Represents a game session with XP assignments
- **Events**: Domain events for async processing

#### Character Stats
- **Basic Stats**: Physical, Mental, Coordination + Life
- **Special Stats**: Skills (Physical, Mental, Energy) + Energy Tank
- **Supernatural Stats**: Transformation abilities (optional)
- **Talents**: 6 independent talent flags affecting XP costs
- **XP**: Three categories (Basic, Special, Supernatural)

### Workflow Example

1. **Master creates campaign** → Campaign entity created
2. **Master invites players** → Invitation entities created
3. **Players create characters (PJs)** → PJ entities created, invitations accepted
4. **Master records session** → Session created, XP assigned via events
5. **Players spend XP** → Character stats updated

## Event-Driven Architecture

The system publishes domain events to RabbitMQ for asynchronous processing:

**Published Events:**
- `UserCreated` - New user registered
- `CampaignCreated` - New campaign created
- `UserInvited` - Player invited to campaign
- `PJCreated` - Character created
- `SessionCreated` - Game session recorded
- `XPAssigned` - XP awarded to character
- `StatsUpdated` - Character stats modified

**Event Store:**
All events are persisted to the `domain_events` table for:
- Audit trail
- Event replay capabilities
- Analytics and reporting
- Eventual consistency

**Worker Process:**
The worker consumes events from RabbitMQ and can:
- Update read models
- Send notifications
- Trigger workflows
- Integrate with external systems

## Security

### Authentication
- JWT Bearer tokens (obtained via `/api/v1/users/login`)
- Tokens expire after configurable duration (default: 1h)
- Include in requests: `Authorization: Bearer <token>`

### Password Security
- Bcrypt hashing with cost factor 12
- Minimum 8 characters, maximum 72 characters
- Validated on entity creation

### Authorization
- Role-based access control (RBAC)
- Middleware enforces role requirements per endpoint
- Resource ownership validation (e.g., user can only update own character)

### Best Practices
- Never commit `.env` to version control
- Rotate JWT secrets regularly in production
- Use strong passwords for database and RabbitMQ
- Enable SSL/TLS in production
- Keep dependencies updated

## Testing

The project includes comprehensive tests at multiple levels:

### Domain Tests
```bash
go test ./internal/domain/campaign/... -v
```
- XP calculation algorithms
- Stat progression logic
- Business rule validation

## Contributing

### Adding a New Endpoint

1. **Define use case** in `internal/application/[domain]/ports.go`
2. **Implement use case** in `internal/application/[domain]/[usecase]/`
3. **Add handler method** in `internal/infrastructure/api/handler/`
4. **Define DTOs** in `internal/infrastructure/api/handler/dto/`
5. **Register route** in `internal/infrastructure/api/api.go`
6. **Update OpenAPI spec** in `openapi.yml`
7. **Write tests**

See [Project Context](./.claude/project_context.md) for detailed patterns and examples.

### Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Write meaningful commit messages
- Add tests for new features
- Update documentation

## Support

For questions or issues:
- Review the [OpenAPI documentation](./openapi.yml)
- Contact the development team

---

**Built with** ❤️ **using Go, Clean Architecture, and Domain-Driven Design**