package models

import (
	"database/sql/driver"
	"errors"
)

type ConceptReplaceByEquivalence string

const (
	ReplaceByRelatedTo   ConceptReplaceByEquivalence = "relatedto"
	ReplaceByEquivalent  ConceptReplaceByEquivalence = "equivalent"
	ReplaceByEqual       ConceptReplaceByEquivalence = "equal"
	ReplaceByWider       ConceptReplaceByEquivalence = "wider"
	ReplaceBySubsumes    ConceptReplaceByEquivalence = "subsumes"
	ReplaceByNarrower    ConceptReplaceByEquivalence = "narrower"
	ReplaceBySpecializes ConceptReplaceByEquivalence = "specializes"
	ReplaceByInexact     ConceptReplaceByEquivalence = "inexact"
	ReplaceByUnmatched   ConceptReplaceByEquivalence = "unmatched"
	ReplaceByDisjoint    ConceptReplaceByEquivalence = "disjoint"
)

func (e *ConceptReplaceByEquivalence) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*e = ConceptReplaceByEquivalence(v)
	case string:
		*e = ConceptReplaceByEquivalence([]byte(v))
	default:
		return errors.New("invalid type for ConceptReplaceByEquivalence")
	}
	return nil
}

func (e ConceptReplaceByEquivalence) Value() (driver.Value, error) {
	return string(e), nil
}

type ConceptReplaceBy struct {
	Code         string                       `gorm:"index;primarykey"`
	MapTo        string                       `gorm:"index;primarykey"`
	CodeSystemID int32                        `gorm:"index;type:integer;primarykey"`
	Equivalence  *ConceptReplaceByEquivalence `gorm:"type:ConceptReplaceByEquivalence"`
	Comment      *string
}
