# Features

## Roles

### Application Roles

There are roles which can be set in KeyCloak to set permissions concerning the whole Miracum Mapper application. The roles are `admin` and `user`.

- `user`: The user role is the default role for all users. It allows to access the application and see all projects where the user is a member of. What the user can do in the project is defined by the project roles which are described later on
- `admin`: The admin role is a special role which has all permissions of the user and can additionally create new projects and assign users to projects. Respectively the admin can also adjust or remove users from projects or edit or delete projects. The creation, editing and deletion of Code Systems is also only allowed for users with the admin role.

When the roles should have other names in the keycloak instance, there are two places where the name has to be adjusted:

1. In [internal/config/constants.go](../internal/config/constants.go) `KeycloakAdminScope` and `KeycloakUserScope` have to be adjusted to the new role names
2. In the [api/openapi.yaml](../api/openapi.yaml) the `securitySchemes` have to be adjusted to the new role names

Then a new docker container has to be built in order for the changes to take effect.

### Project Roles

There are roles which can be set for each project individually to restrict the actions a user can do in the project. The roles are `reviewer`, `editor` and `project_owner`.

- `reviewer`: The reviewer role is the role with least privileges. The reviewer can only view the project and the mappings within them. The creation of new mapping and the deletion of existing mappings is not allowed for the reviewer. Only comments, status and equivalence of a mapping can be adjusted by the reviewer. Therefore the concepts can't be changed. This role should get extended further in the future to allow more actions for the reviewer when proper version management of the mappings is implemented. A look at MIRACUM Mapper 1 can be helpful to see what actions the reviewer should be able to do.
- `editor`: The editor role can do everything the reviewer can do and additionally create new mappings and delete existing mappings. He is also able to fully edit mappings. The editor can not delete the project or adjust the project settings.
- `project_owner`: The project owner is the role with the most privileges. The project owner can do everything the editor can do and additionally adjust the project settings. The project owner can also assign roles to other users in the project.

### Additional Information

The rules are just an initial implementation for the mapper and should be extended further in the future. When a proper workflow management is implemented, new capabilities or new roles in general can be added.

## Code System Import

The service is able to import codeSystemVersions and their corresponding codes. For some codesystems (currently SNOMED CT, LOINC and ICD-10-GM), the downloaded files from the internet can be imported directly without any conversion. Other codesystems or laboratory codesystems (type GENERIC) can be imported from `CSV` files. (For more information, see the [swagger file](/api/openapi.yaml)). The needed files and formats are listed below:
- LOINC:
    - Format: `CSV`
    - Files: `Loinc.csv` for the concepts and `MapTo.csv` for the replace by hints (both in the directory `LoincTable`)
    - Link: https://loinc.org/downloads/archive/ (login required)
- SNOMED CT:
    - Format: `TSV` (.txt file using tabs as separators)
    - Files: `sct2_Concept_Snapshot_<version>.txt` (directory Snapshot/Terminology) for the concepts, `sct2_Description_Snapshot_<version>.txt` (directory Snapshot/Terminology) for the descriptions and `der2_cRefset_AssociationSnapshot_<version>;.txt` (directory Snapshot/Refset/Association) for the replace by hints.
    - Link: https://www.nlm.nih.gov/healthit/snomedct/international.html (login required)
- ICD-10-GM:
    - Format: `FHIR (JSON)` 
    - Files: `Codesystem-icd10gm-<version>.json` for the concepts and `ConceptMap-icd10gm-ueberleitung-alt-neu-<version>.json` for the replace by hints
    - Link: https://terminologien.bfarm.de/CodeSystem-icd10gm-2025.download.html (no login required) - Download the FHIR Package
- GENERIC codesystems:
    - Format: `CSV`
    - Files: A concept file with required columns `code`, `display` and `status` [`active` | `trial`| `deprecated`| `discouraged`] and an optional column `description`. A replace by file with required columns `code` and `map_to` (contains the replace by / map to hints for deleted / deprecated concepts) and the optional columns `comment` and `equivalence` [`relatedto` | `equivalent` | `equal` | `wider` | `subsumes` | `narrower` | `specializes` | `inexact` | `unmatched` | `disjoint`]

If the API is used for the import, the correct endpoint has to be selected, depending on the type of the codesystem.

The upload can be done once per code system version. If the version already has codes, the import will fail so it is not possible (and not intended) to update the concepts of a version after an import. If for example a wrong file was uploaded, a version can be deleted and re-imported (if it is not in use already).

Only one import can run at once to keep the database consistent and to not overload the server.

## Project Migration

A Project Migration is the process of updating a version of a codesystem used in the project. The following endpoints (located at `/projects/{project_id}/migration`) can be used for the process:

- GET /status: check if there is currently an ongoing migration. Only one codeSystemRole can be migrated at once.

Endpoints available when no migration is running:
- GET /options: get the available newer versions for each codeSystemRole of the project
- POST /start: start the migration to a newer version

Endpoints available for a running migration:
- GET /changes: get the changes between the current and the future version for all concepts used in at least one mapping in the project. The changes are grouped by the type of change.
- POST /migrate: post a list of migrations on the mappings. The migrations are temporarily stored in the database and made permanent when the migration is finished. For each mapping it can be chosen between the options: keep the mapping as it is (if the used concept was not deleted), delete the mapping, replace the used concept with another one.
- POST /finish: finish the migration and return to the normal project state. Can only be called if all changes are reviewed / migrated.
- POST /cancel: cancel the migration. All migrations made on the mappings are reversed, except for the deletion of a mapping (cannot be undone).
