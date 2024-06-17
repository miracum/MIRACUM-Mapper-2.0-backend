# MIRACUM Mapper Backend

This is the backend repository for the Miracum Mapper project. It provides a RESTful API to allow mappings from different codeSystems like `Snomed` and `Loinc` to be associated in `n:m` relationships with each other.

## Disclaimer - Not ready for use yet

This project is currently in the development phase and not reliably functioning. It is under very active development which means that the codebase is changing rapidly and may not be stable. As soon as this changes and the project is ready for use, this disclaimer will be removed.

## Features

- Projects can be created to group mappings together which map codes from a specific set of codeSystems to another set of codeSystems. The access to projects for users can be managed with projectPermissions
- CodeSystems and their corrrespondiong Codes can be imported in different formats (e.g CSV)
- The User Management is handled by Keycloak and the backend is able to authenticate users against it
- Effecient Querying of the codes from the codeSystems is possible with the use of filters

## Architecture

The backend is a RESTful API that provides endpoints for the frontend to interact with the database. The service is written in [Go](https://go.dev) and uses the [Gin](https://gin-gonic.com) framework to handle HTTP requests. It is packaged as a Docker Container for easy deployment. The API is sepcified using the [OpenAPI Specification](https://swagger.io/resources/open-api/). Boilerplate Go server code for e.g. validation query and path parameters, the request and response objects as well as defining the Gin router handlers is autmatically generated using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) code generation tool. The services uses a a [PostgreSQL](https://www.postgresql.org) database for storing data which is accessed using the [GORM](https://gorm.io) ORM. User authentication is handled by integrating [Keycloak](https://www.keycloak.org).

## Quick Start

The application can be build and run using Docker. The following command will build the Docker image and run the container next to a postgres database container:

```bash
docker-compose up --build
```

The Docker container will be available at `http://localhost:8080`.

The Service can be configured using a `config.yaml` file. The following default values are used:

```yaml
# Version of the config file
version: 1.0.0
database:
  # How many times should the service try to connect to the database when starting
  retry: 30
  # How long should the service sleep between each try to connect to the database
  sleep: 5
rate_limit:
  # enables the rate limiting of the service
  enabled: true
  limit: 100
  burst: 10
```

## Further Resources

More Information about the development can be found in the [Development Guide](./docs/development.md). There is also a [Database Documentation](./docs/database.md) available which describes the structure of the database. The api is defined in the [OpenAPI Specification](./api/openapi.yaml) file. More information about why the api was designed as it is can be read in the [API Readme](./api/README.md).
