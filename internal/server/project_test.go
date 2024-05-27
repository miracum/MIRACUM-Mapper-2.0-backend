package server

// import (
// 	"log"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn:       db,
// 		DriverName: "postgres",
// 	}), &gorm.Config{})

// 	if err != nil {
// 		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
// 	}

// 	return gormDB, mock
// }

// func TestPostProject(t *testing.T) {
// 	db, mock := NewMockDB()
// 	defer db.Close()

// 	// mock.ExpectBegin()
// 	// mock.ExpectExec(`INSERT INTO "projects"`).
// 	// 	WithArgs("
// }
