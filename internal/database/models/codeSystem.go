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
	Concepts           []Concept           `gorm:"constraint:OnDelete:CASCADE"` // not preloaded on get
	CodeSystemRoles    []CodeSystemRole    // not preloaded on get
	CodeSystemVersions []CodeSystemVersion `gorm:"constraint:OnDelete:CASCADE"`
}
