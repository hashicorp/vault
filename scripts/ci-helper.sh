#!/usr/bin/env bash

# The ci-helper is used to determine build metadata, build Vault binaries,
# package those binaries into artifacts, and execute tests with those artifacts.

set -euo pipefail

# We don't want to get stuck in some kind of interactive pager
export GIT_PAGER=cat

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

# Determine the artifact basename based on metadata
function artifact_basename() {
  : "${PKG_NAME:="vault"}"
  : "${GOOS:=$(go env GOOS)}"
  : "${GOARCH:=$(go env GOARCH)}"

  : "${VERSION:=""}"
  if [ -z "$VERSION" ]; then
    echo "You must specify the VERSION variable for this command" >&2
    exit 1
  fi

  echo "${PKG_NAME}_${VERSION}_${GOOS}_${GOARCH}"
}

# Bundle the dist directory into a zip
function bundle() {
  : "${BUNDLE_PATH:=$(repo_root)/vault.zip}"
  echo "--> Bundling dist/* to $BUNDLE_PATH"
  zip -r -j "$BUNDLE_PATH" dist/
}

# Determine the root directory of the repository
function repo_root() {
  git rev-parse --show-toplevel
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
  local revision
  local build_date
  local ldflags
  local msg

  # Get or set our basic build metadata
  revision=$(build_revision)
  build_date=$(build_date) #
  : "${BIN_PATH:="dist/"}" #if not run by actions-go-build (enos local) then set this explicitly
  : "${GO_TAGS:=""}"
  : "${REMOVE_SYMBOLS:=""}"

  (unset GOOS; unset GOARCH; go generate ./...)

  # Build our ldflags
  msg="--> Building Vault revision $revision, built $build_date"

  # Keep the symbol and dwarf information by default
  if [ -n "$REMOVE_SYMBOLS" ]; then
    ldflags="-s -w "
  else
    ldflags=""
  fi

  ldflags="${ldflags} -X github.com/hashicorp/vault/version.GitCommit=$revision -X github.com/hashicorp/vault/version.BuildDate=$build_date"

  if [[ ${VERSION_METADATA+x} ]]; then
    msg="${msg}, metadata ${VERSION_METADATA}"
    ldflags="${ldflags} -X github.com/hashicorp/vault/version.VersionMetadata=$VERSION_METADATA"
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

# Package version converts a vault version string into a compatible representation for system
# packages.
function version_package() {
  awk '{ gsub("-","~",$1); print $1 }' <<< "$VAULT_VERSION"
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
  revision)
    build_revision
  ;;
  version-package)
    version_package
  ;;
  *)
    echo "unknown sub-command" >&2
    exit 1
  ;;
  esac
}

main "$@"
