# MIRACUM Mapper Backend

This is the backend repository for the Miracum Mapper project. 

## Disclaimer - Not ready for use yet

This project is currently in the development phase and not reliably functioning. It is under very active development which means that the codebase is changing rapidly and may not be stable. As soon as this changes and the project is ready for use, this disclaimer will be removed.


## Architecture

The backend is a RESTful API that provides endpoints for the frontend to interact with the database. The service is written in [Go](https://go.dev) and uses the [Gin](https://gin-gonic.com) framework to handle HTTP requests. It is packaged as a Docker Container for easy deployment. The services uses a a [PostgreSQL](https://www.postgresql.org) database that is accessed using the [GORM](https://gorm.io) ORM. User authentication is handled using [Keycloak](https://www.keycloak.org). 

## Quick Start

The application can be build and run using Docker. The following command will build the Docker image and run the container:

```bash
docker-compose up --build
```

The Docker container will be available at `http://localhost:8080`.

## Development

The application can be run locally using the following command:

```bash
go run cmd/miracummapper/main.go
```

The binary can be build locally with make for a linux environment using the following command:

```bash
make
```

Also the docker image can be build locally using the following command:

```bash
docker build -t miracummapper .
```
