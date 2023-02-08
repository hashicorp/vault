#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.
set -e

GO_CMD=${GO_CMD:-go}

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
SOURCE_DIR=$( dirname "$SOURCE" )
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$SOURCE_DIR/.." && pwd )"

# Change into that directory
cd "$DIR"

# Set build tags
BUILD_TAGS="${BUILD_TAGS:-"vault"}"

# Get the git commit
GIT_COMMIT="$("$SOURCE_DIR"/ci-helper.sh revision)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

BUILD_DATE="$("$SOURCE_DIR"/ci-helper.sh date)"

# If its dev mode, only build for ourself
if [ "${VAULT_DEV_BUILD}x" != "x" ] && [ "${XC_OSARCH}x" == "x" ]; then
    XC_OS=$(${GO_CMD} env GOOS)
    XC_ARCH=$(${GO_CMD} env GOARCH)
    XC_OSARCH=$(${GO_CMD} env GOOS)/$(${GO_CMD} env GOARCH)
elif [ "${XC_OSARCH}x" != "x" ]; then
    IFS='/' read -ra SPLITXC <<< "${XC_OSARCH}"
	DEV_PLATFORM="./pkg/${SPLITXC[0]}_${SPLITXC[1]}"
fi

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64"}
XC_OS=${XC_OS:-linux darwin windows freebsd openbsd netbsd solaris}
XC_OSARCH=${XC_OSARCH:-"linux/386 linux/amd64 linux/arm linux/arm64 darwin/386 darwin/amd64 darwin/arm64 windows/386 windows/amd64 freebsd/386 freebsd/amd64 freebsd/arm openbsd/386 openbsd/amd64 openbsd/arm netbsd/386 netbsd/amd64 solaris/amd64"}

GOPATH=${GOPATH:-$(${GO_CMD} env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac

# Delete the old dir
echo "==> Removing old directory..."
rm -f bin/*
rm -rf pkg/*
mkdir -p bin/

# Build!
# If GOX_PARALLEL_BUILDS is set, it will be used to add a "-parallel=${GOX_PARALLEL_BUILDS}" gox parameter
echo "==> Building..."
gox \
    -osarch="${XC_OSARCH}" \
    -gcflags "${GCFLAGS}" \
    -ldflags "${LD_FLAGS}-X github.com/hashicorp/vault/sdk/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X github.com/hashicorp/vault/sdk/version.BuildDate=${BUILD_DATE}" \
    -output "pkg/{{.OS}}_{{.Arch}}/vault" \
    ${GOX_PARALLEL_BUILDS+-parallel="${GOX_PARALLEL_BUILDS}"} \
    -tags="${BUILD_TAGS}" \
    -gocmd="${GO_CMD}" \
    .

# Move all the compiled things to the $GOPATH/bin
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Copy our OS/Arch to the bin/ directory
DEV_PLATFORM=${DEV_PLATFORM:-"./pkg/$(${GO_CMD} env GOOS)_$(${GO_CMD} env GOARCH)"}
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 1 -type f); do
    cp ${F} bin/
    rm -f ${MAIN_GOPATH}/bin/vault
    cp ${F} ${MAIN_GOPATH}/bin/
done

if [ "${VAULT_DEV_BUILD}x" = "x" ]; then
    # Zip and copy to the dist dir
    echo "==> Packaging..."
    for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
        OSARCH=$(basename ${PLATFORM})
        echo "--> ${OSARCH}"

        pushd $PLATFORM >/dev/null 2>&1
        zip ../${OSARCH}.zip ./*
        popd >/dev/null 2>&1
    done
fi

# Done!
echo
echo "==> Results:"
ls -hl bin/
