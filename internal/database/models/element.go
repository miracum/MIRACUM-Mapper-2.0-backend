package models

// Element defines model for Element.
type Element struct {
	MappingID        uint64 `gorm:"primaryKey"`
	CodeSystemRoleID uint32 `gorm:"primaryKey"`
	ConceptID        *uint64
	Concept          Concept
}
