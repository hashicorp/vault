#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Perform Go formatting checks with gofumpt.
check_fmt() {
  echo "==> Checking code formatting..."

  declare -a malformed=()
  IFS=" " read -r -a files <<< "$(tr '\n' ' ' <<< "$@")"
  if [ -n "${files+set}" ] && [[ "${files[0]}" != "" ]]; then
    echo "--> Checking changed files..."
    for file in "${files[@]}"; do
      if [ ! -f "$file" ]; then
        echo "--> $file no longer exists ⚠"
        continue
      fi

      if echo "$file" | grep -v pb.go | grep -v vendor > /dev/null; then
        local output
        if ! output=$(gofumpt -l "$file") || [ "$output" != "" ]; then
          echo "--> ${file} ✖"
          malformed+=("$file")
          continue
        fi
      fi

      echo "--> ${file} ✔"
    done
  else
    echo "--> Checking all files..."
    IFS=" " read -r -a malformed <<< "$(find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gofumpt -l)"
  fi

  if [ "${#malformed[@]}" -ne 0 ] && [ -n "${malformed[0]}" ] ; then
    echo "--> The following files need to be reformatted with gofumpt"
    printf '%s\n' "${malformed[@]}"
    echo "Run \`make fmt\` to reformat code."
    for file in "${malformed[@]}"; do
      gofumpt -w "$file"
      git diff --no-color "$file"
    done
    exit 1
  fi
}

# Check that the Go toolchain meets minimum version requiremets.
check_version() {
  GO_CMD=${GO_CMD:-go}

  GO_VERSION_MIN=$1
  echo "==> Checking that build is using go version >= $1..."

  if $GO_CMD version | grep -q devel; then
    GO_VERSION="devel"
  else
    GO_VERSION=$($GO_CMD version | grep -o 'go[0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?' | tr -d 'go')

    IFS="." read -r -a GO_VERSION_ARR <<< "$GO_VERSION"
    IFS="." read -r -a GO_VERSION_REQ <<< "$GO_VERSION_MIN"

    if [[ ${GO_VERSION_ARR[0]} -lt ${GO_VERSION_REQ[0]} ||
      ( ${GO_VERSION_ARR[0]} -eq ${GO_VERSION_REQ[0]} &&
      ( ${GO_VERSION_ARR[1]} -lt ${GO_VERSION_REQ[1]} ||
      ( ${GO_VERSION_ARR[1]} -eq ${GO_VERSION_REQ[1]} && ${GO_VERSION_ARR[2]} -lt ${GO_VERSION_REQ[2]} )))
    ]]; then
      echo "Vault requires go $GO_VERSION_MIN to build; found $GO_VERSION."
      exit 1
    fi
  fi

  echo "--> Using go version $GO_VERSION..."
}

# Download all the modules for all go.mod's defined in the project.
mod_download() {
  while IFS= read -r -d '' mod; do
    echo "==> Downloading Go modules for $mod to $(go env GOMODCACHE)..."
    pushd "$(dirname "$mod")" > /dev/null || (echo "failed to push into module dir" && exit 1)
      GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/hashicorp go mod download -x
    popd > /dev/null || (echo "failed to pop out of module dir" && exit 1)
  done < <(find . -type f -name go.mod -not -path "./tools/pipeline/*" -print0 )
}

# Tidy all the go.mod's defined in the project.
mod_tidy() {
  while IFS= read -r -d '' mod; do
    echo "==> Tidying $mod..."
    pushd "$(dirname "$mod")" > /dev/null || (echo "failed to push into module dir" && exit 1)
      GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/hashicorp go mod tidy
    popd > /dev/null || (echo "failed to pop out of module dir" && exit 1)
  done < <(find . -type f -name go.mod -print0)
}

main() {
  case $1 in
  mod-download)
    mod_download
  ;;
  mod-tidy)
    mod_tidy
  ;;
  check-fmt)
    check_fmt "${@:2}"
  ;;
  check-version)
    check_version "$2"
  ;;
  *)
    echo "unknown sub-command" >&2
    exit 1
  ;;
  esac
}

main "$@"
