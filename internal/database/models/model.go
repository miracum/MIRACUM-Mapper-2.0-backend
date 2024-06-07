package models

import "time"

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
//
//	type User struct {
//	  models.Model
//	}
type Model struct {
	ID        uint32 `gorm:"primarykey"` // implicitly autoIncrement
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
}

type ModelBigId struct {
	ID        uint64 `gorm:"primarykey"` // implicitly autoIncrement
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt DeletedAt `gorm:"index"`
}
