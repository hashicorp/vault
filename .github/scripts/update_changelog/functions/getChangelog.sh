#!/bin/bash

function getChangelog {

  if [[ -z "${1}" ]] ; then return ; fi
  if [[ -z "${2}" ]] ; then return ; fi

  local SOURCE="${1}"
  local TARGET="${2}"
  local URL_REGEX='(https?|ftp|file)://[-[:alnum:]\+&@#/%?=~_|!:,.;]*[-[:alnum:]\+&@#/%=~_|]'

  if [[ ${SOURCE} =~ ${URL_REGEX} ]]
  then # Download the changelog from the web
    wget ${SOURCE} --output-document=${TARGET}
  else # Copy the file from a local source
    cp ${SOURCE} ${TARGET}
  fi
}