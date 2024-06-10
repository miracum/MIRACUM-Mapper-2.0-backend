package database

import (
	"fmt"
	"log"
	"miracummapper/internal/config"
	"miracummapper/internal/database/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormConnection(config *config.Config) *gorm.DB {
	db, err := getGormConnection(config)
	if err != nil {
		panic(err)
	}
	return db
}

func getGormConnection(config *config.Config) (*gorm.DB, error) {
	// db, err := createGormConnection(config)
	// if err != nil {
	// 	return nil, err
	// }

	db, err := connectToDb(config)
	if err != nil {
		return nil, err
	}

	initEnums(db)

	db.AutoMigrate(&models.CodeSystem{}, &models.Concept{}, &models.User{}, &models.Project{}, &models.Mapping{}, &models.Element{}, &models.CodeSystemRole{}, &models.ProjectPermission{})

	createTestData(db)

	return db, nil
}

func initEnums(db *gorm.DB) {
	db.Exec(`
	DO $$ BEGIN
		CREATE TYPE Equivalence AS ENUM ('related-to', 'equivalent', 'source-is-narrower-than-target', 'source-is-broader-than-target', 'not-related');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	db.Exec(`
	DO $$ BEGIN
		CREATE TYPE Status AS ENUM ('active', 'inactive', 'pending');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	db.Exec(`
	DO $$ BEGIN
		CREATE TYPE CodeSystemRoleType AS ENUM ('source', 'target');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	db.Exec(`
	DO $$ BEGIN
		CREATE TYPE ProjectPermissionRole AS ENUM ('reviewer', 'projectOwner', 'editor');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)
}

func createGormConnection(config *config.Config) (*gorm.DB, error) {
	DSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		// Conn: db,
		DSN: DSN,
	}), &gorm.Config{}) // Use db.Driver() instead of db.DriverName()
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

// func CreateEnum(db *gorm.DB, enumType interface{}) error {
// 	t := reflect.TypeOf(enumType)
// 	if t.Kind() != reflect.String {
// 		return fmt.Errorf("enumType must be a string")
// 	}

// 	enumName := t.Name()
// 	values := []string{}

// 	v := reflect.ValueOf(enumType)
// 	for i := 0; i < v.NumField(); i++ {
// 		values = append(values, fmt.Sprintf("'%s'", v.Field(i).String()))
// 	}

// 	query := fmt.Sprintf(`
//     DO $$ BEGIN
//         CREATE TYPE %s AS ENUM (%s);
//     EXCEPTION
//         WHEN duplicate_object THEN null;
//     END $$;`, enumName, strings.Join(values, ", "))

// 	return db.Exec(query).Error
// }

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
	gormDB.FirstOrCreate(&testUser, models.User{Id: testUser.Id})

	// Create a test user
	testUser2 := models.User{
		Id:          uuid.MustParse("b1ffcd99-9c0b-4ef8-bb6d-6bb9bd380a20"), // Generate a new UUID
		UserName:    "Test Use2",
		LogName:     "testuse2",
		Affiliation: "Test Affiliatio2",
		// Initialize ProjectPermissions if needed
	}

	// Save the test user 2 to the database
	gormDB.FirstOrCreate(&testUser2, models.User{Id: testUser2.Id})

	description := "Example Code System 1"
	codeSystem := models.CodeSystem{
		Uri:             "http://example.com/codesystem",
		Version:         "1.0",
		Name:            "Example Code System",
		Description:     &description,
		Author:          nil,
		Concepts:        nil,
		CodeSystemRoles: nil,
	}

	gormDB.FirstOrCreate(&codeSystem)

	description2 := "Example Code System 2"

	codeSystem2 := models.CodeSystem{
		Model:       models.Model{ID: 2},
		Uri:         "http://example.com/codesystem2",
		Version:     "1.0",
		Name:        "Example Code System 2",
		Description: &description2,
	}

	gormDB.FirstOrCreate(&codeSystem2)

	// create concept for code system 1
	concept1 := models.Concept{
		CodeSystemID: codeSystem.ID,
		Code:         "1",
		Display:      "Concept 1",
	}

	gormDB.FirstOrCreate(&concept1)

	// create concept for code system 2

	concept2 := models.Concept{
		ModelBigId:   models.ModelBigId{ID: 2},
		CodeSystemID: codeSystem2.ID,
		Code:         "2",
		Display:      "Concept 2",
	}

	gormDB.FirstOrCreate(&concept2)

	// words := []string{"nein", "awesome", "42", "Pills", "Nina", "Loinc", "word7", "word8", "word9", "Very", "and", "for", "some", "Boomer", "Go", "hallo", "blub", "egal", "buch", "katze", "hund", "henrik", "computer", "geben", "halten", "tastatur", "applaudieren", "kontrolle", "schlüssel", "schlange", "schlafen", "schlüsselbund"}

	// rand.Seed(1)

	// count := 2
	// for j := 0; j < 1000; j++ {
	// 	concepts := make([]models.Concept, 1000)
	// 	for i := 0; i < 1000; i++ {
	// 		//code := fmt.Sprintf("%d", rand.Intn(10000)) // generate a random code
	// 		// select 3 random words from the list
	// 		meaning := words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))]

	// 		count = count + rand.Intn(10) + 1
	// 		concepts[i] = models.Concept{
	// 			ModelBigId:   models.ModelBigId{ID: uint64(count)},
	// 			CodeSystemID: 1,
	// 			Code:         fmt.Sprint(count),
	// 			Display:      meaning,
	// 		}
	// 	}

	// 	gormDB.Create(&concepts)
	// }

	// for j := 0; j < 1; j++ {
	// 	concepts := make([]models.Concept, 100)
	// 	for i := 0; i < 100; i++ {
	// 		//code := fmt.Sprintf("%d", rand.Intn(10000)) // generate a random code
	// 		// select 3 random words from the list
	// 		meaning := words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))]

	// 		count = count + rand.Intn(10) + 1
	// 		concepts[i] = models.Concept{
	// 			ModelBigId:   models.ModelBigId{ID: uint64(count)},
	// 			CodeSystemID: 2,
	// 			Code:         fmt.Sprint(count),
	// 			Display:      meaning,
	// 		}
	// 	}

	// 	gormDB.Create(&concepts)
	// }
}
