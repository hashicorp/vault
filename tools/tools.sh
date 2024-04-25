#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Determine the root directory of the repository
repo_root() {
  git rev-parse --show-toplevel
}

# Install an external Go tool.
go_install() {
  if go install "$1"; then
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
  # install action ./github/actions/install-external-tools.yml
  tools=(
    github.com/bufbuild/buf/cmd/buf@v1.25.0
    github.com/favadi/protoc-go-inject-tag@latest
    github.com/golangci/misspell/cmd/misspell@latest
    github.com/golangci/revgrep/cmd/revgrep@latest
    github.com/rinchsan/gosimports/cmd/gosimports@latest
    golang.org/x/tools/cmd/goimports@latest
    google.golang.org/protobuf/cmd/protoc-gen-go@latest
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    gotest.tools/gotestsum@latest
    honnef.co/go/tools/cmd/staticcheck@latest
    mvdan.cc/gofumpt@latest
    github.com/loggerhead/enumer@latest
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
  pushd "$(repo_root)" &> /dev/null
  for tool in "${tools[@]}"; do
    go_install ./tools/"$tool"
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
  check-external)
    check_external
  ;;
  check-internal)
    check_internal
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
