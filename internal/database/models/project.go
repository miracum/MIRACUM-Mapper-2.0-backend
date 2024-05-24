package models

// Project defines model for Project.
type Project struct {
	Model
	Description         string
	EquivalenceRequired bool
	Name                string
	StatusRequired      bool
	Version             string
	Mappings            []Mapping
	CodeSystemRoles     []CodeSystemRole
	Permissions         []ProjectPermission
}
