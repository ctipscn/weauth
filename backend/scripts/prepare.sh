#!/bin/bash

set -euo pipefail

# Define container runtime information.
declare -A CONTAINER_INFO
CONTAINER_INFO["NAME"]='pg-test'
CONTAINER_INFO["IMAGE"]='postgres:14.24-trixie'
CONTAINER_INFO["PORT_MAPPING"]='5432:5432'
CONTAINER_INFO["POSTGRES_USER"]='weauth'
CONTAINER_INFO["POSTGRES_PASSWORD"]='admin123'
CONTAINER_INFO["POSTGRES_DB"]='weauth'


####################################
# Check whether Docker is running.
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Returns:
#   None
####################################
check_docker_running() {
    if ! sudo docker info >/dev/null 2>&1; then
        echo "ERROR: Docker is not running or cannot be accessed." >&2
        exit 1
    fi
    echo "Docker is running."
}


####################################
# Check whether the target container exists.
#
# Globals:
#   CONTAINER_INFO
#
# Arguments:
#   None
#
# Returns:
#   0 if the container exists, 1 otherwise.
####################################
container_exists() {
    local name="${CONTAINER_INFO["NAME"]}"
    sudo docker container inspect "$name" >/dev/null 2>&1
}


####################################
# Stop and remove the target container if it exists.
#
# Globals:
#   CONTAINER_INFO
#
# Arguments:
#   None
#
# Returns:
#   None
####################################
remove_existing_container() {
    local name="${CONTAINER_INFO["NAME"]}"

    if container_exists; then
        echo "Container '$name' exists. Stopping and removing it..."
        sudo docker stop "$name" >/dev/null 2>&1 || true
        sudo docker rm -f "$name" >/dev/null 2>&1 || true
        echo "Container '$name' has been stopped and removed."
    else
        echo "Container '$name' does not exist. No removal needed."
    fi
}


####################################
# Create the target PostgreSQL container.
#
# Globals:
#   CONTAINER_INFO
#
# Arguments:
#   None
#
# Returns:
#   None
####################################
create_container() {
    local name="${CONTAINER_INFO["NAME"]}"
    local image="${CONTAINER_INFO["IMAGE"]}"
    local port_mapping="${CONTAINER_INFO["PORT_MAPPING"]}"
    local pg_user="${CONTAINER_INFO["POSTGRES_USER"]}"
    local pg_password="${CONTAINER_INFO["POSTGRES_PASSWORD"]}"
    local pg_db="${CONTAINER_INFO["POSTGRES_DB"]}"

    echo "Creating container '$name' from image '$image'..."
    sudo docker run -d \
        --name "$name" \
        --rm \
        -p "$port_mapping" \
        -e POSTGRES_USER="$pg_user" \
        -e POSTGRES_PASSWORD="$pg_password" \
        -e POSTGRES_DB="$pg_db" \
        "$image"
    echo "Container '$name' created successfully."
}


check_docker_running
remove_existing_container
create_container
