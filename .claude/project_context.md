# Meye-Core Project Context

## Architecture Overview

This project follows **Hexagonal Architecture** (Ports and Adapters pattern) with clear separation of concerns across three layers.

## Project Structure

```
internal/
├── domain/                          # Domain Layer (Business Logic & Entities)
│   ├── user/
│   │   ├── user.entity.go           # Domain entities
│   │   ├── user.repository.port.go  # Repository port interfaces
│   │   ├── hash.service.port.go     # Service port interfaces
│   │   └── jwt.service.port.go
│   └── shared/
│       └── identification.service.port.go
│
├── application/                     # Application Layer (Use Cases)
│   └── user/
│       ├── ports.go                 # Use case port interfaces + Input DTOs
│       ├── user.output.go           # Shared output DTOs
│       ├── errors.go                # Application-level errors
│       ├── createuser/
│       │   └── create_user.usecase.go  # Concrete use case implementation
│       └── login/
│           └── login.usecase.go
│
└── infrastructure/                  # Infrastructure Layer (Adapters)
    ├── api/
    │   ├── handler/
    │   │   ├── auth.handler.go      # Authentication & authorization middleware
    │   │   ├── user.handler.go      # HTTP handlers
    │   │   └── dto/                 # HTTP-specific DTOs
    │   └── router.go
    ├── repository/
    │   └── user/
    │       └── postgres/
    │           └── user.repository.go  # Repository implementation
    ├── hash/
    │   └── hash.service.go          # Service implementations
    ├── jwt/
    │   └── jwt.service.go
    └── identification/
        └── identification.service.go
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)
**Purpose**: Core business logic, entities, and port interfaces

**What belongs here**:
- Domain entities (e.g., `User`)
- Business rules and validations
- Port interfaces (contracts for repositories and services)
- Domain-specific errors

**Dependencies**:
- ✅ No dependencies on other layers
- ✅ Pure Go with minimal external dependencies

**File naming conventions**:
- Entities: `*.entity.go`
- Port interfaces: `*.port.go` or `*.repository.port.go`

**Example**:
```go
// user.repository.port.go
package user

type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}
```

### 2. Application Layer (`internal/application/`)
**Purpose**: Orchestrate business workflows (use cases)

**What belongs here**:
- Use case implementations
- Use case port interfaces (at the boundary)
- Input/Output DTOs
- Application-level errors

**Dependencies**:
- ✅ Can depend on domain layer (entities and domain ports)
- ❌ MUST NOT depend on infrastructure layer

**File organization**:
```
application/
└── user/
    ├── ports.go              # Use case interfaces + Input DTOs (at boundary)
    ├── user.output.go        # Shared output DTOs
    ├── errors.go             # Application errors
    └── createuser/
        └── create_user.usecase.go  # Concrete implementation
```

**Key Pattern - Use Case Ports**:
```go
// ports.go - Define use case interfaces at the application boundary
package user

type CreateUserInput struct {
    Username string
    Password string
    Role     domainuser.UserRole
}

type CreateUserUseCase interface {
    Execute(ctx context.Context, input CreateUserInput) (UserOutput, error)
}
```

**Key Pattern - Use Case Implementation**:
```go
// createuser/create_user.usecase.go
package createuser

// Compile-time check to ensure implementation matches port
var _ applicationuser.CreateUserUseCase = (*UseCase)(nil)

type UseCase struct {
    userRepository domainuser.Repository  // Depends on domain ports
    hashService    domainuser.HashService
}

func (c *UseCase) Execute(ctx context.Context, input applicationuser.CreateUserInput) (applicationuser.UserOutput, error) {
    // Orchestrate domain objects and services
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)
**Purpose**: Implement adapters for external systems

**What belongs here**:
- HTTP handlers (implement REST API)
- Repository implementations (Postgres, Redis, etc.)
- Service implementations (JWT, hashing, external APIs)
- Database models
- HTTP-specific DTOs

