#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR:?} ]
[ ${FAILED_JOBS:?} ]

echo "$FAILED_JOBS"

#failed_tests_ids=$(gh api \
#  -H "Accept: application/vnd.github+json" \
#  -H "X-GitHub-Api-Version: 2022-11-28" \
#  /repos/hashicorp/"$REPO"/actions/runs/"$RUN_ID"/attempts/1/jobs | jq -r '[.jobs[] | select(.name | startswith("Darwin") or startswith("Linux") or startswith("Other") ) | select(.conclusion == "failure") | .id][0]')
#
#extracted_error=$(gh run view --job $failed_tests_ids --log | grep "\#\#\[error\]")

gh pr comment "$PR" --body "build failed: please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"
