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
	ID                      int32              `gorm:"primarykey;type:integer"`
	Type                    CodeSystemRoleType `gorm:"type:CodeSystemRoleType"`
	Name                    string
	Position                int32     `gorm:"type:integer"`
	ProjectID               int32     `gorm:"index;type:integer"`
	CodeSystemID            int32     `gorm:"type:integer"`
	Elements                []Element `gorm:"constraint:OnDelete:CASCADE"`
	CodeSystem              CodeSystem
	CodeSystemVersionID     int32 `gorm:"type:integer"`
	CodeSystemVersion       CodeSystemVersion
	NextCodeSystemVersionID *int32 `gorm:"type:integer"`
	NextCodeSystemVersion   CodeSystemVersion
}
