package server

import (
	"context"
	"database/sql"
	"log"
	"miracummapper/internal/api"
	"miracummapper/internal/utilities"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Project struct {
	Id                  int
	Name                string
	Version             string
	Description         string
	EquivalenceRequired bool
	StatusRequired      bool
	Created             sql.NullTime
	Modified            sql.NullTime
}

// GetProjects implements codegen.ServerInterface.
func (s *Server) GetProjects(c *gin.Context, params api.GetProjectsParams) {

	// sortBy := DefaultGetProjectsParamsSortBy

	sortBy := utilities.GetOrDefault(params.SortBy, DefaultGetProjectsParamsSortBy)

	sortOrder := utilities.GetOrDefault(params.SortOrder, DefaultGetProjectsParamsSortOrder)
	page := utilities.GetOrDefault(params.Page, DefaultGetProjectsParamsPage)
	pageSize := utilities.GetOrDefault(params.PageSize, DefaultGetProjectsParamsPageSize)

	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`SELECT * FROM "Project" ORDER BY %s %s LIMIT %d OFFSET %d`, sortBy, sortOrder, pageSize, offset)

	rows, err := s.Database.Query(query)
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.Id, &p.Name, &p.Version, &p.Description, &p.EquivalenceRequired, &p.StatusRequired, &p.Created, &p.Modified)
		if err != nil {
			log.Print(err)
			return
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		log.Print(err)
		return
	}

	// return projects
}

func (s *StrictServer) GetProjects(ctx context.Context, request api.GetProjectsRequestObject) (api.GetProjectsResponseObject, error) {
	// query := `SELECT Id, Name, Version, Description, EquivalenceRequired, StatusRequired FROM "Project" ORDER BY ? ? LIMIT ? OFFSET ?`
	// stmt, err := s.Database.Prepare(query)
	// if err != nil {
	// 	log.Print(err)
	// 	return api.GetProjects500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	// }

	pageSize := *request.Params.PageSize
	offset := (*request.Params.Page - 1) * pageSize
	sortBy := *request.Params.SortBy
	switch sortBy {
	case "dateCreated":
		sortBy = "created"
	}
	sortOrder := *request.Params.SortOrder
	switch sortOrder {
	case "asc":
		sortOrder = "ASC"
	case "desc":
		sortOrder = "DESC"
	}

	//TODO prüfen ob alles richtig geht und eingabe validiert
	allowedSortBy := map[string]bool{
		"id":      true,
		"name":    true,
		"version": true,
		// Add other allowed columns here
	}

	allowedSortOrder := map[string]bool{
		"ASC":  true,
		"DESC": true,
	}

	if !allowedSortBy[string(sortBy)] {
		return nil, fmt.Errorf("Invalid sortBy value: %s", sortBy)
	}

	if !allowedSortOrder[string(sortOrder)] {
		return nil, fmt.Errorf("Invalid sortOrder value: %s", sortOrder)
	}

	query := fmt.Sprintf(`SELECT Id, Name, Version, Description, Equivalence_required, status_required FROM "Project" ORDER BY %s %s LIMIT %d OFFSET %d`,
		sortBy, sortOrder, pageSize, offset)

	stmt, err := s.Database.Prepare(query)
	if err != nil {
		log.Print(err)
		return api.GetProjects500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Print(err)
		return api.GetProjects500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}
	defer rows.Close()

	var projects []api.Project = []api.Project{}
	for rows.Next() {
		var p api.Project
		err := rows.Scan(&p.Id, &p.Name, &p.Version, &p.Description, &p.EquivalenceRequired, &p.StatusRequired)
		if err != nil {
			log.Print(err)
			return api.GetProjects500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
		}
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		log.Print(err)
		return api.GetProjects500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}
	// TODO seitenanzahl überprüfen sonst 404 (vorher client sagen wie viele seiten es überhaupt gibt)

	//  type GetProjects200JSONResponse []Project
	return api.GetProjects200JSONResponse(projects), nil
}

