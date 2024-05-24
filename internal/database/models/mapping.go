package models

import (
	"database/sql/driver"
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
	*e = Equivalence(value.([]byte))
	return nil
}

func (e Equivalence) Value() (driver.Value, error) {
	return string(e), nil
}

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Pending  Status = "pending"
)

func (e *Status) Scan(value interface{}) error {
	*e = Status(value.([]byte))
	return nil
}

func (e Status) Value() (driver.Value, error) {
	return string(e), nil
}

// Mapping defines model for Mapping.
type Mapping struct {
	Model
	ProjectID   uint32
	Equivalence *Equivalence `gorm:"type:Equivalence"`
	Status      *Status      `gorm:"type:Status"`
	Comment     *string
	Elements    []Element
}
