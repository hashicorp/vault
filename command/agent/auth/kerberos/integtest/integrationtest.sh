#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Instructions
# This integration test is for the Vault Kerberos agent.
# Before running, execute:
#   pip install --quiet requests-kerberos
# Then run this test from Vault's home directory.
#   ./command/agent/auth/kerberos/integtest/integrationtest.sh

if [[ "$OSTYPE" == "darwin"* ]]; then
  base64cmd="base64 -D"
else
  base64cmd="base64 -d"
fi

VAULT_PORT=8200
SAMBA_VER=4.8.12

export VAULT_TOKEN=${VAULT_TOKEN:-myroot}
DOMAIN_ADMIN_PASS=Pa55word!
DOMAIN_VAULT_ACCOUNT=vault_svc
DOMAIN_VAULT_PASS=vaultPa55word!
DOMAIN_USER_ACCOUNT=grace
DOMAIN_USER_PASS=gracePa55word!

SAMBA_CONF_FILE=/srv/etc/smb.conf
DOMAIN_NAME=matrix
DNS_NAME=host
REALM_NAME=MATRIX.LAN
DOMAIN_DN=DC=MATRIX,DC=LAN
TESTS_DIR=/tmp/vault_plugin_tests

function add_user() {

  username="${1}"
  password="${2}"

  if [[ $(check_user ${username}) -eq 0 ]]
  then
    echo "add user '${username}'"

    docker exec $SAMBA_CONTAINER \
      /usr/bin/samba-tool user create \
      ${username} \
      ${password}\
      --configfile=${SAMBA_CONF_FILE}
  fi
}

function check_user() {

  username="${1}"

  docker exec $SAMBA_CONTAINER \
    /usr/bin/samba-tool user list \
    --configfile=${SAMBA_CONF_FILE} \
    | grep -c ${username}
}

function create_keytab() {

  username="${1}"
  password="${2}"

  user_kvno=$(docker exec $SAMBA_CONTAINER \
    bash -c "ldapsearch -H ldaps://localhost -D \"Administrator@${REALM_NAME}\"  -w \"${DOMAIN_ADMIN_PASS}\" -b \"CN=Users,${DOMAIN_DN}\" -LLL \"(&(objectClass=user)(sAMAccountName=${username}))\" msDS-KeyVersionNumber | sed -n 's/^[ \t]*msDS-KeyVersionNumber:[ \t]*\(.*\)/\1/p'")

  docker exec $SAMBA_CONTAINER \
    bash -c "printf \"%b\" \"addent -password -p \"${username}@${REALM_NAME}\" -k ${user_kvno} -e rc4-hmac\n${password}\nwrite_kt ${username}.keytab\" | ktutil"

  docker exec $SAMBA_CONTAINER \
    bash -c "printf \"%b\" \"read_kt ${username}.keytab\nlist\" | ktutil"

  docker exec $SAMBA_CONTAINER \
    base64 ${username}.keytab > ${TESTS_DIR}/integration/${username}.keytab.base64

  docker cp $SAMBA_CONTAINER:/${username}.keytab ${TESTS_DIR}/integration/
}

