package models

import "github.com/google/uuid"

type User struct {
	Id                 uuid.UUID
	FullName           string
	UserName           string
	Email              string
	ProjectPermissions []ProjectPermission
}
