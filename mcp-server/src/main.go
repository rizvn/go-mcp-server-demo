package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load environment variableÓs from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Printf("Warning: No  env file found, %v\n", err)
	}
	configureLogging()
	
	issuerURL := os.Getenv("ISSUER_URL")
	if issuerURL == "" {
		log.Fatal("ISSUER_URL environment variable is required")
	}
	mcpServerUrl := os.Getenv("MCP_SERVER_URL")
	if mcpServerUrl == "" {
		log.Fatal("MCP_SERVER_URL environment variable is required")
	}

	server := &McpServer{
		IssuerURL:    issuerURL,
		McpServerURL: mcpServerUrl,
	}

	server.Start()
}

func configureLogging() {
	// set slog level from environment variable LOG_LEVEL
	level := os.Getenv("LOG_LEVEL")

	if level == "" {
		level = "INFO"
	}
	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	//logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	slog.SetDefault(slog.New(logger))

}
