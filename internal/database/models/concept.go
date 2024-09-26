package models

import "gorm.io/gorm"

type Concept struct {
	ID                  uint64 `gorm:"primarykey"`
	Code                string `gorm:"index,gin"`
	Display             string `gorm:"index,gin"`
	CodeSystemID        uint32 `gorm:"index"`
	Elements            []Element
	CodeSystem          CodeSystem
	DisplaySearchVector string `gorm:"type:tsvector"` // Correctly map as tsvector for PostgreSQL
}

func (c *Concept) BeforeSave(tx *gorm.DB) (err error) {
	tx.Exec("UPDATE concepts SET display_search_vector = to_tsvector(?) WHERE id = ?", c.Display, c.ID)
	return
}
