package models

import "database/sql/driver"

type CodeSystemRoleType string

const (
	Source CodeSystemRoleType = "source"
	Target CodeSystemRoleType = "target"
)

func (e *CodeSystemRoleType) Scan(value interface{}) error {
	*e = CodeSystemRoleType(value.([]byte))
	return nil
}

func (e CodeSystemRoleType) Value() (driver.Value, error) {
	return string(e), nil
}

// CodeSystemRole defines model for CodeSystemRole.
type CodeSystemRole struct {
	ID           uint32             `gorm:"primarykey"`
	Type         CodeSystemRoleType `gorm:"type:CodeSystemRoleType"`
	Name         string
	Position     uint32
	ProjectID    uint32
	CodeSystemID uint32
	Elements     []Element
}
