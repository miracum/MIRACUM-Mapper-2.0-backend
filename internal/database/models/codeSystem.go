package models

type CodeSystemType string

const (
	GENERIC CodeSystemType = "GENERIC"
	LOINC   CodeSystemType = "LOINC"
)

type CodeSystem struct {
	Model
	Uri                string
	Name               string
	Type               CodeSystemType `gorm:"type:CodeSystemType;default:GENERIC"`
	Title              *string
	Description        *string
	Author             *string
	Concepts           []Concept `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystemRoles    []CodeSystemRole
	CodeSystemVersions []CodeSystemVersion `gorm:"constraint:OnDelete:CASCADE"`
}
