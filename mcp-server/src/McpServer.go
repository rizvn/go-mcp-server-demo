package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
	"github.com/rizvn/go-mcp/auth"
)

type McpServer struct {
	IssuerURL    string
	McpServerURL string
}

func (r *McpServer) Start() {

	// Initialize OAuth provider
	provider := &auth.OauthProvider{}
	provider.IssuerUrl = r.IssuerURL
	provider.Init()

	// Initialize OAuth middleware
	oauthMiddleWare := &OAuthMiddleware{}
	oauthMiddleWare.Init(provider, r.McpServerURL)

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "simple-mcp-server",
		Version: "1.0.0",
	}, nil)

	// add tool to server
	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the input message",
	}, Echo)

	// create streamable HTTP handler
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	// Setup routing
	mux := http.NewServeMux()

	// OAuth 2.1 metadata endpoint (no authorization required)
	// clients will use this to discover the resource server's authorization server and scopes
	mux.HandleFunc("/.well-known/oauth-protected-resource", r.HandleProtectedResourceMetadata)

	// MCP endpoint (OAuth authorization required, with logging)
	mux.Handle("/", LoggingMiddleware(oauthMiddleWare.Handler(mcpHandler)))

	log.Println("Starting MCP server on :8000")
	log.Printf("Authorization Server URL: %s", r.IssuerURL)
	log.Printf("JWKS URL: %s", provider.JwksUri)
	log.Printf("Resource URL: %s", r.McpServerURL)
	log.Println("Tool available: echo")
	log.Println("OAuth2.1 endpoint:")
	log.Println("- /.well-known/oauth-protected-resource")

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Printf("Server failed: %v", err)
	}
}

func (r *McpServer) HandleProtectedResourceMetadata(w http.ResponseWriter, rq *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if rq.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// the resource identifier (what resource this server protects),
	// which OAuth scopes the resource supports,
	// which authorization servers (issuer URLs) are authoritative for access tokens for this resourc
	metadata := oauthex.ProtectedResourceMetadata{
		Resource:             r.McpServerURL,
		ScopesSupported:      []string{"mcp:tools"},
		AuthorizationServers: []string{r.IssuerURL},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(metadata)
}
