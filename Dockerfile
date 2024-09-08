# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-app ./cmd/server

# Final stage
FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /app/todo-app .

CMD ["./todo-app"]

