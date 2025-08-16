# Testing the application

The following sections will provide information on how to run and test the application locally.


## Prerequisites

A Docker compose file is provided to run the project in a containerized environment. To use it, you need to have Docker installed on your machine.


## Running the application

To run the application using docker, use the following command:
```bash
docker compose -f docker-compose-test.yaml up -d
```
At this point, the backend api should be available which can be tested by accessing `http://localhost:8080/ping`.

Keycloak can be reached at `http://localhost:8081`. Please refer to the [Keycloak Setup Guide](./keycloak.md) for setting up the miracum-mapper client.

If you also want to develop someting in the backend, it is recommended to use the provided Dev-Container. More information can be found in the [Development Guide](./development.md).
