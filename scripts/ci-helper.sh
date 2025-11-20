#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

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
    : "${VERSION_METADATA:="ce"}"

    : "${VERSION:=""}"
    if [ -z "$VERSION" ]; then
        echo "You must specify the VERSION variable for this command" >&2
        exit 1
    fi

    local version
    version="$VERSION"
    if [ "$VERSION_METADATA" != "ce" ]; then
        version="${VERSION}+${VERSION_METADATA}"
    fi

    echo "${PKG_NAME}_${version}_${GOOS}_${GOARCH}"
}

# Copy binaries from BIN_PATH to TARGET_BIN_PATH
function copy_binary() {
    : "${BIN_PATH:="dist/"}"

    if [ -z "${TARGET_BIN_PATH:-}" ]; then
        echo "TARGET_BIN_PATH not specified, skipping binary copy" >&2
        return 0
    fi

    echo "--> Copying binary from $BIN_PATH to $TARGET_BIN_PATH"
    mkdir -p "$TARGET_BIN_PATH"
    if [ -d "$BIN_PATH" ]; then
        cp -r "$BIN_PATH"/* "$TARGET_BIN_PATH/"
    else
        echo "Warning: Source directory $BIN_PATH does not exist" >&2
        return 1
    fi
}

# Bundle the dist directory into a zip
# Note: This always bundles from dist/, regardless of TARGET_BIN_PATH
function bundle() {
    : "${BUNDLE_PATH:=$(repo_root)/vault.zip}"
    : "${BIN_PATH:="dist/"}"

    if [ ! -d "$BIN_PATH" ] || [ -z "$(ls -A "$BIN_PATH" 2>/dev/null)" ]; then
        echo "Warning: $BIN_PATH is empty or does not exist, bundle will be empty" >&2
    fi

    echo "--> Bundling $BIN_PATH/* to $BUNDLE_PATH..."
    zip -r -j "$BUNDLE_PATH" "$BIN_PATH/"
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
    yarn install
    npm rebuild node-sass
    yarn run build
    popd
}

# Build Vault
function build() {
    local revision
    local base_version
    local build_date
    local ldflags
    local msg

    # Get or set our basic build metadata
    revision=$(build_revision)
    build_date=$(build_date)
    base_version=$(version_base)
    version=$(version)
    : "${BIN_PATH:="dist/"}" #if not run by actions-go-build (enos local) then set this explicitly
    : "${GO_TAGS:=""}"
    : "${REMOVE_SYMBOLS:=""}"

    # Generate code but make sure we don't slurp in cross compilation env vars
    (
        unset GOOS
        unset GOARCH
        unset CC
        unset CC_FOR_TARGET
        go generate ./...
    )

    # Build our ldflags
    msg="--> Building Vault v$version revision $revision, built $build_date..."

    # Keep the symbol and dwarf information by default
    if [ -n "$REMOVE_SYMBOLS" ]; then
        ldflags="-s -w "
    else
        ldflags=""
    fi

    # If you read what happens in the "version" package you will see that the
    # "version.Version" symbol is automatically set from the embedded VERSION
    # file. So why are we setting it again with linker flags?
    #
    # Well, some third party security scanners like Trivy attempt to determine a
    # Go binaries "version" by reading the embedded debug build info. The main
    # module "version" reported there has little to do with what we consider
    # Vaults version and is instead what the Go module system considers the
    # vault modules "pseudo-version"[0].
    #
    # What Go determines as the pseudo-version can be pretty complicated. If you
    # tag a commit with a semver-ish tag and push it before you build the binary,
    # the "pseudo-version" will be the tag value. But what if you build the binary
    # before a commit has an associated tag like we do? Well, it depends. If you
    # build a Go binary with "-buildvcs" enabled, the "pseudo-version" reported
    # here looks something like: "<prior release tag>-<timestamp>-<sha>+dirty".
    # If Go cannot resolve a prior tag you'll get "v0.0.0" in place of
    # "<prior release tag>". If you disable "-buildvcs" you'll get "devel".
    #
    # As we can see, there's quite a lot of variance in this system and a modules
    # "version" is an unreliable way to reason about a softwares "version". But
    # that doesn't stop tools from using it and reporting CVEs based on it!
    #
    # That's right. If you publish a binary with the "+dirty" style pseudo-version,
    # and the prior tag that is resolves is associated with a CVE, your binary will
    # be flagged for the same CVE even if it has nothing to do with the prior tag.
    # If you disable "buildvcs" (we do) these tools cannot determine a "version"
    # (because it's always "devel"). When this happens these scanners also fail
    # because they can't determine a version. Cool.
    #
    # So that brings us back to our original query: what's going on with the
    # ldflags. To work around this problem, Trivy *reads arbitrary ldflags in
    # the binary build info* to determine the "version"![1] when the main module
    # does not report a version. And it is because of that, dear reader, that we
    # inject our version again via linker flags, to please tooling that relies on
    # the unreliable.
    #
    # [1]: https://go.dev/doc/modules/version-numbers#pseudo-version-number
    # [0]: https://trivy.dev/v0.62/docs/coverage/language/golang/#main-module
    ldflags="${ldflags} -X github.com/hashicorp/vault/version.GitCommit=$revision -X github.com/hashicorp/vault/version.BuildDate=$build_date -X github.com/hashicorp/vault/version.Version=$base_version"

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
    go env
    go build -v -buildvcs=false -tags "$GO_TAGS" -ldflags "$ldflags" -o dist/
    set +x
    popd
}

# ENT: Prepare legal requirements for packaging
function prepare_ent_legal() {
    : "${PKG_NAME:="vault"}"

    if [ -z "${LICENSE_DIR:-}" ]; then
        echo "You must set LICENSE_DIR; example: export LICENSE_DIR=.release/ibm-pao/license/default" 1>&2
        return 1
    fi

    pushd "$(repo_root)"
    mkdir -p dist
    cp -R "$LICENSE_DIR" dist/
    mkdir -p ".release/linux/package/usr/share/doc/$PKG_NAME"
    cp -R "$LICENSE_DIR" ".release/linux/package/usr/share/doc/$PKG_NAME/"
    popd
}

# CE: Prepare legal requirements for packaging
function prepare_ce_legal() {
    : "${PKG_NAME:="vault"}"

    pushd "$(repo_root)"

    mkdir -p dist
    cp LICENSE dist/LICENSE.txt

    mkdir -p ".release/linux/package/usr/share/doc/$PKG_NAME"
    cp LICENSE ".release/linux/package/usr/share/doc/$PKG_NAME/LICENSE.txt"

    popd
}

# version returns the $VAULT_VERSION env variable or reads the VERSION file.
function version() {
    if [[ -n "${VAULT_VERSION+x}" ]]; then
        echo "${VAULT_VERSION}"
        return 0
    fi

    cat "$(readlink -f "$(dirname "$0")/../version/VERSION")"
}

# Base version converts a vault version string into the base version, which omits
# any prerelease or edition metadata.
function version_base() {
    local ver
    ver=$(version)
    echo "${ver%%-*}"
}

# Package version converts a vault version string into a compatible representation for system
# packages.
function version_package() {
    awk '{ gsub("-","~",$1); print $1 }' <<<"$(version)"
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
    copy-binary)
        copy_binary
        ;;
    date)
        build_date
        ;;
    prepare-ent-legal)
        prepare_ent_legal
        ;;
    prepare-ce-legal)
        prepare_ce_legal
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
