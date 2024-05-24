package models

// Element defines model for Element.
type Element struct {
	MappingID        uint32 `gorm:"primaryKey"`
	CodeSystemRoleID uint32 `gorm:"primaryKey"`
	ConceptID        *uint32
}
