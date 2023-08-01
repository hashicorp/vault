#!/bin/bash

license='${license}'
if test $license = "none"; then
  exit 0
fi

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    count=$((count + 1))

    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN='${root_token}'

# Temporary hack until we can make the unseal resource handle legacy license
# setting. If we're running 1.8 and above then we shouldn't try to set a license.
ver=$(${bin_path} version)
if [[ "$(echo "$ver" |awk '{print $2}' |awk -F'.' '{print $2}')" -ge 8 ]]; then
  exit 0
fi

retry 5 ${bin_path} write /sys/license text="$license"
