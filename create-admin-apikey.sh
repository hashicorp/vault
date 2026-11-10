#!/bin/bash

# Create the admin service ID and API key used by the Vault ibmcloud auth method,
# with the least access needed for login verification and root rotation
# (`vault write -f auth/ibmcloud/config/rotate-root`).
#
# Access granted to the new service ID (IBM Cloud IAM):
#   1. Viewer   on iam-identity, all service IDs
#        Login: GET /v1/serviceids/{id} to confirm a service ID belongs to the account.
#   2. Operator on iam-identity, resource type "serviceid", only the new service ID
#        Rotation: create, read and delete API keys of this service ID only.
#   3. Viewer   on user-management
#        Login: GET /v2/accounts/{account}/users/{id} for "User" identities.
#        Skip with --skip-user-lookup if only service IDs will log in.
#
# Requires: ibmcloud CLI (already logged in to the target account), jq, curl.
#
# The API key is written to a file (mode 600) and is never printed.

set -euo pipefail
umask 077

NAME="vault-ibmcloud-admin"
OUT=""
SKIP_USERS=0
SKIP_VERIFY=0

usage() {
    cat <<EOF
Usage: $0 [-n NAME] [-o FILE] [--skip-user-lookup] [--skip-verify]

Options:
  -n, --name NAME       Service ID name (default: ${NAME}). Must not already exist.
  -o, --output FILE     Where to save the API key JSON (default: \$HOME/<NAME>.apikey.json).
                        Must not already exist.
  --skip-user-lookup    Do not grant Viewer on user-management (service IDs only).
  --skip-verify         Do not test the new key against IAM.

Environment:
  ACCOUNT_ID            Account to use (default: the account targeted by ibmcloud).
  IAM_ENDPOINT          IAM endpoint (default: derived from 'ibmcloud api').
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
        -o|--output)
            [ $# -ge 2 ] || die "$1 needs a value"
            OUT="$2"; shift 2 ;;
        --skip-user-lookup) SKIP_USERS=1; shift ;;
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

echo "Account:      ${ACCOUNT_ID}"
echo "IAM endpoint: ${IAM_ENDPOINT}"
echo "Service ID:   ${NAME}"
echo "Key file:     ${OUT}"
echo ""

if ibmcloud iam service-id "${NAME}" --output json >/dev/null 2>&1; then
    die "service ID '${NAME}' already exists; choose another name with -n"
fi

# --- Create service ID, policies and key (rolled back on failure) -------------

SID=""
KEY_CREATED=0

