// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
	OAuth2Scopes     = "OAuth2.Scopes"
)

// Defines values for CodeSystemRoleType.
const (
	CodeSystemRoleTypeSource CodeSystemRoleType = "source"
	CodeSystemRoleTypeTarget CodeSystemRoleType = "target"
)

// Defines values for CreateCodeSystemRoleType.
const (
	CreateCodeSystemRoleTypeSource CreateCodeSystemRoleType = "source"
	CreateCodeSystemRoleTypeTarget CreateCodeSystemRoleType = "target"
)

// Defines values for CreateMappingEquivalence.
const (
	CreateMappingEquivalenceEquivalent                 CreateMappingEquivalence = "equivalent"
	CreateMappingEquivalenceNotRelated                 CreateMappingEquivalence = "not-related"
	CreateMappingEquivalenceRelatedTo                  CreateMappingEquivalence = "related-to"
	CreateMappingEquivalenceSourceIsBroaderThanTarget  CreateMappingEquivalence = "source-is-broader-than-target"
	CreateMappingEquivalenceSourceIsNarrowerThanTarget CreateMappingEquivalence = "source-is-narrower-than-target"
)

// Defines values for CreateMappingStatus.
const (
	CreateMappingStatusActive   CreateMappingStatus = "active"
	CreateMappingStatusInactive CreateMappingStatus = "inactive"
	CreateMappingStatusPending  CreateMappingStatus = "pending"
)

// Defines values for MappingEquivalence.
const (
	MappingEquivalenceEquivalent                 MappingEquivalence = "equivalent"
	MappingEquivalenceNotRelated                 MappingEquivalence = "not-related"
	MappingEquivalenceRelatedTo                  MappingEquivalence = "related-to"
	MappingEquivalenceSourceIsBroaderThanTarget  MappingEquivalence = "source-is-broader-than-target"
	MappingEquivalenceSourceIsNarrowerThanTarget MappingEquivalence = "source-is-narrower-than-target"
)

// Defines values for MappingStatus.
const (
	MappingStatusActive   MappingStatus = "active"
	MappingStatusInactive MappingStatus = "inactive"
	MappingStatusPending  MappingStatus = "pending"
)

// Defines values for Role.
const (
	Editor       Role = "editor"
	ProjectOwner Role = "project_owner"
	Reviewer     Role = "reviewer"
)

// Defines values for UpdateCodeSystemRoleType.
const (
	Source UpdateCodeSystemRoleType = "source"
	Target UpdateCodeSystemRoleType = "target"
)

// Defines values for UpdateMappingEquivalence.
const (
	Equivalent                 UpdateMappingEquivalence = "equivalent"
	NotRelated                 UpdateMappingEquivalence = "not-related"
	RelatedTo                  UpdateMappingEquivalence = "related-to"
	SourceIsBroaderThanTarget  UpdateMappingEquivalence = "source-is-broader-than-target"
	SourceIsNarrowerThanTarget UpdateMappingEquivalence = "source-is-narrower-than-target"
)

// Defines values for UpdateMappingStatus.
const (
	Active   UpdateMappingStatus = "active"
	Inactive UpdateMappingStatus = "inactive"
	Pending  UpdateMappingStatus = "pending"
)

// Defines values for SortOrder.
const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// Defines values for GetAllConceptsParamsSortBy.
const (
	Code    GetAllConceptsParamsSortBy = "code"
	Meaning GetAllConceptsParamsSortBy = "meaning"
)

// Defines values for GetAllConceptsParamsSortOrder.
const (
	GetAllConceptsParamsSortOrderAsc  GetAllConceptsParamsSortOrder = "asc"
	GetAllConceptsParamsSortOrderDesc GetAllConceptsParamsSortOrder = "desc"
)

// Defines values for GetAllProjectsParamsSortBy.
const (
	GetAllProjectsParamsSortByDateCreated GetAllProjectsParamsSortBy = "dateCreated"
	GetAllProjectsParamsSortById          GetAllProjectsParamsSortBy = "id"
	GetAllProjectsParamsSortByName        GetAllProjectsParamsSortBy = "name"
)

