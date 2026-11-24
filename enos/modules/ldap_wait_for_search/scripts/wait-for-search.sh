#!/usr/bin/env bash
# Copyright IBM Corp. 2026
# SPDX-License-Identifier: BUSL-1.1

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$LDAP_BASE_DN" ]] && fail "LDAP_BASE_DN env variable has not been set"
[[ -z "$LDAP_BIND_DN" ]] && fail "LDAP_BIND_DN env variable has not been set"
[[ -z "$LDAP_HOST" ]] && fail "LDAP_HOST env variable has not been set"
[[ -z "$LDAP_PORT" ]] && fail "LDAP_PORT env variable has not been set"
[[ -z "$LDAP_PASSWORD" ]] && fail "LDAP_PASSWORD env variable has not been set"
[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"

# NOTE: An LDAP_QUERY is not technically required here as long as we have a base
# DN and bind DN to search and bind to.

search="ldap://${LDAP_HOST}:${LDAP_PORT} -b ${LDAP_BASE_DN} -D ${LDAP_BIND_DN} -w ${LDAP_PASSWORD} -s base ${LDAP_QUERY}"
safe_search="ldap://${LDAP_HOST}:${LDAP_PORT} -b ${LDAP_BASE_DN} -D ${LDAP_BIND_DN} -s base ${LDAP_QUERY}"
echo "Running test search: ${safe_search}"

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
test_search_out=""
test_search_res=""
declare -i tries=0
while [ "$(date +%s)" -lt "$end_time" ]; do
  tries+=1
  if test_search_out=$(eval "ldapsearch -x -H ${search}" 2>&1); then
    test_search_res=0
    break
  fi

  test_search_res=$?
  echo "Test search failed!, search: ${safe_search}, attempt: ${tries}, exit code: ${test_search_res}, error: ${test_search_out}, retrying..."
  sleep "$RETRY_INTERVAL"
done

if [ "$test_search_res" -ne 0 ]; then
  echo "Timed out waiting for search!, search: ${safe_search}, attempt: ${tries}, exit code: ${test_search_res}, error: ${test_search_out}" 2>&1
  # Exit with the ldapsearch exit code so we bubble that up to error diagnostic
  exit "$test_search_res"
fi
