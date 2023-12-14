#!/bin/bash

function getCurrentVersion {

  if [[ -z "${1}" ]] ; then return ; fi

  local CHANGELOG_FILE="${1}"
  local VERSION_HEADER_REGEX="^## [0-9]\.[0-9]*"
  local TARGET_VERSION
  
  while read -r CURR_LINE
  do
    if [[ "${CURR_LINE}" =~ ${VERSION_HEADER_REGEX} ]] ; then
      TARGET_VERSION=${CURR_LINE//#/}
      TARGET_VERSION=${TARGET_VERSION//[[:space:]]/}
      TARGET_VERSION=$(cut -d '.' -f 1 <<< "${TARGET_VERSION}")"."$(cut -d '.' -f 2 <<< "${TARGET_VERSION}")
      break
    fi
  done < ${CHANGELOG_FILE}

  echo ${TARGET_VERSION}
}