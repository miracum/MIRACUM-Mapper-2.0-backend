package database

import (
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

type Datastore interface {
	// Project
	GetAllProjectsQuery(projects *[]models.Project, userID *uuid.UUID, roles *[]models.ProjectPermissionRole, pageSize int, offset int, sortBy string, sortOrder string) error
	CreateProjectQuery(project *models.Project) error
	GetProjectQuery(project *models.Project, projectId int32) error
	UpdateProjectQuery(project *models.Project, checkFunc func(oldProject, newProject *models.Project) error) error
	DeleteProjectQuery(project *models.Project, projectId int32) error

	// ProjectPermission
	GetAllProjectPermissionsQuery(projectPermissions *[]models.ProjectPermission, projectId int32) error
	CreateProjectPermissionQuery(projectPermission *models.ProjectPermission) error
	GetProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error
	UpdateProjectPermissionQuery(projectPermission *models.ProjectPermission) error
	DeleteProjectPermissionQuery(projectPermission *models.ProjectPermission, projectId int32, userId uuid.UUID) error

	// CodeSystemRole
	GetAllCodeSystemRolesQuery(codeSystemRoles *[]models.CodeSystemRole, projectId int32) error
	GetCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32, codeSystemRoleId int32) error
	UpdateCodeSystemRoleQuery(codeSystemRole *models.CodeSystemRole, projectId int32) error

	// Mapping
	GetAllMappingsQuery(mappings *[]models.Mapping, projectId int, pageSize int, offset int, sortBy string, sortOrder string) error
	CreateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error)) error
	GetMappingQuery(mapping *models.Mapping, projectId int, mappingId int64) error
	UpdateMappingQuery(mapping *models.Mapping, checkFunc func(mapping *models.Mapping, project *models.Project) ([]uint32, error), deleteMissingElements bool) error
	DeleteMappingQuery(mapping *models.Mapping) error

	// User
	GetAllUsersQuery(users *[]models.User) error
	DeleteUserQuery(user *models.User, userId uuid.UUID) error
	CreateOrUpdateUserQuery(user *models.User) error

	// CodeSystem
	GetAllCodeSystemsQuery(codeSystems *[]models.CodeSystem) error
	CreateCodeSystemQuery(codeSystem *models.CodeSystem) error
	GetCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error
	DeleteCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32) error
	UpdateCodeSystemQuery(codeSystem *models.CodeSystem) error
	CheckHasNoConceptsQuery(codeSystemId int32, codeSystemVersionId int32) error
	// GetFirstElementCodeSystemQuery(codeSystem *models.CodeSystem, codeSystemId int32, concept *models.Concept) error

	// CodeSystemVersion
	CreateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error
	//GetCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemVersionId int32) error
	UpdateCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion) error
	DeleteCodeSystemVersionQuery(codeSystemVersion *models.CodeSystemVersion, codeSystemVersionId int32) error

	// Concept
	GetAllConceptsQuery(concepts *[]models.Concept, codeSystemId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error
	GetAllConceptsByVersionQuery(concepts *[]models.Concept, codeSystemId int32, codeSystemVersionId int32, pageSize int, offset int, sortBy string, sortOrder string, meaning string, code string) error
	// CreateConceptsQuery(concepts *[]models.Concept) error
}

type ErrorType int

const (
	NotFound            ErrorType = iota
	InternalServerError ErrorType = iota
	ClientError         ErrorType = iota
)

type Table int

const (
	ProjectTable Table = iota
)

// The ID allows for tracing the error in the logs. In the future it could be possible to additionally set a value if the code should be printed as an api response and in the Error() function a check could be made to only return the message if the value is not set.
type DatabaseError struct {
	ID      uuid.UUID
	Type    ErrorType
	Table   Table
	Message string
}

func NewDBError(t ErrorType, message string) *DatabaseError {
	return &DatabaseError{
		ID:      uuid.New(),
		Type:    t,
		Message: message,
		Table:   -1,
	}
}

func NewDBErrorWithTable(t ErrorType, message string, table Table) *DatabaseError {
	return &DatabaseError{
		ID:      uuid.New(),
		Type:    t,
		Message: message,
		Table:   table,
	}
}

const (
	InternalServerErrorMessage = "An internal server error occurred"
)

func (e DatabaseError) Error() string {
	return e.Message
}

func (e DatabaseError) Is(target error) bool {
	t, ok := target.(DatabaseError)
	if !ok {
		return false
	}
	if e.Table == -1 || t.Table == -1 {
		return e.Type == t.Type
	} else {
		return e.Type == t.Type && e.Table == t.Table
	}
}

var (
	ErrProjectNotFound = DatabaseError{
		Type:  NotFound,
		Table: ProjectTable,
	}
	ErrNotFound = DatabaseError{
		Type: NotFound,
	}
	ErrInternalServerError = DatabaseError{
		Type: InternalServerError,
	}
	ErrClientError = DatabaseError{
		Type: ClientError,
	}
)
