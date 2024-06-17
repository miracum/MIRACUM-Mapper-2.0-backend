package models

// Corrected import for datatypes

// Concept defines model for Concept.
type Concept struct {
	ID                  uint64 `gorm:"primarykey"`
	Code                string
	Display             string
	CodeSystemID        uint32
	Elements            []Element
	CodeSystem          CodeSystem
	DisplaySearchVector string `gorm:"type:tsvector"` // Correctly map as tsvector for PostgreSQL
}
