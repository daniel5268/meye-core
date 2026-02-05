# Meye-Core Project Context

## Project Overview

Meye-Core is an RPG campaign management system built with **Hexagonal Architecture** (Ports and Adapters) and **Domain-Driven Design** principles. The system manages users, campaigns, invitations, and player characters (PJs) with sophisticated character progression mechanics.

## Architecture Overview

This project follows **Hexagonal Architecture** with clear separation across three layers:

```
Domain Layer (Core Business Logic)
    ↑
Application Layer (Use Cases)
    ↑
Infrastructure Layer (Adapters: HTTP, Database, Services)
```

**Key Principle**: Dependencies point inward. Infrastructure depends on Application and Domain, but never the reverse.

## Project Structure

```
meye-core/
├── cmd/
│   └── server/
│       ├── main.go                    # Application entry point
│       └── dependencies.go            # Dependency injection container
│
├── internal/
│   ├── domain/                        # Domain Layer (Business Logic)
│   │   ├── user/
│   │   │   ├── user.entity.go
│   │   │   ├── user.repository.port.go
│   │   │   ├── hash.service.port.go
│   │   │   ├── jwt.service.port.go
│   │   │   └── errors.go
│   │   ├── campaign/
│   │   │   ├── campaign.aggregate.go         # Aggregate root
│   │   │   ├── invitation.entity.go          # Child entity
│   │   │   ├── pj.entity.go                  # Child entity (complex)
│   │   │   ├── campaign.repository.port.go
│   │   │   ├── calculations.go               # Shared calculation logic
│   │   │   ├── calculations_basic.go         # Basic stats XP
│   │   │   ├── calculations_special.go       # Special stats XP
│   │   │   ├── calculations_supernatural.go  # Supernatural stats XP
│   │   │   ├── errors.go
│   │   │   └── *_test.go                     # Domain tests
│   │   └── shared/
│   │       └── identification.service.port.go
│   │
│   ├── application/                   # Application Layer (Use Cases)
│   │   ├── user/
│   │   │   ├── ports.go               # Use case interfaces
│   │   │   ├── dto.go                 # Input + Output DTOs
│   │   │   ├── errors.go              # Application errors
│   │   │   ├── createuser/
│   │   │   │   └── create_user.usecase.go
│   │   │   └── login/
│   │   │       └── login.usecase.go
│   │   └── campaign/
│   │       ├── ports.go
│   │       ├── dto.go
│   │       ├── errors.go
│   │       ├── createcampaign/
│   │       │   └── create_campaign.usecase.go
│   │       └── inviteuser/
│   │           └── invite_user.usecase.go
│   │
│   ├── infrastructure/                # Infrastructure Layer (Adapters)
│   │   ├── api/
│   │   │   ├── api.go                 # Router setup
│   │   │   ├── handler/
│   │   │   │   ├── user.handler.go    # HTTP handlers
│   │   │   │   ├── campaign.handler.go
│   │   │   │   ├── auth.handler.go    # Auth middleware
│   │   │   │   ├── errors_mapper.go   # Error to HTTP mapping
│   │   │   │   └── dto/               # HTTP request/response DTOs
│   │   │   │       ├── user.dto.go
│   │   │   │       └── campaign.dto.go
│   │   │   └── validator/
│   │   │       └── validators.go       # Custom validators
│   │   ├── repository/
│   │   │   ├── user/postgres/
│   │   │   │   ├── user.repository.go
│   │   │   │   └── user.model.go      # Database model
│   │   │   └── campaign/postgres/
│   │   │       ├── campaign.repository.go
│   │   │       ├── campaign.model.go
│   │   │       ├── campaign_invitation.model.go
│   │   │       └── campaign_pj.model.go
│   │   ├── hash/
│   │   │   └── hash.service.go        # Bcrypt implementation
│   │   ├── jwt/
│   │   │   └── jwt.service.go         # JWT implementation
│   │   └── identification/
│   │       └── identification.service.go  # UUID generation
│   │
│   └── config/
│       └── config.go                  # Configuration loading
│
├── tests/
│   ├── mocks/                         # Generated mocks (gomock)
│   └── data/                      # Test data builders
│       ├── user.data.go
│       ├── campaign.data.go
│       └── pj.data.go
│
├── migrations/                        # Database migrations
│   ├── 000_create_users_table.up.sql
│   ├── 000_create_users_table.down.sql
│   ├── 001_create_campaign_table.up.sql
│   ├── 001_create_campaign_table.down.sql
│   ├── 002_create_campaign_invitations_table.up.sql
│   ├── 002_create_campaign_invitations_table.down.sql
│   ├── 003_create_pjs_table.up.sql
│   ├── 003_create_pjs_table.down.sql
│   ├── 004_update_pjs_talents_system.up.sql
│   └── 004_update_pjs_talents_system.down.sql
│
├── .env.example                       # Environment variables template
├── go.mod
└── go.sum
```

