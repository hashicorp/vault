#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

EVENT_TYPE="${1}" # GH event type (pull_request)
REF_NAME="${2}"   # branch reference that triggered the workflow
BASE_REF="${3}"   # PR branch base ref

if [[ "$EVENT_TYPE" == "pull_request" ]]; then

  # The following environment variables must be set to continue
  [ ${GITHUB_TOKEN:?} ]
  [ ${REPO:?} ]
  [ ${PR_NUMBER:?} ]

  # https://github.com/hashicorp/vault/pull/{PR_NUMBER}/files --> files in the PR
  # https://github.com/hashicorp/vault/tree/{REPO}/{BRANCH}/CHANGELOG.md ---> base URL for the files

  # Pull in the constant and function definitions
  source "constants.sh"
  for FNAME in $(ls functions) ; do source "functions/$FNAME"; done

  # Step 1: Pull in changelog file
  getChangelog "${CHANGELOG_URL}" "${LOCAL_CHANGELOG}"

  # Step 2: Find the most recent major version number if none was provided
  if [[ -z ${1} ]] ; then 
    TARGET_VERSION=$( getCurrentVersion ${LOCAL_CHANGELOG} )
  else
    TARGET_VERSION="${1}"
  fi

  # Step 3: Parse entries related to the current version
  buildPartials ${TARGET_VERSION}

  # Cleanup
  for FNAME in $(ls output) ; do rm "output/$FNAME"; done
  for FNAME in $(ls downloads) ; do rm "downloads/$FNAME"; done

  unset FNAME
  unset TARGET_VERSION

fi