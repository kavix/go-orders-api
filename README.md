# Go Orders API

A REST API for managing orders using Go, Chi router, and Redis.

## Features

- Order management API
- Redis-backed order storage
- HTTP routing with Chi
- Modular project structure

## Tech Stack

- Go
- Chi
- Redis

## Project Structure

```text
application/   # application setup and services
handeler/      # HTTP handlers
model/         # data models
order/         # order domain/repository logic
main.go        # API entrypoint
```

## Prerequisites

- Go
- Redis running locally

## Run Locally

```bash
go mod tidy
go run main.go
```

## Notes

- Make sure Redis is running before starting the API.
- Update Redis connection settings if needed.
