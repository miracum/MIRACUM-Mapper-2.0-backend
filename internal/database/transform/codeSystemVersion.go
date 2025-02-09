package transform

import (
	"miracummapper/internal/api"
	"miracummapper/internal/database/models"

	"github.com/oapi-codegen/runtime/types"
)

func GormCodeSystemVersionToApiCodeSystemVersion(codeSystemVersion *models.CodeSystemVersion) *api.CodeSystemVersion {
	return &api.CodeSystemVersion{
		Id:          int32(codeSystemVersion.ID),
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
		CodeSystemID: uint32(codeSystemId),
		VersionName:  codeSystemVersion.VersionName,
		ReleaseDate:  codeSystemVersion.ReleaseDate.Time,
	}
}

func ApiCodeSystemVersionToGormCodeSystemVersion(codeSystemVersion *api.CodeSystemVersion, codeSystemId int32) *models.CodeSystemVersion {
	return &models.CodeSystemVersion{
		ID:           uint32(codeSystemVersion.Id),
		CodeSystemID: uint32(codeSystemId),
		VersionName:  codeSystemVersion.VersionName,
		ReleaseDate:  codeSystemVersion.ReleaseDate.Time,
	}
}
