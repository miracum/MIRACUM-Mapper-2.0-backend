package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"
	"slices"

	"github.com/oapi-codegen/runtime/types"
)

func GormCodeSystemVersionToApiCodeSystemVersion(codeSystemVersion *models.CodeSystemVersion, isAdmin bool) *api.CodeSystemVersion {
	if codeSystemVersion.ID == 0 {
		return nil
	}
	projectUses := []string{}

	if isAdmin {
		for _, role := range codeSystemVersion.CodeSystemRoles {
			projectName := role.Project.Name
			if !slices.Contains(projectUses, projectName) {
				projectUses = append(projectUses, projectName)
			}
		}
		for _, nextRole := range codeSystemVersion.NextCodeSystemRoles {
			projectName := nextRole.Project.Name
			if !slices.Contains(projectUses, projectName) {
				projectUses = append(projectUses, projectName)
			}
		}
	}

	return &api.CodeSystemVersion{
		Id:          codeSystemVersion.ID,
		VersionName: codeSystemVersion.VersionName,
		ReleaseDate: types.Date{Time: codeSystemVersion.ReleaseDate},
		Imported:    codeSystemVersion.Imported,
		ProjectUses: projectUses,
	}
}

func GormCodeSystemVersionsToApiCodeSystemVersions(codeSystemVersions *[]models.CodeSystemVersion, isAdmin bool) *[]api.CodeSystemVersion {
	apiCodeSystemVersions := []api.CodeSystemVersion{}
	for _, version := range *codeSystemVersions {
		apiCodeSystemVersions = append(apiCodeSystemVersions, *GormCodeSystemVersionToApiCodeSystemVersion(&version, isAdmin))
	}
	return &apiCodeSystemVersions
}

func ApiBaseCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion *api.BaseCodeSystemVersion, codeSystemId int32) *models.CodeSystemVersion {
	return &models.CodeSystemVersion{
		CodeSystemID: codeSystemId,
		VersionName:  codeSystemVersion.VersionName,
		ReleaseDate:  codeSystemVersion.ReleaseDate.Time,
	}
}

func ApiCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion *api.CodeSystemVersion, codeSystemId int32) *models.CodeSystemVersion {
	return &models.CodeSystemVersion{
		ID:           codeSystemVersion.Id,
		CodeSystemID: codeSystemId,
		VersionName:  codeSystemVersion.VersionName,
		ReleaseDate:  codeSystemVersion.ReleaseDate.Time,
	}
}

func ApiUpdateCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion *api.UpdateCodeSystemVersion, codeSystemId int32) *models.CodeSystemVersion {
	return &models.CodeSystemVersion{
		ID:           codeSystemVersion.Id,
		CodeSystemID: codeSystemId,
		VersionName:  codeSystemVersion.VersionName,
	}
}
