package models

// Concept defines model for Concept.
type Concept struct {
	Model
	Code         string
	Display      string
	CodeSystemID uint32
	Elements     []Element
}
