package database

import (
	"gorm.io/gorm"
)

type CodeSystem struct {
	gorm.Model
	Uri         string
	Version     string
	Name        string
	Title       *string
	Description *string
	Author      *string
	Concepts    *[]Concept `gorm:"foreignKey:CodeSystemID"`
}

// Concept defines model for Concept.
type Concept struct {
	Id           uint `gorm:"primaryKey"`
	CodeSystemID uint
	Code         string
	Display      string
	Elements     *[]Element `gorm:"foreignKey:ConceptId"`
}

type Equivalence string

const (
	RelatedTo                  Equivalence = "related-to"
	Equivalent                 Equivalence = "equivalent"
	SourceIsNarrowerThanTarget Equivalence = "source-is-narrower-than-target"
	SourceIsBroaderThanTarget  Equivalence = "source-is-broader-than-target"
	NotRElated                 Equivalence = "not-related"
)

// Mapping defines model for Mapping.
type Mapping struct {
	gorm.Model
	ProjectId   uint
	Equivalence Equivalence `gorm:"type:enum('related-to', 'equivalent', 'source-is-narrower-than-target', 'source-is-broader-than-target', 'not-related')"`
	Status      Status      `gorm:"type:enum('Active', 'Inactive', 'Pending')"`
	Comment     *string
	Elements    []Element `gorm:"foreignKey:MappingId"`
}

// Element defines model for Element.
type Element struct {
	MappingId        uint `gorm:"primaryKey"`
	CodeSystemRoleId uint `gorm:"primaryKey"`
	ConceptId        uint
}

type Status string

const (
	Active   Status = "Active"
	Inactive Status = "Inactive"
	Pending  Status = "Pending"
)

type ProjectPermissionRole string

const (
	Admin ProjectPermissionRole = "Reviewer"
	User  ProjectPermissionRole = "ProjectOwner"
	Guest ProjectPermissionRole = "Editor"
)

// func (e *ProjectPermissionRole) Scan(value interface{}) error {
// 	str, ok := value.(string)
// 	if !ok {
// 		return errors.New("Scan source was not string")
// 	}

// 	*e = ProjectPermissionRole(str)
// 	return nil
// }

// func (e ProjectPermissionRole) Value() (driver.Value, error) {
// 	return string(e), nil
// }

type CodeSystemRoleType string

const (
	Source CodeSystemRoleType = "Source"
	Target CodeSystemRoleType = "Target"
)

// CodeSystemRole defines model for CodeSystemRole.
type CodeSystemRole struct {
	Id           uint `gorm:"primarykey"`
	Name         string
	Position     int32
	Type         CodeSystemRoleType `gorm:"type:enum('Source', 'Target')"`
	CodeSystemID *int32             `gorm:"foreignKey:SystemID;references:Id"`
	// CodeSystem   uint         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Elements *[]Element `gorm:"foreignKey:CodeSystemRoleId"`
}

// Project defines model for Project.
type Project struct {
	gorm.Model
	Description         string
	EquivalenceRequired bool
	Modified            *string
	Name                string
	StatusRequired      bool
	Version             string
	Mappings            []Mapping           `gorm:"foreignKey:ProjectId"`
	CodeSystemRoles     []CodeSystemRole    `gorm:"foreignKey:ProjectId"`
	Permissions         []ProjectPermission `gorm:"foreignKey:ProjectId"`
}

// ProjectPermission defines model for ProjectPermission.
type ProjectPermission struct {
	Role      ProjectPermissionRole `gorm:"type:enum('Reviewer', 'ProjectOwner', 'Editor')"`
	UserId    string                `gorm:"primaryKey;foreignKey:UserId"`
	ProjectId string                `gorm:"primaryKey;foreignKey:ProjectId"`
}
