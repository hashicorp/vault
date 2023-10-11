#!/bin/bash
#  This script validates if the git diff contains only docs/ui changes

event_type=$1 # GH event type (pull_request)
ref_name=$2 # branch reference that triggered the workflow
base_ref=$3 # PR branch base ref

contains() {
    target=$1; shift
    for i; do
        if [[ "$i" == "$target" ]]; then
            return 0
        fi
    done
    return 1
}

if [[ "$event_type" == "pull_request" ]]; then
  git fetch --no-tags --prune origin $base_ref
  head_commit="HEAD"
  base_commit="origin/$base_ref"
else
  git fetch --no-tags --prune origin $ref_name
  head_commit=$(git log origin/$ref_name --oneline | head -1 | awk '{print $1}')
  base_commit=$(git log origin/$ref_name --oneline | head -2 | awk 'NR==2 {print $1}')
fi

# git diff with ... shows the differences between base_commit and head_commit starting at the last common commit
# excluding the changelog directory
changed_dir=$(git diff $base_commit...$head_commit --name-only | awk -F"/" '{ print $1}' | uniq | sed '/changelog/d')
change_count=$(git diff $base_commit...$head_commit --name-only | awk -F"/" '{ print $1}' | uniq | sed '/changelog/d' | wc -l)

# There are 4 main conditions to check:
#
# 1. more than two changes found, set the flags to false
# 2. doc only change
# 3. ui only change
# 4. two changes found, if either doc or ui does not exist in the changes, set both flags to false

if [[ $change_count -gt 2 ]]; then
  echo "is_docs_change=false" >> "$GITHUB_OUTPUT"
  echo "is_ui_change=false" >> "$GITHUB_OUTPUT"
elif [[ $change_count -eq 1 && "$changed_dir" == "website" ]]; then
  echo "is_docs_change=true" >> "$GITHUB_OUTPUT"
  echo "is_ui_change=false" >> "$GITHUB_OUTPUT"
elif [[ $change_count -eq 1 && "$changed_dir" == "ui" ]]; then
  echo "is_ui_change=true" >> "$GITHUB_OUTPUT"
  echo "is_docs_change=false" >> "$GITHUB_OUTPUT"
else
  if ! contains "website" ${changed_dir[@]} || ! contains "ui" ${changed_dir[@]}; then
   echo "is_docs_change=false" >> "$GITHUB_OUTPUT"
   echo "is_ui_change=false" >> "$GITHUB_OUTPUT"
  else
    echo "is_docs_change=true" >> "$GITHUB_OUTPUT"
    echo "is_ui_change=true" >> "$GITHUB_OUTPUT"
  fi
fi
