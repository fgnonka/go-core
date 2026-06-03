# go-core

A minimal Go web API for practicing core concepts (handlers, data layer, simple store).

## Overview

This repository provides a small HTTP API implemented in Go. The server exposes endpoints under `/api/v1/` and uses a lightweight internal data layer in `internal/data`.

## Prerequisites

- Go 1.20+ installed (latest stable version recommended)
- (Optional) Docker if you prefer containerized runs

## Quick start

From the repository root, run the server locally:

```bash
go run ./cmd/api/main.go
```

By default the server listens on port `8080`. Example request (from this workspace):

```bash
curl -i http://localhost:8080/api/v1/entities
```

## Build

Build a binary:

```bash
go build -o bin/api ./cmd/api
./bin/api
```

## Project structure

- `cmd/api/` — application entrypoint and HTTP wiring (`main.go`, `handlers.go`)
- `internal/data/` — data models and store implementation (`models.go`, `store.go`)
- `go.mod` — module definition

## API

Endpoints are mounted under `/api/v1/`. See `cmd/api/handlers.go` for handler implementations and available routes. The repository includes a simple entities endpoint used for manual testing.

## Development notes

- Follow Go idioms: keep packages small and testable.
- If you add a Dockerfile, you can run the app in a container. Example (if a Dockerfile exists):

```bash
docker build -t go-core .
docker run -p 8080:8080 go-core
```

## Testing

No test framework is included by default. Add tests under the packages you modify and run `go test ./...`.

## Contributing

Open issues or PRs with clear descriptions. Keep changes small and focused.

## License

This project does not include a license by default — add one if you plan to open-source it.

## Contact

For questions about this workspace, reach out to the repository owner.
