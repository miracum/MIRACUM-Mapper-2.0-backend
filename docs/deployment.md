# Deployment

All services are containerized and can be deployed using Docker. A `Dockerfile` for the Backend and Frontend are provided to build the containers and `docker-compose.yaml` files are provided to run these containers with additional services like a PostgreSQL database, Keycloak and Nginx. The Frontend and Backend itself can be configured to work with these external services.

## Backend Configuration

The Backend can be configured using a `config.yaml` (see `README.md`)file and environment variables.The following environment variables can be used:

```yaml
environment:
  - PORT=8080
  - KEYCLOAK_URL=http://keycloak:8080/auth
  - KEYCLOAK_REALM=master
  - KEYCLOAK_CLIENT_ID=miracum-mapper
  - DB_HOST=miracum-postgres
  - DB_NAME=miracum_db
  - DB_USER=miracum_user
  - DB_PASSWORD=miracum_password
```

## Frontend Configuration

The Frontend can be configured using a `.env` file. The following environment variables can be used:

```.env
VITE_KEYCLOAK_URL=http://localhost:8081/
VITE_KEYCLOAK_CLIENT_ID=miracum-mapper
VITE_KEYCLOAK_REALM=master
VITE_APP_URL=http://localhost:5173
VITE_API_URL=http://localhost:8080
```

The `VITE_KEYCLOAK_URL` is the url of the KeyCloak server. The `VITE_KEYCLOAK_CLIENT_ID` is the client id of the KeyCloak client. The `VITE_KEYCLOAK_REALM` is the realm of the KeyCloak server. The `VITE_APP_URL` is the url of the frontend. The `VITE_API_URL` is the url of the backend.

## Service Overview

To get an overview of how the services work together, see the following graphic:

![Deployment Overview](images/ArchitectureDeployment-dark.svg#gh-dark-mode-only)
![Deployment Overview](images/ArchitectureDeployment-light.svg#gh-light-mode-only)

The two `docker-compose.yaml` files can be used to start up the frontend and backend services. The frontend includes the nginx server which serves the frontend and provides SSL. The backend gets proxied by the nginx server to also provide SSL. Currently, also the Keycloak server is started in the frontend `docker-compose.yaml` file. For a production deployment, most likely an external Keycloak server is getting used so this service can be removed from the `docker-compose.yaml` file and the config files of the services need to be adjusted accordingly. A uniform docker-compose file should be created to start all services together. To use the current two `docker-compose.yaml` files, a volume for the Backend and the Keycloak Postgres Databases have to be created (`miracum_postgres_data` and `keycloak_postgres_data`) and a shared network called `shared_network` is needed.
