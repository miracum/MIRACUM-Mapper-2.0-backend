package server

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"miracummapper/internal/api"
	"miracummapper/internal/database"
	"miracummapper/internal/database/models"
	"miracummapper/internal/database/transform"
	"slices"
)

// GetMigrationOptions implements api.StrictServerInterface.
func (s *Server) GetMigrationOptions(ctx context.Context, request api.GetMigrationOptionsRequestObject) (api.GetMigrationOptionsResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetMigrationOptions404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetMigrationOptions500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MigrationPermission, permissions) {
		return api.GetMigrationOptions403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to get the migration changes for the project with ID %d", projectId))}, nil
	}

	var project models.Project
	if err := s.Database.GetProjectQuery(&project, projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetMigrationOptions404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetMigrationOptions500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project"}, nil
		}
	}

	var newerCodeSystemVersions map[int32][]models.CodeSystemVersion = make(map[int32][]models.CodeSystemVersion)
	for _, codeSystemRole := range project.CodeSystemRoles {
		if codeSystemRole.NextCodeSystemVersionID != nil {
			return api.GetMigrationOptions400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse("There is already a migration in progress for this project.")}, nil
		}
		newerCodeSystemVersions[codeSystemRole.ID] = []models.CodeSystemVersion{}
		var codeSystem models.CodeSystem
		if err := s.Database.GetCodeSystemQuery(&codeSystem, codeSystemRole.CodeSystemID); err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				return api.GetMigrationOptions404JSONResponse(fmt.Sprintf("CodeSystem with ID %d couldn't be found.", codeSystemRole.CodeSystemID)), nil
			default:
				return api.GetMigrationOptions500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the code system"}, nil
			}
		}
		for _, codeSystemVersion := range codeSystem.CodeSystemVersions {
			if codeSystemVersion.VersionID > codeSystemRole.CodeSystemVersion.VersionID {
				newerCodeSystemVersions[codeSystemRole.ID] = append(newerCodeSystemVersions[codeSystemRole.ID], codeSystemVersion)
			}
		}
	}
	return api.GetMigrationOptions200JSONResponse(*transform.GormMigrationOptionsToApiMigrationOptions(&project, &newerCodeSystemVersions)), nil
}

// GetMigrationStatus implements api.StrictServerInterface.
func (s *Server) GetMigrationStatus(ctx context.Context, request api.GetMigrationStatusRequestObject) (api.GetMigrationStatusResponseObject, error) {
	projectId := request.ProjectId

	codeSystemRole, err := s.Database.GetMigrationCodeSystemRoleQuery(projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMigrationStatus404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.GetMigrationStatus200JSONResponse{Running: false, CodeSystemRole: nil}, nil
		default:
			return api.GetMigrationStatus500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration code system role"}, err
		}
	}
	return api.GetMigrationStatus200JSONResponse{Running: true, CodeSystemRole: transform.GormCodeSystemRoleToApiCodeSystemRole(codeSystemRole)}, nil
}

// StartMigration implements api.StrictServerInterface.
func (s *Server) StartMigration(ctx context.Context, request api.StartMigrationRequestObject) (api.StartMigrationResponseObject, error) {
	projectId := request.ProjectId
	codeSystemRoleId := request.Body.CodeSystemRoleId
	versionId := request.Body.VersionId

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.StartMigration404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.StartMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(StartMigrationPermission, permissions) {
		return api.StartMigration403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update the project with ID %d", projectId))}, nil
	}

	if err := s.Database.StartMigrationQuery(projectId, codeSystemRoleId, versionId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.StartMigration404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.StartMigration400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.StartMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to start the migration"}, err
		}
	}

	return api.StartMigration200JSONResponse("Migration started successfully"), nil
}

// CancelMigration implements api.StrictServerInterface.
func (s *Server) CancelMigration(ctx context.Context, request api.CancelMigrationRequestObject) (api.CancelMigrationResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.CancelMigration404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.CancelMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(StartMigrationPermission, permissions) {
		return api.CancelMigration403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update the project with ID %d", projectId))}, nil
	}

	if err := s.Database.CancelMigrationQuery(projectId); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.CancelMigration404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.CancelMigration400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.CancelMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to cancel the migration"}, err
		}
	}
	return api.CancelMigration200JSONResponse("Migration cancelled successfully"), nil
}

