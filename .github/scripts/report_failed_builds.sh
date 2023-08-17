#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


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
if [ ${#failed_jobs[@]} -eq 0 ]; then
  new_body="Build Results:
All builds succeeded! :white_check_mark:"
else
  new_body="Build Results:
Build failed for these jobs: ${failed_jobs[*]}. Please refer to this workflow to learn more: https://github.com/hashicorp/vault/actions/runs/$RUN_ID"
fi


source ./.github/scripts/gh_comment.sh

update_or_create_comment "$REPO" "$PR_NUMBER" "Build Results:" "$new_body"