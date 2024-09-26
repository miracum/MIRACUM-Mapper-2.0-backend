package models

import "time"

type Model struct {
	ID        uint32    `gorm:"primarykey"` // implicitly autoIncrement
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type ModelBigId struct {
	ID        uint64    `gorm:"primarykey"` // implicitly autoIncrement
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}