// GetMigrationChanges implements api.StrictServerInterface.
func (s *Server) GetMigrationChanges(ctx context.Context, request api.GetMigrationChangesRequestObject) (api.GetMigrationChangesResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.GetMigrationChanges404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.GetMigrationChanges500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MigrationPermission, permissions) {
		return api.GetMigrationChanges403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to get the migration changes for the project with ID %d", projectId))}, nil
	}

	codeSystemRole, err := s.Database.GetMigrationCodeSystemRoleQuery(projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMigrationChanges404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.GetMigrationChanges400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.GetMigrationChanges500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration code system role"}, err
		}
	}

	validToVersions, err := s.Database.GetMigrationValidToVersionIdsQuery(codeSystemRole.CodeSystemID, *codeSystemRole.NextCodeSystemVersionID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMigrationChanges404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.GetMigrationChanges400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.GetMigrationChanges500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration changes"}, err
		}
	}

	var mappings []models.Mapping = []models.Mapping{}

	if err := s.Database.GetAllMappingsQuery(&mappings, projectId, int(^uint(0)>>1), 0, "ID", "ASC"); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMigrationChanges404JSONResponse(err.Error()), nil
		default:
			return api.GetMigrationChanges500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mappings"}, nil
		}
	}

	changes, err := computeMigrationChanges(s.Database, codeSystemRole, validToVersions, &mappings)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.GetMigrationChanges404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.GetMigrationChanges400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.GetMigrationChanges500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration changes"}, err
		}
	}
	return api.GetMigrationChanges200JSONResponse(*changes), nil

}

func computeMigrationChanges(db database.Datastore, codeSystemRole *models.CodeSystemRole, validToVersions []int32, mappings *[]models.Mapping) (*api.MigrationChanges, error) {
	var changeDescription = make(map[int32]api.MigrationChangeOldAndNewConcept)
	var changeDisplay = make(map[int32]api.MigrationChangeOldAndNewConcept)
	var deleted = make(map[int32]api.MigrationChangeOldConcept)
	var deprecated = make(map[int32]api.MigrationChangeOldConcept)
	var discouraged = make(map[int32]api.MigrationChangeOldConcept)

	codeSystemRoleId := codeSystemRole.ID
	codeSystemId := codeSystemRole.CodeSystemID
	//codeSystemVersionId := codeSystemRole.CodeSystemVersion.ID
	nextCodeSystemVersionId := *codeSystemRole.NextCodeSystemVersionID

	for _, mapping := range *mappings {
		var element *models.Element
		for _, elementLoop := range mapping.Elements {
			if elementLoop.CodeSystemRoleID == codeSystemRoleId {
				element = &elementLoop
				break
			}
		}
		if element == nil {
			continue
		}
		if element.ConceptID == nil {
			continue
		}
		if element.NextConceptID != nil {
			continue
		}
		if slices.Contains(validToVersions, element.Concept.ValidToVersionID) {
			continue
		}

		var concept models.Concept = models.Concept{}
		if err := db.GetConceptQuery(&concept, element.Concept.Code, codeSystemId, nextCodeSystemVersionId); err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				// if concept is not found, it means it was deleted
				deletedConcept, ok := deleted[*element.ConceptID]
				if ok {
					deletedConcept.Mappings = append(deletedConcept.Mappings, transform.GormMappingToApiMapping(mapping))
					deleted[*element.ConceptID] = deletedConcept
				} else {
					deleted[*element.ConceptID] = api.MigrationChangeOldConcept{
						OldConcept: *transform.GormConceptToApiConcept(&element.Concept),
						Mappings:   []api.Mapping{transform.GormMappingToApiMapping(mapping)},
					}
				}
				continue
			default:
				return nil, err
			}
		}

		if concept.Status == models.Deprecated {
			deprecatedConcept, ok := deprecated[*element.ConceptID]
			if ok {
				deprecatedConcept.Mappings = append(deprecatedConcept.Mappings, transform.GormMappingToApiMapping(mapping))
				deprecated[*element.ConceptID] = deprecatedConcept
			} else {
				deprecated[*element.ConceptID] = api.MigrationChangeOldConcept{
					OldConcept: *transform.GormConceptToApiConcept(&element.Concept),
					Mappings:   []api.Mapping{transform.GormMappingToApiMapping(mapping)},
				}
			}
		} else if concept.Status == models.Discouraged {
			discouragedConcept, ok := discouraged[*element.ConceptID]
			if ok {
				discouragedConcept.Mappings = append(discouragedConcept.Mappings, transform.GormMappingToApiMapping(mapping))
				discouraged[*element.ConceptID] = discouragedConcept
			} else {
				discouraged[*element.ConceptID] = api.MigrationChangeOldConcept{
					OldConcept: *transform.GormConceptToApiConcept(&element.Concept),
					Mappings:   []api.Mapping{transform.GormMappingToApiMapping(mapping)},
				}
			}
		} else if concept.Display != element.Concept.Display {
			changeDisplayConcept, ok := changeDisplay[*element.ConceptID]
			if ok {
				changeDisplayConcept.Mappings = append(changeDisplayConcept.Mappings, transform.GormMappingToApiMapping(mapping))
				changeDisplay[*element.ConceptID] = changeDisplayConcept
			} else {
				changeDisplay[*element.ConceptID] = api.MigrationChangeOldAndNewConcept{
					OldConcept: *transform.GormConceptToApiConcept(&element.Concept),
					NewConcept: *transform.GormConceptToApiConcept(&concept),
					Mappings:   []api.Mapping{transform.GormMappingToApiMapping(mapping)},
				}
			}
		} else if ((concept.Description == nil || element.Concept.Description == nil) && concept.Description != element.Concept.Description) ||
			(concept.Description != nil && element.Concept.Description != nil && *concept.Description != *element.Concept.Description) {
			changeDescriptionConcept, ok := changeDescription[*element.ConceptID]
			if ok {
				changeDescriptionConcept.Mappings = append(changeDescriptionConcept.Mappings, transform.GormMappingToApiMapping(mapping))
				changeDescription[*element.ConceptID] = changeDescriptionConcept
			} else {
				changeDescription[*element.ConceptID] = api.MigrationChangeOldAndNewConcept{
					OldConcept: *transform.GormConceptToApiConcept(&element.Concept),
					NewConcept: *transform.GormConceptToApiConcept(&concept),
					Mappings:   []api.Mapping{transform.GormMappingToApiMapping(mapping)},
				}
			}
		}

	}
	changes := api.MigrationChanges{
		ChangeDescription: slices.Collect(maps.Values(changeDescription)),
		ChangeDisplay:     slices.Collect(maps.Values(changeDisplay)),
		Deleted:           slices.Collect(maps.Values(deleted)),
		Deprecated:        slices.Collect(maps.Values(deprecated)),
		Discouraged:       slices.Collect(maps.Values(discouraged)),
	}
	if changes.ChangeDescription == nil {
		changes.ChangeDescription = []api.MigrationChangeOldAndNewConcept{}
	}
	if changes.ChangeDisplay == nil {
		changes.ChangeDisplay = []api.MigrationChangeOldAndNewConcept{}
	}
	if changes.Deleted == nil {
		changes.Deleted = []api.MigrationChangeOldConcept{}
	}
	if changes.Deprecated == nil {
		changes.Deprecated = []api.MigrationChangeOldConcept{}
	}
	if changes.Discouraged == nil {
		changes.Discouraged = []api.MigrationChangeOldConcept{}
	}
	return &changes, nil
}