function main() {
  # make and start vault
  make dev
  vault server -dev -dev-root-token-id=root &

  # start our domain controller
  SAMBA_CONTAINER=$(docker run --net=${DNS_NAME} -d -ti --privileged -e "SAMBA_DC_ADMIN_PASSWD=${DOMAIN_ADMIN_PASS}" -e "KERBEROS_PASSWORD=${DOMAIN_ADMIN_PASS}" -e SAMBA_DC_DOMAIN=${DOMAIN_NAME} -e SAMBA_DC_REALM=${REALM_NAME} "bodsch/docker-samba4:${SAMBA_VER}")
  sleep 15

  # set up users
  add_user $DOMAIN_VAULT_ACCOUNT $DOMAIN_VAULT_PASS
  create_keytab $DOMAIN_VAULT_ACCOUNT $DOMAIN_VAULT_PASS

  add_user $DOMAIN_USER_ACCOUNT $DOMAIN_USER_PASS
  create_keytab $DOMAIN_USER_ACCOUNT $DOMAIN_USER_PASS

  # add the service principals we'll need
  docker exec $SAMBA_CONTAINER \
    samba-tool spn add HTTP/localhost ${DOMAIN_VAULT_ACCOUNT} --configfile=${SAMBA_CONF_FILE}
  docker exec $SAMBA_CONTAINER \
    samba-tool spn add HTTP/localhost:${VAULT_PORT} ${DOMAIN_VAULT_ACCOUNT} --configfile=${SAMBA_CONF_FILE}
  docker exec $SAMBA_CONTAINER \
    samba-tool spn add HTTP/localhost.${DNS_NAME} ${DOMAIN_VAULT_ACCOUNT} --configfile=${SAMBA_CONF_FILE}
  docker exec $SAMBA_CONTAINER \
    samba-tool spn add HTTP/localhost.${DNS_NAME}:${VAULT_PORT} ${DOMAIN_VAULT_ACCOUNT} --configfile=${SAMBA_CONF_FILE}

  # enable and configure the kerberos plugin in Vault
  vault auth enable -passthrough-request-headers=Authorization -allowed-response-headers=www-authenticate kerberos
  vault write auth/kerberos/config keytab=@${TESTS_DIR}/integration/vault_svc.keytab.base64 service_account="vault_svc"
  vault write auth/kerberos/config/ldap binddn=${DOMAIN_VAULT_ACCOUNT}@${REALM_NAME} bindpass=${DOMAIN_VAULT_PASS} groupattr=sAMAccountName groupdn="${DOMAIN_DN}" groupfilter="(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))" insecure_tls=true starttls=true userdn="CN=Users,${DOMAIN_DN}" userattr=sAMAccountName upndomain=${REALM_NAME} url=ldaps://localhost:636

  mkdir -p ${TESTS_DIR}/integration

  echo "
[libdefaults]
  default_realm = ${REALM_NAME}
  dns_lookup_realm = false
  dns_lookup_kdc = true
    ticket_lifetime = 24h
    renew_lifetime = 7d
    forwardable = true
    rdns = false
  preferred_preauth_types = 23
[realms]
  ${REALM_NAME} = {
    kdc = localhost
    admin_server = localhost
    master_kdc = localhost
    default_domain = localhost
  }
" > ${TESTS_DIR}/integration/krb5.conf

  echo "
auto_auth {
        method \"kerberos\" {
                mount_path = \"auth/kerberos\"
                config = {
                        username = \"$DOMAIN_USER_ACCOUNT\"
                        service = \"HTTP/localhost:8200\"
                        realm = \"$REALM_NAME\"
                        keytab_path = \"$TESTS_DIR/integration/grace.keytab\"
                        krb5conf_path = \"$TESTS_DIR/integration/krb5.conf\"
                }
        }
        sink \"file\" {
                config = {
                        path = \"$TESTS_DIR/integration/agent-token.txt\"
                }
        }
}
" > ${TESTS_DIR}/integration/agent.conf

  vault agent -config=${TESTS_DIR}/integration/agent.conf &
  sleep 10
  token=$(cat $TESTS_DIR/integration/agent-token.txt)

  # clean up: kill vault and stop the docker container we started
  kill -9 $(ps aux | grep vault | awk '{print $2}' | head -1) # kill vault server
  kill -9 $(ps aux | grep vault | awk '{print $2}' | head -1) # kill vault agent
  docker rm -f ${SAMBA_CONTAINER}

  # a valid Vault token starts with "s.", check for that
  if [[ $token !=  s.* ]]; then
    echo "received invalid token: $token"
    return 1
  fi
  
  echo "vault kerberos agent obtained auth token: $token"
  echo "exiting successfully!"
  return 0
}
main
