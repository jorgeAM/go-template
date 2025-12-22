# Go Template

A comprehensive boilerplate template for building production-ready Golang APIs with clean architecture principles. This template provides a complete foundation with domain-driven design, extensive utility packages, and enterprise-grade features for scalable API development.

## Features

### Core Architecture

- **Clean Architecture**: Domain-driven design with clear separation of concerns (Domain, Application, Infrastructure layers)
- **Structured Project Layout**: Organized folder structure following Go best practices
- **Configuration Management**: Type-safe environment variable loading with defaults
- **Logging**: Structured logging with Zap integration
- **HTTP Routing**: Chi v5 router with comprehensive middleware stack
- **Database Integration**: PostgreSQL with `sqlx` and `goqu` query builder
- **Migration System**: Database migrations using golang-migrate

### Security & Authentication

- **Password Hashing**: Secure bcrypt password hashing
- **JWT Authentication**: Token generation and validation
- **CORS Support**: Configurable cross-origin resource sharing
- **Request Security**: Request ID tracking, real IP detection, timeout protection

### Utility Packages (15+ packages)

- **Collections**: Generic utilities for data manipulation (chunking, key-by operations)
- **Criteria**: Advanced query filtering, pagination, and ordering system
- **Events**: Event bus system with in-memory and AWS SNS/SQS support
- **Mailer**: Multi-provider email sending (SendGrid, AWS SES, in-memory)
- **Storage**: Cloud storage integration (Cloudflare R2 presigned URLs)
- **HTTP Utilities**: Response helpers, REST client, comprehensive middleware
- **Database**: Transaction management and connection handling
- **Error Handling**: Custom error types with metadata support
- **Crypto**: JWT and password utilities
- **Model**: Common value objects (Country, Currency, Email, etc.)

### AWS Integration

- **S3**: File storage operations
- **SES**: Email sending service
- **SNS**: Pub/sub messaging
- **SQS**: Message queue processing

### Development Tools

- **Mock Generation**: Automated mock generation for testing
- **Test Coverage**: Built-in coverage reporting
- **Docker Support**: Containerization ready
- **Code Generation**: Go generate integration

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

    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_DB=mydb
    POSTGRES_USER=admin
    POSTGRES_PASSWORD=passwd123
    POSTGRES_MAX_IDLE_CONNECTIONS=10
    POSTGRES_MAX_OPEN_CONNECTIONS=30
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

**Create User Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

## Available Make Commands

- `make generate` - Run go generate for mock generation
- `make test` - Run tests with coverage
- `make show-cover` - Display test coverage in browser
- `make tidy` - Tidy and vendor dependencies
- `make run` - Start the application
- `make new_migration MIGRATION_NAME=<name>` - Create new migration
- `make migration_up` - Run database migrations
- `make migration_down` - Rollback database migrations

## Project Structure

```
├── cmd/app/                    # Application entry point
├── internal/user/              # User domain module
│   ├── application/            # Use cases and commands
│   ├── domain/                 # Business logic and entities
│   ├── infrastructure/         # External concerns (HTTP, persistence)
│   └── mock/                   # Generated mocks
├── pkg/                        # Reusable packages
│   ├── collections/           # Generic utilities
│   ├── criteria/              # Query filtering system
│   ├── crypto/                # JWT and password utilities
│   ├── db/                    # Database utilities
│   ├── env/                   # Environment variable loading
│   ├── errors/                # Error handling
│   ├── events/                # Event bus system
│   ├── http/                  # HTTP utilities and middleware
│   ├── log/                   # Structured logging
│   ├── mailer/                # Email sending
│   ├── model/                 # Value objects
│   ├── pin/                   # PIN generation
│   ├── ref/                   # Pointer utilities
│   └── storage/               # Cloud storage
├── database/migration/         # Database migrations
└── cfg/                       # Configuration management
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/jorgeAM/go-template/blob/main/LICENCE) file for details.
