# Concurrent URL Shortener API

A highly concurrent, lightning-fast URL shortener backend constructed in Go. The application leverages Go's native routing infrastructure, interfaces for decoupled dependency injection, asynchronous background goroutines for non-blocking telemetry tracking, and `sqlc` for compile-time safe database interactions with MySQL.

## Features
* **Asynchronous Telemetry:** Click logging, IP capturing, and User-Agent recording run inside an isolated fire-and-forget Goroutine pipeline, preventing database writes from blocking user redirection.
* **Type-Safe Persistence:** Employs `sqlc` to generate type-safe Go source code directly from schema SQL files.
* **Decoupled Architecture:** Built completely around interfaces (`DBQuerier`) to support seamless storage layer mocking and switching.

## Tech Stack
* **Language:** Go (Golang)
* **Database:** MySQL
* **Migration Utility:** Goose
* **SQL Compiler:** sqlc

---

## Getting Started

### Prerequisites
* Go 1.22+ installed
* MySQL server instance active
* Goose CLI binary installed

### 1. Environment Setup
Clone the repository and replicate the sample configuration parameters file:
```bash
cp .env.example .env

Open .env and fill in your explicit local database credentials.

### 2.Run Database Migrations
Initialize your database layout scheme using the Goose utility driver:

goose mysql "YOUR_DB_USER:YOUR_DB_PASS@tcp(127.0.0.1:3306)/YOUR_DB_NAME" up

### 3. Ignition
Launch the application server:

go run cmd/server/main.go

The server will boot and begin listening on port 8080

4. Running Tests
Execute the table-driven unit test suite using Go's built-in testing toolchain:

go test -v ./internal/shortener/...

API Specification
1. Create Shortened URL
Endpoint: POST /api/v1/shorten

Content-Type: application/json

Payload:

{
  "long_url": "[https://github.com/alx-tools](https://github.com/alx-tools)"
}

Success Response (201 Created):

{
  "code": "a4TYIK",
  "short_code": "[http://127.0.0.1:8080/r/a4TYIK](http://127.0.0.1:8080/r/a4TYIK)"
}

2. Execute Redirection
Endpoint: GET /r/{short_code}

Action: Instantly triggers an HTTP 302 Found redirection straight to the registered destination while capturing metrics asynchronously in the background.
