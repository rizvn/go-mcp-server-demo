# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY . /src
RUN mkdir dist

COPY .env dist/.env


# Install dependencies
RUN go mod download

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/mcp-go .

# Copy .env file to the dist directory
RUN cp .env dist/.env

# Final stage
FROM busybox:latest
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /src/dist .

# run the application
ENTRYPOINT ["./mcp-go"]