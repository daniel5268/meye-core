# Meye Core - Project Context

## Overview

**Meye Core** is a comprehensive RPG (Role-Playing Game) campaign management system built with Go and Clean Architecture principles. The system enables game masters to create campaigns, invite players, manage game sessions, and track character progression through a sophisticated stat and XP system.

## Technology Stack

- **Language**: Go 1.24.0
- **Web Framework**: Gin v1.11.0
- **Database**: PostgreSQL with GORM ORM
- **Message Queue**: RabbitMQ (amqp091-go) for domain events
- **Authentication**: JWT (golang-jwt v5)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator v10
- **Logging**: Logrus v1.9.4
- **Migrations**: golang-migrate v4
- **Testing**: testify v1.11.1, go.uber.org/mock

## Architecture

### Clean Architecture Layers

The project follows Clean Architecture with clear separation of concerns:

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

### Key Architectural Patterns

1. **Domain-Driven Design (DDD)**
   - Aggregates: User, Campaign, PJ (Player Character), Session
   - Domain Events: Published for all significant business operations
   - Value Objects: Stats, XP, Physical/Mental/Coordination attributes

2. **Dependency Injection**
   - Centralized in `cmd/server/dependencies.go` for the API and `cmd/worker/dependencies.go` for the worker
   - Interface-based design for testability
   - Clear dependency flow: Handlers → Use Cases → Domain → Repositories

3. **Event-Driven Architecture**
   - Domain events published to RabbitMQ
   - Worker process consumes events asynchronously
   - Events: UserCreated, CampaignCreated, PJCreated, XPConsumed, StatsUpdated, etc.

4. **CQRS Patterns**
   - Use Cases implement command and query operations
   - Clear input/output DTOs for each use case

## API Documentation

### OpenAPI v3 Specification

**Location**: `openapi.yml` (root directory)

The API is fully documented using OpenAPI v3.0.3 specification. This file contains:
- All endpoint definitions with detailed descriptions
- Request/response schemas with validation rules
- Authentication requirements (JWT Bearer tokens)
- Authorization rules (role-based access control)
- Complete examples for all operations
- Error codes and status codes
- Business logic explanations

### Using the OpenAPI Documentation

**View in Swagger UI**:
   - Visit https://editor.swagger.io/
   - Copy/paste the content of `openapi.yml`
   - Interactive API exploration with "Try it out" feature


### API Endpoints Summary

**Base URL**: `http://localhost:3000`

#### Health Check
- `GET /health` - Service health status

#### User Management
- `POST /api/v1/users/login` - Authenticate user → Returns JWT token
- `POST /api/v1/users` - Create user (Admin only)

#### Campaign Management
- `POST /api/v1/campaigns` - Create campaign (Master role)
- `GET /api/v1/campaigns/{campaignID}` - Get campaign details (Master only)
- `GET /api/v1/campaigns` - Get campaigns basic information details (Master only)
- `POST /api/v1/campaigns/{campaignID}/invitations` - Invite user (Master only)
- `POST /api/v1/campaigns/{campaignID}/pjs` - Create player character (Player role)
- `POST /api/v1/campaigns/{campaignID}/sessions` - Create session (Master only)

#### Player Character Management
- `GET /api/v1/pjs/{pjID}` - Get character details (Owner only)
- `PUT /api/v1/pjs/{pjID}/stats` - Update character stats (Owner only)

## Domain Model

### User Roles

- **admin**: System administrator, can create users
- **master**: Campaign master, can create campaigns and manage sessions
- **player**: Regular player, can create and manage their own characters

### Character Stats System

The RPG system uses a three-tier stat structure:

#### 1. Basic Stats
- **Physical**: strength, agility, speed, resistance
- **Mental**: intelligence, wisdom, concentration, will
- **Coordination**: precision, calculation, range, reflexes
- **Life**: Hit points

#### 2. Special Stats (Skills)
- **Physical Skills**: empowerment, vital_control
- **Mental Skills**: illusion, mental_control
- **Energy Skills**: object_handling, energy_handling
- **Energy Tank**: Total energy/mana capacity

#### 3. Supernatural Stats (Supernatural Characters Only)
- **Skills**: Array of skills with transformation levels
- Each skill has multiple transformation tiers with power levels

### Talent System

Each stat category can have a "talent" flag (`is_*_talented`):
- Provides bonuses in that category
- Reduces XP costs for improvements
- Set during character creation, permanent

### XP System

Characters earn and spend XP in three categories:

1. **Basic XP**: For improving basic stats (physical, mental, coordination)
2. **Special XP**: For improving special skills
3. **Supernatural XP**: For improving supernatural abilities

**XP Flow**:
1. Master creates session with XP assignments
2. XP added to character's available pools
3. Player spends XP to increase stats (via PUT /api/v1/pjs/{pjID}/stats)
4. XP automatically deducted based on stat increases

**Business Rules**:
- Stats can only increase, never decrease
- Higher stats cost more XP to increase further
- Must have sufficient XP before updating
- Atomic transactions - all or nothing updates

### Character Types

- **human**: Normal human with basic and special stats
- **supernatural**: Supernatural being with basic, special, AND supernatural stats

## Authentication & Authorization

### JWT Authentication

**Flow**:
1. User logs in via `POST /api/v1/users/login`
2. Receives JWT token in response
3. Includes token in subsequent requests: `Authorization: Bearer <token>`

