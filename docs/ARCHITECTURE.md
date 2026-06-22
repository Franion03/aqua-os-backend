# Architecture: Environment Parity Strategy

## Overview

The backend uses a **repository pattern** where the `internal/db/` package provides the data access layer. All HTTP handlers call db package functions — they never interact with the storage engine directly.

## Storage Engines

| Environment | Engine | Notes |
|-------------|--------|-------|
| Development | SQLite | File-based, zero config, instant startup |
| Production | DynamoDB | Managed by Terraform in `aqua-os-infrastructure` |

## Configuration

The switch is controlled by the `DB_DRIVER` environment variable:

- `sqlite` (default) — uses a local `.db` file, no external dependencies
- `dynamodb` — connects to AWS DynamoDB tables provisioned by Terraform

## Repository Interface

All handlers call these functions, which form the contract between handlers and storage:

- `GetAllLevels()`
- `GetLevel(id)`
- `GetExercises(levelID)`
- `AddExercise(req)`
- `DeleteExercise(id)`

The interface is defined in `internal/db/repository.go`.

## Implementation Files

| File | Purpose |
|------|---------|
| `internal/db/repository.go` | Interface definition |
| `internal/db/sqlite.go` | SQLite implementation (dev) |
| `internal/db/dynamo.go` | DynamoDB implementation (prod) |

## Adding DynamoDB Support

1. Implement the `Repository` interface in `internal/db/dynamo.go` using the AWS SDK for Go v2.
2. At init time, read `DB_DRIVER` and instantiate the corresponding implementation.
3. DynamoDB tables are created by Terraform in the `aqua-os-infrastructure` repo — the backend only reads/writes, never creates tables.
