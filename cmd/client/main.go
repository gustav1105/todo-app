package main

import (
    "context"
    "log"
    "os"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
    "todo-app/proto"
)

func main() {
    rootCmd := &cobra.Command{Use: "todo-client"}
    
    addCmd := &cobra.Command{
        Use:   "add",
        Short: "Add a new task",
        Run: func(cmd *cobra.Command, args []string) {
            title, _ := cmd.Flags().GetString("title")
            description, _ := cmd.Flags().GetString("description")
            addTask(title, description)
        },
    }
    
    getCmd := &cobra.Command{
        Use:   "get",
        Short: "Get all tasks",
        Run: func(cmd *cobra.Command, args []string) {
            getTasks()
        },
    }
    
    addCmd.Flags().String("title", "", "Title of the task")
    addCmd.Flags().String("description", "", "Description of the task")
    
    rootCmd.AddCommand(addCmd, getCmd)
    
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("Error starting CLI: %v", err)
        os.Exit(1)
    }
}

func addTask(title, description string) {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := proto.NewTodoServiceClient(conn)
    _, err = client.AddTask(context.Background(), &proto.TaskRequest{Title: title, Description: description})
    if err != nil {
        log.Fatalf("Error adding task: %v", err)
    }
    log.Println("Task added successfully")
}

func getTasks() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := proto.NewTodoServiceClient(conn)
    response, err := client.GetTasks(context.Background(), &proto.Empty{})
    if err != nil {
        log.Fatalf("Error getting tasks: %v", err)
    }
    
    for _, task := range response.Tasks {
        log.Printf("Task: %v, Completed: %v", task.Title, task.Completed)
    }
}
