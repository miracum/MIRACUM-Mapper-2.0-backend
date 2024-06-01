package database

import (
	"miracummapper/internal/database/models"

	"github.com/google/uuid"
)

type Datastore interface {
	GetProjectQuery(project *models.Project, projectId int32) error
	GetProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error
	AddProjectQuery(project *models.Project) error
	DeleteProjectQuery(project *models.Project, projectId int32) error
	UpdateProjectQuery(project *models.Project, checkFunc func(oldProject, newProject *models.Project) error) error
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