// Defines values for GetAllProjectsParamsSortOrder.
const (
	GetAllProjectsParamsSortOrderAsc  GetAllProjectsParamsSortOrder = "asc"
	GetAllProjectsParamsSortOrderDesc GetAllProjectsParamsSortOrder = "desc"
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

// BaseCodeSystem defines model for BaseCodeSystem.
type BaseCodeSystem struct {
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
	Name        string  `json:"name"`
	Title       *string `json:"title,omitempty"`
	Uri         string  `json:"uri"`
}

// BaseCodeSystemVersion defines model for BaseCodeSystemVersion.
type BaseCodeSystemVersion struct {
	ReleaseDate openapi_types.Date `json:"release_date"`
	VersionName string             `json:"version_name"`
}

// BaseProject defines model for BaseProject.
type BaseProject struct {
	Description         string `json:"description"`
	EquivalenceRequired bool   `json:"equivalence_required"`
	Name                string `json:"name"`
	StatusRequired      bool   `json:"status_required"`
	Version             string `json:"version"`
}

// CodeSystem defines model for CodeSystem.
type CodeSystem struct {
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          int32   `json:"id"`
	Name        string  `json:"name"`
	Title       *string `json:"title,omitempty"`
	Uri         string  `json:"uri"`
}

// CodeSystemRole defines model for CodeSystemRole.
type CodeSystemRole struct {
	Id     int32  `json:"id"`
	Name   string `json:"name"`
	System struct {
		Id      int32  `json:"id"`
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"system"`
	Type CodeSystemRoleType `json:"type"`
}

// CodeSystemRoleType defines model for CodeSystemRole.Type.
type CodeSystemRoleType string

// CodeSystemVersion defines model for CodeSystemVersion.
type CodeSystemVersion struct {
	Id          int32              `json:"id"`
	ReleaseDate openapi_types.Date `json:"release_date"`
	VersionName string             `json:"version_name"`
}

// Concept defines model for Concept.
type Concept struct {
	Code    string `json:"code"`
	Id      int64  `json:"id"`
	Meaning string `json:"meaning"`
}

// CreateCodeSystem defines model for CreateCodeSystem.
type CreateCodeSystem = BaseCodeSystem

// CreateCodeSystemRole defines model for CreateCodeSystemRole.
type CreateCodeSystemRole struct {
	Name   string                   `json:"name"`
	System int32                    `json:"system"`
	Type   CreateCodeSystemRoleType `json:"type"`
}

// CreateCodeSystemRoleType defines model for CreateCodeSystemRole.Type.
type CreateCodeSystemRoleType string

// CreateMapping defines model for CreateMapping.
type CreateMapping struct {
	Comment     *string                   `json:"comment,omitempty"`
	Elements    *[]Element                `json:"elements,omitempty"`
	Equivalence *CreateMappingEquivalence `json:"equivalence,omitempty"`
	Status      *CreateMappingStatus      `json:"status,omitempty"`
}

// CreateMappingEquivalence defines model for CreateMapping.Equivalence.
type CreateMappingEquivalence string

// CreateMappingStatus defines model for CreateMapping.Status.
type CreateMappingStatus string

// CreateProjectDetails defines model for CreateProjectDetails.
type CreateProjectDetails struct {
	CodeSystemRoles     []CreateCodeSystemRole  `json:"code_system_roles"`
	Description         string                  `json:"description"`
	EquivalenceRequired bool                    `json:"equivalence_required"`
	Name                string                  `json:"name"`
	ProjectPermissions  []SendProjectPermission `json:"project_permissions"`
	StatusRequired      bool                    `json:"status_required"`
	Version             string                  `json:"version"`
}

// Element defines model for Element.
type Element struct {
	CodeSystemRole *int32 `json:"codeSystemRole,omitempty"`
	Concept        *int64 `json:"concept,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = string

// FullElement defines model for FullElement.
type FullElement struct {
	CodeSystemRole *int32   `json:"codeSystemRole,omitempty"`
	Concept        *Concept `json:"concept,omitempty"`
}

// GetCodeSystem defines model for GetCodeSystem.
type GetCodeSystem struct {
	Author      *string             `json:"author,omitempty"`
	Description *string             `json:"description,omitempty"`
	Id          int32               `json:"id"`
	Name        string              `json:"name"`
	Title       *string             `json:"title,omitempty"`
	Uri         string              `json:"uri"`
	Versions    []CodeSystemVersion `json:"versions"`
}

// Mapping defines model for Mapping.
type Mapping struct {
	Comment     *string             `json:"comment,omitempty"`
	Created     string              `json:"created"`
	Elements    []FullElement       `json:"elements"`
	Equivalence *MappingEquivalence `json:"equivalence,omitempty"`
	Id          int64               `json:"id"`
	Modified    string              `json:"modified"`
	Status      *MappingStatus      `json:"status,omitempty"`
}

// MappingEquivalence defines model for Mapping.Equivalence.
type MappingEquivalence string

// MappingStatus defines model for Mapping.Status.
type MappingStatus string

// Project defines model for Project.
type Project struct {
	Created             string `json:"created"`
	Description         string `json:"description"`
	EquivalenceRequired bool   `json:"equivalence_required"`
	Id                  int32  `json:"id"`
	Modified            string `json:"modified"`
	Name                string `json:"name"`
	StatusRequired      bool   `json:"status_required"`
	Version             string `json:"version"`
}

// ProjectDetails defines model for ProjectDetails.
type ProjectDetails struct {
	CodeSystemRoles     []CodeSystemRole     `json:"code_system_roles"`
	Created             string               `json:"created"`
	Description         string               `json:"description"`
	EquivalenceRequired bool                 `json:"equivalence_required"`
	Id                  int32                `json:"id"`
	Modified            string               `json:"modified"`
	Name                string               `json:"name"`
	ProjectPermissions  *[]ProjectPermission `json:"project_permissions,omitempty"`
	StatusRequired      bool                 `json:"status_required"`
	Version             string               `json:"version"`
}

// ProjectPermission defines model for ProjectPermission.
type ProjectPermission struct {
	Role Role `json:"role"`
	User User `json:"user"`
}

// Role defines model for Role.
type Role string

// SendProjectPermission defines model for SendProjectPermission.
type SendProjectPermission struct {
	Role   Role   `json:"role"`
	UserId string `json:"user_id"`
}

// UpdateCodeSystemRole defines model for UpdateCodeSystemRole.
type UpdateCodeSystemRole struct {
	Id   int32                    `json:"id"`
	Name string                   `json:"name"`
	Type UpdateCodeSystemRoleType `json:"type"`
}

// UpdateCodeSystemRoleType defines model for UpdateCodeSystemRole.Type.
type UpdateCodeSystemRoleType string

// UpdateMapping defines model for UpdateMapping.
type UpdateMapping struct {
	Comment     *string                   `json:"comment,omitempty"`
	Elements    *[]Element                `json:"elements,omitempty"`
	Equivalence *UpdateMappingEquivalence `json:"equivalence,omitempty"`
	Id          int64                     `json:"id"`
	Status      *UpdateMappingStatus      `json:"status,omitempty"`
}

// UpdateMappingEquivalence defines model for UpdateMapping.Equivalence.
type UpdateMappingEquivalence string

// UpdateMappingStatus defines model for UpdateMapping.Status.
type UpdateMappingStatus string

// UpdateProject defines model for UpdateProject.
type UpdateProject struct {
	Description         string `json:"description"`
	EquivalenceRequired bool   `json:"equivalence_required"`
	Id                  int32  `json:"id"`
	Name                string `json:"name"`
	StatusRequired      bool   `json:"status_required"`
	Version             string `json:"version"`
}

// User defines model for User.
type User struct {
	Email    *string `json:"email,omitempty"`
	Fullname *string `json:"fullname,omitempty"`
	Id       string  `json:"id"`
	Username string  `json:"username"`
}

// CodesystemRoleId defines model for codesystem-role_id.
type CodesystemRoleId = int32

// CodesystemVersionId defines model for codesystem-version_id.
type CodesystemVersionId = int32

// CodesystemId defines model for codesystem_id.
type CodesystemId = int32

// MappingId defines model for mapping_id.
type MappingId = int64

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

// ForbiddenError defines model for ForbiddenError.
type ForbiddenError = string

// InternalServerError defines model for InternalServerError.
type InternalServerError = string

// UnauthorizedError defines model for UnauthorizedError.
type UnauthorizedError = string

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

	// CodeSearch search for the code
	CodeSearch *string `form:"codeSearch,omitempty" json:"codeSearch,omitempty"`

	// MeaningSearch search for meaning
	MeaningSearch *string `form:"meaningSearch,omitempty" json:"meaningSearch,omitempty"`
}

// GetAllConceptsParamsSortBy defines parameters for GetAllConcepts.
type GetAllConceptsParamsSortBy string

// GetAllConceptsParamsSortOrder defines parameters for GetAllConcepts.
type GetAllConceptsParamsSortOrder string

// GetAllProjectsParams defines parameters for GetAllProjects.
type GetAllProjectsParams struct {
	// Page Page number (must be a positive integer)
	Page *Page `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page (minimum 1, maximum 100)
	PageSize *PageSize `form:"pageSize,omitempty" json:"pageSize,omitempty"`

	// SortBy Field to sort sortBy
	SortBy *GetAllProjectsParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortOrder Order of sorting (asc or desc)
	SortOrder *GetAllProjectsParamsSortOrder `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// GetAllProjectsParamsSortBy defines parameters for GetAllProjects.
type GetAllProjectsParamsSortBy string

// GetAllProjectsParamsSortOrder defines parameters for GetAllProjects.
type GetAllProjectsParamsSortOrder string

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

// CreateCodeSystemJSONRequestBody defines body for CreateCodeSystem for application/json ContentType.
type CreateCodeSystemJSONRequestBody = CreateCodeSystem

// UpdateCodeSystemJSONRequestBody defines body for UpdateCodeSystem for application/json ContentType.
type UpdateCodeSystemJSONRequestBody = CodeSystem

// CreateCodeSystemVersionJSONRequestBody defines body for CreateCodeSystemVersion for application/json ContentType.
type CreateCodeSystemVersionJSONRequestBody = BaseCodeSystemVersion

// UpdateCodeSystemVersionJSONRequestBody defines body for UpdateCodeSystemVersion for application/json ContentType.
type UpdateCodeSystemVersionJSONRequestBody = CodeSystemVersion

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = CreateProjectDetails

// UpdateProjectJSONRequestBody defines body for UpdateProject for application/json ContentType.
type UpdateProjectJSONRequestBody = UpdateProject

// UpdateCodeSystemRoleJSONRequestBody defines body for UpdateCodeSystemRole for application/json ContentType.
type UpdateCodeSystemRoleJSONRequestBody = UpdateCodeSystemRole

// PatchMappingJSONRequestBody defines body for PatchMapping for application/json ContentType.
type PatchMappingJSONRequestBody = UpdateMapping

// CreateMappingJSONRequestBody defines body for CreateMapping for application/json ContentType.
type CreateMappingJSONRequestBody = CreateMapping

// UpdateMappingJSONRequestBody defines body for UpdateMapping for application/json ContentType.
type UpdateMappingJSONRequestBody = UpdateMapping

// CreatePermissionJSONRequestBody defines body for CreatePermission for application/json ContentType.
type CreatePermissionJSONRequestBody = SendProjectPermission

// UpdatePermissionJSONRequestBody defines body for UpdatePermission for application/json ContentType.
type UpdatePermissionJSONRequestBody = SendProjectPermission
