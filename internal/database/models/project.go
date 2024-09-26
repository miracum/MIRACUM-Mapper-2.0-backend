package models

type Project struct {
	Model
	Description         string
	EquivalenceRequired bool
	Name                string
	StatusRequired      bool
	Version             string
	Mappings            []Mapping           `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystemRoles     []CodeSystemRole    `gorm:"constraint:OnDelete:CASCADE"`
	Permissions         []ProjectPermission `gorm:"constraint:OnDelete:CASCADE"`
}
