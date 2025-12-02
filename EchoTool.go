package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EchoInput struct {
	Message string `json:"message" jsonschema:"The message to echo back"`
}

type EchoOutput struct {
	Response string `json:"response" jsonschema:"The echoed response"`
}

func Echo(ctx context.Context, req *mcp.CallToolRequest, input *EchoInput) (*mcp.CallToolResult, EchoOutput, error) {
	return nil, EchoOutput{
		Response: input.Message,
	}, nil
}
