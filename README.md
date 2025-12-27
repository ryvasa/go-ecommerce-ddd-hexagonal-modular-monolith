# Go DDD Hexagonal Modular Monolith

Boilerplate backend Go menggunakan:
- DDD (Domain Driven Design)
- Hexagonal Architecture
- Modular Monolith
- GORM
- Transaction Manager
- Dependency Injection (Wire)

## Architecture Overview

- `domain`      : pure business entity
- `application` : use case & port
- `adapter/in`  : HTTP (Gin)
- `adapter/out` : persistence (GORM)
- `cmd/server`  : composition root

## Modules
- User
- Task

Task module bergantung ke User module melalui **port (`UserReader`)**, bukan implementasi.

## Transaction
Transaction dikelola di application service menggunakan `TransactionManager`.

- Commit → jika use case return `nil`
- Rollback → jika use case return `error`

## Run
```bash
wire ./cmd/server
go run ./cmd/server
