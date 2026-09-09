#!/bin/bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -euo pipefail

# Replaces (or inserts) a marker-delimited section within a PR's description,
# leaving the rest of the body untouched. Unlike a comment, the description
# is always visible at the top of the PR and isn't buried by later activity.
#
# Always strips any existing marked section first and re-appends fresh at the
# bottom, rather than editing it in place — otherwise the section stays stuck
# wherever it first landed (e.g. still at the top from an old run), instead
# of tracking wherever "the bottom" currently is as the description grows.
#
# SECTION_FILE is optional: omit it to just remove the section (e.g. when a
# trigger label is removed) without replacing it with anything.
#
# When given, SECTION_FILE (not the section text itself) is taken as an
# argument because BSD awk (macOS's default) rejects embedded newlines in a
# `-v` value; reading the replacement via `getline` from a file works on both
# BSD awk and gawk.
function update_pr_body_section {
  REPO=$1
  PR_NUMBER=$2
  MARKER_START=$3
  MARKER_END=$4
  SECTION_FILE=${5:-}

  current_body=$(gh api \
                    -H "Accept: application/vnd.github+json" \
                    -H "X-GitHub-Api-Version: 2022-11-28" \
                    /repos/"$GITHUB_REPOSITORY_OWNER"/"$REPO"/pulls/"$PR_NUMBER" --jq '.body // ""')

  # Strip \r first: if the body ever has CRLF line endings, an exact-line
  # match against $MARKER_START/$MARKER_END would otherwise silently fail,
  # leaving the old section in place and letting sections pile up over time.
  body_without_section=$(awk -v start="$MARKER_START" -v end="$MARKER_END" '
    { gsub(/\r/, "") }
    $0 == start { skipping = 1; next }
    $0 == end { skipping = 0; next }
    skipping { next }
    { print }
  ' <<<"$current_body")

  if [[ -n "$SECTION_FILE" ]]; then
    new_body=$(printf '%s\n\n%s' "$body_without_section" "$(cat "$SECTION_FILE")")
  else
    new_body="$body_without_section"
  fi

  gh api \
    --method PATCH \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    /repos/"$GITHUB_REPOSITORY_OWNER"/"$REPO"/pulls/"$PR_NUMBER" \
    -f body="$new_body"
}
