#!/usr/bin/env bash
set -e

# Get the version from the command line
VERSION=$1
if [ -z $VERSION ]; then
  echo "Please specify a version."
  exit 1
fi

# Make sure we have AWS API keys
if ([ -z $AWS_ACCESS_KEY_ID ] || [ -z $AWS_SECRET_ACCESS_KEY ]) && [ ! -z $HC_RELEASE ]; then
  echo "Please set your AWS access key information in the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars."
  exit 1
fi

if [ -z $NOBUILD ] && [ -z $DOCKER_CROSS_IMAGE ]; then
  echo "Please set the Docker cross-compile image in DOCKER_CROSS_IMAGE"
  exit 1
fi

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that dir because we expect that
cd $DIR

if [ -z $RELBRANCH ]; then
  RELBRANCH=master
fi

# Tag, unless told not to
if [ -z $NOTAG ]; then
  echo "==> Tagging..."
  git commit --allow-empty --gpg-sign=348FFC4C -m "Cut version $VERSION"
  git tag -a -m "Version $VERSION" -s -u 348FFC4C "v${VERSION}" $RELBRANCH
fi

# Build the packages
if [ -z $NOBUILD ]; then
  # This should be a local build of the Dockerfile in the cross dir
  docker run --rm -v "$(pwd)":/gopath/src/github.com/hashicorp/vault -w /gopath/src/github.com/hashicorp/vault ${DOCKER_CROSS_IMAGE}
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
if [ ! -z $HC_RELEASE ]; then
  hc-releases upload $DIR/pkg/dist
  hc-releases publish

  curl -X PURGE https://releases.hashicorp.com/vault/${VERSION}
  for FILENAME in $(find $DIR/pkg/dist -type f); do
    FILENAME=$(basename $FILENAME)
    curl -X PURGE https://releases.hashicorp.com/vault/${VERSION}/${FILENAME}
  done
fi

exit 0
