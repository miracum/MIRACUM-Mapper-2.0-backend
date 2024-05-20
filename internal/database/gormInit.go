package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitGorm(db *sql.DB) {

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{}) // Use db.Driver() instead of db.DriverName()
	if err != nil {
		log.Fatalf("Failed to create GORM database connection: %v", err)
		return
	}

	gormDB.AutoMigrate(&CodeSystem{})
	gormDB.AutoMigrate(&Concept{})
	gormDB.AutoMigrate(&Mapping{})
	gormDB.AutoMigrate(&Element{})
	gormDB.AutoMigrate(&CodeSystemRole{})
	gormDB.AutoMigrate(&Project{})
	gormDB.AutoMigrate(&ProjectPermission{})
}