// FinishMigration implements api.StrictServerInterface.
func (s *Server) FinishMigration(ctx context.Context, request api.FinishMigrationRequestObject) (api.FinishMigrationResponseObject, error) {
	projectId := request.ProjectId

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.FinishMigration404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.FinishMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(StartMigrationPermission, permissions) {
		return api.FinishMigration403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to update the project with ID %d", projectId))}, nil
	}

	codeSystemRole, err := s.Database.GetMigrationCodeSystemRoleQuery(projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.FinishMigration404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.FinishMigration400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.FinishMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration code system role"}, err
		}
	}

	validToVersions, err := s.Database.GetMigrationValidToVersionIdsQuery(codeSystemRole.CodeSystemID, *codeSystemRole.NextCodeSystemVersionID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.FinishMigration404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.FinishMigration400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.FinishMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to finish the migration"}, err
		}
	}

	var mappings []models.Mapping = []models.Mapping{}

	if err := s.Database.GetAllMappingsQuery(&mappings, projectId, int(^uint(0)>>1), 0, "ID", "ASC"); err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.FinishMigration404JSONResponse(err.Error()), nil
		default:
			return api.FinishMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the mappings"}, nil
		}
	}

	err = tryFinishMigration(s.Database, codeSystemRole, validToVersions, &mappings)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.FinishMigration404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.FinishMigration400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.FinishMigration500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to finish the migration"}, err
		}
	}
	return api.FinishMigration200JSONResponse("Migration finished successfully"), nil
}

