package database

import (
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

type Datastore interface {
	// Project
	GetProjectQuery(project *models.Project, projectId int32) error
	GetAllProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error
	CreateProjectQuery(project *models.Project) error
	UpdateProjectQuery(project *models.Project, checkFunc func(oldProject, newProject *models.Project) error) error
	DeleteProjectQuery(project *models.Project, projectId int32) error

	// CodeSystemRole
	GetAllCodeSystemRolesQuery(codeSystemRoles *[]models.CodeSystemRole, projectId int32) error
	GetCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32) error
	UpdateCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32, checkFunc func(oldCodeSystemRole, newCodeSystemRole *models.CodeSystemRole) error) error

	// ProjectPermission
	GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error
	GetProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error
	CreateProjectPermissionQuery(projectPermission *models.ProjectPermission) error
	UpdateProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32) error
	DeleteProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error

	// Mapping
	GetAllMappingsQuery(mappings *[]models.Mapping, projectId int, pageSize int, offset int, sortBy string, sortOrder string) error
	GetMappingQuery(mapping *models.Mapping, projectId int, mappingId int32) error
	CreateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) error) error
	UpdateMappingQuery(mapping *models.Mapping, checkFunc func(oldMapping, newMapping *models.Mapping) error) error
	DeleteMappingQuery(mapping *models.Mapping, mappingId int32) error
	//Add other methods here...
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
