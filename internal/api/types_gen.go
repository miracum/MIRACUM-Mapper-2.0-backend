// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

const (
	BearerAuthScopes = "BearerAuth.Scopes"
	OAuth2Scopes     = "OAuth2.Scopes"
)

// Defines values for CodeSystemRoleType.
const (
	Destination CodeSystemRoleType = "destination"
	Source      CodeSystemRoleType = "source"
)

// Defines values for MappingEquivalence.
const (
	Equivalent                 MappingEquivalence = "equivalent"
	NotRelated                 MappingEquivalence = "not-related"
	RelatedTo                  MappingEquivalence = "related-to"
	SourceIsBroaderThanTarget  MappingEquivalence = "source-is-broader-than-target"
	SourceIsNarrowerThanTarget MappingEquivalence = "source-is-narrower-than-target"
)

// Defines values for ProjectPermissionRole.
const (
	Editor       ProjectPermissionRole = "editor"
	ProjectOwner ProjectPermissionRole = "project_owner"
	Reviewer     ProjectPermissionRole = "reviewer"
)

// Defines values for SortOrder.
const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// Defines values for GetAllConceptsParamsSortBy.
const (
	Code    GetAllConceptsParamsSortBy = "code"
	Display GetAllConceptsParamsSortBy = "display"
)

// Defines values for GetAllConceptsParamsSortOrder.
const (
	GetAllConceptsParamsSortOrderAsc  GetAllConceptsParamsSortOrder = "asc"
	GetAllConceptsParamsSortOrderDesc GetAllConceptsParamsSortOrder = "desc"
)

// Defines values for GetAllMappingsParamsSortBy.
const (
	GetAllMappingsParamsSortByComment     GetAllMappingsParamsSortBy = "comment"
	GetAllMappingsParamsSortByCreated     GetAllMappingsParamsSortBy = "created"
	GetAllMappingsParamsSortByEquivalence GetAllMappingsParamsSortBy = "equivalence"
	GetAllMappingsParamsSortById          GetAllMappingsParamsSortBy = "id"
	GetAllMappingsParamsSortByModified    GetAllMappingsParamsSortBy = "modified"
	GetAllMappingsParamsSortByStatus      GetAllMappingsParamsSortBy = "status"
)

// Defines values for GetAllMappingsParamsSortOrder.
const (
	GetAllMappingsParamsSortOrderAsc  GetAllMappingsParamsSortOrder = "asc"
	GetAllMappingsParamsSortOrderDesc GetAllMappingsParamsSortOrder = "desc"
)

// Defines values for GetProjectsParamsSortBy.
const (
	GetProjectsParamsSortByDateCreated GetProjectsParamsSortBy = "dateCreated"
	GetProjectsParamsSortById          GetProjectsParamsSortBy = "id"
	GetProjectsParamsSortByName        GetProjectsParamsSortBy = "name"
)

// Defines values for GetProjectsParamsSortOrder.
const (
	GetProjectsParamsSortOrderAsc  GetProjectsParamsSortOrder = "asc"
	GetProjectsParamsSortOrderDesc GetProjectsParamsSortOrder = "desc"
)

// CodeSystem defines model for CodeSystem.
type CodeSystem struct {
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *int32  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Title       *string `json:"title,omitempty"`
	Uri         string  `json:"uri"`
	Version     string  `json:"version"`
}

// CodeSystemRole defines model for CodeSystemRole.
type CodeSystemRole struct {
	Id       *int32 `json:"id,omitempty"`
	Name     string `json:"name"`
	Position int32  `json:"position"`
	System   struct {
		Id      *int32  `json:"id,omitempty"`
		Name    *string `json:"name,omitempty"`
		Version *string `json:"version,omitempty"`
	} `json:"system"`
	Type CodeSystemRoleType `json:"type"`
}

// CodeSystemRoleType defines model for CodeSystemRole.Type.
type CodeSystemRoleType string

// Concept defines model for Concept.
type Concept struct {
	Code    *string `json:"code,omitempty"`
	Id      *int    `json:"id,omitempty"`
	Meaning *string `json:"meaning,omitempty"`
}