## Domain Model

### User Domain

**Entity: User**
- Properties: `id`, `username`, `hashedPassword`, `role`
- Roles: `admin`, `master` (game master), `player`
- Business rules:
  - Password is hashed on entity creation
  - Only players can be invited to campaigns

**Repository Port**: `Save()`, `FindByUsername()`, `FindByID()`

**Service Ports**: `HashService`, `JWTService`

---

### Campaign Domain (Complex Aggregate)

The campaign domain is the core of the application with rich business logic for RPG campaign and character management.

#### Aggregate Root: Campaign

**Entity: Campaign** (`campaign.aggregate.go`)
- Properties: `id`, `masterID`, `name`, `invitations[]`, `pjs[]`
- Aggregate children: Invitations, PJs
- Methods:
  - `InviteUser(user, identificationService) -> Invitation` - creates invitation for player
  - `AddPJ(userID, params, identificationService) -> PJ` - creates PJ after invitation
  - `GetPendingUserInvitation(userID) -> Invitation` - finds pending invitation

**Business Rules**:
- Only players can be invited
- User must have pending invitation to create PJ
- Invitation automatically accepts when PJ is created

---

#### Child Entity: Invitation

**Entity: Invitation** (`invitation.entity.go`)
- Properties: `id`, `campaignID`, `userID`, `state`
- States: `pending`, `accepted`
- Internal method: `accept()` - transitions to accepted state

---

#### Child Entity: PJ (Player Character)

**Entity: PJ** (`pj.entity.go`)

The most complex entity in the system, representing a player character with detailed stats and progression.

**Core Properties**:
- `id`, `userID`, `name`
- Physical attributes: `weight`, `height`, `age`, `look`
- Character traits: `charisma`, `villainy`, `heroism`
- Type: `pjType`
- Stats: `basicStats`, `specialStats`, `supernaturalStats`

**Enums**:
- `PJType`: `human`, `supernatural`

**Value Objects** (Immutable stat structures with talent flags):

1. **Basic Stats** (for all PJs):
   - `Physical`: strength, agility, speed, resistance, **isTalented**
   - `Mental`: intelligence, wisdom, concentration, will, **isTalented**
   - `Coordination`: precision, calculation, range, reflexes, **isTalented**
   - `BasicStats`: physical, mental, coordination, life

2. **Special Stats** (for all PJs):
   - `PhysicalSkills`: empowerment, vitalControl, **isTalented**
   - `MentalSkills`: illusion, mentalControl, **isTalented**
   - `EnergySkills`: objectHandling, energyHandling, **isTalented**
   - `SpecialStats`: physical, mental, energy, energyTank

3. **Supernatural Stats** (only for supernatural PJs):
   - `Skill`: transformations[] (array of uint)
   - `SupernaturalStats`: skills[] (nullable)

**Talent System**:
- Each stat group has an `isTalented` boolean flag
- Multiple talents can be selected (e.g., physical + mental + energy skills)
- Talents reduce XP costs for that specific stat category

---

#### Domain Logic: XP Calculation System

The campaign domain includes a sophisticated XP (experience points) calculation system for character progression.

**Files**:
- `calculations.go` - Base calculation utilities
- `calculations_basic.go` - Basic stats XP calculation
- `calculations_special.go` - Special stats XP calculation
- `calculations_supernatural.go` - Supernatural stats XP calculation

**Key Concepts**:

1. **Level Steps**: Different progression speeds for stat categories
   - Basic stats: level step = 10 (faster progression)
   - Special stats: level step = 100 (slower progression)
   - Supernatural stats: level step = 100

