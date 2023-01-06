#!/usr/bin/env bash

# The ci-helper is used to determine build metadata, build Vault binaries,
# package those binaries into artifacts, and execute tests with those artifacts.

set -euo pipefail

# We don't want to get stuck in some kind of interactive pager
export GIT_PAGER=cat

# Get the full version information
function version() {
  local version
  local prerelease
  local metadata

  version=$(version_base)
  prerelease=$(version_pre)
  metadata=$(version_metadata)

  if [ -n "$metadata" ] && [ -n "$prerelease" ]; then
    echo "$version-$prerelease+$metadata"
  elif [ -n "$metadata" ]; then
    echo "$version+$metadata"
  elif [ -n "$prerelease" ]; then
    echo "$version-$prerelease"
  else
    echo "$version"
  fi
}

# Get the base version
function version_base() {
  : "${VAULT_VERSION:=""}"

  if [ -n "$VAULT_VERSION" ]; then
    echo "$VAULT_VERSION"
    return
  fi

  : "${VERSION_FILE:=$(repo_root)/sdk/version/version_base.go}"
  awk '$1 == "Version" && $2 == "=" { gsub(/"/, "", $3); print $3 }' < "$VERSION_FILE"
}

# Get the version major
function version_major() {
  version_base | cut -d '.' -f 1
}

# Get the version minor
function version_minor() {
  version_base | cut -d '.' -f 2
}

# Get the version patch
function version_patch() {
  version_base | cut -d '.' -f 3
}

# Get the version pre-release
function version_pre() {
  : "${VAULT_PRERELEASE:=""}"

  if [ -n "$VAULT_PRERELEASE" ]; then
    echo "$VAULT_PRERELEASE"
    return
  fi

  : "${VERSION_FILE:=$(repo_root)/sdk/version/version_base.go}"
  awk '$1 == "VersionPrerelease" && $2 == "=" { gsub(/"/, "", $3); print $3 }' < "$VERSION_FILE"
}

# Get the version metadata, which is commonly the edition
function version_metadata() {
  : "${VAULT_METADATA:=""}"

  if [[ (-n "$VAULT_METADATA") && ("$VAULT_METADATA" != "oss") ]]; then
    echo "$VAULT_METADATA"
    return
  fi

  : "${VERSION_FILE:=$(repo_root)/sdk/version/version_base.go}"
  awk '$1 == "VersionMetadata" && $2 == "=" { gsub(/"/, "", $3); print $3 }' < "$VERSION_FILE"
}

# Get the version formatted for Debian and RHEL packages
function version_package() {
  version | awk '{ gsub("-","~",$1); print $1 }'
}

# Get the build date from the latest commit since it can be used across all
# builds
function build_date() {
  # It's tricky to do an RFC3339 format in a cross platform way, so we hardcode UTC
  : "${DATE_FORMAT:="%Y-%m-%dT%H:%M:%SZ"}"
  git show --no-show-signature -s --format=%cd --date=format:"$DATE_FORMAT" HEAD
}

# Get the revision, which is the latest commit SHA
function build_revision() {
  git rev-parse HEAD
}

# Determine our repository by looking at our origin URL
function repo() {
  basename -s .git "$(git config --get remote.origin.url)"
}

# Determine the root directory of the repository
function repo_root() {
  git rev-parse --show-toplevel
}

# Determine the artifact basename based on metadata
function artifact_basename() {
  : "${PKG_NAME:="vault"}"
  : "${GOOS:=$(go env GOOS)}"
  : "${GOARCH:=$(go env GOARCH)}"

  echo "${PKG_NAME}_$(version)_${GOOS}_${GOARCH}"
}

# Build the UI
function build_ui() {
  local repo_root
  repo_root=$(repo_root)

  pushd "$repo_root"
  mkdir -p http/web_ui
  popd
  pushd "$repo_root/ui"
  yarn install --ignore-optional
  npm rebuild node-sass
  yarn --verbose run build
  popd
}