**Dependencies**:
- ✅ Can depend on domain ports (implements them)
- ✅ Can depend on application ports (handlers use them)
- ❌ MUST NOT depend on concrete application or domain implementations

**Key Pattern - Handlers Depend on Use Case Ports**:
```go
// user.handler.go
package handler

import "meye-core/internal/application/user"

type UserHandler struct {
    createUserUseCase user.CreateUserUseCase  // Interface, not concrete type
    loginUseCase      user.LoginUseCase
}

func NewUserHandler(createUserUC user.CreateUserUseCase, loginUseCase user.LoginUseCase) *UserHandler {
    return &UserHandler{
        createUserUseCase: createUserUC,
        loginUseCase:      loginUseCase,
    }
}
```

## Common Patterns

### Adding a New Use Case

1. **Define Input DTO** in `internal/application/[domain]/ports.go`:
```go
type MyNewUseCaseInput struct {
    Field1 string
    Field2 int
}
```

2. **Define Use Case Port** in `internal/application/[domain]/ports.go`:
```go
type MyNewUseCase interface {
    Execute(ctx context.Context, input MyNewUseCaseInput) (OutputDTO, error)
}
```

3. **Implement Use Case** in `internal/application/[domain]/myfeature/`:
```go
package myfeature

var _ applicationdomain.MyNewUseCase = (*UseCase)(nil)

type UseCase struct {
    // Dependencies (domain ports only)
}

func (u *UseCase) Execute(ctx context.Context, input applicationdomain.MyNewUseCaseInput) (applicationdomain.OutputDTO, error) {
    // Implementation
}
```

4. **Use in Handler** - depend on the port interface:
```go
type MyHandler struct {
    myUseCase applicationdomain.MyNewUseCase  // Interface
}
```

### Adding a New Domain Service

1. **Define Port** in `internal/domain/[domain]/[service].port.go`:
```go
package user

type EmailService interface {
    SendEmail(to, subject, body string) error
}
```

2. **Implement in Infrastructure** in `internal/infrastructure/email/`:
```go
package email

type Service struct {
    // SMTP config, etc.
}

func (s *Service) SendEmail(to, subject, body string) error {
    // Implementation
}
```

3. **Wire in Dependency Injection** in `cmd/server/dependencies.go`

### Authentication & Authorization Middleware

**Location**: `internal/infrastructure/api/handler/auth.handler.go`

**Middleware Chain**:
```go
// 1. AuthMiddleware - validates JWT, sets UserID in context
router.Use(authHandler.AuthMiddleware())

// 2. RequireRole - checks user role from database
router.Use(authHandler.RequireMasterRole())        // Single role
router.Use(authHandler.RequireRole(Master, Admin))  // Multiple roles (OR logic)
```

**Available Middleware**:
- `InternalAPIKeyMiddleware()` - Check API key header
- `AuthMiddleware()` - Validate JWT, set AuthContext
- `RequireRole(...roles)` - Generic role checker (OR logic)
- `RequireMasterRole()` - Master only
- `RequireAdminRole()` - Admin only
- `RequirePlayerRole()` - Player only

### Dependency Injection

**Location**: `cmd/server/dependencies.go`

**Pattern**:
```go
type DependencyContainer struct {
    Config       *config.Config
    Database     *gorm.DB
    Services     *Services
    Repositories *Repositories
    UseCases     *UseCases
    Handlers     *Handlers
    APIRouter    *api.Router
}

// Initialization order:
1. loadEnvironment()
2. loadConfig()
3. connectDatabase()
4. initializeServices()    // Infrastructure services
5. initializeRepositories() // Infrastructure repositories
6. initializeUseCases()     // Application use cases (depend on domain ports)
7. initializeHandlers()     // Infrastructure handlers (depend on use case ports)
8. initializeRouter()       // API routes
```

**Key Principle**: Use case ports in the UseCases struct:
```go
type UserUseCases struct {
    CreateUser user.CreateUserUseCase  // Interface type
    Login      user.LoginUseCase       // Interface type
}
```