cleanup() {
    local rc=$?
    rm -f "${BODY:-}"
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
BODY=$(mktemp)
trap cleanup EXIT

echo "Creating service ID..."
SID=$(ibmcloud iam service-id-create "${NAME}" \
    -d "Vault ibmcloud auth admin (login verification and root rotation)" \
    --output json | jq -r '.id // empty')
[ -n "${SID}" ] || die "service ID creation returned no ID"
echo "  ${SID}"

echo "Granting Viewer on iam-identity (all service IDs)..."
ibmcloud iam service-policy-create "${SID}" --roles Viewer \
    --service-name iam-identity -q -f >/dev/null

echo "Granting Operator on iam-identity, restricted to ${SID}..."
ibmcloud iam service-policy-create "${SID}" --roles Operator \
    --service-name iam-identity --resource-type serviceid --resource "${SID}" -q -f >/dev/null

if [ "${SKIP_USERS}" -eq 0 ]; then
    echo "Granting Viewer on user-management..."
    ibmcloud iam service-policy-create "${SID}" --roles Viewer \
        --service-name user-management -q -f >/dev/null
fi

echo "Creating API key (saved to ${OUT}, not printed)..."
ibmcloud iam service-api-key-create "${NAME}-key" "${SID}" \
    -d "Vault ibmcloud auth admin key" --file "${OUT}" -f -q >/dev/null
[ -s "${OUT}" ] || die "API key file ${OUT} was not written"
chmod 600 "${OUT}"
APIKEY=$(jq -r '.apikey // empty' "${OUT}")
[ -n "${APIKEY}" ] || die "no apikey field in ${OUT}"
KEY_CREATED=1

# --- Verify (kept even on failure: policies can take a while to propagate) ----

TOKEN=""

# curl with the secret kept off the command line (-K - reads config from stdin)
iam_call() {
    local method=$1 path=$2 data=${3:-}
    if [ -n "${data}" ]; then
        printf 'header = "Authorization: Bearer %s"\n' "${TOKEN}" | \
            curl -s -o "${BODY}" -w '%{http_code}' -K - -X "${method}" "${IAM_ENDPOINT}${path}" \
                -H 'Content-Type: application/json' -d "${data}" || true
    else
        printf 'header = "Authorization: Bearer %s"\n' "${TOKEN}" | \
            curl -s -o "${BODY}" -w '%{http_code}' -K - -X "${method}" "${IAM_ENDPOINT}${path}" || true
    fi
}

key_payload() {
    jq -nc --arg n "vault-permission-check" --arg i "$1" \
        '{name:$n, description:"temporary key, safe to delete", iam_id:$i, account_id:""}'
}

FAILED=0
check() { # description expected actual
    if [ "$2" = "$3" ]; then
        echo "  PASS  $1 (HTTP $3)"
    else
        echo "  FAIL  $1 (expected HTTP $2, got $3)"
        FAILED=1
    fi
}

verify() {
    echo ""
    echo "Verifying against ${IAM_ENDPOINT} ..."

    TOKEN=$(printf 'data-urlencode = "apikey=%s"\n' "${APIKEY}" | \
        curl -s -K - -X POST "${IAM_ENDPOINT}/identity/token" \
            --data-urlencode "grant_type=urn:ibm:params:oauth:grant-type:apikey" | \
        jq -r '.access_token // empty')
    [ -n "${TOKEN}" ] || { echo "  FAIL  could not exchange the new key for an IAM token"; FAILED=1; return; }
    echo "  PASS  key exchanges for an IAM token"

    # Policies can take a few seconds to take effect.
    local code="" i
    for i in 1 2 3 4 5 6 7 8; do
        code=$(iam_call GET "/v1/serviceids/${SID}")
        [ "${code}" = "200" ] && break
        sleep 5
    done
    check "read own service ID (login lookup)" 200 "${code}"

    local other
    other=$(ibmcloud iam service-ids --output json 2>/dev/null | \
        jq -r --arg me "${SID}" '[.[] | select(.id != $me)][0].id // empty' 2>/dev/null || true)
    if [ -n "${other}" ]; then
        code=$(iam_call GET "/v1/serviceids/${other}")
        check "read another service ID (login lookup)" 200 "${code}"
    else
        echo "  SKIP  no other service ID in the account to read"
    fi

    # Rotation: create, read and delete a key for the admin service ID. Vault sends
    # the service ID exactly as configured (no iam- prefix), so test that form.
    local created_id=""
    code=$(iam_call POST "/v1/apikeys" "$(key_payload "${SID}")")
    if [ "${code}" != "201" ]; then
        local code_iam
        code_iam=$(iam_call POST "/v1/apikeys" "$(key_payload "iam-${SID}")")
        if [ "${code_iam}" = "201" ]; then
            echo "  WARN  key creation worked only with iam_id=iam-${SID}, but Vault sends '${SID}' (HTTP ${code})"
            echo "        rotate-root will fail until CreateAPIKey in client.go sends the iam- prefixed ID"
            code="${code_iam}"
        fi
    fi
    check "create an API key for own service ID (rotation)" 201 "${code}"
    if [ "${code}" = "201" ]; then
        created_id=$(jq -r '.id // empty' "${BODY}")
        code=$(iam_call GET "/v1/apikeys/${created_id}")
        check "read that API key by ID (rotation)" 200 "${code}"
        code=$(iam_call DELETE "/v1/apikeys/${created_id}")
        if [ "${code}" = "204" ] || [ "${code}" = "200" ]; then
            echo "  PASS  delete that API key (rotation) (HTTP ${code})"
        else
            echo "  FAIL  delete that API key (rotation) (expected 204 or 200, got ${code})"
            echo "        remove key ${created_id} manually"
            FAILED=1
        fi
    fi
    # Least privilege: creating a key for a different service ID must be refused.
    if [ -n "${other}" ]; then
        code=$(iam_call POST "/v1/apikeys" "$(key_payload "${other}")")
        if [ "${code}" = "201" ]; then
            iam_call DELETE "/v1/apikeys/$(jq -r '.id // empty' "${BODY}")" >/dev/null || true
            echo "  FAIL  create a key for ANOTHER service ID was allowed; access is broader than intended"
            echo "        check other policies that apply to ${SID} or to the account"
            FAILED=1
        else
            check "create a key for another service ID is refused" 403 "${code}"
        fi
    fi
}

if [ "${SKIP_VERIFY}" -eq 0 ]; then
    verify
fi

echo ""
if [ "${FAILED}" -ne 0 ]; then
    echo "Service ID ${SID} and its key were created, but verification failed (see above)."
    echo "Policies can take a minute to propagate; re-test, or delete with:"
    echo "  ibmcloud iam service-id-delete ${SID} -f"
    exit 1
fi

cat <<EOF
Done.

  Service ID: ${SID}
  Account:    ${ACCOUNT_ID}
  Key file:   ${OUT}

Configure Vault (run with VAULT_ADDR and VAULT_TOKEN set):

  ./demo-setup.sh "\$(jq -r .apikey '${OUT}')" ${ACCOUNT_ID} ${SID}

Rotate the key from Vault:

  ./bin/vault write -f auth/ibmcloud/config/rotate-root

Delete the key file once Vault has it; after rotation Vault holds the only valid key.
EOF
