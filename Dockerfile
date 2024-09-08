# Build stage
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Tidy dependencies
RUN go mod tidy

# Build statically linked binaries for gRPC server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-server ./cmd/server

# Final stage (using distroless static Debian 11)
FROM gcr.io/distroless/static-debian11

WORKDIR /

# Copy the binaries from the build stage
COPY --from=builder /app/todo-server .

# Default CMD to run the gRPC server
CMD ["./todo-server"]

