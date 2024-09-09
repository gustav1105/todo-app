package main

import (
    "os"

    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "todo-app/internal/commands"
)

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    rootCmd := &cobra.Command{Use: "todo-client"}

    // Register commands from the tasks and commands package
    rootCmd.AddCommand(
        commands.NewAddCommand(logger),
        commands.NewGetCommand(logger),
        commands.NewServeCommand(logger), // Add serve command here
    )

    if err := rootCmd.Execute(); err != nil {
        logger.Fatal("Error executing CLI command", zap.Error(err))
        os.Exit(1)
    }
}

