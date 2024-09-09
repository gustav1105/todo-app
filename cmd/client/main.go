package main

import (
    "os"

    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "todo-app/internal/tasks"
)

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    rootCmd := &cobra.Command{Use: "todo-client"}

    // Register commands from the tasks package
    rootCmd.AddCommand(tasks.NewAddCommand(logger), tasks.NewGetCommand(logger))

    if err := rootCmd.Execute(); err != nil {
        logger.Fatal("Error executing CLI command", zap.Error(err))
        os.Exit(1)
    }
}

