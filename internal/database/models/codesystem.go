package models

type CodeSystem struct {
	Model
	Uri             string
	Version         string
	Name            string
	Title           *string
	Description     *string
	Author          *string
	Concepts        []Concept
	CodeSystemRoles []CodeSystemRole
}
