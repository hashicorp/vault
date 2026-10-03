#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Determine the root directory of the repository
repo_root() {
  git rev-parse --show-toplevel
}

# Install an external Go tool.
go_install() {
  local tags=""
  if [ "$(go env GOOS)" == "darwin" ]; then
    tags="netcgo"
  fi
  if eval CGO_ENABLED=0 go install "-tags=${tags}" \"-ldflags=-w -s\" "$1"; then
    echo "--> $1 ✔"
  else
    echo "--> $1 ✖"
    return 1
  fi
}

# Check for a tool binary in the path.
check_tool() {
  if builtin type -P "$2" &> /dev/null; then
    echo "--> $2 ✔"
  else
    echo "--> $2 ✖"
    echo "Could not find required $1 tool $2. Run 'make tools-$1' to install it." 1>&2
    return 1
  fi
}

# Install external tools.
install_external() {
  local tools
  # If you update this please update check_external below as well as our external tools
  # install action .github/actions/install-external-tools/action.yml
  #
  # All tool versions should match the versions in .github/actions/install-external-tools/action.yml
  # Protobuf tool versions should match what's in Vault's go.mod.
  tools=(
    honnef.co/go/tools/cmd/staticcheck@v0.6.0
    github.com/bufbuild/buf/cmd/buf@v1.45.0
    github.com/favadi/protoc-go-inject-tag@v1.4.0
    github.com/golangci/misspell/cmd/misspell@v0.6.0
    github.com/golangci/revgrep/cmd/revgrep@v0.8.0
    github.com/stevendpclark/enumer@v0.0.0-20250122154818-a42b666c3cd3
    github.com/rinchsan/gosimports/cmd/gosimports@v0.3.8
    golang.org/x/tools/cmd/goimports@v0.30.0
    google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
    gotest.tools/gotestsum@v1.12.3
    mvdan.cc/gofumpt@v0.8.0
    mvdan.cc/sh/v3/cmd/shfmt@v3.10.0
  )

  echo "==> Installing external tools..."
  for tool in "${tools[@]}"; do
    go_install "$tool"
  done
}

# Check that all tools are installed
check_external() {
  # Ensure that all external tools are available. In CI we'll prefer installing pre-built external
  # tools for speed instead of go install so that we don't require downloading Go modules and
  # compiling tools from scratch in every CI job.
  # See .github/actions/install-external-tools.yml for that workflow.
  local tools
  tools=(
    buf
    enumer
    gofumpt
    goimports
    gosimports
    gotestsum
    misspell
    protoc-gen-go
    protoc-gen-go-grpc
    protoc-go-inject-tag
    revgrep
    shfmt
    staticcheck
  )

  echo "==> Checking for external tools..."
  for tool in "${tools[@]}"; do
    check_tool external "$tool"
  done
}

# Install internal tools.
install_internal() {
  local tools
  # If you update this please update check tools below.
  tools=(
    codechecker
    stubmaker
  )

  echo "==> Installing internal tools..."
  pushd "$(repo_root)/tools" &> /dev/null
  for tool in "${tools[@]}"; do
    go_install ./"$tool"
  done
  popd &> /dev/null
}

# Check internal that all tools are installed
check_internal() {
  # Ensure that all required internal tools are available.
  local tools
  tools=(
    codechecker
    stubmaker
  )

  echo "==> Checking for internal tools..."
  for tool in "${tools[@]}"; do
    check_tool internal "$tool"
  done
}

# Install our pipeline tools. In some cases these may require access to internal repositories so
# they are excluded from our baseline toolset.
install_pipeline() {
  echo "==> Installing pipeline tools..."
  pushd "$(repo_root)/tools/pipeline" &> /dev/null
  if env GOPRIVATE=github.com/hashicorp go install ./...; then
    echo "--> pipeline ✔"
  else
    echo "--> pipeline ✖"
    popd &> /dev/null
    return 1
  fi
  popd &> /dev/null
}

# Check that all required pipeline tools are installed
check_pipeline() {
  echo "==> Checking for pipeline tools..."
  check_tool pipeline pipeline
}

# Install tools.
install() {
  install_internal
  install_external
}

# Check tools.
check() {
  check_internal
  check_external
}

main() {
  case $1 in
    install-external)
      install_external
      ;;
    install-internal)
      install_internal
      ;;
    install-pipeline)
      install_pipeline
      ;;
    check-external)
      check_external
      ;;
    check-internal)
      check_internal
      ;;
    check-pipeline)
      check_pipeline
      ;;
    install)
      install
      ;;
    check)
      check
      ;;
    *)
      echo "unknown sub-command" >&2
      exit 1
      ;;
  esac
}

main "$@"
