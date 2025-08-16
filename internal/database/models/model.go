package models

import "time"

type Model struct {
	ID        int32     `gorm:"primarykey;type:integer"` // implicitly autoIncrement
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}
