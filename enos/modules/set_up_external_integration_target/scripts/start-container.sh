#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
    echo "$1" 1>&2
    exit 1
}

show_help() {
    cat << EOF
Usage: $0 --image IMAGE [OPTIONS]

Required:
  --image IMAGE         Docker image to run (e.g., osixia/openldap:latest)

Optional:
  --name NAME           Container name (default: auto-generated)
  --port PORT[:HOST_PORT]  Port mapping (can be used multiple times)
  --env KEY=VALUE       Environment variable (can be used multiple times)
  --volume SRC:DEST     Volume mount (can be used multiple times)
  --container-cmd CMD   Container command (default: sudo podman)
  --args ARGS           Additional arguments to pass to container run command
  --help               Show this help message

Examples:
  # Basic LDAP setup
  $0 --image osixia/openldap:latest --port 389 --port 636 --name openldap \\
     --env LDAP_ORGANISATION="My Org" --env LDAP_DOMAIN="example.com"

  # KMIP/Percona setup
  $0 --image percona/percona-server:8.0 --name kmip \\
     --volume "\$(pwd)":/TEMP_DIR --env MYSQL_ROOT_PASSWORD=testpassword \\
     --args "--port 3306"

EOF
}

# Default values
CONTAINER_CMD="sudo podman"
NAME=""
DOCKER_IMAGE=""
PORTS=()
ENVS=()
VOLUMES=()
ADDITIONAL_ARGS=""

# Check for environment variable configuration (Terraform style)
if [[ -n "${CONTAINER_IMAGE}" ]]; then
    DOCKER_IMAGE="${CONTAINER_IMAGE}"
fi

if [[ -n "${CONTAINER_NAME}" ]]; then
    NAME="${CONTAINER_NAME}"
fi

if [[ -n "${CONTAINER_PORTS}" ]]; then
    IFS=',' read -ra PORT_ARRAY <<< "${CONTAINER_PORTS}"
    PORTS=("${PORT_ARRAY[@]}")
fi

if [[ -n "${CONTAINER_ENVS}" ]]; then
    IFS=',' read -ra ENV_ARRAY <<< "${CONTAINER_ENVS}"
    ENVS=("${ENV_ARRAY[@]}")
fi

if [[ -n "${CONTAINER_VOLUMES}" ]]; then
    IFS=',' read -ra VOL_ARRAY <<< "${CONTAINER_VOLUMES}"
    VOLUMES=("${VOL_ARRAY[@]}")
fi

if [[ -n "${CONTAINER_ARGS}" ]]; then
    ADDITIONAL_ARGS="${CONTAINER_ARGS}"
fi

# Parse command line arguments (these will override environment variables)
while [[ $# -gt 0 ]]; do
    case $1 in
    --image)
        DOCKER_IMAGE="$2"
        shift 2
        ;;
    --name)
        NAME="$2"
        shift 2
        ;;
    --port)
        PORTS+=("$2")
        shift 2
        ;;
    --env)
        ENVS+=("$2")
        shift 2
        ;;
    --volume)
        VOLUMES+=("$2")
        shift 2
        ;;
    --container-cmd)
        CONTAINER_CMD="$2"
        shift 2
        ;;
    --args)
        ADDITIONAL_ARGS="$2"
        shift 2
        ;;
    --help | -h)
        show_help
        exit 0
        ;;
    *)
        fail "Unknown option: $1. Use --help for usage information."
        ;;
  esac
done

# Validate required parameters
[[ -z "${DOCKER_IMAGE}" ]] && fail "Docker image is required. Use --image to specify."

# Generate container name if not provided
if [[ -z "${NAME}" ]]; then
    NAME=$(echo "${DOCKER_IMAGE}" | sed 's/.*\///' | sed 's/:.*$//')
    echo "Using auto-generated container name: ${NAME}"
fi

# Pull the Docker image
echo "Pulling image: ${DOCKER_IMAGE}"
${CONTAINER_CMD} pull "${DOCKER_IMAGE}"

# Build the run command
RUN_CMD="${CONTAINER_CMD} run -d --name ${NAME}"

# Add port mappings
for port in "${PORTS[@]}"; do
  if [[ "${port}" == *":"* ]]; then
        # Port mapping format: host_port:container_port
        RUN_CMD="${RUN_CMD} -p ${port}"
  else
        # Same port for host and container
        RUN_CMD="${RUN_CMD} -p ${port}:${port}"
  fi
done

# Add environment variables
for env in "${ENVS[@]}"; do
    RUN_CMD="${RUN_CMD} -e ${env}"
done

# Add volume mounts
for volume in "${VOLUMES[@]}"; do
    RUN_CMD="${RUN_CMD} --volume ${volume}"
done

# Add the image
RUN_CMD="${RUN_CMD} ${DOCKER_IMAGE}"

# Add any additional arguments
if [[ -n "${ADDITIONAL_ARGS}" ]]; then
    RUN_CMD="${RUN_CMD} ${ADDITIONAL_ARGS}"
fi

# Execute the run command
echo "Starting container with command:"
echo "${RUN_CMD}"
echo ""

eval "${RUN_CMD}"

echo "${NAME} container is now running!"
