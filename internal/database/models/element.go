package models

type Element struct {
	MappingID        uint64 `gorm:"primaryKey;index"`
	CodeSystemRoleID uint32 `gorm:"primaryKey"`
	ConceptID        *uint64
	Concept          Concept
	NextConceptID    *uint64
	NextConcept      Concept
}
