#!/usr/bin/env bash

binpath=${VAULT_INSTALL_DIR}/vault

IFS="," read -a keys <<< ${UNSEAL_KEYS}

function fail() {
	echo "$1" 1>&2
	exit 1
}
count=0
retries=5
while :; do
   for key in ${keys[@]}; do

    # Check the Vault seal status
    seal_status=$($binpath status -format json | jq '.sealed')
    
    if [[ "$seal_status" == "true" ]]; then
      echo "running unseal with $key count $count with retry $retry" >> /tmp/unseal_script.out
      $binpath operator unseal $key > /dev/null 2>&1
    else
      exit 0
    fi
  done

  wait=$((1 ** count))
  count=$((count + 1))
  if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
  else
      fail "failed to unseal node"
  fi
done
