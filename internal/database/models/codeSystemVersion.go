package models

import "time"

type CodeSystemVersion struct {
	ID                  uint32 `gorm:"primarykey"`
	CodeSystemID        uint32
	VersionID           uint32 `gorm:"unique"`
	VersionName         string
	ReleaseDate         time.Time
	CodeSystemRoles     []CodeSystemRole `gorm:"foreignKey:CodeSystemVersionID"`
	NextCodeSystemRoles []CodeSystemRole `gorm:"foreignKey:NextCodeSystemVersionID"`
	ValidFromConcepts   []Concept        `gorm:"foreignKey:ValidFromVersionID"`
	ValidToConcepts     []Concept        `gorm:"foreignKey:ValidToVersionID"`
}
