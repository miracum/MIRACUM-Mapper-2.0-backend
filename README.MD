# MIRACUM Mapper Backend

This is the backend repository for the Miracum Mapper project. It provides a RESTful API to allow mappings from different codeSystems like `Snomed` and `Loinc` to be associated in `n:m` relationships with each other. The corresponding frontend repository can be found [here](https://github.com/miracum/MIRACUM-Mapper-2.0-frontend)

## Disclaimer - Not ready for use yet

This project is currently in the development phase and not reliably functioning. It is under very active development which means that the codebase is changing rapidly and may not be stable. As soon as this changes and the project is ready for use, this disclaimer will be removed.

## Features

- Projects can be created to group mappings together which map codes from a specific set of codeSystems to another set of codeSystems. The access to projects for users can be managed with projectPermissions
- CodeSystems and their corresponding Codes can be imported in different formats (e.g CSV)
- The User Management is handled by KeyCloak and the backend is able to authenticate users against it
- Efficient Querying of the codes from the codeSystems is possible with the use of filters

## Architecture

The backend is a RESTful API that provides endpoints for the frontend to interact with the database. The service is written in [Go](https://go.dev) and uses the [Gin](https://gin-gonic.com) framework to handle HTTP requests. It is packaged as a Docker Container for easy deployment. The API is specified using the [OpenAPI Specification](https://swagger.io/resources/open-api/). Boilerplate Go server code for e.g. validation query and path parameters, the request and response objects as well as defining the Gin router handlers is automatically generated using the [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) code generation tool. The service uses a [PostgreSQL](https://www.postgresql.org) database for storing data which is accessed using the [GORM](https://gorm.io) ORM. User authentication is handled by integrating [Keycloak](https://www.keycloak.org).

## Quick Start

The application can be build and run using Docker. The following command will build the Docker image and run the container next to a postgres database container. An Keycloak instance with the corresponding postgres database is also started. The Keycloak instances need to persist data, thats why two volumes have to be created in advance.

```bash
docker volume create miracum_postgres_data
docker volume create keycloak_postgres_data
docker-compose up --build
```

At this point, the backend api should be available which can be tested by accessing `http://localhost:8080/ping`. But most endpoints require authentication so Keycloak has to be configured in the next step. Keycloak should now be available under `http://localhost:8081`. Please login using `admin` as the username and `admin` as the password. Afterwards, import the miracum-mapper client by clicking on `Client` and then `Import Client`. Select the file which can be found [here](tools/setup/keycloak-client-miracum-mapper.json). Afterward, roles need to be created in the section `Roles` within the newly created `miracum-mapper` client. Two roles should be created, called `admin` and `normal`. These roles need to be assign to either the admin user or a newly created user in oder for them to access the miracum mapper. A new user can be created by clicking on `Users`, `Add User` and following the dialog. Afterwards, click on the user, then on `Role Mapping`, `Assign Role`, in the search bar search for either `admin` or `normal`, select the role and click on `assign`. Afterwards, when hosting the frontend, you should be able to login to the frontend and then copy the access token for the api in the `Profile` section. A guide on how to host the frontend can be seen in the [Frontend Readme](https://github.com/miracum/MIRACUM-Mapper-2.0-frontend). Any client like [Insomnia](https://insomnia.rest), the VS Code Extension [Swagger Preview](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer),... can be used to access the api. By default access tokens are only valid for one minute and get refreshed automatically by the frontend. This can be frustrating when developing as this means you have to copy a fresh access token every minute. To extend the duration of an access token, login to keycloak, click on `Realm Settings` > `Tokens`, scroll down to `Access Tokens` and adjust the `Access Token Lifespan` accordingly. For more information to get started with the development and e.g. debug the `miracum-mapper`, have a look at [development.md](./docs/deployment.md).

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
```

## Code System Import

The service is able to import codeSystems and their corresponding codes from a `CSV` file. The `CSV` file should have the two columns `code` and `meaning`. The UI doesn't feature an import yet so the `CSV` file has to be uploaded to the `/codesystems/{codesystem_id}/import` endpoint directly (For more information, see the [swagger file](api/openapi.yaml)). In order to create a new CodeSystem, follow these steps:

1. Create a new CodeSystem using the POST `/codesystems` endpoint
2. Prepare a `CSV` file for import. To achieve it, [this python script](tools/codesystem-import/parser/README.md) can be used.
3. Upload the `CSV` file to the `/codesystems/{codesystem_id}/import` endpoint. This upload can be done once per code system. If the code system already has codes, the import will fail so it is not possible via the API to update codes within a codesystem after they were created

Important: Only one input is allowed to run at a time for a codeSystem. When two imports are triggered simultaneously for one codeSystem, this can lead to potentially duplicated codes.

## Test System

Currently the code is deployed on a Test System which can be found at [https://131.188.43.66](https://131.188.43.66). The backend and the database work together with the frontend, a own keycloak instance for authentication as well as a nginx reverse proxy to provide SSL encryption. All service are running in docker containers. An overview of the current deployment can be seen in the graphic below:

![Deployment Test System](docs/images/ArchitectureTestSystem-dark.svg#gh-dark-mode-only)
![Deployment Test System](docs/images/ArchitectureTestSystem-light.svg#gh-light-mode-only)

## Further Resources

More Information about the development can be found in the [Development Guide](./docs/development.md). There is also a [Database Documentation](./docs/database.md) available which describes the structure of the database. The api is defined in the [OpenAPI Specification](./api/openapi.yaml) file. More information about why the api was designed as it is can be read in the [API Readme](./api/README.md). More information about the features of the project like the roles and permissions can be found in the [Features Documentation](./docs/features.md). How the deployment process takes place can be read in the [Deployment Documentation](./docs/deployment.md).
