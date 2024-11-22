# Database

The structure of the database looks like this:

![Test](../db/v1.0.0/miracum-mapper-database.png)

## General structure

The database is structured in a way so the data is normalized as much as possible. Indexes are used to speed up queries on large tables like the Concepts. `GORM` is used to access the Database so the definition of the tables as well as the relations between them are defined in the `GORM` Models.

The `CodeSystem` table defines the different codeSystems that can be used in the project. The `Concept` table holds the individual codes of a `CodeSystem`. The `Project` table defines a central entity in the database. A project can have multiple `ProjectPermissions` which define the access rights of a `User` to a project. A `Project` contains `Mappings` which describe an `n:m` relationship of codes from one codeSystem to others. The `n:m` Mappings are internally represented in the database using the `Element` table which maps a concept to a specific code_system_role within a `Mapping`. Which kind of `Concepts` of a `CodeSystem` are in a relationship to other `Concepts` of another `CodeSystems` is described in the `CodeSystemRole` table. They are defined individually for each project and each can have a type of `source` or `destination` and an association to one specific `CodeSystem`.

## Indexes

Indexes are created to allow efficient searches on the database. In most cases default `B-Tree` indexes are used.

- `codeSystemID` in `Concept`: When providing a suggestion `Concept` for the user while mapping, only concepts in the current codeSystem are relevant so it has to be searched for them
- `projectID` in `Mapping`: To show all mappings for a project this index is used
- `projectID` in `code_system_role`: When requesting project details the CodeSystemRoles for the project need to be queried
- `projectID` in `project_permission`: When requesting project details, the Project Permission of all users for the project need to be queried
- `mappingID` in `element`: To get all mappings of a project, all elements of a mapping have to be queried

The `Concepts` table contains a huge amount of elements as CodeSystems like Loinc can have over `50.000` elements. Therefore efficiently searching in this table by utilizing indexes is crucial. Concepts have the fields `code` and `display` which are searchable.

`code` is usually a number. When a user searches, the values must start with the specified code and can end with anything. The search is implemented in `Go` as follows (see [conceptQuery.go](../internal/database/gormQuery/conceptQuery.go)):

```go
query.Where("LOWER(code) ILIKE ?", strings.ToLower(code)+"%")
```

For searching the `meaning`, the `pg_trgm` extension is used. The implementation can be found in the `setupFullTextSearch` function in the [gormInit.go](../internal/database/gormInit.go) file. The search is implemented in go as follows (see [conceptQuery.go](../internal/database/gormQuery/conceptQuery.go)):

```go
formattedMeaning := strings.Join(strings.Fields(meaning), ":* & ") + ":*" // Adjust for partial matches
query = query.Where("display_search_vector @@ to_tsquery('english', ?)", formattedMeaning)
```
