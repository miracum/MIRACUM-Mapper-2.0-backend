#!/bin/sh

# Get the Docker group ID from the host system
DOCKER_GROUP_ID=$(stat -c '%g' /var/run/docker.sock)

# Create the Docker group inside the container with the same ID
if ! getent group docker >/dev/null; then
    addgroup -g $DOCKER_GROUP_ID docker
fi

# Add the vscode user to the Docker group
adduser vscode docker

# Execute the original command
exec "$@"