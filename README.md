# Go Template

A comprehensive boilerplate template for building production-ready Golang APIs with clean architecture principles. This template provides a complete foundation with domain-driven design, extensive utility packages, and enterprise-grade features for scalable API development.

## Features

### Core Architecture

- **Modular Monolith**: Each business module lives in `internal/<module>/` and implements `module.Module` (`Name`, `Init`, `RegisterHttp`, `MigrationFS`); modules are registered in `internal/bootstrap`
- **Hexagonal Architecture**: Each module splits into `domain/`, `app/`, `api/` (driving adapters) and `adapters/` (driven adapters)
- **CQRS Pattern**: Separate command and query handlers for better separation of concerns
- **Repository Pattern**: Domain-driven repository interfaces with clean abstractions
- **Module-owned Wiring**: Each module builds its own dependencies in `module.go` `Init()` from its own `config.go`
- **OpenAPI-first HTTP**: Per-module `openapi.yaml` compiled to a chi strict server by `oapi-codegen`
- **Structured Project Layout**: Organized folder structure following Go best practices
- **Configuration Management**: Type-safe environment variable loading with defaults and validation
- **Logging**: Structured logging with Zap integration and context-aware request tracking
- **HTTP Routing**: Chi v5 router with comprehensive middleware stack
- **Database Integration**: PostgreSQL with `pgx` and `sqlc`-generated, type-safe queries
- **Migration System**: Database migrations using golang-migrate with up/down support

### Security & Authentication

- **Password Hashing**: Secure bcrypt password hashing and comparison
- **JWT Authentication**: Token generation, validation, and type checking; failures map to client-safe sentinel errors (`crypto.ErrTokenExpired`, `ErrTokenMalformed`, `ErrTokenType`, `ErrTokenIssuer`)
- **Authentication Middleware**: Bearer token validation with principal (JWT `sub` + claims) context injection (available in `middleware.Authenticate`, not mounted by default)
- **Cookie Management**: Refresh token cookie with HttpOnly and SameSite protection; the caller passes the `secure` flag from its own config
- **CORS Support**: Configurable cross-origin resource sharing with flexible origin/method control
- **Request Security**: Request ID tracking, real IP detection, configurable timeout protection

### Advanced Query System

- **Dynamic Filtering**: Support for multiple operators (EQ, GT, LT) with PostgreSQL conversion
- **Pagination**: Built-in page and page_size validation with offset calculation
- **Ordering**: ASC/DESC ordering with field validation
- **Query Parameter Parsing**: Automatic conversion from HTTP query parameters to criteria objects
- **Type-safe SQL Generation**: hand-written SQL compiled to Go by `sqlc`

### Utility Packages (15+ packages)

- **Collections**: Generic utilities for data manipulation (chunking for batch processing)
- **Criteria**: Advanced query filtering, pagination, and ordering system
- **Events**: Complete event bus system with in-memory and AWS SNS/SQS implementations
- **Mailer**: Multi-provider email sending (SendGrid, AWS SES, in-memory for testing)
- **Storage**: Cloud storage integration (Cloudflare R2 presigned URLs with content type validation)
- **HTTP Utilities**: Response helpers, REST client with retry support, comprehensive middleware
- **Database**: Transaction management interface and connection handling
- **Error Handling**: Custom error types with error codes, cause tracking, and metadata support
- **Crypto**: JWT utilities and secure password hashing
- **Value objects**: Email with regex validation, UUID ID, timestamps
- **Environment**: Type-safe environment variable loading for int, string, bool types
- **PIN**: Cryptographically secure 4-digit PIN generation
- **Reference**: Generic pointer utility functions
- **Log**: Structured logging with multiple levels and cloneable loggers

### Event Bus Architecture

- **Event Model**: Complete event system with ID, topic, payload, timestamp, and versioning
- **Multiple Publishers**:
  - In-memory publisher (for testing/local development)
  - AWS SNS publisher with batch support (max 10 events per batch)
- **Multiple Listeners**:
  - In-memory listener with handler registration
  - AWS SQS listener with message polling and automatic deletion
- **Event Handlers**: Interface-based event handling with topic routing
- **Event Collector**: Testing utilities for event verification

### AWS Integration

