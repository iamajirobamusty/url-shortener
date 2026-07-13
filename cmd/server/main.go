package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog" // Added for structured logging
	"net/http"
	"os" // Added to access system environment variables

	"url-shortener/internal/db"
	"url-shortener/internal/shortener"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, reading from system environment")
	}

	// FIX 1: Fixed the spelling of "parseTime"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 1. Establish structural infrastructure connection
	conn, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("Database connection failure here: %v", err)
	}
	defer conn.Close()

	// Verify the database is actually reachable
	if err := conn.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// Initialize sqlc execution engine
	queries := db.New(conn)

	// Inject our sqlc engine directly into the shortener constructor interface
	shortenerHandler := shortener.NewHandler(queries)

	// Declare clean routing paths
	mux := http.NewServeMux()

	// FIX 2: Ensure your route patterns start with a forward slash "/"
	mux.HandleFunc("POST /api/v1/shorten", shortenerHandler.Shorten)
	mux.HandleFunc("GET /r/", shortenerHandler.Redirect)

	// Determine server port (default to 8080 if not explicitly defined)
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("URL Shortener application ignition started on port %s...\n", serverPort)

	// FIX 3: Start the server so the application stays alive and listens for requests
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatalf("Server failed to spin up: %v", err)
	}
}
