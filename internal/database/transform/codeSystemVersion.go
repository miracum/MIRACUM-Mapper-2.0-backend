package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"

	"github.com/oapi-codegen/runtime/types"
)

func GormCodeSystemVersionToApiCodeSystemVersion(codeSystemVersion *models.CodeSystemVersion) *api.CodeSystemVersion {
	return &api.CodeSystemVersion{
		Id:          codeSystemVersion.ID,
		VersionName: codeSystemVersion.VersionName,
		ReleaseDate: types.Date{Time: codeSystemVersion.ReleaseDate},
	}
}

func GormCodeSystemVersionsToApiCodeSystemVersions(codeSystemVersions *[]models.CodeSystemVersion) *[]api.CodeSystemVersion {
	apiCodeSystemVersions := []api.CodeSystemVersion{}
	for _, version := range *codeSystemVersions {
		apiCodeSystemVersions = append(apiCodeSystemVersions, *GormCodeSystemVersionToApiCodeSystemVersion(&version))
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
