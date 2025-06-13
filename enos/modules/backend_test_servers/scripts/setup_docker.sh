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
    sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update
    sudo apt install docker-ce docker-ce-cli containerd.io -y
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
sudo docker --version