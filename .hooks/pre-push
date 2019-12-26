#!/usr/bin/env bash

remote="$1"
remote_url=$(git remote get-url $remote)

if [[ $remote_url == *"vault-enterprise"* ]]; then
    exit 0
fi

if [ "$remote" = "enterprise" ]; then
    exit 0
fi

if [ "$remote" = "ent" ]; then
    exit 0
fi

if [ -f command/version_ent.go ]; then
    echo "Found enterprise version file while pushing to oss remote"
    exit 1
fi

exit 0
