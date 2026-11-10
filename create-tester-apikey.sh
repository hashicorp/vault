#!/bin/bash

# Create a tester service ID and API key for logging in to Vault's ibmcloud auth method.
#
# Without --access-group the service ID gets NO IAM policies and NO group
# membership. That is enough for Vault to identify it: the token carries the
# account, and the service ID can be matched with bound_service_ids.
#
# With --access-group the service ID is added to that access group, so it gets
# whatever access the group grants and can log in to roles that set
# bound_access_group_ids to the same group.
#
# Requires: ibmcloud CLI (already logged in to the target account), jq, curl.
#
# The API key is written to a file (mode 600) and is never printed.

set -euo pipefail
umask 077

NAME="vault-ibmcloud-tester"
GROUP=""
OUT=""
SKIP_VERIFY=0

usage() {
    cat <<EOF
Usage: $0 [-n NAME] [-g ACCESS_GROUP] [-o FILE] [--skip-verify]

Options:
  -n, --name NAME            Service ID name (default: ${NAME}). Must not already exist.
  -g, --access-group GROUP   Access group ID (AccessGroupId-...) or name. The service ID is
                             added to it. Without this the service ID has no access at all.
  -o, --output FILE          Where to save the API key JSON (default: \$HOME/<NAME>.apikey.json).
                             Must not already exist.
  --skip-verify              Do not test the new key against IAM.

Environment:
  ACCOUNT_ID                 Account to use (default: the account targeted by ibmcloud).
  IAM_ENDPOINT               IAM endpoint (default: derived from 'ibmcloud api').
EOF
}

die() {
    echo "ERROR: $*" >&2
    exit 1
}

while [ $# -gt 0 ]; do
    case "$1" in
        -n|--name)
            [ $# -ge 2 ] || die "$1 needs a value"
            NAME="$2"; shift 2 ;;
        -g|--access-group)
            [ $# -ge 2 ] || die "$1 needs a value"
            GROUP="$2"; shift 2 ;;
        -o|--output)
            [ $# -ge 2 ] || die "$1 needs a value"
            OUT="$2"; shift 2 ;;
        --skip-verify) SKIP_VERIFY=1; shift ;;
        -h|--help) usage; exit 0 ;;
        *) usage >&2; die "unknown argument: $1" ;;
    esac
done

[ -n "${OUT}" ] || OUT="${HOME}/${NAME}.apikey.json"

for tool in ibmcloud jq curl; do
    command -v "${tool}" >/dev/null 2>&1 || die "${tool} is required but not installed"
done

[ ! -e "${OUT}" ] || die "output file ${OUT} already exists; choose another with -o"

# --- Target account and IAM endpoint -----------------------------------------

if [ -z "${ACCOUNT_ID:-}" ]; then
    ACCOUNT_ID=$(ibmcloud target --output json 2>/dev/null | jq -r '.account.guid // empty' 2>/dev/null || true)
fi
if [ -z "${ACCOUNT_ID}" ]; then
    ACCOUNT_ID=$(ibmcloud account show --output json 2>/dev/null | jq -r '.account_id // empty' 2>/dev/null || true)
fi
[ -n "${ACCOUNT_ID}" ] || die "could not determine the account; run 'ibmcloud login' first or set ACCOUNT_ID"

if [ -z "${IAM_ENDPOINT:-}" ]; then
    API_ENDPOINT=$(ibmcloud api 2>/dev/null | tr -d '\r' | awk '/API endpoint:/{print $3}')
    [ -n "${API_ENDPOINT}" ] || die "could not read the API endpoint; set IAM_ENDPOINT"
    IAM_ENDPOINT="https://iam.${API_ENDPOINT#https://}"
fi

# --- Resolve the access group (ID or name) before creating anything -----------

GROUP_ID=""
GROUP_NAME=""
if [ -n "${GROUP}" ]; then
    AG_JSON=$(ibmcloud iam access-groups --output json 2>/dev/null) || die "could not list access groups"
    GROUP_ID=$(printf '%s' "${AG_JSON}" | jq -r --arg g "${GROUP}" \
        '[.[] | select(.id == $g or .name == $g)][0].id // empty')
    GROUP_NAME=$(printf '%s' "${AG_JSON}" | jq -r --arg g "${GROUP}" \
        '[.[] | select(.id == $g or .name == $g)][0].name // empty')
    [ -n "${GROUP_ID}" ] || die "access group '${GROUP}' not found in account ${ACCOUNT_ID}"
fi

echo "Account:      ${ACCOUNT_ID}"
echo "IAM endpoint: ${IAM_ENDPOINT}"
echo "Service ID:   ${NAME}"
if [ -n "${GROUP_ID}" ]; then
    echo "Access group: ${GROUP_NAME} (${GROUP_ID})"
else
    echo "Access group: none (the service ID will have no IAM access)"
fi
echo "Key file:     ${OUT}"
echo ""

if ibmcloud iam service-id "${NAME}" --output json >/dev/null 2>&1; then
    die "service ID '${NAME}' already exists; choose another name with -n"
fi

# --- Create service ID, group membership and key (rolled back on failure) -----

SID=""
KEY_CREATED=0
BODY=$(mktemp)