- **S3 Compatible**: Cloudflare R2 integration with presigned URL generation
- **SES**: Email sending service with HTML content support
- **SNS**: Pub/sub messaging with batch publishing and message attributes
- **SQS**: Message queue processing with configurable wait times and batch processing

### HTTP Middleware Stack

- **Request ID**: Automatic request ID generation and tracking
- **Structured Logging**: Request/response logging with method, path, status, duration
- **Panic Recovery**: Graceful panic handling with stack trace logging
- **Real IP Detection**: Accurate client IP detection behind proxies
- **CORS**: Configurable cross-origin resource sharing
- **Response Headers**: Automatic Content-Type and Accept header injection
- **Timeout Management**: Per-request timeout with X-Timeout header support (default 15s)

`middleware.Authenticate` (JWT Bearer validation with principal context injection) is available but not mounted in `cmd/app/router.go`; add it to the routes that need it.

### Development Tools

- **Mock Generation**: Automated mock generation using `go.uber.org/mock` with `//go:generate` directives
- **SQL & API Codegen**: `sqlc` and `oapi-codegen`, pinned as Go tools (`go tool`)
- **Test Coverage**: Built-in coverage reporting with browser display
- **Docker Support**: Multi-stage Docker build with distroless base image for security
- **Code Generation**: Go generate integration for mocks and other generated code
- **REST Client**: Built-in HTTP client with retry support and error handling

## Getting Started

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/jorgeAM/go-template.git
    cd go-template
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```
3.  **Set up environment variables:**
    Create a .env file in the root directory with the following variables (you can copy .env.example):

    ```env
    # Application Configuration
    APP_ENV=local
    APP_NAME=go-template
    PORT=8080

    # Database Configuration
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_DB=mydb
    POSTGRES_USER=admin
    POSTGRES_PASSWORD=passwd123
    POSTGRES_MAX_OPEN_CONNECTIONS=30

    # JWT Configuration
    JWT_KEY=your-secret-jwt-key
    JWT_ISSUER=your-app-name
    ```

4.  **Set up the database:**
    Run database migrations:

    ```sh
    make migrate
    ```

5.  **Run the application:**
    Start the server:

    ```sh
    make run
    ```

    Or run directly:

    ```sh
    go run cmd/app/main.go
    ```

## API Endpoints

### Health Check

- `GET /health` - Health check endpoint

### User Management

- `POST /api/v1/user` - Create a new user (`201` with the created user, `400` on invalid data)
- `GET /api/v1/user/{id}` - Get user by ID (`404` when it doesn't exist)

**Create User Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Get User Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Query Parameters Support

The `criteria` package (`criteria.CriteriaInput`) can bind these URL parameters for list endpoints. No endpoint uses it yet:

- `order_by` - Field to order by
- `order_type` - ASC or DESC
- `page` - Page number (starts from 1)
- `page_size` - Number of items per page

Example for a list endpoint: `?order_by=created_at&order_type=DESC&page=1&page_size=10`

## Available Make Commands

- `make generate` - Run go generate (sqlc, oapi-codegen, mockgen)
- `make unit-tests` - Run tests with coverage reporting
- `make show-cover` - Display test coverage in browser
- `make tidy` - Tidy and vendor dependencies
- `make run` - Start the application server
- `make migrate` - Apply every module's pending migrations (`cmd/migrate`)
- `make new_migration MODULE=<module> MIGRATION_NAME=<name>` - Create new migration files for a module
- `make migration_down MODULE=<module>` - Rollback the last migration of a module

## Usage Examples

### Event System

```go
// Create and publish an event
event, err := events.NewEvent("user.created", map[string]interface{}{
    "user_id": "123",
    "email": "user@example.com",
})
if err != nil {
    return err
}

