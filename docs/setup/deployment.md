# Deployment

The following sections will provide information on how to deploy the application on a server.


## Service Overview

To get an overview of how the services work together, see the following graphic:

![Deployment Overview](/docs/images/ArchitectureDeployment-dark.svg#gh-dark-mode-only)
![Deployment Overview](/docs/images/ArchitectureDeployment-light.svg#gh-light-mode-only)


## Prerequisites

All services are containerized and can be deployed using Docker. A `Dockerfile` for the Backend and Frontend are provided to build the containers and `docker-compose.yaml` files are provided to run these containers with additional services like a PostgreSQL database, Keycloak and Nginx. The Frontend and Backend itself can be configured to work with these external services. To use the containers, Docker has to be installed on the machine.


## Running the application

First, a shared docker network for the frontend and the backend has to be created:
```bash
docker network create miracum_public_network
```

To run the application, use the following command:
```bash
docker compose up -d
```

The `docker-compose.yaml` file by default pulls the image `ghcr.io/miracum/miracum-mapper-2.0-backend:1.1.0` form the GitHub Container registry (In the docker compose, the image section can be replaced with the build section which is currently commented out to build the image locally).

Please note that no ports are exposed to the local machine. The frontend includes the nginx server which serves the frontend and provides SSL. The backend gets proxied by the nginx server to also provide SSL.

Currently, also the Keycloak server is started in the frontend `docker-compose.yaml` file. For a production deployment, most likely an external Keycloak server is getting used so this service can be removed from the `docker-compose.yaml` file and the config files of the services need to be adjusted accordingly.

It can be useful to create a uniform docker-compose file to start all services together. The shared network can then be declared in the docker-compose file.

For using authentication with Keycloak, a miracum-mapper client has to be created. Please refer to the [Keycloak Setup Guide](./keycloak.md) for setting up the miracum-mapper client.


## Configuration

The Backend can be configured using a `config.yaml` (see `README.md`) file and environment variables. The following environment variables can be used:

```yaml
environment:
  - PORT=8080
  - KEYCLOAK_URL=http://keycloak:8080
  - KEYCLOAK_REALM=master
  - KEYCLOAK_CLIENT_ID=miracum-mapper
  - DB_HOST=miracum-postgres
  - DB_NAME=miracum_db
  - DB_USER=miracum_user
  - DB_PASSWORD=miracum_password
```

There could be problems with the scratch image when communicating over https with other services as it does not have certificates. In this case, either manually copy the certificates over or use a different base image like Alpine.
