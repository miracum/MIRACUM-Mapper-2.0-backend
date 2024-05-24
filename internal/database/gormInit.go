package database

import (
	"fmt"
	"log"
	"miracummapper/internal/config"
	"miracummapper/internal/database/models"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGorm(config *config.Config) {

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		// Conn: db,
		DSN: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Env.DBHost,
			config.Env.DBPort,
			config.Env.DBUser,
			config.Env.DBPassword,
			config.Env.DBName),
	}), &gorm.Config{}) // Use db.Driver() instead of db.DriverName()
	if err != nil {
		log.Fatalf("Failed to create GORM database connection: %v", err)
		return
	}

	// Assuming gormDB is your *gorm.DB instance
	// enums := []interface{}{
	// 	ProjectPermissionRole(""),
	// 	CodeSystemRoleType(""),
	// 	Status(""),
	// 	Equivalence(""),
	// }

	// for _, enum := range enums {
	// 	err := CreateEnum(gormDB, enum)
	// 	if err != nil {
	// 		log.Fatalf("Failed to create enum for %T: %v", enum, err)
	// 	}
	// }

	gormDB.Exec(`
	DO $$ BEGIN
		CREATE TYPE Equivalence AS ENUM ('related-to', 'equivalent', 'source-is-narrower-than-target', 'source-is-broader-than-target', 'not-related');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	gormDB.Exec(`
	DO $$ BEGIN
		CREATE TYPE Status AS ENUM ('active', 'inactive', 'pending');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	gormDB.Exec(`
	DO $$ BEGIN
		CREATE TYPE CodeSystemRoleType AS ENUM ('source', 'target');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	gormDB.Exec(`
	DO $$ BEGIN
		CREATE TYPE ProjectPermissionRole AS ENUM ('reviewer', 'projectOwner', 'editor');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	gormDB.AutoMigrate(&models.CodeSystem{})
	gormDB.AutoMigrate(&models.Concept{})
	gormDB.AutoMigrate(&models.Mapping{})
	gormDB.AutoMigrate(&models.Element{})
	gormDB.AutoMigrate(&models.CodeSystemRole{})
	gormDB.AutoMigrate(&models.Project{})
	gormDB.AutoMigrate(&models.ProjectPermission{})
	gormDB.AutoMigrate(&models.User{})

	createTestData(gormDB)
}

func CreateEnum(db *gorm.DB, enumType interface{}) error {
	t := reflect.TypeOf(enumType)
	if t.Kind() != reflect.String {
		return fmt.Errorf("enumType must be a string")
	}

	enumName := t.Name()
	values := []string{}

	v := reflect.ValueOf(enumType)
	for i := 0; i < v.NumField(); i++ {
		values = append(values, fmt.Sprintf("'%s'", v.Field(i).String()))
	}

	query := fmt.Sprintf(`
    DO $$ BEGIN
        CREATE TYPE %s AS ENUM (%s);
    EXCEPTION
        WHEN duplicate_object THEN null;
    END $$;`, enumName, strings.Join(values, ", "))

	return db.Exec(query).Error
}

func createTestData(gormDB *gorm.DB) {
	// Create a test user
	testUser := models.User{
		Id:          uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), // Generate a new UUID
		UserName:    "Test User",
		LogName:     "testuser",
		Affiliation: "Test Affiliation",
		// Initialize ProjectPermissions if needed
	}

	// Save the test user to the database
	result := gormDB.Create(&testUser)
	if result.Error != nil {
		log.Fatalf("Failed to create test user: %v", result.Error)
	}
}
