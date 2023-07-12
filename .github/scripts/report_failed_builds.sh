#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR:?} ]
[ ${BUILD_OTHER:?} ]
[ ${BUILD_LINUX:?} ]
[ ${BUILD_DARWIN:?} ]
[ ${BUILD_DOCKER:?} ]
[ ${BUILD_UBI:?} ]
[ ${TEST:?} ]
[ ${TEST_DOCKER_K8S:?} ]

#failed_tests_ids=$(gh api \
#  -H "Accept: application/vnd.github+json" \
#  -H "X-GitHub-Api-Version: 2022-11-28" \
#  /repos/hashicorp/"$REPO"/actions/runs/"$RUN_ID"/attempts/1/jobs | jq -r '[.jobs[] | select(.name | startswith("Darwin") or startswith("Linux") or startswith("Other") ) | select(.conclusion == "failure") | .id][0]')
#
#extracted_error=$(gh run view --job $failed_tests_ids --log | grep "\#\#\[error\]")
jobs=( "build-other:$BUILD_OTHER" "build-linux:$BUILD_LINUX" "build-darwin:$BUILD_DARWIN" "build-docker:$BUILD_DOCKER" "build-ubi:$BUILD_UBI" "test:$TEST" "test-docker-k8s:$TEST_DOCKER_K8S" )
failed_jobs=()
for job in "${jobs[@]}";do
  if [[ "$job" == *"failure"* ]]; then
    failed_jobs+=("$job")
  fi
done

gh pr comment "$PR" --body "build failed for these jobs: $failed_jobs. Please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"