# Build Vault
function build() {
  local version
  local revision
  local prerelease
  local build_date
  local ldflags
  local msg

  # Get or set our basic build metadata
  version=$(version_base)
  revision=$(build_revision)
  metadata=$(version_metadata)
  prerelease=$(version_pre)
  build_date=$(build_date)
  : "${GO_TAGS:=""}"
  : "${KEEP_SYMBOLS:=""}"

  # Build our ldflags
  msg="--> Building Vault v$version, revision $revision, built $build_date"

  # Strip the symbol and dwarf information by default
  if [ -n "$KEEP_SYMBOLS" ]; then
    ldflags=""
  else
    ldflags="-s -w "
  fi

  ldflags="${ldflags}-X github.com/hashicorp/vault/sdk/version.Version=$version -X github.com/hashicorp/vault/sdk/version.GitCommit=$revision -X github.com/hashicorp/vault/sdk/version.BuildDate=$build_date"

  if [ -n "$prerelease" ]; then
    msg="${msg}, prerelease ${prerelease}"
    ldflags="${ldflags} -X github.com/hashicorp/vault/sdk/version.VersionPrerelease=$prerelease"
  fi

  if [ -n "$metadata" ]; then
    msg="${msg}, metadata ${metadata}"
    ldflags="${ldflags} -X github.com/hashicorp/vault/sdk/version.VersionMetadata=$metadata"
  fi

  # Build vault
  echo "$msg"
  pushd "$(repo_root)"
  mkdir -p dist
  mkdir -p out
  set -x
  go build -v -tags "$GO_TAGS" -ldflags "$ldflags" -o dist/
  set +x
  popd
}

# Bundle the dist directory into a zip
function bundle() {
  : "${BUNDLE_PATH:=$(repo_root)/vault.zip}"
  echo "--> Bundling dist/* to $BUNDLE_PATH"
  zip -r -j "$BUNDLE_PATH" dist/
}

# Prepare legal requirements for packaging
function prepare_legal() {
  : "${PKG_NAME:="vault"}"

  pushd "$(repo_root)"
  mkdir -p dist
  curl -o dist/EULA.txt https://eula.hashicorp.com/EULA.txt
  curl -o dist/TermsOfEvaluation.txt https://eula.hashicorp.com/TermsOfEvaluation.txt
  mkdir -p ".release/linux/package/usr/share/doc/$PKG_NAME"
  cp dist/EULA.txt ".release/linux/package/usr/share/doc/$PKG_NAME/EULA.txt"
  cp dist/TermsOfEvaluation.txt ".release/linux/package/usr/share/doc/$PKG_NAME/TermsOfEvaluation.txt"
  popd
}

# Determine the matrix group number that we'll select for execution. If the
# MATRIX_TEST_GROUP environment variable has set then it will always return
# that value. If has not been set, we will randomly select a number between 1
# and the value of MATRIX_MAX_TEST_GROUPS.
function matrix_group_id() {
  : "${MATRIX_TEST_GROUP:=""}"
  if [ -n "$MATRIX_TEST_GROUP" ]; then
    echo "$MATRIX_TEST_GROUP"
    return
  fi

  : "${MATRIX_MAX_TEST_GROUPS:=1}"
  awk -v min=1 -v max=$MATRIX_MAX_TEST_GROUPS 'BEGIN{srand(); print int(min+rand()*(max-min+1))}'
}

# Filter matrix file reads in the contents of MATRIX_FILE and filters out
# scenarios that are not in the current test group and/or those that have not
# met minimux or maximum version requirements.
function matrix_filter_file() {
  : "${MATRIX_FILE:=""}"
  if [ -z "$MATRIX_FILE" ]; then
    echo "You must specify the MATRIX_FILE variable for this command" >&2
    exit 1
  fi

  : "${MATRIX_TEST_GROUP:=$(matrix_group_id)}"

  local path
  local matrix
  path=$(readlink -f $MATRIX_FILE)
  matrix=$(cat "$path" | jq ".include |
    map(. |
      select(
        ((.min_minor_version == null) or (.min_minor_version <= $(version_minor))) and
        ((.max_minor_version == null) or (.max_minor_version >= $(version_minor))) and
        ((.test_group == null) or (.test_group == $MATRIX_TEST_GROUP))
      )
    )"
  )

  echo "{\"include\":$matrix}" | jq -c .
}

# Run the CI Helper
function main() {
  case $1 in
  artifact-basename)
    artifact_basename
  ;;
  build)
    build
  ;;
  build-ui)
    build_ui
  ;;
  bundle)
    bundle
  ;;
  date)
    build_date
  ;;
  prepare-legal)
    prepare_legal
  ;;
  matrix-filter-file)
    matrix_filter_file
  ;;
  matrix-group-id)
    matrix_group_id
  ;;
  revision)
    build_revision
  ;;
  version)
    version
  ;;
  version-base)
    version_base
  ;;
  version-pre)
    version_pre
  ;;
  version-major)
    version_major
  ;;
  version-meta)
    version_metadata
  ;;
  version-minor)
    version_minor
  ;;
  version-package)
    version_package
  ;;
  version-patch)
    version_patch
  ;;
  *)
    echo "unknown sub-command" >&2
    exit 1
  ;;
  esac
}

main "$@"
