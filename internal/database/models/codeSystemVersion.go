package models

import "time"

type CodeSystemVersion struct {
	ID           int32 `gorm:"primarykey;type:integer"`
	CodeSystemID int32 `gorm:"type:integer"`
	VersionID    int32 `gorm:"index;type:integer"`
	VersionName  string
	ReleaseDate  time.Time
	Imported     bool `gorm:"default:false"`
	// CodeSystemRoles     []CodeSystemRole `gorm:"foreignKey:CodeSystemVersionID"`
	// NextCodeSystemRoles []CodeSystemRole `gorm:"foreignKey:NextCodeSystemVersionID"`
	// ValidFromConcepts   []Concept        `gorm:"foreignKey:ValidFromVersionID"`
	// ValidToConcepts     []Concept        `gorm:"foreignKey:ValidToVersionID"`
}
