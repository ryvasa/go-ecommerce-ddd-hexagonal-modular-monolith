# Go DDD Hexagonal Modular Monolith

A robust and scalable e-commerce backend boilerplate built with Go, following Domain-Driven Design (DDD), Hexagonal Architecture, and Modular Monolith principles.

## Tech Stack

- **Language**: Go 1.25.1
- **Framework**: HTTP Router ([Gin](https://github.com/gin-gonic/gin))
- **ORM**: [GORM](https://gorm.io/) with MySQL
- **Dependency Injection**: [Google Wire](https://github.com/google/wire)
- **Authentication**: JWT ([golang-jwt](https://github.com/golang-jwt/jwt))
- **Validation**: [validator/v10](https://github.com/go-playground/validator)
- **Configuration**: [godotenv](https://github.com/joho/godotenv)

## Architecture Overview

The project follows a modular monolith structure where each module is self-contained and communicates through well-defined ports or events.

- `domain` : Pure business logic, entities, and value objects. No external dependencies.
- `application` : Use cases, input ports (interfaces), and service implementations.
- `adapter/in` : Primary adapters (e.g., HTTP Handlers, CLI).
- `adapter/out` : Secondary adapters (e.g., Persistence with GORM, External APIs).
- `cmd/server` : Composition root where dependencies are wired together.

## Modules

1.  **User**: Management of user accounts, authentication (JWT), and security.
2.  **Task**: Project management functionality.
3.  **Cart**: E-commerce shopping cart management.

### Inter-Module Communication

Modules are decoupled. For example, the **Task** module interacts with the **User** module via **ports (`UserReader`)** or through the **Event System**, ensuring no direct implementation dependency.

## Shared & Infrastructure

- `internal/shared`: Reusable components used across modules.
  - `transaction`: Transaction Manager for atomic operations.
  - `apperror`: Unified error handling.
  - `response`: Standardized API responses.
  - `middleware`: RequestID, Logger, Auth, etc.
- `pkg/`: Infrastructure-level packages.
  - `database`: GORM/MySQL connection management.
  - `event`: In-memory event publisher/subscriber.
  - `logger`: Structured logging.

## Database & Transactions

Transactions are managed at the **Application Service** layer using a `TransactionManager`.

- **Commit**: Automatic if the use case returns `nil`.
- **Rollback**: Automatic if the use case returns an `error`.

## Quick Start

1.  **Clone & Install**:
    ```bash
    go mod download
    ```
2.  **Setup Environment**: Create a `.env` file in the root directory.
3.  **Generate DI**:
    ```bash
    wire ./cmd/server
    ```
4.  **Run Server**:
    ```bash
    go run ./cmd/server
    ```

## Development

- **Migrations**: Database migrations are handled automatically or via script (see `migrations/`).
- **Composition**: Use `wire_gen.go` for dependency injection. Run `wire` command after changing providers.
