# Future Development

Ideas on how the project can be further developed are listed here. This also includes the frontend.

## Functionality / Extensions

- **Mapping Version Management**: Currently there is no version management implemented for the mappings. This would allow to see the workflow of a mapping, see who made which changes, makes it possible to revert changes and to compare different versions of a mapping.
- **Concept Relationships**: Currently the concepts are stored in a plain table. Some codesystems like SNOMED CT or ICD-10-GM include relationships between the concepts. These relationships could also be implemented / importet to a new table and displayed in the concept browser or the concept search.
- **Configurable Validation-Workflow**: This would allow to define a process per project on how a mapping is processed in which steps until a it reaches a final state.
- **Export Mapping to FHIR**: The mappings could be exported to FHIR resources to be used in other applications.
- **Support for a external Terminilogy Server**: Currently the Code Systems are managed in the application itself.
- **Define mapping subsets**: Process to set specific mappings for specific users so they map or review only this subset per project
- **Improve concept browser**: Maybe the URL to the online browser could be stored for each codesystem / version and then a link to the browser could be displayed in the concept browser. Also the concept browser takes a while to load the concepts. The search, sorting and filtering is also very slow. Maybe a lazy table could be used to speed it up and to reduce memory footprint.
- **Mark versio as deprecated**: A flag (boolean) could be stored for every version to indicate that it is deprecate and shouldn't be used in new / existing projects.
- **Import overview**: Currently, when a user started an import and leaves the page, there is no way to see the progress of the import again. Maybe a page could be added which is always accassible to see the progress of a currently running import.

## Technical

- **Automated Tests**: Currently tests were done manually and there are no automated tests implemented for the project. Go has good support for testing built in so unit and integration tests could be implemented withoput using any external testing framework.
- **Faster codesystem import**: The import of a codesystemversion can take some minutes for large codesystems like SNOMED CT. The limiting factor are probably the database queries. Currently for each imported concept, one or more queries are executed to find out if the concept already exists and has changed. This could maybe be improved by using batches for the queries.

# Known Issues

Things which can be improved and could not be implemented yet for time reasons are listed here.

## Insufficient time for fix

- **ProjectPermissions**: Currently, the api only accepts to post, put or patch a single permission. When a user changes multiple permissions in a project in the UI, this leads to the problem that the UI has to determine which permissions changed and then send multiple requests to the api. This could be improved by adjusting the body types from a single permission to a list of permissions so the frontend can easily call a patch request with all permissions and they get adjusted accordingly. The logic for correctly updating the permissions is then implemented centrally in the backend and not every client using the api has to implement this logic.
- **Pagination**: Currently, the project and Mappings are returned in a paginated way, but the client doesn't know how many elements/pages are available. A Meta object should be added to the response to provide this information (current page, total pages, total elements). Also filtering and sorting should be implemented. This has to be done for every field in the response object (comment, status, equivalence and also the dynamic concept fields). Then the DataTable in the frontend has to be adjusted to call the backend whenever the user changes the page, filters or sorts.
- **Unused replace by / map to hints**: In the ConceptReplaceBy table, the entries aren't linked to a codesystemVersion. So the replace by hints remain in the database, even if a version is deleted. Maybe a function could be implemented that searches for replace bies with concepts / codes that don't exist anymore and delete them. The function can then be called manually or regularly. It should be discussed if it is good to delete the entries because they cannot be restored anymore.

## Couldn't be solved / figured out

- **CodeSystemImport**: when importing a csv file for a CodeSystem, the time until the file reaches the backend code is extremely long. The import process is fast but something before this takes very long. Maybe it is part of the code generator, maybe part of gin but we couldn't figure it out
- **Mapping table in Migration View**: The mapping tables for each concept in the accordion of the project migration view should be horizontally scrollable if the content is to large to get displayed, but it does not. In the mapping table of a project it works. Could not be figured out why it behaves different.
