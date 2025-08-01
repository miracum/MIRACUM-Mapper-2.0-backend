# MIRACUM Mapper Backend

This is the backend repository for the Miracum Mapper project. It provides a RESTful API to allow mappings from different codeSystems like `SNOMED CT` or `LOINC` to be associated in `n:m` relationships with each other. The corresponding frontend repository can be found [here](https://github.com/miracum/MIRACUM-Mapper-2.0-frontend)

<!-- ## Disclaimer - Not ready for use yet

This project is currently in the development phase and not reliably functioning. It is under very active development which means that the codebase is changing rapidly and may not be stable. As soon as this changes and the project is ready for use, this disclaimer will be removed. -->

## Features

- **Projects** can be created to group **Mappings** together which map codes from a specific set of codeSystems to another set of codeSystems. The access to projects for users can be managed with projectPermissions
- **CodeSystems** and their corresponding **Codes** can be imported in different formats (e.g CSV, FHIR/JSON). The changes between imported versions of codeSystems are computed automatically and stored efficiently in the database
- The used version of a codeSystem in a project can be upgraded in a process called *Migration*, where all changes of concepts which are used in mappings have to be reviewed
- The **User Management** is handled by **KeyCloak** and the backend is able to authenticate users against it
- Efficient Querying of the codes from the codeSystems is possible with the use of filters

## Architecture

The backend is a RESTful API that provides endpoints for the frontend to interact with the database. The service is written in [Go](https://go.dev) and uses the [Gin](https://gin-gonic.com) framework to handle HTTP requests. It is packaged as a Docker Container for easy deployment. The API is specified using the [OpenAPI Specification](https://swagger.io/resources/open-api/). Boilerplate Go server code for e.g. validation query and path parameters, the request and response objects as well as defining the Gin router handlers is automatically generated using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) code generation tool. The service uses a [PostgreSQL](https://www.postgresql.org) database for storing data which is accessed using the [GORM](https://gorm.io) ORM. User authentication is handled by integrating [Keycloak](https://www.keycloak.org).

## Quick Start

The application can be build and run using Docker. A `Dockerfile` and `docker-compose.yaml` files are provided for this purpose. The backend can be run either in the docker container or in a Dev-Container which is also provided for development. The backend uses a postgres database service that stores its data in a volume. For authentication during development and testing, a Keycloak service is provided that also uses a postgres database service and a volume.

There are different guides depending on how the application should be set up and run:

- If you want to contribute to the project and actively develop something, see the [Development Guide](/docs/setup/development.md)
- If you want to test the application (e. g. import codesystems, create projects and mappings) without changing the code, see the [Test Guide](/docs/setup/testing.md)
- If you want to deploy the application for production or on a test system, see the [Deployment Guide](/docs/setup/deployment.md)

## Configuration

The Service can be configured using a `config.yaml` file. The following default values are used:

```yaml
# Version of the config file
version: 1.0.0
database:
  # How many times should the service try to connect to the database when starting
  retry: 30
  # How long should the service sleep between each try to connect to the database
  sleep: 5
keycloak:
  # How many times should the service try to connect to the database when starting
  retry: 30
  # How long should the service sleep between each try to connect to the database
  sleep: 5
cors:
  # Allowed origins for the CORS policy
  allowed_origins:
    - "*"
# Debug is a flag to enable debug mode (e.g. more verbose logging)
debug: true
```


## Further Resources

- More Information about the **development** can be found in the [Development Guide](/docs/setup/development.md). The **code structure** is described in the [Code Structure Documentation](/docs/implementation/code.md). There is also a [Database Documentation](/docs/implementation/database.md) available which describes the structure of the **database**. The **api** is defined in the [OpenAPI Specification](/api/openapi.yaml) file. More information about why the api was designed as it is can be read in the [API Documentation](/docs/implementation/api.md).
- How the deployment process takes place can be read in the [Deployment Guide](/docs/setup/deployment.md).
- More information about the features of the project like the codesystem import and the roles and permissions can be found in the [Features Documentation](/docs/features.md).
- Ideas for future development and known issues can be found in the [Roadmap](/docs/roadmap.md).


## Test System

UPDATE 07/2025: The Test System is currently not working / not reachable.

Currently the code is deployed on a Test System which can be found at [https://131.188.43.66](https://131.188.43.66). The backend and the database work together with the frontend, a own keycloak instance for authentication as well as a nginx reverse proxy to provide SSL encryption. All service are running in docker containers. An overview of the current deployment can be seen in the graphic below:

![Deployment Test System](docs/images/ArchitectureTestSystem-dark.svg#gh-dark-mode-only)
![Deployment Test System](docs/images/ArchitectureTestSystem-light.svg#gh-light-mode-only)
