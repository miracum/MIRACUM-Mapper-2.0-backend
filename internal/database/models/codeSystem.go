package models

type CodeSystem struct {
	Model
	Uri             string
	Version         string
	Name            string
	Title           *string
	Description     *string
	Author          *string
	Concepts        []Concept `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystemRoles []CodeSystemRole
}
