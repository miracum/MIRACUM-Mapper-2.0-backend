package database

import (
	"errors"
	"fmt"
	"miracummapper/internal/database/models"
)

type Datastore interface {
	GetProjectQuery(project *models.Project, projectId int32) error
	GetProjectsQuery(projects *[]models.Project, pageSize int, offset int, sortBy string, sortOrder string) error
	AddProjectQuery(project *models.Project) error
	DeleteProjectQuery(project *models.Project, projectId int32) error
	UpdateProjectQuery(project *models.Project) error
	//Add other methods here...
}

type DatabaseError struct {
	Message string
	Code    int
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

var (
	ErrRecordNotFound = errors.New("record not found")
	// Define other errors here...
)