**Token Claims**:
- User ID
- Username
- Role
- Issuer: `meye-core`
- Expiration: Configurable (default 1h)

### Role-Based Access Control (RBAC)

**Middleware Chain**:
1. `AuthMiddleware()` - Validates JWT, extracts user context
2. Role middleware - Checks user role:
   - `RequireAdminRole()` - Admin only
   - `RequireMasterRole()` - Master only
   - `RequirePlayerRole()` - Player only
   - `RequireCampaignMaster()` - Must be master of specific campaign
   - `RequirePjUser()` - Must be owner of specific character

## Configuration

### Environment Variables

See `.env.example` for complete list:

```bash
# API Configuration
API_PORT=3000
API_KEY=supersecretapikey  # For internal services

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_ISSUER=meye-core
JWT_EXPIRATION_TIME=1h

# Database
DATABASE_DSN=postgresql://user:pass@localhost:5432/meye_core?sslmode=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EVENTS_QUEUE=domain_events
```

### Docker Compose

Start dependencies:
```bash
docker-compose up -d  # Starts PostgreSQL and RabbitMQ
```

## Development Workflow

### Setup

1. Install Go 1.24.0
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Start dependencies:
   ```bash
   docker-compose up -d
   ```
4. Copy environment file:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

### Running the Server

```bash
GO_ENV=development go run ./cmd/server/...

# Variable GO_ENV=development should be passed in order to read from .env

```

### Running the Worker

```bash
GO_ENV=development go run ./cmd/worker/...

```

### Running Tests

```bash
# All tests
go test ./...
```

## Database Migrations

Located in `migrations/` directory.

**Running migrations**:
- Automatic on server startup

## Domain Events

All significant business operations publish domain events to RabbitMQ:

- `UserCreated` - New user registered
- `CampaignCreated` - New campaign created
- `UserInvited` - User invited to campaign
- `PJCreated` - Player character created
- `SessionCreated` - Game session recorded
- `XPConsumed` - XP awarded to character
- `StatsUpdated` - Character stats modified

**Event Structure**:
```go
type DomainEvent interface {
    EventType() string
    AggregateID() string
    OccurredOn() time.Time
    EventData() map[string]interface{}
}
```

**Consumer**: Worker process in `cmd/worker/main.go`

## Validation Rules

Custom validators defined in `internal/infrastructure/api/validator/`:

- `pjtype` - Validates character type (human/supernatural)
- `userrole` - Validates user role (admin/master/player)

Standard validators from go-playground/validator:
- `required` - Field required
- `min/max` - String length, numeric ranges
- `alphanum` - Alphanumeric only
- `uuid` - Valid UUID format
- `gt/gte/lt/lte` - Numeric comparisons

## Error Handling

Centralized error mapping in `internal/infrastructure/api/handler/errors_mapper.go`

**Status Codes**:
- `200 OK` - Successful operation
- `201 Created` - Resource created
- `400 Bad Request` - Validation failures
- `401 Unauthorized` - Missing/invalid auth
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `406 Not Acceptable` - Business logic violations
- `409 Conflict` - Resource conflict (e.g., username exists)
- `500 Internal Server Error` - Server errors

**Error Response Format**:
```json
{
  "error": "Human-readable error message",
  "code": "MACHINE_READABLE_CODE"
}
```

## Security Considerations

1. **Password Security**:
   - Minimum 8 characters
   - Maximum 72 characters (bcrypt limitation)
   - Hashed with bcrypt before storage
   - Never returned in API responses

2. **JWT Security**:
   - Tokens expire (configurable, default 1h)
   - Secret key stored in environment variables
   - Claims include minimal data (ID, role)

3. **Authorization**:
   - Middleware enforces role-based access
   - Resource ownership validated (e.g., user can only update own character)
   - Campaign master privileges enforced

4. **Validation**:
   - All inputs validated using go-playground/validator
   - Custom validators for domain-specific rules
   - SQL injection prevented by GORM parameterization

## Testing Strategy

1. **Unit Tests**:
   - Domain logic in `internal/domain/*/`
   - Use case logic in `internal/application/*/`
   - Mock external dependencies

## Logging

Using Logrus for structured logging:

- Request/response logging
- Error logging with stack traces
- Business event logging
- Configurable log levels (debug, info, warn, error)

## Graceful Shutdown

Server implements graceful shutdown:
- Catches SIGINT/SIGTERM signals
- Finishes processing current requests
- Closes database connections
- Closes RabbitMQ connections
- Timeout: 5 seconds

## Performance Considerations

1. **Database**:
   - Indexes on foreign keys
   - Connection pooling via GORM
   - Prepared statements for repeated queries

2. **Caching**:
   - JWT tokens cached in client for reuse
   - Consider adding Redis for session caching (future)

3. **Async Processing**:
   - Domain events processed asynchronously via RabbitMQ
   - Non-blocking event publishing

## Additional Resources

- **OpenAPI Documentation**: `openapi.yml` - Complete API reference
- **Environment Template**: `.env.example` - Configuration template
- **Docker Compose**: `compose.yml` - Local development setup
- **Migrations**: `migrations/` - Database schema evolution

## Support & Contribution

For questions or contributions:
1. Review the OpenAPI documentation for API details
2. Follow Clean Architecture patterns
3. Write tests for new features
4. Update OpenAPI spec when adding/modifying endpoints
5. Follow Go best practices and project conventions

---

**Version**: 1.0.0
**Last Updated**: 2026-02-07
**Maintained By**: Meye Core Team