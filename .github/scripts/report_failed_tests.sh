#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR_NUMBER:?} ]
if [ -z "$TABLE_DATA" ]; then
  echo "Invalid TABLE_DATA"
  exit 1
fi

# Create a comment to be posted on the PR
# This comment reports failed jobs and the url to the failed workflow
BODY="CI Test Results:
${TABLE_DATA}"

source ./.github/scripts/gh_comment.sh

update_or_create_comment "$REPO" "$PR_NUMBER" "CI Test Results:" "$BODY"