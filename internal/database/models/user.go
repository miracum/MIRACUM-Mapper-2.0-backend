package models

import "github.com/google/uuid"

type User struct {
	Id                 uuid.UUID
	UserName           string
	LogName            string
	Affiliation        string
	ProjectPermissions []ProjectPermission
}
