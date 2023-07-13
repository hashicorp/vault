#!/bin/bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${RUN_ID:?} ]
[ ${REPO:?} ]
[ ${PR_NUMBER:?} ]
# list of build jobs
[ ${BUILD_OTHER:?} ]
[ ${BUILD_LINUX:?} ]
[ ${BUILD_DARWIN:?} ]
[ ${BUILD_DOCKER:?} ]
[ ${BUILD_UBI:?} ]
[ ${TEST:?} ]
[ ${TEST_DOCKER_K8S:?} ]

# listing out all of the jobs with the status
jobs=( "build-other:$BUILD_OTHER" "build-linux:$BUILD_LINUX" "build-darwin:$BUILD_DARWIN" "build-docker:$BUILD_DOCKER" "build-ubi:$BUILD_UBI" "test:$TEST" "test-docker-k8s:$TEST_DOCKER_K8S" )

# there is a case where even if a job is failed, it reports as cancelled. So, we look for both.
failed_jobs=()
for job in "${jobs[@]}";do
  if [[ "$job" == *"failure"* || "$job" == *"cancelled"* ]]; then
    failed_jobs+=("$job")
  fi
done

# Create a comment to be posted on the PR
# This comment reports failed jobs and the url to the failed workflow
new_body="build failed for these jobs: ${failed_jobs[*]}. Please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"

# We only want for the GH bot to place one comment to report build failures
# and if we rerun a job, that comment needs to be updated.
# Let's try to find if the GH bot has placed a similar comment
comment_id=$(gh api \
               -H "Accept: application/vnd.github+json" \
               -H "X-GitHub-Api-Version: 2022-11-28" \
               /repos/hashicorp/"$REPO"/issues/"$PR_NUMBER"/comments | jq -r '.[] | select (.body | contains("build failed for these job")) | .id')

if [[ "$comment_id" != "" ]]; then
  # update the comment with the new body
  gh api \
    --method PATCH \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    /repos/hashicorp/"$REPO"/issues/comments/"$comment_id" \
    -f body="$new_body"
else
  # create a comment with the new body
  gh api \
    --method POST \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    /repos/hashicorp/"$REPO"/issues/"$PR_NUMBER"/comments \
    -f body="$new_body"
fi
