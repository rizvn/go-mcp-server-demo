package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rizvn/go-mcp/auth"
	"github.com/rizvn/go-mcp/echo"
)

type McpServer struct {
	IssuerURL      string
	TargetAudience string
	Scope          string
}

func (r *McpServer) Start() {

	// Initialize OAuth middleware
	oauthMiddleWare := &auth.OAuthMiddleware{
		IssuerUrl:      r.IssuerURL,
		TargetAudience: r.TargetAudience,
		Scope:          r.Scope,
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
	mux.Handle("/", oauthMiddleWare.Handler(mcpHandler))

	slog.Info(fmt.Sprintf("Starting MCP server on :8000"))
	slog.Info(fmt.Sprintf("Authorization Server URL: %s", r.IssuerURL))
	slog.Info(fmt.Sprintf("Resource URL: %s", r.TargetAudience))
	slog.Info(fmt.Sprintf("Tool available: echo"))
	slog.Info(fmt.Sprintf("OAuth2.1 endpoint:"))
	slog.Info(fmt.Sprintf("- /.well-known/oauth-protected-resource"))

	if err := http.ListenAndServe(":8000", mux); err != nil {
		slog.Error(fmt.Sprintf("Server failed: %v", err))
	}
}
