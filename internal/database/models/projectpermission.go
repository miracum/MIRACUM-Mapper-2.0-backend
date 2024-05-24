package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type ProjectPermissionRole string

const (
	AdminRole ProjectPermissionRole = "reviewer"
	UserRole  ProjectPermissionRole = "projectOwner"
	GuestRole ProjectPermissionRole = "editor"
)

func (e *ProjectPermissionRole) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = ProjectPermissionRole(v)
	case string:
		*e = ProjectPermissionRole(v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (e ProjectPermissionRole) Value() (driver.Value, error) {
	return string(e), nil
}

// ProjectPermission defines model for ProjectPermission.
type ProjectPermission struct {
	Role      ProjectPermissionRole `gorm:"type:ProjectPermissionRole"`
	UserID    uuid.UUID             `gorm:"primarykey"`
	ProjectID uint32                `gorm:"primarykey"`
	User      User
}
