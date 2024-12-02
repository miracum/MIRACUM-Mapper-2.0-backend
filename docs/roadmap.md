# Future Development

Ideas on how the project can be further developed are listed here.

- **Automated Tests**: Currently tests were done manually and there are no automated tests implemented for the project. Go has good support for testing built in so unit and integration tests could be implemented withoput using any external testing framework.
- **Version Management**: Currently there is no version management implemented for the mappings. This would allow to see the workflow of a mapping, see who made which changes, makes it possible to revert changes and to compare different versions of a mapping.
- **Configurable Validation-Workflow**: This would allow to defined a process per project on how a mapping is processed in which steps until a it reaches a final state.
- **Export Mapping to FHIR**: The mappings could be exported to FHIR resources to be used in other applications.
- **Support for a external Terminilogy Server**: Currently the Code Systems are managed in the application itself.
- **Define mapping subsets**: Process to set specific mappings for specific users so they map or review only this subset per project

# Known Issues

Things which can be improved and could not be implemented yet for time reasons are listed here.

- **ProjectPermissions**: Currently, the api only accepts to post, put or patch a single permission. When a user changes multiple permissions in a project in the UI, this leads to the problem that the UI has to determine which permissions changed and then send multiple requests to the api. This could be improved by adjusting the body types from a single permission to a list of permissions so the frontend can easily call a patch request with all permissions and they get adjusted accordingly. The logic for correctly updating the permissions is then implemented centrally in the backend and not every client using the api has to implement this logic.
- **Pagination**: Currently, the project and Mappings are returned in a paginated way, but the client doesn't know how many elements/pages are available. A Meta object should be added to the response to provide this information (current page, total pages, total elements). Also filtering and sorting should be implemented. This has to be done for every field in the response object (comment, status, equivalence and also the dynamic concept fields). Then the DataTable in the frontend has to be adjusted to call the backend whenever the user changes the page, filters or sorts.
- **CodeSystemImport**: when importing a csv file for a CodeSystem, the time until the file reaches the backend code is extremely long. The import process is fast but something before this takes very long. Maybe it is part of the code generator, maybe part of gin but we couldn't figure it out
