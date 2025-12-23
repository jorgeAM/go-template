# Go Template

A comprehensive boilerplate template for building production-ready Golang APIs with clean architecture principles. This template provides a complete foundation with domain-driven design, extensive utility packages, and enterprise-grade features for scalable API development.

## Features

### Core Architecture

- **Clean Architecture**: Domain-driven design with clear separation of concerns (Domain, Application, Infrastructure layers)
- **CQRS Pattern**: Separate command and query handlers for better separation of concerns
- **Repository Pattern**: Domain-driven repository interfaces with clean abstractions
- **Dependency Injection**: Configuration-based dependency management
- **Structured Project Layout**: Organized folder structure following Go best practices
- **Configuration Management**: Type-safe environment variable loading with defaults and validation
- **Logging**: Structured logging with Zap integration and context-aware request tracking
- **HTTP Routing**: Chi v5 router with comprehensive middleware stack
- **Database Integration**: PostgreSQL with `sqlx` and `goqu` query builder for type-safe SQL
- **Migration System**: Database migrations using golang-migrate with up/down support

### Security & Authentication

- **Password Hashing**: Secure bcrypt password hashing and comparison
- **JWT Authentication**: Complete JWT system with token generation, validation, and type checking
- **Authentication Middleware**: Bearer token validation with automatic user context injection
- **Cookie Management**: Secure refresh token handling with HttpOnly and SameSite protection
- **CORS Support**: Configurable cross-origin resource sharing with flexible origin/method control
- **Request Security**: Request ID tracking, real IP detection, configurable timeout protection

### Advanced Query System

- **Dynamic Filtering**: Support for multiple operators (EQ, GT, LT) with PostgreSQL conversion
- **Pagination**: Built-in page and page_size validation with offset calculation
- **Ordering**: ASC/DESC ordering with field validation
- **Query Parameter Parsing**: Automatic conversion from HTTP query parameters to criteria objects
- **Type-safe SQL Generation**: Goqu-based query building with prepared statements

### Utility Packages (15+ packages)

- **Collections**: Generic utilities for data manipulation (chunking for batch processing)
- **Criteria**: Advanced query filtering, pagination, and ordering system with PostgreSQL converter
- **Events**: Complete event bus system with in-memory and AWS SNS/SQS implementations
- **Mailer**: Multi-provider email sending (SendGrid, AWS SES, in-memory for testing)
- **Storage**: Cloud storage integration (Cloudflare R2 presigned URLs with content type validation)
- **HTTP Utilities**: Response helpers, REST client with retry support, comprehensive middleware
- **Database**: Transaction management interface and connection handling
- **Error Handling**: Custom error types with error codes, cause tracking, and metadata support
- **Crypto**: JWT utilities and secure password hashing
- **Model**: Value objects (Country, Currency, Email with regex validation, UUID, timestamps)
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
- **Authentication**: JWT Bearer token validation with user context injection

### Development Tools

- **Mock Generation**: Automated mock generation using `go.uber.org/mock` with `//go:generate` directives
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
    ENV=local
    PORT=8080

    # Database Configuration
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_DB=mydb
    POSTGRES_USER=admin
    POSTGRES_PASSWORD=passwd123
    POSTGRES_MAX_IDLE_CONNECTIONS=10
    POSTGRES_MAX_OPEN_CONNECTIONS=30

    # JWT Configuration
    JWT_KEY=your-secret-jwt-key
    JWT_ISSUER=your-app-name

    # Application Configuration
    APP_ENV=local
    APP_NAME=go-template
    ```

4.  **Set up the database:**
    Run database migrations:

    ```sh
    make migration_up
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

- `POST /api/v1/user` - Create a new user
- `GET /api/v1/user/{id}` - Get user by ID

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

The API supports advanced querying through URL parameters:

- `order_by` - Field to order by
- `order_type` - ASC or DESC
- `page` - Page number (starts from 1)
- `page_size` - Number of items per page

Example: `GET /api/v1/users?order_by=created_at&order_type=DESC&page=1&page_size=10`

## Available Make Commands