2. **Talent Multipliers**: Cost varies based on stat group's `isTalented` flag
   - Basic stats (Physical, Mental, Coordination):
     - isTalented = true: 1x cost (cheaper)
     - isTalented = false: 3x cost (expensive)
   - Special stats (PhysicalSkills, MentalSkills, EnergySkills):
     - isTalented = true: 1x cost (cheaper)
     - isTalented = false: 2x cost (expensive)
   - Energy tank:
     - Basic Coordination isTalented = true: 50% cost reduction
     - Otherwise: Standard cost

3. **XP Formula**: `XP = completeLevels * (completeLevels + firstLevelCost * 2 - 1) / 2 * levelStep + partialPoints * firstLevelCost`
   - Calculates total XP based on stat points
   - Different costs for first level based on isTalented flags

**Methods**:
- `BasicStats.GetRequiredXP() -> int` - uses isTalented from each stat group
- `SpecialStats.GetRequiredXP(isEnergyTalented bool) -> int` - uses isTalented from skill groups + energy talent flag
- `SupernaturalStats.GetRequiredXP() -> int`

**Use Case**: These calculations determine character progression costs and enable balanced character development in the RPG system.

---

#### Repository Port

**CampaignRepository** (`campaign.repository.port.go`)
- `Save(ctx, campaign)` - atomically saves aggregate with all children
- `FindByID(ctx, id)` - loads campaign with all related entities

**Aggregate Persistence Pattern**:
The repository handles the entire aggregate atomically:
1. Upserts campaign
2. Upserts/deletes invitations (synchronizes with domain state)
3. Upserts/deletes PJs (synchronizes with domain state)
4. All in a single transaction

---

## Application Layer

### Use Case Pattern

**Structure**:
```
application/
└── [domain]/
    ├── ports.go        # Use case interfaces + Input DTOs
    ├── dto.go          # Output DTOs + mappers
    ├── errors.go       # Application-level errors
    └── [usecase]/
        └── [usecase].usecase.go  # Implementation
```

**Key Pattern**: Use case ports defined at application boundary
```go
// ports.go
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

**Implementation Pattern**:
```go
// createuser/create_user.usecase.go
package createuser

// Compile-time interface check
var _ applicationuser.CreateUserUseCase = (*UseCase)(nil)

type UseCase struct {
    userRepository       domainuser.Repository  // Domain port
    identificationService shared.IdentificationService
    hashService          domainuser.HashService
}

func (u *UseCase) Execute(ctx context.Context, input applicationuser.CreateUserInput) (applicationuser.UserOutput, error) {
    // Use case orchestration logic
}
```

---

### Current Use Cases

#### User Use Cases

1. **CreateUser** (`createuser/create_user.usecase.go`)
   - Dependencies: userRepository, identificationService, hashService
   - Flow:
     1. Create user entity (validates & hashes password)
     2. Check username uniqueness
     3. Save to repository
     4. Return user output
   - Errors: `ErrUsernameAlreadyExists`

2. **Login** (`login/login.usecase.go`)
   - Dependencies: userRepository, hashService, jwtService
   - Flow:
     1. Find user by username
     2. Verify password hash
     3. Generate JWT token
     4. Return token
   - Errors: `ErrInvalidCredentials`, `ErrUserNotFound`

---

#### Campaign Use Cases

1. **CreateCampaign** (`createcampaign/create_campaign.usecase.go`)
   - Dependencies: campaignRepository, identificationService
   - Flow:
     1. Create campaign aggregate
     2. Save to repository
     3. Return campaign output

2. **InviteUser** (`inviteuser/invite_user.usecase.go`)
   - Dependencies: campaignRepository, userRepository, identificationService
   - Flow:
     1. Find campaign by ID
     2. Find user by ID
     3. Call campaign.InviteUser() (domain logic validates user is player)
     4. Save campaign aggregate (with new invitation)
     5. Return invitation output
   - Errors: `ErrCampaignNotFound`, `ErrUserNotFound`, `ErrUserNotPlayer`

---

## Infrastructure Layer

### API Layer

**Framework**: Gin (HTTP web framework)

**Router Structure** (`api.go`):
```go
type Router struct {
    router   *gin.Engine
    handlers *Handlers
}
```

#### Endpoints

**Public Endpoints**:
- `GET /health` - Health check
- `POST /api/v1/users/login` - User login

**Protected Endpoints**:
- `POST /api/v1/users` - Create user (requires admin role)
- `POST /api/v1/campaigns` - Create campaign (requires master role)
- `POST /api/v1/campaigns/:campaignID/invitations` - Invite user (requires campaign master)

---

### Handlers

**Handler Pattern**: Depend on use case ports, not concrete implementations

```go
type UserHandler struct {
    createUserUseCase user.CreateUserUseCase  // Interface, not concrete type
    loginUseCase      user.LoginUseCase
}
```

**Handler DTOs** (`handler/dto/`):
- Request DTOs: validation tags (`required`, `min`, `max`, `alphanum`)
- Response DTOs: consistent JSON structure
- Separate from application DTOs (different layer responsibilities)

---

### Authentication & Authorization

**Auth Handler** (`handler/auth.handler.go`)

**Middleware Functions**:
1. `InternalAPIKeyMiddleware()` - Validates API key header
2. `AuthMiddleware()` - Validates JWT Bearer token, sets user context
3. `RequireRole(...roles)` - Generic role checker (OR logic)
4. `RequireAdminRole()` - Admin only
5. `RequireMasterRole()` - Master only
6. `RequirePlayerRole()` - Player only
7. `RequireCampaignMaster()` - Validates user is the campaign master (resource ownership)

**Auth Context**:
```go
type AuthContext struct {
    UserID string
}
```
- Stored in Gin context after successful JWT validation
- Retrieved by handlers for authorization checks

**Middleware Chain Example**:
```go
protected := router.Group("/api/v1")
protected.Use(authHandler.AuthMiddleware())  // Validate JWT
protected.Use(authHandler.RequireMasterRole())  // Check role

