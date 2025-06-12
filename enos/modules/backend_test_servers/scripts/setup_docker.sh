#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

# Function to check if Docker is already installed
check_docker_installed() {
  if command -v docker &> /dev/null; then
    echo "Docker is already installed: $(docker --version)"
    exit 0
  fi
}

# Function to detect the OS
detect_os() {
  if [ -f /etc/os-release ]; then
    # shellcheck disable=SC1091
    . /etc/os-release
    echo "$ID"
  else
    echo "Unknown OS: /etc/os-release not found"
  fi
}

# Main logic
check_docker_installed

os_id=$(detect_os)
echo "Installing Docker for: ${os_id}"
case "$os_id" in
  amzn)
    sudo dnf upgrade --refresh -y
    sudo dnf install -y docker
    ;;
  ubuntu)
    sudo apt update -y
    sudo apt install -y docker.io
    ;;
  rhel | centos)
    sudo yum update -y
    sudo yum install -y docker
    ;;
  *)
    echo "Unsupported or unknown OS: $os_id"
    exit 1
    ;;
esac

echo "Successfully installed Docker."
sudo systemctl enable --now docker
docker version
