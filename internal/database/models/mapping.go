package models

import (
	"database/sql/driver"
	"errors"
)

type Equivalence string

const (
	RelatedTo                  Equivalence = "related-to"
	Equivalent                 Equivalence = "equivalent"
	SourceIsNarrowerThanTarget Equivalence = "source-is-narrower-than-target"
	SourceIsBroaderThanTarget  Equivalence = "source-is-broader-than-target"
	NotRelated                 Equivalence = "not-related"
)

func (e *Equivalence) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = Equivalence(v)
	case string:
		*e = Equivalence([]byte(v))
	default:
		return errors.New("invalid type for Equivalence")
	}
	return nil
}

func (e Equivalence) Value() (driver.Value, error) {
	return string(e), nil
}

type MappingStatus string

const (
	ActiveMapping MappingStatus = "active"
	Inactive      MappingStatus = "inactive"
	Pending       MappingStatus = "pending"
)

func (e *MappingStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = MappingStatus(v)
	case string:
		*e = MappingStatus([]byte(v))
	default:
		return errors.New("invalid type for Status")
	}
	return nil
}

func (e MappingStatus) Value() (driver.Value, error) {
	return string(e), nil
}

type Mapping struct {
	ModelBigId
	ProjectID   uint32         `gorm:"index"`
	Equivalence *Equivalence   `gorm:"type:Equivalence"`
	Status      *MappingStatus `gorm:"type:MappingStatus"`
	Comment     *string
	Elements    []Element `gorm:"constraint:OnDelete:CASCADE"`
}
