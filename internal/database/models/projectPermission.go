package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type ProjectPermissionRole string

const (
	ReviewerRole     ProjectPermissionRole = "reviewer"
	ProjectOwnerRole ProjectPermissionRole = "project_owner"
	EditorRole       ProjectPermissionRole = "editor"
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

type ProjectPermission struct {
	Role      ProjectPermissionRole `gorm:"type:ProjectPermissionRole"`
	UserID    uuid.UUID             `gorm:"primarykey"` // ;constraint:OnDelete:CASCADE
	ProjectID uint32                `gorm:"primarykey;index"`
	User      User                  //`gorm:"foreignKey:UserID"`
}
