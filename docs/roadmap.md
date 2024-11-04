# Future Development

Ideas on how the project can be further developed are listed here.

- **Automated Tests**: Currently tests were done manually and there are no automated tests implemented for the project. Go has good support for testing built in so unit and integration tests could be implemented withoput using any external testing framework.
- **Version Management**: Currently there is no version management implemented for the mappings. This would allow to see the workflow of a mapping, see who made which changes, makes it possible to revert changes and to compare different versions of a mapping.
- **Configurable Validation-Workflow**: This would allow to defined a process per project on how a mapping is processed in which steps until a it reaches a final state.
- **Export Mapping to FHIR**: The mappings could be exported to FHIR resources to be used in other applications.
- **Support for a external Terminilogy Server**: Currently the Code Systems are managed in the application itself.
- **Define mapping subsets**: Process to set specific mappings for specific users so they map or review only this subset per project
