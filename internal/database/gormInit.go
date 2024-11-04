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

	if err := db.AutoMigrate(&models.CodeSystem{}, &models.Concept{}, &models.User{}, &models.Project{}, &models.Mapping{}, &models.Element{}, &models.CodeSystemRole{}, &models.ProjectPermission{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %v", err)
	}

	if err := setupFullTextSearchOld(db); err != nil {
		return nil, err
	}

	createTestData(db)

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
		"CREATE TYPE Status AS ENUM ('active', 'inactive', 'pending');",
		"CREATE TYPE CodeSystemRoleType AS ENUM ('source', 'target');",
		"CREATE TYPE ProjectPermissionRole AS ENUM ('reviewer', 'project_owner', 'editor');",
	}

	if err := executeSQLWithExceptionHandling(db, enumStatements); err != nil {
		return err
	}
	return nil
}

func setupFullTextSearch(db *gorm.DB) error {
	sqlStatements := []string{
		"CREATE EXTENSION IF NOT EXISTS pg_trgm;",
		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1
		        FROM   pg_class c
		        JOIN   pg_namespace n ON n.oid = c.relnamespace
		        WHERE  c.relname = 'idx_code_trigram'
		        AND    n.nspname = 'public'
		    ) THEN
		        EXECUTE 'CREATE INDEX idx_code_trigram ON concepts USING gin (code gin_trgm_ops);';
		    END IF;
		END
		$$;`,
		// `DO $$
		// BEGIN
		//     ALTER TABLE concepts ADD COLUMN IF NOT EXISTS display_search_vector tsvector
		// 	GENERATED ALWAYS AS (to_tsvector('english', display)) STORED;
		// 	CREATE INDEX IF NOT EXISTS idx_display_search_vector ON concepts USING gin (display_search_vector);
		// END
		// $$;`,
		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1
		        FROM   pg_class c
		        JOIN   pg_namespace n ON n.oid = c.relnamespace
		        WHERE  c.relname = 'idx_display_fulltext'
		        AND    n.nspname = 'public'
		    ) THEN
		        EXECUTE 'CREATE INDEX idx_display_fulltext ON concepts USING gin (to_tsvector(''english'', display));';
		    END IF;
		END
		$$;`,
		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1
		        FROM   pg_class c
		        JOIN   pg_namespace n ON n.oid = c.relnamespace
		        WHERE  c.relname = 'idx_code_lower'
		        AND    n.nspname = 'public'
		    ) THEN
		        EXECUTE 'CREATE INDEX idx_code_lower ON concepts (LOWER(code));';
		    END IF;
		END
		$$;`,
		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1
		        FROM   pg_class c
		        JOIN   pg_namespace n ON n.oid = c.relnamespace
		        WHERE  c.relname = 'idx_display_lower'
		        AND    n.nspname = 'public'
		    ) THEN
		        EXECUTE 'CREATE INDEX idx_display_lower ON concepts (LOWER(display));';
		    END IF;
		END
		$$;`,
	}

	if err := executeSQL(db, sqlStatements); err != nil {
		return err
	}
	return nil
}

func setupFullTextSearchOld(db *gorm.DB) error {
	sqlStatements := []string{
		"CREATE EXTENSION IF NOT EXISTS pg_trgm;",

		"UPDATE concepts SET display_search_vector = to_tsvector('english', display);",

		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1
		        FROM   pg_class c
		        JOIN   pg_namespace n ON n.oid = c.relnamespace
		        WHERE  c.relname = 'display_search_vector_idx'
		        AND    n.nspname = 'public'
		    ) THEN
		        EXECUTE 'CREATE INDEX display_search_vector_idx ON concepts USING gin(display_search_vector);';
		    END IF;
		END
		$$;`,

		`CREATE OR REPLACE FUNCTION concepts_display_trigger() RETURNS trigger AS $$
		    begin
		      new.display_search_vector := to_tsvector('english', new.display);
		      return new;
		    end
		$$ LANGUAGE plpgsql;`,

		`DO $$
		BEGIN
		    IF NOT EXISTS (
		        SELECT 1 FROM pg_trigger WHERE tgname = 'update_display_search_vector'
		    ) THEN
		        EXECUTE 'CREATE TRIGGER update_display_search_vector BEFORE INSERT OR UPDATE ON concepts FOR EACH ROW EXECUTE FUNCTION concepts_display_trigger();';
		    END IF;
		END
		$$;`,
		// 	`DO $$
		// BEGIN
		// 	IF NOT EXISTS (
		// 		SELECT 1
		// 		FROM   pg_class c
		// 		JOIN   pg_namespace n ON n.oid = c.relnamespace
		// 		WHERE  c.relname = 'idx_display_trgm'
		// 		AND    n.nspname = 'public'
		// 	) THEN
		// 		EXECUTE 'CREATE INDEX idx_display_trgm ON public.concepts USING gin (display gin_trgm_ops);';
		// 	END IF;
		// END
		// $$;`,

		`DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM   pg_class c
            JOIN   pg_namespace n ON n.oid = c.relnamespace
            WHERE  c.relname = 'idx_code_system_id'
            AND    n.nspname = 'public'
        ) THEN
            EXECUTE 'CREATE INDEX idx_code_system_id ON concepts (code_system_id);';
        END IF;
    END
    $$;`,

		`DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM   pg_class c
            JOIN   pg_namespace n ON n.oid = c.relnamespace
            WHERE  c.relname = 'idx_code_lower'
            AND    n.nspname = 'public'
        ) THEN
            EXECUTE 'CREATE INDEX idx_code_lower ON concepts (LOWER(code));';
        END IF;
    END
    $$;`,
		// 	`DO $$
		// BEGIN
		// 	IF NOT EXISTS (
		//     	SELECT 1
		//     	FROM   pg_class c
		//     	JOIN   pg_namespace n ON n.oid = c.relnamespace
		//     	WHERE  c.relname = 'display_trgm_idx'
		//     	AND    n.nspname = 'public'
		// 	) THEN
		//     	EXECUTE 'CREATE INDEX display_trgm_idx ON concepts USING gin (display gin_trgm_ops);';
		// 	END IF;
		// END
		// $$;`,
	}

	if err := executeSQL(db, sqlStatements); err != nil {
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

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // TODO remove for production. Maybe toggle this with debug flag
	})
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
		Id:       uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), // Generate a new UUID
		FullName: "Test User",
		// LogName:     "testuser",
		// Affiliation: "Test Affiliation",
		UserName: "testuser",
		Email:    "test.user@123.com",
		// Initialize ProjectPermissions if needed
	}

	// Save the test user to the database
	gormDB.FirstOrCreate(&testUser, models.User{Id: testUser.Id})

	// Create a test user
	testUser2 := models.User{
		Id:       uuid.MustParse("b1ffcd99-9c0b-4ef8-bb6d-6bb9bd380a20"), // Generate a new UUID
		FullName: "Test User 2",
		// LogName:     "testuse2",
		// Affiliation: "Test Affiliatio2",
		UserName: "testuser2",
		Email:    "test.user2@123.com",
		// Initialize ProjectPermissions if needed
	}

	// Save the test user 2 to the database
	gormDB.FirstOrCreate(&testUser2, models.User{Id: testUser2.Id})

	// description := "Example Code System 1"
	// codeSystem := models.CodeSystem{
	// 	Uri:             "http://example.com/codesystem",
	// 	Version:         "1.0",
	// 	Name:            "Example Code System",
	// 	Description:     &description,
	// 	Author:          nil,
	// 	Concepts:        nil,
	// 	CodeSystemRoles: nil,
	// }

	// gormDB.FirstOrCreate(&codeSystem)

	// description2 := "Example Code System 2"

	// codeSystem2 := models.CodeSystem{
	// 	Model:       models.Model{ID: 2},
	// 	Uri:         "http://example.com/codesystem2",
	// 	Version:     "1.0",
	// 	Name:        "Example Code System 2",
	// 	Description: &description2,
	// }

	// gormDB.FirstOrCreate(&codeSystem2)

	// // create concept for code system 1
	// concept1 := models.Concept{
	// 	CodeSystemID: codeSystem.ID,
	// 	Code:         "1",
	// 	Display:      "Concept 1",
	// }

	// gormDB.FirstOrCreate(&concept1)

	// // create concept for code system 2

	// concept2 := models.Concept{
	// 	ID:           2,
	// 	CodeSystemID: codeSystem2.ID,
	// 	Code:         "2",
	// 	Display:      "Concept 2",
	// }

	// gormDB.FirstOrCreate(&concept2)

	// var sampleLoincCodes = []struct {
	// 	Code    string
	// 	Meaning string
	// }{
	// 	// artificially generated test data which is not valid
	// 	{"1000-9", "Hemoglobin A1c/Hemoglobin.total in Blood"},
	// 	{"1001-7", "Hemoglobin A1c in Blood by HPLC"},
	// 	{"1002-5", "Glucose level in Blood"},
	// 	{"1003-3", "Potassium level in Serum or Plasma"},
	// 	{"1004-1", "Sodium level in Serum or Plasma"},
	// 	{"1005-8", "Cholesterol in Serum or Plasma"},
	// 	{"1006-6", "Triglycerides in Serum or Plasma"},
	// 	{"1007-4", "HDL Cholesterol in Serum or Plasma"},
	// 	{"1008-2", "LDL Cholesterol in Serum or Plasma"},
	// 	{"1009-0", "Creatinine level in Serum or Plasma"},
	// 	{"1010-8", "Urea Nitrogen level in Blood"},
	// 	{"1011-6", "Protein total in Serum or Plasma"},
	// 	{"1012-4", "Albumin level in Serum or Plasma"},
	// 	{"1013-2", "Calcium level in Serum or Plasma"},
	// 	{"1014-0", "Phosphorus level in Serum or Plasma"},
	// 	{"1015-7", "Iron level in Serum or Plasma"},
	// 	{"1016-5", "Bilirubin total in Serum or Plasma"},
	// 	{"1017-3", "Alkaline Phosphatase level in Serum or Plasma"},
	// 	{"1018-1", "Alanine Aminotransferase level in Serum or Plasma"},
	// 	{"1019-9", "Aspartate Aminotransferase level in Serum or Plasma"},
	// 	{"1020-7", "Gamma Glutamyl Transferase level in Serum or Plasma"},
	// 	{"1021-5", "Blood Uric acid level"},
	// 	{"1022-3", "C-Reactive Protein level in Serum or Plasma"},
	// 	{"1023-1", "Thyroid Stimulating Hormone level in Serum or Plasma"},
	// 	{"1024-9", "Free T4 level in Serum or Plasma"},
	// 	{"1025-6", "Total T3 level in Serum or Plasma"},
	// 	{"1026-4", "Prostate Specific Antigen in Serum or Plasma"},
	// 	{"1027-2", "Rheumatoid Factor in Serum or Plasma"},
	// 	{"1028-0", "Hepatitis C Virus Antibody in Serum or Plasma"},
	// 	{"1029-8", "HIV 1+2 Antibodies in Serum or Plasma"},
	// 	{"1030-6", "Hemoglobin level in Blood"},
	// 	{"1031-4", "Erythrocyte Sedimentation Rate"},
	// 	{"1032-2", "White Blood Cell count in Blood"},
	// 	{"1033-0", "Platelet count in Blood"},
	// 	{"1034-8", "Red Blood Cell count in Blood"},
	// 	{"1035-5", "Mean Corpuscular Volume"},
	// 	{"1036-3", "Mean Corpuscular Hemoglobin"},
	// 	{"1037-1", "Mean Corpuscular Hemoglobin Concentration"},
	// 	{"1038-9", "Red Cell Distribution Width"},
	// 	{"1039-7", "Neutrophils.auto count in Blood"},
	// 	{"1040-5", "Lymphocytes.auto count in Blood"},
	// 	{"1041-3", "Monocytes.auto count in Blood"},
	// 	{"1042-1", "Eosinophils.auto count in Blood"},
	// 	{"1043-9", "Basophils.auto count in Blood"},
	// 	{"1044-7", "Blood Type in Blood"},
	// 	{"1045-4", "Rh(D) Typing in Blood"},
	// 	{"1046-2", "Antibody Screen in Blood"},
	// 	{"1047-0", "Direct Antiglobulin Test in Blood"},
	// 	{"1048-8", "Indirect Antiglobulin Test in Blood"},
	// 	{"1049-6", "Blood Culture for Bacteria"},
	// }

	// for i, sampleLoincCode := range sampleLoincCodes {
	// 	concept := models.Concept{
	// 		ID:           uint64(i + 3),
	// 		CodeSystemID: 2,
	// 		Code:         sampleLoincCode.Code,
	// 		Display:      sampleLoincCode.Meaning,
	// 	}

	// 	gormDB.FirstOrCreate(&concept)
	// }

	// words := []string{"nein", "awesome", "42", "Pills", "Nina", "Loinc", "word7", "word8", "word9", "Very", "and", "for", "some", "Boomer", "Go", "hallo", "blub", "egal", "buch", "katze", "hund", "henrik", "computer", "geben", "halten", "tastatur", "applaudieren", "kontrolle", "schlüssel", "schlange", "schlafen", "schlüsselbund"}

	// rand.Seed(1)

	// count := 2
	// for j := 0; j < 1000; j++ {
	// 	concepts := make([]models.Concept, 100)
	// 	for i := 0; i < 100; i++ {
	// 		//code := fmt.Sprintf("%d", rand.Intn(10000)) // generate a random code
	// 		// select 3 random words from the list
	// 		// meaning := words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))] + " " + words[rand.Intn(len(words))]
	// 		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 		s := make([]rune, 20)
	// 		for i := range s {
	// 			s[i] = letters[rand.Intn(len(letters))]
	// 		}
	// 		meaning := string(s)

	// 		count = count + rand.Intn(10) + 1
	// 		concepts[i] = models.Concept{
	// 			ID:           uint64(count),
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
	// 			ID:           uint64(count),
	// 			CodeSystemID: 2,
	// 			Code:         fmt.Sprint(count),
	// 			Display:      meaning,
	// 		}
	// 	}

	// 	gormDB.Create(&concepts)
	// }
}