## Testing Patterns

### Use Case Tests

**Location**: Same package as use case with `_test` suffix

**Pattern**:
```go
package createuser_test

func TestUseCase_Execute(t *testing.T) {
    // 1. Setup mocks for domain ports
    ctrl := gomock.NewController(t)
    repoMock := mocks.NewMockRepository(ctrl)

    // 2. Setup expectations
    repoMock.EXPECT().FindByID(ctx, "123").Return(user, nil)

    // 3. Create use case with mocks
    useCase := createuser.NewUseCase(repoMock, ...)

    // 4. Execute and assert
    output, err := useCase.Execute(ctx, input)
    assert.NoError(t, err)
}
```

### Generating Mocks

Use `go:generate` directives in port files:
```go
//go:generate mockgen -destination=../../../tests/mocks/user_repository_mock.go -package=mocks meye-core/internal/domain/user Repository
```

Run: `go generate ./...`

## Error Handling

**Domain Errors**: Define in domain package (e.g., validation errors)

**Application Errors**: Define in `internal/application/[domain]/errors.go`:
```go
var (
    ErrUsernameAlreadyExists = errors.New("username already exists")
    ErrInvalidCredentials    = errors.New("invalid credentials")
)
```

**HTTP Error Mapping**: Handlers map application/domain errors to HTTP status codes

## Configuration

**Location**: `internal/config/`

**Loading**: Environment variables via `godotenv` (development) or system env (production)

**Pattern**:
```go
type Config struct {
    Database DatabaseConfig
    JWT      JWTConfig
    Api      ApiConfig
}
```

## Database

**ORM**: GORM

**Migrations**: Located in `migrations/`

**Pattern**: Separate database models from domain entities:
```go
// Infrastructure layer - Postgres model
type User struct {
    ID             string
    Username       string
    HashedPassword string
    Role           string
}

// Convert to domain entity
func (u *User) ToDomainUser() *user.User {
    return user.CreateUserWithoutValidation(...)
}
```

## Go Module

**Module name**: `meye-core`

**Go version**: Check `go.mod`

**Key dependencies**:
- `github.com/gin-gonic/gin` - HTTP framework
- `gorm.io/gorm` - ORM
- `github.com/golang-jwt/jwt` - JWT
- `go.uber.org/mock/gomock` - Mocking for tests

## Build Commands

```bash
# Build all packages
go build ./...

# Run tests
go test ./...

# Run specific test package
go test ./internal/application/user/createuser/...

# Generate mocks
go generate ./...

# Run server
go run cmd/server/main.go
```

## Key Principles

1. **Dependency Inversion**: Infrastructure depends on domain/application, not the other way around
2. **Port Interfaces**: Define at layer boundaries (domain ports, use case ports)
3. **Input DTOs**: Live with use case ports at application boundary
4. **No Circular Dependencies**: Input DTOs moved from child packages to parent to avoid cycles
5. **Explicit Implementation Checks**: Use `var _ Interface = (*Concrete)(nil)` pattern
6. **Clean Boundaries**: Handlers never import concrete use case implementations

## Common Pitfalls to Avoid

❌ **Don't** make handlers depend on concrete use case types
✅ **Do** make handlers depend on use case port interfaces

❌ **Don't** define Input DTOs in use case subdirectories (causes circular imports)
✅ **Do** define Input DTOs in `ports.go` at the boundary

❌ **Don't** have use cases depend on infrastructure
✅ **Do** have use cases depend on domain ports

❌ **Don't** put business logic in handlers
✅ **Do** put business logic in domain entities and use cases

## When to Create New Layers

**New Domain Package**: When adding a distinct business capability (e.g., `order`, `payment`)

**New Use Case**: When adding a new business workflow

**New Handler**: When adding a new API endpoint or middleware

**New Service Port**: When needing integration with external systems