func tryFinishMigration(db database.Datastore, codeSystemRole *models.CodeSystemRole, validToVersions []int32, mappings *[]models.Mapping) error {
	codeSystemRoleId := codeSystemRole.ID
	codeSystemId := codeSystemRole.CodeSystemID
	nextCodeSystemVersionId := *codeSystemRole.NextCodeSystemVersionID

	for _, mapping := range *mappings {
		var element *models.Element
		for _, elementLoop := range mapping.Elements {
			if elementLoop.CodeSystemRoleID == codeSystemRoleId {
				element = &elementLoop
				break
			}
		}
		if element == nil {
			continue
		}
		if element.ConceptID == nil {
			continue
		}
		if element.NextConceptID != nil {
			continue
		}
		if slices.Contains(validToVersions, element.Concept.ValidToVersionID) {
			continue
		}

		var concept models.Concept = models.Concept{}
		if err := db.GetConceptQuery(&concept, element.Concept.Code, codeSystemId, nextCodeSystemVersionId); err != nil {
			switch {
			case errors.Is(err, database.ErrNotFound):
				// if concept is not found, it means it was deleted
				return database.NewDBError(database.ClientError, "There is at least one mapping that uses a concept that was deleted / changed but not migrated yet. Please review all changes before finishing the migration.")
			default:
				return err
			}
		}

		if concept.Status == models.Deprecated {
			return database.NewDBError(database.ClientError, "There is at least one mapping that uses a concept that was deleted / changed but not migrated yet. Please review all changes before finishing the migration.")
		} else if concept.Status == models.Discouraged {
			return database.NewDBError(database.ClientError, "There is at least one mapping that uses a concept that was deleted / changed but not migrated yet. Please review all changes before finishing the migration.")
		} else if concept.Display != element.Concept.Display {
			return database.NewDBError(database.ClientError, "There is at least one mapping that uses a concept that was deleted / changed but not migrated yet. Please review all changes before finishing the migration.")
		} else if ((concept.Description == nil || element.Concept.Description == nil) && concept.Description != element.Concept.Description) ||
			(concept.Description != nil && element.Concept.Description != nil && *concept.Description != *element.Concept.Description) {
			return database.NewDBError(database.ClientError, "There is at least one mapping that uses a concept that was deleted / changed but not migrated yet. Please review all changes before finishing the migration.")
		}

	}
	return db.FinishMigrationQuery(codeSystemRole, mappings)
}

