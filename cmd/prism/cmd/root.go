package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	dbType  string
	dbHost  string
	dbPort  int
	dbName  string
	dbUser  string
	dbPass  string
	dbFile  string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "prism",
	Short: "Prism - Terraform Execution Platform",
	Long: `Prism is a platform for managing and executing Terraform configurations
with support for multi-provider, real-time output streaming, and distributed locking.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

	// Database flags
	rootCmd.PersistentFlags().StringVar(&dbType, "db-type", "sqlite", "database type: mysql, postgres, sqlite")
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", "localhost", "database host")
	rootCmd.PersistentFlags().IntVar(&dbPort, "db-port", 3306, "database port")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", "prism", "database name")
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", "root", "database user")
	rootCmd.PersistentFlags().StringVar(&dbPass, "db-pass", "", "database password")
	rootCmd.PersistentFlags().StringVar(&dbFile, "db-file", "prism", "sqlite database file path (without .db extension)")
}

// GetDBConfig returns database config from flags
func GetDBConfig() map[string]interface{} {
	return map[string]interface{}{
		"type": dbType,
		"host": dbHost,
		"port": dbPort,
		"name": dbName,
		"user": dbUser,
		"pass": dbPass,
		"file": dbFile,
	}
}

func exitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}
