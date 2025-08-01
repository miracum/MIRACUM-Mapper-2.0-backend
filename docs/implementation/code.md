# Code structure

The following sections will provide information on how the code is structured.

## Code structure overview

The `api` folder contains the OpenAPI specification file which documents the api and is also used to generate the server boilerplate code. For more information on the API see [here](./api.md)

The `cmd` folder contains the main entrypoint for the application with the `main.go` file. There is also a `default-config.yaml` file containing default values to configure the service. These values are applied if the provided config file of the user doesn't include a value.

Documentation files for the development or database specific information is located in the `docs` folder.

The `internal` folder contains the main application code (see below).

The `test` folder contains the integration tests for the application.

The `tools` folder contains code generation tools like the `oapi-codegen` to generate the server boilerplate code.

## Implementation go code

The `iternal` folder is further subdivided into the following packages:

- `api`: contains the generated code from the OpenAPI specification file divided into `server_gen.go` and the `types_gen.go` code. The types define the request and response objects with structs while the server file includes the Gin router handlers and validation of the query and path parameters.

- `config`: contains the `Config` struct which holds all values related to configure the service externally by environment variables or the config file. The logic for loading the configuration file and environment variables as well as using default values if non are provided is handled here.

- `utilities`: hold different functions which are widely used accros all packages like generating UUIDs, loading environment variables, etc.

- `server`: `server.go` creates the server struct which implements the `StrictServerInterface` form the `api` package. Middlewares like authentication and rate limiting are implemented in the `middleware` directory. The other files hold the business logic of the endpoints. The files are divided in the same way as they are structured in the OpenAPI specification files. The files only hold the business logic and return the corresponding api responses (success or different errors depending on the specific endpoins). Database operations aren't implemented here in order to isolate the the database logic from the api business logic. The functions in the api endpoints only call generic operations defined in the `Datastore` interface from the `database` package or transform methods in order to transfrom API structs to their corresponding database ones. This also makes it easier to test the api endpoints without the need of a database. And the database and transform logic can be tested separately as well.

- `database`: The directory contains all database related operations. `datastore.go` defines the `Datastore` interface as well as `DatabaseErrors`. The interface capsulates all database logic so the defined functions can be called from the endpoints. The folder `gormQuery` is an implementation of the `Datastore` interface using GORM asn an ORM. Inside the `models` folder the GORM models are defined which specify the database tables and relations between them. `gormInit` is used to create a connection to the database and autoMigrate these models. In the `transform` directory the transformation between the API models and database models are defined at one central place.