// AddProject implements api.StrictServerInterface.
func (s *StrictServer) AddProject(ctx context.Context, request api.AddProjectRequestObject) (api.AddProjectResponseObject, error) {
	// Extract the project details from the request
	projectDetails := request.Body

	// Validate the project details, must contain at least one code system role
	if len(projectDetails.CodeSystemRoles) == 0 {
		return api.AddProject422JSONResponse("CodeSystemRoles are required"), nil
	}

	// Start a new transaction
	tx, err := s.Database.Begin()
	if err != nil {
		log.Print(err)
		return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}

	// Prepare the SQL query
	query := `INSERT INTO "Project" (name, version, description, equivalence_required, status_required, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $6) RETURNING Id`

	// Execute the query
	// stmt, err := tx.Prepare(query)
	// if err != nil {
	// 	log.Print(err)
	// 	return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	// }

	// get current time
	currentTime := time.Now()

	var id int32
	err = tx.QueryRow(query, projectDetails.Name, projectDetails.Version, projectDetails.Description, projectDetails.EquivalenceRequired, projectDetails.StatusRequired, currentTime).Scan(&id)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}

	// TDOO database trigger oder hier validieren, dass keine position zweimal oder position nicht im endpunkt übergeben und nur über reihenfolge geregetl wird
	// Insert the code system roles
	for _, role := range projectDetails.CodeSystemRoles {
		query := `INSERT INTO "code_system_role" (project, name, system, type, position) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.Exec(query, id, role.Name, role.System.Id, role.Type, role.Position)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
		}
	}

	// Insert permissions for the project
	for _, permission := range *projectDetails.ProjectPermissions {
		query := `INSERT INTO "project_permission" ("project", "role", "user") VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, id, permission.Role, permission.UserId)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return api.AddProject500JSONResponse{InternalServerErrorJSONResponse: "Internal Server Error"}, err
	}

	// Return the ID of the newly created project
	return api.AddProject200JSONResponse(*projectDetails), nil
}

func (s *StrictServer) GetProject(ctx context.Context, request api.GetProjectRequestObject) (api.GetProjectResponseObject, error) {
	// TODO check that projectID is valid, eg no sql injection possible

	projectId := request.ProjectId

	tx, err := s.Database.Begin()
	if err != nil {
		log.Print(err)
		// TODO 500
	}

	query := fmt.Sprintf(`SELECT name, version, description, equivalence_required, status_required, modified FROM "Project" WHERE id = %d`, projectId)

	var pd api.ProjectDetails = api.ProjectDetails{Id: &projectId}
	err = tx.QueryRow(query).Scan(&pd.Name, &pd.Version, &pd.Description, &pd.EquivalenceRequired, &pd.StatusRequired, &pd.Modified)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return api.GetProject404Response{}, nil
	} else if err != nil {
		log.Println(err)
		tx.Rollback()
		if err == sql.ErrNoRows {
			return api.GetProject404Response{}, err
		} else {
			// TODO Handle other errors
			// ...
		}
		// TODO 500
	}

	query = fmt.Sprintf(`SELECT csr.id, csr.type, csr.name, csr.position, cs.id, cs.name, cs.version FROM "code_system_role" AS csr JOIN "CodeSystem" AS cs ON csr.system = cs.id WHERE project = %d`, projectId)
	rows, err := tx.Query(query)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		//TODO 500
	}

	var csrs []api.CodeSystemRole = []api.CodeSystemRole{}
	for rows.Next() {
		var csr api.CodeSystemRole
		err = rows.Scan(&csr.Id, &csr.Type, &csr.Name, &csr.Position, &csr.System.Id, &csr.System.Name, &csr.System.Version)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			// TODO 500
		}
		csrs = append(csrs, csr)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		tx.Rollback()
		//TODO 500
	}

	pd.CodeSystemRoles = csrs

	query = fmt.Sprintf(`SELECT pp.user, pp.role, u.user_name FROM "project_permission" AS pp JOIN "User" AS u ON pp.user = u.id WHERE pp.project = %d`, projectId)
	var pps []api.ProjectPermission = []api.ProjectPermission{}
	rows, err = tx.Query(query)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		tx.Rollback()
		// TODO 500
	} else if err == nil {
		for rows.Next() {
			var pp api.ProjectPermission
			err = rows.Scan(&pp.UserId, &pp.Role, &pp.UserName)
			if err != nil {
				log.Println(err)
				tx.Rollback()
				// TODO 500
			}
			pps = append(pps, pp)
		}
	}
	pd.ProjectPermissions = &pps

	tx.Commit()
	return api.GetProject200JSONResponse(pd), nil
}
