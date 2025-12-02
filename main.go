package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load environment variableÓs from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Printf("Warning: No  env file found, %v\n", err)
	}
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
