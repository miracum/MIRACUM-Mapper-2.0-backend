package server

import "miracummapper/internal/api"

const (
	DefaultGetProjectsParamsSortBy    api.GetProjectsParamsSortBy    = api.GetProjectsParamsSortByDateCreated
	DefaultGetProjectsParamsSortOrder api.GetProjectsParamsSortOrder = api.GetProjectsParamsSortOrderAsc
	DefaultGetProjectsParamsPage      api.Page                       = 1
	DefaultGetProjectsParamsPageSize  api.Page                       = 10
)