// Element defines model for Element.
type Element struct {
	Concept  *Concept `json:"concept,omitempty"`
	SystemId *int32   `json:"system-id,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = string

// Mapping defines model for Mapping.
type Mapping struct {
	Comment     *string             `json:"comment,omitempty"`
	Created     *string             `json:"created,omitempty"`
	Elements    *[]Element          `json:"elements,omitempty"`
	Equivalence *MappingEquivalence `json:"equivalence,omitempty"`
	Id          *int64              `json:"id,omitempty"`
	Modified    *string             `json:"modified,omitempty"`
	Status      *string             `json:"status,omitempty"`
}

// MappingEquivalence defines model for Mapping.Equivalence.
type MappingEquivalence string

// Project defines model for Project.
type Project struct {
	Created             *string `json:"created,omitempty"`
	Description         string  `json:"description"`
	EquivalenceRequired bool    `json:"equivalence_required"`
	Id                  *int32  `json:"id,omitempty"`
	Modified            *string `json:"modified,omitempty"`
	Name                string  `json:"name"`
	StatusRequired      bool    `json:"status_required"`
	Version             string  `json:"version"`
}

// ProjectDetails defines model for ProjectDetails.
type ProjectDetails struct {
	CodeSystemRoles     []CodeSystemRole     `json:"code_system_roles"`
	Created             *string              `json:"created,omitempty"`
	Description         string               `json:"description"`
	EquivalenceRequired bool                 `json:"equivalence_required"`
	Id                  *int32               `json:"id,omitempty"`
	Modified            *string              `json:"modified,omitempty"`
	Name                string               `json:"name"`
	ProjectPermissions  *[]ProjectPermission `json:"project_permissions,omitempty"`
	StatusRequired      bool                 `json:"status_required"`
	Version             string               `json:"version"`
}

// ProjectPermission defines model for ProjectPermission.
type ProjectPermission struct {
	Role     ProjectPermissionRole `json:"role"`
	UserId   string                `json:"user_id"`
	UserName *string               `json:"user_name,omitempty"`
}

// ProjectPermissionRole defines model for ProjectPermission.Role.
type ProjectPermissionRole string

// CodeSystemRoleId defines model for code-system-role_id.
type CodeSystemRoleId = int32

// CodeSystemId defines model for code-system_id.
type CodeSystemId = int32

// Limit defines model for limit.
type Limit = int

// MappingId defines model for mapping_id.
type MappingId = int32

// Page defines model for page.
type Page = int

// PageSize defines model for pageSize.
type PageSize = int

// ProjectId defines model for project_id.
type ProjectId = int32

// SortOrder defines model for sortOrder.
type SortOrder string

// UserId defines model for user_id.
type UserId = string

// BadRequestError defines model for BadRequestError.
type BadRequestError = string

// InternalServerError defines model for InternalServerError.
type InternalServerError = string

// GetAllConceptsParams defines parameters for GetAllConcepts.
type GetAllConceptsParams struct {
	// Page Page number (must be a positive integer)
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page (minimum 1, maximum 100)
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// SortBy Field to sort sortBy
	SortBy *GetAllConceptsParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortOrder Order of sorting (asc or desc)
	SortOrder *GetAllConceptsParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// GetAllConceptsParamsSortBy defines parameters for GetAllConcepts.
type GetAllConceptsParamsSortBy string

// GetAllConceptsParamsSortOrder defines parameters for GetAllConcepts.
type GetAllConceptsParamsSortOrder string

// FindConceptByCodeParams defines parameters for FindConceptByCode.
type FindConceptByCodeParams struct {
	// Limit maximum number of items to return
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// CodeSearch search string for the code field
	CodeSearch *string `form:"codeSearch,omitempty" json:"codeSearch,omitempty"`

	// MeaningSearch search string for the meaning field
	MeaningSearch *string `form:"meaningSearch,omitempty" json:"meaningSearch,omitempty"`
}

// GetAllMappingsParams defines parameters for GetAllMappings.
type GetAllMappingsParams struct {
	// Page Page number (must be a positive integer)
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page (minimum 1, maximum 100)
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// SortBy Field to sort by
	SortBy *GetAllMappingsParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortOrder Order of sorting (asc or desc)
	SortOrder *GetAllMappingsParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// GetAllMappingsParamsSortBy defines parameters for GetAllMappings.
type GetAllMappingsParamsSortBy string

// GetAllMappingsParamsSortOrder defines parameters for GetAllMappings.
type GetAllMappingsParamsSortOrder string

// GetProjectsParams defines parameters for GetProjects.
type GetProjectsParams struct {
	// Page Page number (must be a positive integer)
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page (minimum 1, maximum 100)
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// SortBy Field to sort sortBy
	SortBy *GetProjectsParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortOrder Order of sorting (asc or desc)
	SortOrder *GetProjectsParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// GetProjectsParamsSortBy defines parameters for GetProjects.
type GetProjectsParamsSortBy string

// GetProjectsParamsSortOrder defines parameters for GetProjects.
type GetProjectsParamsSortOrder string

// UpdateCodeSystemJSONRequestBody defines body for UpdateCodeSystem for application/json ContentType.
type UpdateCodeSystemJSONRequestBody = CodeSystem

// AddCodeSystemJSONRequestBody defines body for AddCodeSystem for application/json ContentType.
type AddCodeSystemJSONRequestBody = CodeSystem

// EditProjectJSONRequestBody defines body for EditProject for application/json ContentType.
type EditProjectJSONRequestBody = Project

// UpdateCodeSystemRoleJSONRequestBody defines body for UpdateCodeSystemRole for application/json ContentType.
type UpdateCodeSystemRoleJSONRequestBody = CodeSystemRole

// UpdateMappingJSONRequestBody defines body for UpdateMapping for application/json ContentType.
type UpdateMappingJSONRequestBody = Mapping

// AddMappingJSONRequestBody defines body for AddMapping for application/json ContentType.
type AddMappingJSONRequestBody = Mapping

// AddPermissionJSONRequestBody defines body for AddPermission for application/json ContentType.
type AddPermissionJSONRequestBody = ProjectPermission

// UpdatePermissionJSONRequestBody defines body for UpdatePermission for application/json ContentType.
type UpdatePermissionJSONRequestBody = ProjectPermission

// AddProjectJSONRequestBody defines body for AddProject for application/json ContentType.
type AddProjectJSONRequestBody = ProjectDetails
