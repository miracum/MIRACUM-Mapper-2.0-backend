package models

type Element struct {
	MappingID        int32  `gorm:"primaryKey;index;type:integer"`
	CodeSystemRoleID int32  `gorm:"primaryKey;type:integer"`
	ConceptID        *int32 `gorm:"type:integer"`
	Concept          Concept
	NextConceptID    *int32 `gorm:"type:integer"`
	NextConcept      Concept
}