- `make generate` - Run go generate for mock generation
- `make test` - Run tests with coverage reporting
- `make show-cover` - Display test coverage in browser
- `make tidy` - Tidy and vendor dependencies
- `make run` - Start the application server
- `make new_migration MIGRATION_NAME=<name>` - Create new database migration files
- `make migration_up` - Run all pending database migrations
- `make migration_down` - Rollback the last database migration

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
publisher := events.NewSNSPublisher(snsClient, topicArn)
err = publisher.Publish(ctx, event)
```

### Email Sending

```go
// Using SendGrid
mailer := mailer.NewSendgridMailer(sendgridClient)
payload := &mailer.MailerPayload{
    From:    "noreply@example.com",
    To:      "user@example.com",
    Subject: "Welcome!",
    Body:    "<h1>Welcome to our service!</h1>",
}
err := mailer.Send(ctx, payload)
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

// Convert to PostgreSQL query
converter := criteria.NewCriteriaToPostgresConverter()
sql, args, err := converter.Convert(ctx, "users", criteria)
```

### Cloud Storage

```go
// Generate presigned URL for file upload
signer := storage.NewCloudflareR2Client(bucketName, accessKey, secretKey, endpoint)
url, err := signer.GeneratePresignedURL(ctx, "file.jpg", storage.JPEG)
```

## Project Structure

```
├── cmd/app/                    # Application entry point
│   ├── main.go                # Main application file with graceful shutdown
│   └── router.go              # HTTP router configuration and middleware setup
├── internal/user/              # User domain module (clean architecture)
│   ├── application/            # Use cases (commands and queries)
│   │   ├── command/           # Write operations (CreateUser)
│   │   └── query/             # Read operations (GetUser)
│   ├── domain/                # Business logic and entities
│   ├── infrastructure/        # External concerns (HTTP handlers, persistence)
│   │   ├── http/              # HTTP handlers
│   │   └── persistence/       # Database repositories
│   └── mock/                  # Generated mocks for testing
├── pkg/                       # Reusable packages (15+ utilities)
│   ├── collections/           # Generic utilities (chunking, key-by operations)
│   ├── criteria/              # Advanced query filtering, pagination, ordering
│   ├── crypto/                # JWT and password utilities
│   ├── db/                    # Database utilities and transaction management
│   ├── env/                   # Environment variable loading with type safety
│   ├── errors/                # Custom error types with metadata and error codes
│   ├── events/                # Event bus system (in-memory, SNS, SQS)
│   ├── http/                  # HTTP utilities, middleware, response helpers
│   │   ├── handler/           # Common HTTP handlers (health check)
│   │   ├── middleware/        # Authentication, CORS, logging, timeout
│   │   └── response/          # Response helper functions
│   ├── log/                   # Structured logging with Zap
│   ├── mailer/                # Multi-provider email sending
│   ├── model/                 # Value objects (Country, Currency, Email, etc.)
│   ├── pin/                   # Cryptographically secure PIN generation
│   ├── ref/                   # Pointer utility functions
│   └── storage/               # Cloud storage (Cloudflare R2 presigned URLs)
├── database/migration/         # Database migrations (up/down SQL files)
├── cfg/                       # Configuration management and dependency injection
│   ├── config.go              # Environment-based configuration loading
│   └── dependencies.go        # Dependency injection and service wiring
├── vendor/                    # Vendored dependencies
├── Dockerfile                 # Multi-stage Docker build with distroless base
├── Makefile                   # Development and deployment commands
└── go.mod                     # Go module definition with all dependencies
```

## Key Dependencies

- **Go 1.25.4** - Latest Go version
- **Chi v5** - Lightweight HTTP router
- **Zap** - Structured logging
- **sqlx + goqu** - Database operations and query building
- **AWS SDK v2** - S3, SES, SNS, SQS integration
- **golang-jwt** - JWT token handling
- **bcrypt** - Password hashing
- **SendGrid** - Email service integration
- **go-resty** - HTTP client with retry support
- **testify + gomock** - Testing framework and mocks

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jorgeAM/go-template/blob/main/LICENCE) file for details.