// Publish to SNS
publisher := eventbus.NewSNSPublisher(snsClient, topicArn)
err = publisher.Publish(ctx, event)
```

### Email Sending

```go
// Using SendGrid
sender := mailer.NewSendgridMailer(sendgridClient)
payload := &mailer.MailerPayload{
    From:    "noreply@example.com",
    To:      "user@example.com",
    Subject: "Welcome!",
    Body:    "<h1>Welcome to our service!</h1>",
}
err := sender.Send(ctx, payload)
```

### Advanced Querying

```go
// Build criteria for filtering and pagination
criteria, err := criteria.FromPrimitive(&criteria.CriteriaPrimitive{
    Filters: []*criteria.FilterPrimitive{
        {Field: "status", Operator: "EQ", Value: "active"},
        {Field: "created_at", Operator: "GT", Value: "2023-01-01"},
    },
    OrderBy:   &[]string{"created_at"}[0],
    OrderType: &[]string{"DESC"}[0],
    Page:      1,
    PageSize:  10,
})
```

### Cloud Storage

```go
// Generate presigned URL for file upload
signer, err := storage.NewCloudflareR2Client(bucketName, accessKey, secretKey, endpoint)
if err != nil {
    return err
}

url, err := signer.GeneratePresignedURL(ctx, "file.jpg", storage.IMAGE_JPEG)
```

## Project Structure

```
├── cmd/app/                    # Application entry point
│   ├── main.go                # Initializes every module, graceful shutdown
│   └── router.go              # Router, middleware, each module's RegisterHttp
├── cmd/migrate/                # Applies every module's migrations
├── internal/bootstrap/         # Module list shared by cmd/app and cmd/migrate
├── internal/identity/          # Identity module (users), reference module layout
│   ├── module.go              # Module bootstrap: Init, RegisterHttp, MigrationFS
│   ├── config.go              # Env vars this module needs
│   ├── domain/                # User aggregate, invariants, repository interface, sentinel errors
│   ├── app/                   # Use cases
│   │   ├── command/           # Write operations (CreateUser)
│   │   ├── query/             # Read operations (GetUser)
│   │   └── models/            # Output shapes (UserInfo, never the password hash)
│   ├── api/http/              # openapi.yaml + oapi-codegen strict server (server.go)
│   ├── adapters/db/           # sqlc scaffold: migrations/, queries/, generated sqlc/, repository
│   └── mocks/                 # Generated mocks for testing
├── internal/platform/          # I/O adapters shared across modules
│   ├── db/                    # Transaction management
│   ├── eventbus/              # Event publishers/listeners (in-memory, SNS, SQS)
│   ├── http/                  # HTTP utilities
│   │   ├── handler/           # Common HTTP handlers (health check)
│   │   ├── middleware/        # Authentication, CORS, logging, timeout
│   │   └── restclient/        # REST client with retries
│   ├── log/                   # Structured logging with Zap
│   ├── mailer/                # Multi-provider email sending
│   └── storage/               # Cloud storage (Cloudflare R2 presigned URLs)
├── internal/shared/            # Domain-agnostic types, zero I/O
│   ├── collections/           # Generic utilities (chunking, key-by operations)
│   ├── criteria/              # Query filtering, pagination, ordering
│   ├── crypto/                # JWT and password utilities
│   ├── env/                   # Environment variable loading with type safety
│   ├── errors/                # Custom error types with metadata and error codes
│   ├── events/                # Domain event, topic and collector
│   ├── generator/             # Cryptographically secure PIN generation
│   ├── module/                # module.Module interface every module implements
│   ├── ref/                   # Pointer utility functions
│   └── valueobject/           # Value objects (Email, ID, Timestamps)
├── vendor/                    # Vendored dependencies
├── Dockerfile                 # Multi-stage Docker build with distroless base
├── Makefile                   # Development and deployment commands
└── go.mod                     # Go module definition with all dependencies
```

## Key Dependencies

- **Go 1.26.5** - Latest Go version
- **Chi v5** - Lightweight HTTP router
- **Zap** - Structured logging
- **pgx + sqlc** - PostgreSQL driver and type-safe query generation
- **oapi-codegen** - OpenAPI → chi strict server generation
- **httpin** - HTTP request input binding
- **golang-migrate** - Per-module migrations (`cmd/migrate`)
- **AWS SDK v2** - S3, SES, SNS, SQS integration
- **golang-jwt** - JWT token handling
- **bcrypt** - Password hashing
- **SendGrid** - Email service integration
- **go-resty** - HTTP client with retry support
- **testify + gomock** - Testing framework and mocks

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jorgeAM/go-template/blob/main/LICENCE) file for details.
