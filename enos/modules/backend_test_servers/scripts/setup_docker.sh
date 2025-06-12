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
    . /etc/os-release
    echo "OS: $ID"
  else
    echo "unknown"
  fi
}

# Main logic
check_docker_installed

echo "Installing Docker..."
os_id=$(detect_os)
case "$os_id" in
  amzn)
    sudo dnf install -y docker
    ;;
  ubuntu)
    sudo apt install docker
    ;;
  rhel | centos)
    sudo yum install docker
    ;;
  *)
    echo "Unsupported or unknown OS: $os_id"
    exit 1
    ;;
esac

echo "Successfully installed Docker."
sudo systemctl enable --now docker
docker version
