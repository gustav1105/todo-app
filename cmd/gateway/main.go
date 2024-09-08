package main

import (
    "context"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "todo-app/proto"
)

func main() {
    router := gin.Default()

    log.Println("Starting API Gateway...")

    // Connect to gRPC server
    conn, err := grpc.Dial("grpc-server:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Could not connect to gRPC server: %v", err)
    }
    log.Println("Successfully connected to gRPC server")
    defer conn.Close()

    grpcClient := proto.NewTodoServiceClient(conn)

    // REST Endpoint: POST /tasks
    router.POST("/tasks", func(c *gin.Context) {
        var taskRequest struct {
            Title       string `json:"title"`
            Description string `json:"description"`
        }

        if err := c.BindJSON(&taskRequest); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        // Call gRPC AddTask method
        _, err := grpcClient.AddTask(context.Background(), &proto.TaskRequest{
            Title:       taskRequest.Title,
            Description: taskRequest.Description,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Task added successfully"})
    })

    // REST Endpoint: GET /tasks
    router.GET("/tasks", func(c *gin.Context) {
        tasksResponse, err := grpcClient.GetTasks(context.Background(), &proto.Empty{})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
            return
        }

        c.JSON(http.StatusOK, tasksResponse.Tasks)
    })

    // Start the HTTP server
    router.Run(":8080")
}

