# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Tidy dependencies
RUN go mod tidy

# Build statically linked binaries for the API Gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-gateway ./cmd/gateway

# Final stage (using distroless static Debian 11)
FROM gcr.io/distroless/static-debian11

WORKDIR /

# Copy the binary from the build stage
COPY --from=builder /app/todo-gateway .

# Default CMD to run the API Gateway
CMD ["./todo-gateway"]

