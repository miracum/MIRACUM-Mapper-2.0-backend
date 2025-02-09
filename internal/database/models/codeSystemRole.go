package models

import (
	"database/sql/driver"
	"fmt"
)

type CodeSystemRoleType string

const (
	Source CodeSystemRoleType = "source"
	Target CodeSystemRoleType = "target"
)

func (e *CodeSystemRoleType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = CodeSystemRoleType(v)
	case string:
		*e = CodeSystemRoleType(v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (e CodeSystemRoleType) Value() (driver.Value, error) {
	return string(e), nil
}

type CodeSystemRole struct {
	ID                      uint32             `gorm:"primarykey"`
	Type                    CodeSystemRoleType `gorm:"type:CodeSystemRoleType"`
	Name                    string
	Position                uint32
	ProjectID               uint32 `gorm:"index"`
	CodeSystemID            uint32
	Elements                []Element `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystem              CodeSystem
	CodeSystemVersionID     uint32
	CodeSystemVersion       CodeSystemVersion
	NextCodeSystemVersionID *uint32
	NextCodeSystemVersion   CodeSystemVersion
}
