package main

import (
    "context"
    "os"

    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "todo-app/proto"
)

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    rootCmd := &cobra.Command{Use: "todo-client"}

    addCmd := &cobra.Command{
        Use:   "add",
        Short: "Add a new task",
        Run: func(cmd *cobra.Command, args []string) {
            title, _ := cmd.Flags().GetString("title")
            description, _ := cmd.Flags().GetString("description")
            addTask(logger, title, description)
        },
    }

    getCmd := &cobra.Command{
        Use:   "get",
        Short: "Get all tasks",
        Run: func(cmd *cobra.Command, args []string) {
            getTasks(logger)
        },
    }

    addCmd.Flags().String("title", "", "Title of the task")
    addCmd.Flags().String("description", "", "Description of the task")

    rootCmd.AddCommand(addCmd, getCmd)

    if err := rootCmd.Execute(); err != nil {
        logger.Fatal("Error executing CLI command", zap.Error(err))
        os.Exit(1)
    }
}

func addTask(logger *zap.Logger, title, description string) {
    logger.Info("Attempting to add a new task", zap.String("title", title))

    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        logger.Fatal("Failed to connect", zap.Error(err))
    }
    defer conn.Close()

    client := proto.NewTodoServiceClient(conn)
    _, err = client.AddTask(context.Background(), &proto.TaskRequest{Title: title, Description: description})
    if err != nil {
        logger.Fatal("Error adding task", zap.Error(err))
    }

    logger.Info("Task added successfully", zap.String("title", title))
}

func getTasks(logger *zap.Logger) {
    logger.Info("Fetching tasks")

    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        logger.Fatal("Failed to connect", zap.Error(err))
    }
    defer conn.Close()

    client := proto.NewTodoServiceClient(conn)
    response, err := client.GetTasks(context.Background(), &proto.Empty{})
    if err != nil {
        logger.Fatal("Error getting tasks", zap.Error(err))
    }

    for _, task := range response.Tasks {
        logger.Info("Task found", zap.String("title", task.Title), zap.Bool("completed", task.Completed))
    }
}