cleanup() {
    local rc=$?
    rm -f "${BODY}"
    if [ "${rc}" -ne 0 ] && [ "${KEY_CREATED}" -eq 0 ]; then
        if [ -n "${SID}" ]; then
            echo "Rolling back: deleting service ID ${SID}" >&2
            ibmcloud iam service-id-delete "${SID}" -f >/dev/null 2>&1 || \
                echo "WARNING: could not delete ${SID}; remove it manually" >&2
        fi
        rm -f "${OUT}"
    fi
    return "${rc}"
}
trap cleanup EXIT

echo "Creating service ID..."
SID=$(ibmcloud iam service-id-create "${NAME}" \
    -d "Vault ibmcloud auth tester" --output json | jq -r '.id // empty')
[ -n "${SID}" ] || die "service ID creation returned no ID"
echo "  ${SID}"

if [ -n "${GROUP_ID}" ]; then
    echo "Adding to access group ${GROUP_NAME}..."
    ibmcloud iam access-group-service-id-add "${GROUP_NAME}" "${SID}" -q >/dev/null
fi

echo "Creating API key (saved to ${OUT}, not printed)..."
ibmcloud iam service-api-key-create "${NAME}-key" "${SID}" \
    -d "Vault ibmcloud auth tester key" --file "${OUT}" -f -q >/dev/null
[ -s "${OUT}" ] || die "API key file ${OUT} was not written"
chmod 600 "${OUT}"
APIKEY=$(jq -r '.apikey // empty' "${OUT}")
[ -n "${APIKEY}" ] || die "no apikey field in ${OUT}"
KEY_CREATED=1

# --- Verify what Vault will see at login --------------------------------------

FAILED=0

pass() { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=1; }

verify() {
    echo ""
    echo "Verifying against ${IAM_ENDPOINT} ..."

    local token
    token=$(printf 'data-urlencode = "apikey=%s"\n' "${APIKEY}" | \
        curl -s -K - -X POST "${IAM_ENDPOINT}/identity/token" \
            --data-urlencode "grant_type=urn:ibm:params:oauth:grant-type:apikey" | \
        jq -r '.access_token // empty')
    if [ -z "${token}" ]; then
        fail "could not exchange the new key for an IAM token"
        return
    fi
    pass "key exchanges for an IAM token"

    # What Vault reads at login: token introspection.
    local intro
    intro=$(printf 'data-urlencode = "token=%s"\n' "${token}" | \
        curl -s -K - -X POST "${IAM_ENDPOINT}/identity/introspect")

    [ "$(printf '%s' "${intro}" | jq -r '.active')" = "true" ] && pass "token is active" || fail "token is not active"
    [ "$(printf '%s' "${intro}" | jq -r '.sub_type')" = "ServiceId" ] && pass "identity type is ServiceId" || fail "unexpected identity type"
    [ "$(printf '%s' "${intro}" | jq -r '.account.bss // .account_id')" = "${ACCOUNT_ID}" ] && \
        pass "token is bound to account ${ACCOUNT_ID}" || fail "token account does not match ${ACCOUNT_ID}"
    [ "$(printf '%s' "${intro}" | jq -r '.iam_id')" = "iam-${SID}" ] && pass "iam_id is iam-${SID}" || fail "unexpected iam_id"

    # What Vault reads for bound_access_group_ids: the caller's own group list.
    if [ -n "${GROUP_ID}" ]; then
        local i code found=0
        for i in 1 2 3 4 5 6; do
            code=$(printf 'header = "Authorization: Bearer %s"\n' "${token}" | \
                curl -s -o "${BODY}" -w '%{http_code}' -K - \
                    "${IAM_ENDPOINT}/v2/groups?account_id=${ACCOUNT_ID}&iam_id=iam-${SID}" || true)
            if [ "${code}" = "200" ] && \
               jq -e --arg g "${GROUP_ID}" 'any(.groups[]?; .id == $g)' "${BODY}" >/dev/null 2>&1; then
                found=1; break
            fi
            sleep 3
        done
        if [ "${found}" -eq 1 ]; then
            pass "service ID is a member of ${GROUP_ID} (as Vault will see it)"
        else
            fail "own group list (HTTP ${code:-?}) does not include ${GROUP_ID}"
        fi
    fi
}

if [ "${SKIP_VERIFY}" -eq 0 ]; then
    verify
fi

echo ""
if [ "${FAILED}" -ne 0 ]; then
    echo "Service ID ${SID} and its key were created, but verification failed (see above)."
    echo "Delete with: ibmcloud iam service-id-delete ${SID} -f"
    exit 1
fi

# --- Next steps ----------------------------------------------------------------

ROLE_CMD="./bin/vault write auth/ibmcloud/role/tester-role \\
      type=apikey \\
      bound_account_id=${ACCOUNT_ID} \\
      bound_service_ids=iam-${SID} \\
      token_policies=default"
if [ -n "${GROUP_ID}" ]; then
    ROLE_CMD="${ROLE_CMD} \\
      bound_access_group_ids=${GROUP_ID}"
fi

cat <<EOF
Done.

  Service ID: ${SID}
  IAM ID:     iam-${SID}
  Account:    ${ACCOUNT_ID}
  Key file:   ${OUT}

Create a role that accepts this tester (run with VAULT_ADDR and VAULT_TOKEN set):

  ${ROLE_CMD}

Log in:

  ./bin/vault login -method=ibmcloud role=tester-role type=apikey \\
      apikey="\$(jq -r .apikey '${OUT}')" iam_endpoint=${IAM_ENDPOINT}

Remove the tester when done:

  ibmcloud iam service-id-delete ${SID} -f && rm -f '${OUT}'
EOF
