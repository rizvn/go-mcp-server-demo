# Build stage
FROM golang:1.24-alpine AS builder

COPY . /code
RUN mkdir -p /code/dist

COPY src/.env /code/dist/.env

WORKDIR /code/src

# Install dependencies
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../dist/mcp-go .

# Final stage
FROM busybox:latest
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /code/dist .

# run the application
ENTRYPOINT ["./mcp-go"]