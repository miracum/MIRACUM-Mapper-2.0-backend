package models

// Concept defines model for Concept.
type Concept struct {
	ModelBigId
	Code         string
	Display      string
	CodeSystemID uint32
	Elements     []Element
	CodeSystem   CodeSystem
}
