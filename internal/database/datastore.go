package database

import (
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

type Datastore interface {
	// Project
	GetAllProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error
	CreateProjectQuery(project *models.Project) error
	GetProjectQuery(project *models.Project, projectId int32) error
	UpdateProjectQuery(project *models.Project, checkFunc func(oldProject, newProject *models.Project) error) error
	DeleteProjectQuery(project *models.Project, projectId int32) error

	// ProjectPermission
	GetAllProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error
	CreateProjectPermissionQuery(projectPermission *models.ProjectPermission) error
	GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error
	UpdateProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32) error
	DeleteProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error

	// CodeSystemRole
	GetAllCodeSystemRolesQuery(codeSystemRoles *[]models.CodeSystemRole, projectId int32) error
	GetCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32) error
	UpdateCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32, checkFunc func(oldCodeSystemRole, newCodeSystemRole *models.CodeSystemRole) error) error

	// Mapping
	GetAllMappingsQuery(mappings *[]models.Mapping, projectId int, pageSize int, offset int, sortBy string, sortOrder string) error
	CreateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error)) error
	GetMappingQuery(mapping *models.Mapping, projectId int, mappingId int64) error
	UpdateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error), deleteMissingElements bool) error
	DeleteMappingQuery(mapping *models.Mapping) error

	// CodeSystem
	GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error
	CreateCodeSystemQuery(codeSystem *models.CodeSystem) error
	GetCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error
	DeleteCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error
	UpdateCodeSystemQuery(codeSystem *models.CodeSystem) error

	// Concept
	GetAllConceptsQuery(concepts *[]models.Concept, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error
}

type ErrorType int

const (
	NotFound            ErrorType = iota
	InternalServerError ErrorType = iota
	ClientError         ErrorType = iota
	// Add other error types here...
)

// The ID allows for tracing the error in the logs. In the future it could be possible to additionally set a value if the code should be printed as an api response and in the Error() function a check could be made to only return the message if the value is not set.
type DatabaseError struct {
	ID      uuid.UUID
	Type    ErrorType
	Message string
}

func NewDBError(t ErrorType, message string) *DatabaseError {
	return &DatabaseError{
		ID:      uuid.New(),
		Type:    t,
		Message: message,
	}
}

const (
	InternalServerErrorMessage = "An internal server error occurred"
)

// func NewGenericDDBError() *DatabaseError {
// 	return NewDBError(InternalServerError, InternalServerErrorMessage)
// }

func (e DatabaseError) Error() string {
	//return fmt.Sprintf("Database error %s: %s", e.ID, e.Message)
	return e.Message
	// return fmt.Sprintf("%s (Reference-Code:%s)", e.Message, e.ID)
}

func (e DatabaseError) Is(target error) bool {
	t, ok := target.(DatabaseError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

var (
	ErrNotFound = DatabaseError{
		Type: NotFound,
	}
	ErrInternalServerError = DatabaseError{
		Type: InternalServerError,
	}
	ErrClientError = DatabaseError{
		Type: ClientError,
	}
	// Define other error types here
)

// var (
// 	ErrRecordNotFound = errors.New("record not found")
// 	// Define other errors here...
// )
