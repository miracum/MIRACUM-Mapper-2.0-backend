package models

type Element struct {
	MappingID        uint64  `gorm:"primaryKey"`
	CodeSystemRoleID uint32  `gorm:"primaryKey"`
	ConceptID        *uint64 `gorm:"index"`
	Concept          Concept
}
