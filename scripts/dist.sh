#!/usr/bin/env bash
set -e

# Get the version from the command line
VERSION=$1
if [ -z $VERSION ]; then
    echo "Please specify a version."
    exit 1
fi

# Make sure we have a bintray API key
if [ -z $BINTRAY_API_KEY ] && [ ! -z $BINTRAY ]; then
    echo "Please set your bintray API key in the BINTRAY_API_KEY env var."
    exit 1
fi

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that dir because we expect that
cd $DIR

# Tag, unless told not to
if [ -z $NOTAG ]; then
  echo "==> Tagging..."
  git commit --allow-empty -a --gpg-sign=348FFC4C -m "Cut version $VERSION"
  git tag -a -m "Version $VERSION" -s -u 348FFC4C "v${VERSION}" master
fi

# Build the packages
if [ -z $NOBUILD ]; then
# Yes, jefferai/gox should be parameterized; it's just a local build of the Dockerfile in the cross dir
  docker run --rm -v "$(pwd)":/gopath/src/github.com/hashicorp/vault -w /gopath/src/github.com/hashicorp/vault jefferai/gox:1.5.1
fi

# Zip all the files
rm -rf ./pkg/dist
mkdir -p ./pkg/dist
for FILENAME in $(find ./pkg -mindepth 1 -maxdepth 1 -type f); do
    FILENAME=$(basename $FILENAME)
    cp ./pkg/${FILENAME} ./pkg/dist/vault_${VERSION}_${FILENAME}
done

if [ -z $NOSIGN ]; then
  echo "==> Signing..."
  pushd ./pkg/dist
  rm -f ./vault_${VERSION}_SHA256SUMS*
  shasum -a256 * > ./vault_${VERSION}_SHA256SUMS
  gpg --default-key 348FFC4C --detach-sig ./vault_${VERSION}_SHA256SUMS
  popd
fi

# Upload
if [ ! -z $BINTRAY ]; then
  for ARCHIVE in ./pkg/dist/*; do
    ARCHIVE_NAME=$(basename ${ARCHIVE})

    echo Uploading: $ARCHIVE_NAME
    curl \
        -T ${ARCHIVE} \
        -umitchellh:${BINTRAY_API_KEY} \
        "https://api.bintray.com/content/mitchellh/vault/vault/${VERSION}/${ARCHIVE_NAME}"
  done
fi

exit 0
