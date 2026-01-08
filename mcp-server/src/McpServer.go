package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rizvn/go-mcp/auth"
	"github.com/rizvn/go-mcp/echo"
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
	oauthMiddleWare := &auth.OAuthMiddleware{
		IssuerUrl:         r.IssuerURL,
		TargetAudienceUrl: r.McpServerURL,
	}
	oauthMiddleWare.Init()

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "simple-mcp-server",
		Version: "1.0.0",
	}, nil)

	echoTool := &echo.EchoTool{}
	// add tool to server
	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the input message",
	}, echoTool.Call)

	// create streamable HTTP handler
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	// Setup routing
	mux := http.NewServeMux()

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
