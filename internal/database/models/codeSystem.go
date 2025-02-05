package models

type CodeSystem struct {
	Model
	Uri                string
	Name               string
	Title              *string
	Description        *string
	Author             *string
	Concepts           []Concept `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystemRoles    []CodeSystemRole
	CodeSystemVersions []CodeSystemVersion `gorm:"constraint:OnDelete:CASCADE"`
}