protected.POST("/campaigns", campaignHandler.CreateCampaign)
```

---

### Error Handling

**Error Mapping** (`handler/errors_mapper.go`)

Central function maps domain/application errors to HTTP status codes:
```go
func respondMappedError(c *gin.Context, err error)
```

**Mappings**:
- `ErrUsernameAlreadyExists` → 409 Conflict
- `ErrInvalidCredentials` → 401 Unauthorized
- `ErrUserNotFound` → 404 Not Found
- `ErrCampaignNotFound` → 404 Not Found
- `ErrUserNotPlayer` → 400 Bad Request
- `ErrUserNotInvited` → 400 Bad Request
- Default → 500 Internal Server Error

---

### Repository Implementations

#### User Repository (`repository/user/postgres/`)

**Model**:
```go
type User struct {
    ID             string `gorm:"primaryKey"`
    Username       string `gorm:"unique"`
    HashedPassword string
    Role           string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

**Methods**:
- `Save()` - Upsert by ID with conflict resolution
- `FindByUsername()` - Returns nil if not found
- `FindByID()` - Returns nil if not found

**Mappers**:
- `GetModelFromDomainUser(user) -> Model`
- `(m *User) ToDomainUser() -> domain.User`

---

#### Campaign Repository (`repository/campaign/postgres/`)

**Models**:

1. `Campaign`: id, name, master_id, timestamps
2. `CampaignInvitation`: id, campaign_id, user_id, state, timestamps
3. `CampaignPJ`: id, campaign_id, user_id, [40+ stat columns], supernatural_stats (JSONB), timestamps

**Complex Model: CampaignPJ**
- Uses domain types directly: `campaign.PJType`
- All stat fields stored as separate columns for querying
- Talent flags stored as 6 boolean columns: `is_physical_talented`, `is_mental_talented`, `is_coordination_talented`, `is_physical_skills_talented`, `is_mental_skills_talented`, `is_energy_skills_talented`
- Supernatural stats stored as JSONB (nullable for human PJs)
- Custom GORM type `SupernaturalStatsJSON` implements `driver.Valuer` and `sql.Scanner`

**Aggregate Save Pattern**:
```go
func (r *Repository) Save(ctx context.Context, c *campaign.Campaign) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. Upsert campaign
        // 2. Upsert all invitations
        // 3. Delete removed invitations
        // 4. Upsert all PJs
        // 5. Delete removed PJs
    })
}
```

**Synchronization Logic**:
- Tracks IDs in domain aggregate
- Inserts/updates entities present in aggregate
- Deletes entities removed from aggregate
- Ensures database reflects exact aggregate state

---

### Service Implementations

1. **HashService** (`infrastructure/hash/`)
   - Implementation: bcrypt with cost 12
   - `Hash()`, `Compare()`

2. **JWTService** (`infrastructure/jwt/`)
   - Implementation: golang-jwt/jwt (HS256)
   - Claims: UserID + standard claims (exp, iss)
   - `GenerateSignedToken()`, `ValidateToken()`

3. **IdentificationService** (`infrastructure/identification/`)
   - Implementation: Google UUID (v4)
   - `GenerateID()` - returns UUID as string

---

## Configuration

**Structure** (`internal/config/config.go`):
```go
type Config struct {
    Api      ApiConfig
    Database DatabaseConfig
    JWT      JWTConfig
}
```

**Environment Variables**:
- `API_PORT` - Server port
- `API_KEY` - Internal API key
- `DATABASE_DSN` - Postgres connection string
- `JWT_SECRET` - JWT signing secret
- `JWT_ISSUER` - JWT issuer claim
- `JWT_EXPIRATION_TIME` - Token expiration (duration string)

**Loading**:
- `.env` file for development (godotenv)
- System environment variables for production
- Validation: returns error if required fields missing

---

## Database

### ORM: GORM

**Connection**: PostgreSQL via DSN

**Migration Management**: SQL files in `migrations/`
- Naming: `XXX_description.up.sql` / `XXX_description.down.sql`
- Applied manually or via migration tool

### Tables

1. **users**
   - Columns: id (PK), username (unique), hashed_password, role, timestamps
   - Indexes: unique username

2. **campaigns**
   - Columns: id (PK), master_id (FK), name, timestamps
   - Foreign key: master_id → users(id) CASCADE

3. **campaign_invitations**
   - Columns: id (PK), campaign_id (FK), user_id (FK), state (enum), timestamps
   - Foreign keys: campaign_id → campaigns(id) CASCADE, user_id → users(id) CASCADE
   - Indexes: campaign_id, user_id, state
   - Enum: invitation_state (pending, accepted)

4. **pjs**
   - Columns: id (PK), campaign_id (FK), user_id (FK), character info (name, weight, height, age, look, charisma, villainy, heroism), pj_type, 30+ stat columns, 6 talent boolean columns (is_physical_talented, is_mental_talented, is_coordination_talented, is_physical_skills_talented, is_mental_skills_talented, is_energy_skills_talented), supernatural_stats (JSONB), timestamps
   - Foreign keys: campaign_id → campaigns(id) CASCADE, user_id → users(id) CASCADE
   - Indexes: campaign_id, user_id, pj_type
   - Enums: pj_type

---

## Dependency Injection

**Container** (`cmd/server/dependencies.go`):
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
```

**Initialization Order**:
1. `loadEnvironment()` - Load .env file
2. `loadConfig()` - Parse environment variables
3. `connectDatabase()` - Connect to Postgres
4. `initializeServices()` - Create stateless services
5. `initializeRepositories()` - Create repositories
6. `initializeUseCases()` - Wire use cases with domain ports
7. `initializeHandlers()` - Wire handlers with use case ports
8. `initializeRouter()` - Setup routes and middleware

**Key Pattern**: Use case structs hold interfaces
```go
type UserUseCases struct {
    CreateUser user.CreateUserUseCase  // Interface, not concrete type
    Login      user.LoginUseCase
}
```

---

## Testing

### Test Organization

**Test Types**:
1. **Domain Tests**: In domain packages with `_test` suffix
   - `campaign.aggregate_test.go` - Tests aggregate behavior
   - `calculations_basic_test.go` - Tests XP calculation logic
   - `calculations_special_test.go`
   - `calculations_supernatural_test.go`

2. **Use Case Tests**: In use case packages with `_test` suffix
   - Use mocks for domain ports

**Test Data Builders** (`tests/data/`):
- `user.data.go` - User fixtures
- `campaign.data.go` - Campaign fixtures
- `pj.data.go` - PJ stat fixtures, complex value objects

---

### Mock Generation

**Tool**: `go.uber.org/mock/gomock`

**Pattern**: Use `go:generate` directives in port files
```go
//go:generate mockgen -destination=../../../tests/mocks/user_repository_mock.go -package=mocks meye-core/internal/domain/user Repository
```

**Generated Mocks** (`tests/mocks/`):
- `*_mock.go` files for all domain ports
- Used in use case tests

**Generate Command**: `go generate ./...`

---

## Key Design Patterns

### 1. Hexagonal Architecture (Ports & Adapters)

**Ports** (Interfaces):
- Domain ports: Repository, Service interfaces in domain layer
- Application ports: Use case interfaces in application layer

**Adapters** (Implementations):
- Infrastructure implementations of domain ports
- Handlers are adapters for HTTP (implement HTTP → use case translation)

**Dependency Direction**: Infrastructure → Application → Domain (inward only)

---

### 2. Aggregate Root Pattern

**Aggregate**: Campaign
- Root entity: Campaign
- Child entities: Invitation, PJ
- Consistency boundary: All changes go through aggregate root
- Transactional boundary: Entire aggregate saved atomically

**Rules**:
- External references only to root (by campaign ID)
- Children cannot be modified directly
- Root enforces invariants (e.g., only players invited, PJ requires invitation)

---

### 3. Value Objects Pattern

**Immutable**: Physical, Mental, Coordination, BasicStats, SpecialStats, etc.

**Characteristics**:
- No identity (defined by values)
- Immutable (no setters)
- Created via factory functions (`CreateXWithoutValidation`)
- Used extensively in PJ entity for stat composition

---

### 4. Repository Pattern

**Interface in Domain**: Defines what persistence operations are needed

**Implementation in Infrastructure**: How data is persisted

**Aggregate Persistence**: Repository handles entire aggregate atomically
- Save() synchronizes domain aggregate state with database
- FindByID() reconstructs aggregate with all children

---

### 5. Use Case Pattern (Interactors)

**Single Responsibility**: Each use case = one business workflow

**Interface**: Defined at application boundary
- Input DTO: Simple data structure
- Output DTO: Simple data structure
- Execute method: Orchestrates domain objects

**Dependencies**: Only domain ports (Repository, Services)

---

### 6. Factory Methods

**Purpose**: Entity reconstruction from database

**Pattern**:
```go
// Domain entity
func CreateUserWithoutValidation(id, username, hashedPassword string, role UserRole) *User {
    return &User{...}
}
```

**Usage**: Infrastructure layer uses factories to convert database models to domain entities

**Naming**: `Create[Entity]WithoutValidation` - bypasses domain validation for trusted data sources

---

### 7. Middleware Chain

**Composable Middleware**:
```go
router.Use(authHandler.AuthMiddleware())       // Step 1: Validate JWT
router.Use(authHandler.RequireMasterRole())    // Step 2: Check role
router.Use(authHandler.RequireCampaignMaster()) // Step 3: Check ownership
```

**Context Passing**: Each middleware enriches Gin context
- AuthMiddleware: sets UserID
- Route handlers: retrieve UserID from context

---

## Common Patterns & Conventions

### Adding a New Use Case

1. **Define Input DTO** in `application/[domain]/ports.go`:
```go
type MyUseCaseInput struct {
    Field1 string
    Field2 int
}
```

2. **Define Use Case Port** in `application/[domain]/ports.go`:
```go
type MyUseCase interface {
    Execute(ctx context.Context, input MyUseCaseInput) (OutputDTO, error)
}
```

3. **Implement Use Case** in `application/[domain]/myusecase/`:
```go
package myusecase

var _ applicationdomain.MyUseCase = (*UseCase)(nil)

type UseCase struct {
    // Dependencies (domain ports only)
}

func (u *UseCase) Execute(ctx context.Context, input applicationdomain.MyUseCaseInput) (applicationdomain.OutputDTO, error) {
    // Implementation
}
```

4. **Add to UseCases Struct** in `dependencies.go`:
```go
type DomainUseCases struct {
    MyUseCase domain.MyUseCase  // Interface type
}
```

5. **Wire in Handler**:
```go
type MyHandler struct {
    myUseCase applicationdomain.MyUseCase  // Interface
}
```

6. **Add Route** in `api.go`:
```go
router.POST("/my-endpoint", myHandler.MyEndpoint)
```

---

### Adding a New Domain Service

1. **Define Port** in `domain/[domain]/[service].port.go`:
```go
package domain

type MyService interface {
    DoSomething(param string) error
}
```

2. **Implement in Infrastructure** in `infrastructure/[service]/`:
```go
package service

type Service struct {
    // Configuration
}

func (s *Service) DoSomething(param string) error {
    // Implementation
}
```

3. **Wire in Dependencies** in `dependencies.go`:
```go
services.MyService = service.New()
```

---

### Adding a New Endpoint

1. **Create Handler DTOs** in `infrastructure/api/handler/dto/`:
```go
type MyRequestBody struct {
    Field string `json:"field" binding:"required"`
}

type MyResponseBody struct {
    Result string `json:"result"`
}
```

2. **Create Handler Method**:
```go
func (h *MyHandler) MyEndpoint(c *gin.Context) {
    var req dto.MyRequestBody
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Map to use case input
    input := application.MyUseCaseInput{...}

    // Execute use case
    output, err := h.myUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        respondMappedError(c, err)
        return
    }

    // Map to response DTO
    response := dto.MyResponseBody{...}
    c.JSON(200, response)
}
```

3. **Add Route** in `api.go`:
```go
router.POST("/my-endpoint", myHandler.MyEndpoint)
```

---

## Error Handling Strategy

### Error Types

1. **Domain Errors** (`domain/[domain]/errors.go`):
   - Business rule violations
   - Example: `ErrUserNotPlayer`, `ErrUserNotInvited`

2. **Application Errors** (`application/[domain]/errors.go`):
   - Use case failures
   - Example: `ErrUsernameAlreadyExists`, `ErrInvalidCredentials`

3. **Infrastructure Errors**:
   - Technical failures (database, network)
   - Example: GORM errors, HTTP errors

### Error Propagation

**Flow**: Domain → Application → Infrastructure

**Handler Mapping**:
```go
func respondMappedError(c *gin.Context, err error) {
    switch err {
    case application.ErrUsernameAlreadyExists:
        c.JSON(409, gin.H{"error": err.Error()})
    case application.ErrInvalidCredentials:
        c.JSON(401, gin.H{"error": err.Error()})
    // ... more mappings
    default:
        c.JSON(500, gin.H{"error": "Internal server error"})
    }
}
```

---

## Request Validation

### Gin Binding Validation

**Tags on DTOs**:
```go
type CreateUserInputBody struct {
    Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"required,userrole"`
}
```

**Standard Tags**: `required`, `min`, `max`, `alphanum`, `email`, etc.

**Custom Validators** (`api/validator/validators.go`):
```go
// Custom validator for user role enum
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("userrole", validateUserRole)
}
```

**Validation Failure**: Returns 400 Bad Request with validation errors

---

## Build & Run Commands

```bash
# Build all packages
go build ./...

# Run server
go run cmd/server/main.go

# Run all tests
go test ./...

# Run specific test package
go test ./internal/domain/campaign/...

# Run tests with verbose output
go test -v ./...

# Generate mocks
go generate ./...

# Run with race detector
go test -race ./...
```

---

## Key Dependencies

**Framework & Web**:
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/go-playground/validator/v10` - Request validation

**Database**:
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - Postgres driver

**Security**:
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `golang.org/x/crypto/bcrypt` - Password hashing

**Testing**:
- `github.com/stretchr/testify/assert` - Test assertions
- `go.uber.org/mock/gomock` - Mock generation

**Utilities**:
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading

---

## Key Principles

### Architectural Principles

1. **Dependency Inversion**: Infrastructure depends on domain/application, never reverse
2. **Separation of Concerns**: Each layer has distinct responsibility
3. **Explicit Boundaries**: Port interfaces define layer boundaries
4. **Immutability**: Value objects are immutable
5. **Single Responsibility**: Each use case does one thing
6. **Interface Segregation**: Small, focused interfaces

### Code Principles

1. **Compile-Time Checks**: `var _ Interface = (*Concrete)(nil)` pattern
2. **No Circular Dependencies**: DTOs at layer boundaries to avoid cycles
3. **Explicit Over Implicit**: Clear interfaces and dependencies
4. **Business Logic in Domain**: Handlers and use cases orchestrate, don't implement rules
5. **Factory Methods**: Clear construction patterns for entities
6. **Consistent Error Handling**: Errors flow from domain → application → infrastructure

---

## Common Pitfalls to Avoid

❌ **Don't** make handlers depend on concrete use case types
```go
type Handler struct {
    useCase *createuser.UseCase  // ❌ Concrete type
}
```

✅ **Do** make handlers depend on use case port interfaces
```go
type Handler struct {
    useCase user.CreateUserUseCase  // ✅ Interface
}
```

---

❌ **Don't** define Input DTOs in use case subdirectories
```go
// createuser/input.go
package createuser
type Input struct {...}  // ❌ Creates circular dependency
```

✅ **Do** define Input DTOs in ports.go at the boundary
```go
// ports.go
package user
type CreateUserInput struct {...}  // ✅ At application boundary
```

---

❌ **Don't** have use cases depend on infrastructure
```go
type UseCase struct {
    db *gorm.DB  // ❌ Infrastructure dependency
}
```

✅ **Do** have use cases depend on domain ports
```go
type UseCase struct {
    repository domain.Repository  // ✅ Domain port
}
```

---

❌ **Don't** put business logic in handlers
```go
func (h *Handler) CreateUser(c *gin.Context) {
    // ❌ Validation and business rules in handler
    if len(req.Password) < 8 {
        return errors.New("password too short")
    }
}
```

✅ **Do** put business logic in domain entities and use cases
```go
// Domain entity
func NewUser(...) (*User, error) {
    // ✅ Validation in domain
    if len(password) < 8 {
        return nil, ErrPasswordTooShort
    }
}
```

---

❌ **Don't** modify aggregate children directly
```go
invitation := campaign.Invitations()[0]
invitation.Accept()  // ❌ Direct modification bypasses aggregate
```

✅ **Do** modify through aggregate root
```go
campaign.AddPJ(userID, params)  // ✅ Goes through aggregate root
```

---

## Project-Specific Features

### Campaign XP System

The project includes a comprehensive XP calculation system for RPG character progression:

**Features**:
- Three stat categories: Basic, Special, Supernatural
- Talent-based cost multipliers
- Non-linear progression formulas
- Support for both human and supernatural characters

**Implementation**: Domain logic in `calculations_*.go` files

**Testing**: Extensive tests verify XP calculations for all scenarios

### Character Stat System

**Complexity**: PJ entity has 40+ individual stat fields organized into value objects

**Flexibility**:
- Basic stats: universal attributes (strength, agility, intelligence, etc.)
- Special stats: advanced skills (empowerment, illusion, energy handling)
- Supernatural stats: transformation abilities (JSONB in database)
- Multiple talents: Each stat group can be talented independently, allowing flexible character builds

**Pattern**: Composition of value objects maintains clean domain model while supporting complex game mechanics

### Multi-Talent System (2024 Refactor)

**Evolution**: Refactored from single basic/special talent selection to multiple independent talent flags

**Previous System**:
- `BasicTalentType` enum: one of physical, mental, coordination, energy
- `SpecialTalentType` enum: one of physical, mental, energy
- Limited to 2 total talents per character

**New System**:
- 6 independent `isTalented` boolean flags on stat groups
- Allows multiple talent combinations (e.g., physical + mental + energy skills)
- More flexible character progression system
- Calculation methods updated to use boolean flags instead of enum comparisons

**Migration**: `004_update_pjs_talents_system` converts existing talent enums to boolean flags

### Aggregate Persistence with Children

**Challenge**: Save aggregate root with multiple child collections (invitations, PJs)

**Solution**: Transactional repository pattern with synchronization logic
- Tracks changes in domain aggregate
- Synchronizes database state with aggregate state
- Handles inserts, updates, and deletes in single transaction

---

## When to Create New Layers

**New Domain Package**: When adding a distinct business capability
- Example: `order`, `payment`, `inventory`
- Should have clear bounded context

**New Use Case**: When adding a new business workflow
- Example: `AcceptInvitation`, `UpdateCampaign`, `DeletePJ`
- Each workflow gets its own use case

**New Handler**: When adding a new API endpoint
- One handler per domain/resource
- Methods correspond to endpoints

**New Service Port**: When needing integration with external systems
- Example: `EmailService`, `NotificationService`, `PaymentGateway`
- Define port in domain, implement in infrastructure

**New Repository**: When adding a new aggregate root
- Each aggregate root gets its own repository
- Repository manages aggregate persistence

---

## Project Status

The project is actively developed with focus on campaign management and character progression features. Recent additions include:
- PJ (player character) entity with complex stat system
- XP calculation system for character progression
- Invitation workflow for campaign management
- Aggregate pattern implementation for campaign persistence

The architecture is well-established and follows clean architecture principles with clear boundaries between layers.
