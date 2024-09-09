package commands

import (
    "fmt"
    "go.uber.org/zap"
    "github.com/spf13/cobra"
    "io/ioutil"
)

// NewServeCommand creates the "serve" command to set database credentials
func NewServeCommand(logger *zap.Logger) *cobra.Command {
    var dbUser, dbPassword, dbHost, dbName, dbPort string

    serveCmd := &cobra.Command{
        Use:   "serve",
        Short: "Set database credentials for the gRPC server",
        Run: func(cmd *cobra.Command, args []string) {
            // Check if all necessary flags are provided
            if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
                logger.Fatal("All database credentials must be provided")
            }

            // Prepare the content for the .env file
            envContent := fmt.Sprintf(`DB_USER=%s
DB_PASSWORD=%s
DB_HOST=%s
DB_PORT=%s
DB_NAME=%s
`, dbUser, dbPassword, dbHost, dbPort, dbName)

            // Write the content to a .env file in the root directory (or wherever Docker expects it)
            err := ioutil.WriteFile(".env", []byte(envContent), 0644)
            if err != nil {
                logger.Fatal("Failed to write .env file", zap.Error(err))
            }

            // Inform the user that the credentials were set successfully
            logger.Info("Database credentials saved to .env file.")
            logger.Info("Run 'docker-compose up --build' to start the server with these credentials.")
        },
    }

    // Define flags to accept database credentials via CLI
    serveCmd.Flags().StringVar(&dbUser, "db-user", "", "Database user")
    serveCmd.Flags().StringVar(&dbPassword, "db-password", "", "Database password")
    serveCmd.Flags().StringVar(&dbHost, "db-host", "localhost", "Database host")
    serveCmd.Flags().StringVar(&dbPort, "db-port", "3306", "Database port")
    serveCmd.Flags().StringVar(&dbName, "db-name", "", "Database name")

    // Mark the required flags
    serveCmd.MarkFlagRequired("db-user")
    serveCmd.MarkFlagRequired("db-password")
    serveCmd.MarkFlagRequired("db-host")
    serveCmd.MarkFlagRequired("db-port")
    serveCmd.MarkFlagRequired("db-name")

    return serveCmd
}