// MigrateMapping implements api.StrictServerInterface.
func (s *Server) MigrateMapping(ctx context.Context, request api.MigrateMappingRequestObject) (api.MigrateMappingResponseObject, error) {
	projectId := request.ProjectId
	body := request.Body

	permissions, err := getUserPermissions(ctx, s, projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrProjectNotFound):
			return api.MigrateMapping404JSONResponse(fmt.Sprintf("Project with ID %d couldn't be found.", projectId)), nil
		default:
			return api.MigrateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the project permission for the user"}, nil
		}
	}
	if !checkUserHasPermissions(MigrationPermission, permissions) {
		return api.MigrateMapping403JSONResponse{ForbiddenErrorJSONResponse: api.ForbiddenErrorJSONResponse(fmt.Sprintf("User is not authorized to migrate the mappings for the project with ID %d", projectId))}, nil
	}

	codeSystemRole, err := s.Database.GetMigrationCodeSystemRoleQuery(projectId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.MigrateMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.MigrateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.MigrateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to get the migration code system role"}, err
		}
	}

	validToVersions, err := s.Database.GetMigrationValidToVersionIdsQuery(codeSystemRole.CodeSystemID, *codeSystemRole.NextCodeSystemVersionID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.MigrateMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.MigrateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.MigrateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to migrate the mappings"}, err
		}
	}

	validFromVersions, err := s.Database.GetMigrationValidFromVersionIdsQuery(codeSystemRole.CodeSystemID, codeSystemRole.CodeSystemVersionID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotFound):
			return api.MigrateMapping404JSONResponse(err.Error()), nil
		case errors.Is(err, database.ErrClientError):
			return api.MigrateMapping400JSONResponse{BadRequestErrorJSONResponse: api.BadRequestErrorJSONResponse(err.Error())}, nil
		default:
			return api.MigrateMapping500JSONResponse{InternalServerErrorJSONResponse: "An Error occurred while trying to migrate the mappings"}, err
		}
	}

	response := api.MigrateMapping200JSONResponse{Errors: []api.MigrateMappingError{}}

	for _, migration := range *body {
		switch migration.MigrationType {

		case "none":
			continue

		case "delete":
			var mapping models.Mapping
			if err := s.Database.GetMappingQuery(&mapping, projectId, migration.MappingId); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "get the mapping"))
				continue
			}
			if err := s.Database.DeleteMappingQuery(&mapping); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "delete the mapping"))
				continue
			}
			continue

		case "keep":
			var mapping models.Mapping
			if err := s.Database.GetMappingQuery(&mapping, projectId, migration.MappingId); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "get the mapping"))
				continue
			}
			var element *models.Element
			for _, elementLoop := range mapping.Elements {
				if elementLoop.CodeSystemRoleID == codeSystemRole.ID {
					element = &elementLoop
					break
				}
			}
			if element == nil {
				response.Errors = append(response.Errors, getError(nil, &migration, "get the element"))
				continue
			}

			var concept models.Concept
			if err := s.Database.GetConceptQuery(&concept, element.Concept.Code, codeSystemRole.CodeSystemID, *codeSystemRole.NextCodeSystemVersionID); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "get the concept"))
				continue
			}
			if err := s.Database.MigrationElementSetNextConceptQuery(element.MappingID, element.CodeSystemRoleID, &concept.ID); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "set the next concept"))
				continue
			}
			continue

		case "new":
			if migration.NewConceptId == nil {
				response.Errors = append(response.Errors, api.MigrateMappingError{
					Migration: &migration,
					Code:      "400",
					Error:     "new_concept_id is required for migration_type new.",
				})
				continue
			}
			var mapping models.Mapping
			if err := s.Database.GetMappingQuery(&mapping, projectId, migration.MappingId); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "get the mapping"))
				continue
			}
			var element *models.Element
			for _, elementLoop := range mapping.Elements {
				if elementLoop.CodeSystemRoleID == codeSystemRole.ID {
					element = &elementLoop
					break
				}
			}
			if element == nil {
				response.Errors = append(response.Errors, getError(nil, &migration, "get the element"))
				continue
			}

			var concept models.Concept
			if err := s.Database.GetConceptByIdQuery(&concept, *migration.NewConceptId); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "get the concept"))
				continue
			}
			if concept.CodeSystemID != codeSystemRole.CodeSystemID {
				response.Errors = append(response.Errors, api.MigrateMappingError{
					Migration: &migration,
					Code:      "400",
					Error:     fmt.Sprintf("The concept with ID %d does not belong to the code system with ID %d.", *migration.NewConceptId, codeSystemRole.CodeSystemID),
				})
				continue
			}
			if !slices.Contains(validToVersions, concept.ValidToVersionID) || !slices.Contains(validFromVersions, concept.ValidFromVersionID) {
				response.Errors = append(response.Errors, api.MigrateMappingError{
					Migration: &migration,
					Code:      "400",
					Error:     fmt.Sprintf("The concept with ID %d is not valid in this version of the code system version with ID %d.", *migration.NewConceptId, codeSystemRole.CodeSystemID),
				})
				continue
			}
			if err := s.Database.MigrationElementSetNextConceptQuery(element.MappingID, element.CodeSystemRoleID, migration.NewConceptId); err != nil {
				response.Errors = append(response.Errors, getError(&err, &migration, "set the next concept"))
				continue
			}
			continue

		default:
			response.Errors = append(response.Errors, api.MigrateMappingError{
				Migration: &migration,
				Code:      "400",
				Error:     fmt.Sprintf("Invalid migration type %s. Valid types are none, delete, keep and new.", migration.MigrationType),
			})
			continue
		}
	}

	return response, nil
}

func getError(err *error, migration *api.MigrateMapping, errString string) api.MigrateMappingError {
	if err == nil {
		return api.MigrateMappingError{
			Migration: migration,
			Code:      "500",
			Error:     fmt.Sprintf("An Error occurred while trying to %s.", errString),
		}
	}
	switch {
	case errors.Is(*err, database.ErrClientError):
		return api.MigrateMappingError{
			Migration: migration,
			Code:      "400",
			Error:     (*err).Error(),
		}
	case errors.Is(*err, database.ErrNotFound):
		return api.MigrateMappingError{
			Migration: migration,
			Code:      "404",
			Error:     (*err).Error(),
		}
	default:
		return api.MigrateMappingError{
			Migration: migration,
			Code:      "500",
			Error:     fmt.Sprintf("An Error occurred while trying to %s.", errString),
		}
	}
}
