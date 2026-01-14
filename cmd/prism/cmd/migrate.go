package cmd

import (
	"github.com/cylonchau/prism/pkg/logger"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/cylonchau/prism/pkg/store"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Migrate creates or updates all database tables based on the model definitions.`,
	RunE:  runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	logger.Info("Starting database migration...")

	// Initialize database
	dbConfig := getStoreConfig()
	dbStore := store.GetInstance()

	if err := dbStore.Initialize(dbConfig); err != nil {
		logger.Error("Failed to initialize database", logger.Err(err))
		return err
	}
	defer dbStore.Close()

	logger.Info("Connected to database", logger.String("type", dbType))

	// Run migrations
	allModels := []interface{}{
		&models.ExecutionLock{},
		&models.ExecutionTask{},
		&models.Provider{},
		&models.Plugin{},
		&models.TerraformConfig{},
		&models.TerraformConfigMetadata{},
		&models.TerraformConfigParam{},
		&models.TerraformResource{},
		&models.TerraformResourceAttribute{},
		&models.TerraformResourceOutput{},
	}

	logger.Info("Starting model migration", logger.Int("count", len(allModels)))

	for i, model := range allModels {
		if err := dbStore.AutoMigrate(model); err != nil {
			logger.Error("Migration failed",
				logger.Int("model_index", i),
				logger.String("model_type", getModelName(model)),
				logger.Err(err))
			return err
		}
		logger.Info("Model migrated",
			logger.Int("progress", i+1),
			logger.Int("total", len(allModels)),
			logger.String("model", getModelName(model)))
	}

	logger.Info("Database migration completed successfully!")
	return nil
}

func getModelName(model interface{}) string {
	switch model.(type) {
	case *models.ExecutionLock:
		return "ExecutionLock"
	case *models.ExecutionTask:
		return "ExecutionTask"
	case *models.Provider:
		return "Provider"
	case *models.Plugin:
		return "Plugin"
	case *models.TerraformConfig:
		return "TerraformConfig"
	case *models.TerraformConfigMetadata:
		return "TerraformConfigMetadata"
	case *models.TerraformConfigParam:
		return "TerraformConfigParam"
	case *models.TerraformResource:
		return "TerraformResource"
	case *models.TerraformResourceAttribute:
		return "TerraformResourceAttribute"
	case *models.TerraformResourceOutput:
		return "TerraformResourceOutput"
	default:
		return "Unknown"
	}
}

func getStoreConfig() store.DatabaseConfig {
	var dbTypeEnum store.DBType
	switch dbType {
	case "mysql":
		dbTypeEnum = store.MySQL
	case "postgres", "postgresql":
		dbTypeEnum = store.PostgreSQL
	default:
		dbTypeEnum = store.SQLite
	}

	return store.DatabaseConfig{
		Type:     dbTypeEnum,
		Host:     dbHost,
		Port:     dbPort,
		Database: dbName,
		Username: dbUser,
		Password: dbPass,
		File:     dbFile,
	}
}
