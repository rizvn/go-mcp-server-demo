package main

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rizvn/panics"
)

type EchoInput struct {
	Message string `json:"message" jsonschema:"The message to echo back"`
}

type EchoOutput struct {
	Response string `json:"response" jsonschema:"The echoed response"`
}

func Echo(ctx context.Context, req *mcp.CallToolRequest, input *EchoInput) (*mcp.CallToolResult, EchoOutput, error) {
	authHeader := req.GetExtra().Header.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// dont to verify the token signature, as it would have been validated by middleware
	token, _, err := jwt.NewParser().ParseUnverified(accessToken, jwt.MapClaims{})
	panics.OnError(err, "failed to parse JWT token")

	claims, ok := token.Claims.(jwt.MapClaims)
	panics.OnFalse(ok, "failed to get JWT claims from token")

	email, ok := claims["email"].(string)
	panics.OnFalse(ok, "email claim is missing in JWT token")

	if authHeader == "" {
		panic("Authorization header is missing")
	}

	// of accessing request headers
	return nil, EchoOutput{
		Response: input.Message + " from " + email,
	}, nil
}
