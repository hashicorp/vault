#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Perform Go formatting checks with gofumpt.
check_fmt() {
  echo "==> Checking code formatting..."

  declare -a malformed_imports=()
  declare -a malformed_formatting=()
  local failed=0
  IFS=" " read -r -a files <<< "$(tr '\n' ' ' <<< "$@")"
  if [ -n "${files+set}" ] && [[ "${files[0]}" != "" ]]; then
    echo "--> Checking changed files..."
    for file in "${files[@]}"; do
      if [ ! -f "$file" ]; then
        echo "----> $file no longer exists ⚠"
        continue
      fi

      if echo "$file" | grep -v pb.go | grep -v vendor > /dev/null; then
        local output
        if ! output=$(gosimports -d "$file") || [ "$output" != "" ]; then
          echo "----> ${file} ✖ (gosimports)"
          malformed_imports+=("$file")
          failed=1
        fi
        if ! output=$(gofumpt -l "$file") || [ "$output" != "" ]; then
          echo "----> ${file} ✖ (gofumpt)"
          malformed_formatting+=("$file")
          failed=1
        fi

        if [ "$failed" == 1 ]; then
          continue
        fi
      fi

      echo "----> ${file} ✔"
    done
  else
    echo "--> Checking all files..."
    IFS=" " read -r -a malformed_imports <<< "$(find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gosimports -l | tr '\n' ' ')"
    IFS=" " read -r -a malformed_formatting <<< "$(find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs gofumpt -l)"
  fi

  if [ "${#malformed_imports[@]}" -ne 0 ] && [ -n "${malformed_imports[0]}" ]; then
    failed=1
    echo "--> The following files need to be reformatted with gosimports"
    printf '%s\n' "${malformed_imports[*]}"
    echo "Run \`make fmt\` to reformat code."
    for file in "${malformed_imports[@]}"; do
      gosimports -d "$file"
    done
  fi

  if [ "${#malformed_formatting[@]}" -ne 0 ] && [ -n "${malformed_formatting[0]}" ]; then
    failed=1
    echo "--> The following files need to be reformatted with gofumpt"
    printf '%s\n' "${malformed_formatting[@]}"
    echo "Run \`make fmt\` to reformat code."
    for file in "${malformed_formatting[@]}"; do
      gofumpt -w "$file"
      git diff --no-color "$file"
    done
  fi

  if [ $failed == 1 ]; then
    return 1
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
      (${GO_VERSION_ARR[0]} -eq ${GO_VERSION_REQ[0]} &&
      (${GO_VERSION_ARR[1]} -lt ${GO_VERSION_REQ[1]} ||
      (${GO_VERSION_ARR[1]} -eq ${GO_VERSION_REQ[1]} && ${GO_VERSION_ARR[2]} -lt ${GO_VERSION_REQ[2]}))) ]] \
      ; then
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
    GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/hashicorp go list ./...
    popd > /dev/null || (echo "failed to pop out of module dir" && exit 1)
  done < <(find . -type f -name go.mod -not -path "./tools/pipeline/*" -print0)
}

# Tidy all the go.mod's defined in the project.
mod_tidy() {
  while IFS= read -r -d '' mod; do
    echo "==> Tidying $mod..."
    pushd "$(dirname "$mod")" > /dev/null || (echo "failed to push into module dir" && exit 1)
    GOOS=linux GOARCH=amd64 GOPRIVATE=github.com/hashicorp go mod tidy
    popd > /dev/null || (echo "failed to pop out of module dir" && exit 1)
  done < <(find . -type f -name go.mod ! -path '*/fixtures/*' -print0)
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
