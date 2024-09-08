package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net"
    "os"
    "time"

    "go.uber.org/fx"
    "go.uber.org/zap"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
    "google.golang.org/grpc"
    "todo-app/internal/handlers"
    "todo-app/proto"
)

// NewDB provides a *sql.DB or an error if something goes wrong
func NewDB(lc fx.Lifecycle) (*sql.DB, error) {
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
    log.Printf("Connecting to database using DSN: %s", dsn)

    var db *sql.DB
    var err error

    // Retry connection for a certain duration (e.g., 30 seconds)
    maxRetries := 10
    retryDelay := 3 * time.Second

    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("mysql", dsn)
        if err != nil {
            log.Printf("Error connecting to database: %v. Retrying in %v...", err, retryDelay)
            time.Sleep(retryDelay)
            continue
        }

        err = db.Ping()
        if err == nil {
            log.Println("Successfully connected to database")
            break
        }

        log.Printf("Database not ready yet: %v. Retrying in %v...", err, retryDelay)
        time.Sleep(retryDelay)
    }

    if err != nil {
        return nil, fmt.Errorf("could not connect to database: %w", err)
    }

    lc.Append(fx.Hook{
        OnStop: func(ctx context.Context) error {
            log.Println("Closing the database connection")
            return db.Close()
        },
    })

    return db, nil
}

func runServer(lifecycle fx.Lifecycle, logger *zap.Logger, todoService proto.TodoServiceServer) {
    lifecycle.Append(
        fx.Hook{
            OnStart: func(ctx context.Context) error {
                logger.Info("Starting gRPC server...")

                lis, err := net.Listen("tcp", ":50051")
                if err != nil {
                    logger.Fatal("Failed to listen", zap.Error(err))
                }
                grpcServer := grpc.NewServer()
                proto.RegisterTodoServiceServer(grpcServer, todoService)

                go func() {
                    if err := grpcServer.Serve(lis); err != nil {
                        logger.Fatal("Failed to serve gRPC", zap.Error(err))
                    }
                }()
                return nil
            },
            OnStop: func(ctx context.Context) error {
                logger.Info("Shutting down gRPC server...")
                return nil
            },
        },
    )
}

func main() {
    app := fx.New(
        fx.Provide(
            NewDB,
            handlers.NewTodoServiceHandler,
            NewLogger,
        ),
        fx.Invoke(runServer),
    )

    app.Run()
}

// NewLogger sets up the zap logger
func NewLogger() *zap.Logger {
    logger, _ := zap.NewProduction() // You can use zap.NewDevelopment() for dev mode
    defer logger.Sync()               // Flush logs
    return logger
}
