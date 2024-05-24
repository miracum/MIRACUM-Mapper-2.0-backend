package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type ProjectPermissionRole string

const (
	AdminRole ProjectPermissionRole = "reviewer"
	UserRole  ProjectPermissionRole = "projectOwner"
	GuestRole ProjectPermissionRole = "editor"
)

func (e *ProjectPermissionRole) Scan(value interface{}) error {
	*e = ProjectPermissionRole(value.([]byte))
	return nil
}

func (e ProjectPermissionRole) Value() (driver.Value, error) {
	return string(e), nil
}

// ProjectPermission defines model for ProjectPermission.
type ProjectPermission struct {
	Role      ProjectPermissionRole `gorm:"type:ProjectPermissionRole"`
	UserID    uuid.UUID
	ProjectID uint32
}
