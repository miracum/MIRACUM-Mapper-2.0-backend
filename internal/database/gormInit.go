package database

import (
	"fmt"
	"log"
	"miracummapper/internal/config"
	"miracummapper/internal/database/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormConnection(config *config.Config) *gorm.DB {
	db, err := getGormConnection(config)
	if err != nil {
		panic(err)
	}
	return db
}

func getGormConnection(config *config.Config) (*gorm.DB, error) {

	db, err := connectToDb(config)
	if err != nil {
		return nil, err
	}

	if err = initEnums(db); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.CodeSystem{}, &models.Concept{}, &models.User{}, &models.Project{}, &models.Mapping{}, &models.Element{}, &models.CodeSystemRole{}, &models.ProjectPermission{}, &models.CodeSystemVersion{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %v", err)
	}

	// Create concept index (cant be created with Gorm annotations)
	if err := models.CreateConceptIndex(db); err != nil {
		log.Fatal(err)
	}

	// createTestData(db)

	return db, nil
}

// executeSQLWithExceptionHandling wraps each SQL statement with exception handling and executes them
func executeSQLWithExceptionHandling(db *gorm.DB, sqlStatements []string) error {
	for _, sqlStatement := range sqlStatements {
		wrappedSQL := `
        DO $$ BEGIN
            ` + sqlStatement + `
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;
        `
		if err := db.Exec(wrappedSQL).Error; err != nil {
			return err
		}
	}
	return nil
}

func executeSQL(db *gorm.DB, sqlStatements []string) error {
	for _, sqlStatement := range sqlStatements {
		if err := db.Exec(sqlStatement).Error; err != nil {
			return err
		}
	}
	return nil
}

func initEnums(db *gorm.DB) error {
	enumStatements := []string{
		"CREATE TYPE Equivalence AS ENUM ('related-to', 'equivalent', 'source-is-narrower-than-target', 'source-is-broader-than-target', 'not-related');",
		"CREATE TYPE MappingStatus AS ENUM ('active', 'inactive', 'pending');",
		"CREATE TYPE CodeSystemRoleType AS ENUM ('source', 'target');",
		"CREATE TYPE ProjectPermissionRole AS ENUM ('reviewer', 'project_owner', 'editor');",
		"CREATE TYPE ConceptStatus AS ENUM ('active', 'trial', 'deprecated', 'discouraged');",
		"CREATE TYPE CodeSystemType AS ENUM ('GENERIC', 'LOINC');",
	}

	if err := executeSQLWithExceptionHandling(db, enumStatements); err != nil {
		return err
	}
	return nil
}

func createGormConnection(config *config.Config) (*gorm.DB, error) {
	DSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName)

	var gormConfig *gorm.Config
	if config.File.Debug {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		gormConfig = &gorm.Config{}
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create GORM database connection: %v", err)
	}
	return gormDB, nil
}

func connectToDb(config *config.Config) (*gorm.DB, error) {
	for i := 0; i < config.File.DatabaseConfig.Retry; i++ {
		db, err := createGormConnection(config)
		if err == nil {
			err = pingGormDB(db)
			if err == nil {
				log.Printf("Successfully connected to database: %s", config.Env.DBName)
				return db, nil
			}
		}
		log.Printf("Failed to connect to database: %s. Retrying in %d seconds", config.Env.DBName, config.File.DatabaseConfig.Sleep)
		if i != config.File.DatabaseConfig.Retry-1 {
			time.Sleep(time.Duration(config.File.DatabaseConfig.Sleep) * time.Second)
		}
	}
	return nil, fmt.Errorf("failed to connect to database %s after %d retries", config.Env.DBName, config.File.DatabaseConfig.Retry)
}

func pingGormDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
