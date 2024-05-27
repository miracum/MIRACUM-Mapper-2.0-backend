- How to develop locally (go version, maybe Dev Containers)
- Pipeline Setup
- How to run tests (unit tests and integration tests)
- How to run and debug the app locally


- GitHub Actions workflow
  - Stage: Pre-Checks
    - Linting
    - Run generate for oapi-codegen and check if the generated code is up-to-date with spec file
    - check changelog is up-to-date (next version)
  - Stage Build
    - Build the docker image
  - Stage Test
    - Run unit tests
    - Run integration tests (with postgres-db)
    - (Run security checks)
  - Stage Deploy
    - Upload Docker Image to Registry