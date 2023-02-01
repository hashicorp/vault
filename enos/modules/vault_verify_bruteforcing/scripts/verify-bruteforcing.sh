#!/usr/bin/env bash

set -e

function retry {
  local retries=$1
  shift
  local count=0

  until "$@"; do
    exit=$?
    wait=$((2 ** count))
    echo $wait
    count=$((count + 1))
    if [ "$count" -lt "$retries" ]; then
      sleep "$wait"
    else
      return "$exit"
    fi
  done

  return 0
}

function fail {
	echo "$1" 1>&2
	exit 1
}

binpath=${VAULT_INSTALL_DIR}/vault

fail() {
  echo "$1" 1>&2
  return 1
}

test -x "$binpath" || fail "unable to locate vault binary at $binpath"

# To keep the authentication method and module verification consistent between all
# Enos scenarios we authenticate using testuser created by vault_verify_write_data module
retry 5 $binpath login -method=userpass username=testuser password=passuser1

# login_testuser is to login to testuser using correct password 
login_testuser () {
cat << EOF > payload.json
{
"password": "passuser1"
}
EOF

leader_address=$(curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/sys/leader" | jq '.leader_address' | sed 's/\"//g')
user_login_output=$(curl --header "X-Vault-Token: $VAULT_TOKEN" --data @payload.json "$leader_address/v1/auth/userpass/login/testuser")
  if [[ "$user_login_output" == *"invalid username or password"* ]] || [[ "$user_login_output" == *"permission denied"* ]]; then
    fail "expected user login to be successful"
  fi
}

# login_wrong_password is to login to testuser using incorrect password 
login_wrong_password () {
cat << EOF > payload.json
{
"password": "wrongPassword"
}
EOF

leader_address=$(curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/sys/leader" | jq '.leader_address' | sed 's/\"//g')
  # login using wrong password
  user_login_output=$(curl --header "X-Vault-Token: $VAULT_TOKEN" --data @payload.json "$leader_address/v1/auth/userpass/login/testuser")
  if ! [[ "$user_login_output" == *"invalid username or password"* || *"permission denied"* ]]; then
    fail "expected user to be return invalid credentials or permission denied error"
  fi
}

# check_user_locked checks if the user is locked 
check_user_locked () {
cat << EOF > payload.json
{
"password": "passuser1"
}
EOF

  leader_address=$(curl -H "X-Vault-Request: true" -H "X-Vault-Token: $VAULT_TOKEN" "$VAULT_ADDR/v1/sys/leader" | jq '.leader_address' | sed 's/\"//g')
  # login using correct password
  user_login_output=$(curl --header "X-Vault-Token: $VAULT_TOKEN" --data @payload.json "$leader_address/v1/auth/userpass/login/testuser")
  # if user locked, we get permission denied error 
  if [[ "$user_login_output" != *"permission denied"* ]]; then
    fail "expected user to be locked with permission denied error"
  fi
}

# Test 1: Verifying user gets locked out for lockout duration after multiple 
# failure login attempts (in this case, default lockout threshold of 5 attempts)

# Change lockout duration value for testuser using auth tune for faster testing of this feature
# Default lockout threshold is 5 attempts, lockout duration is 15 minutes, lockout counter reset is 15 minutes 
retry 5 $binpath  auth tune -user-lockout-duration=30s userpass/

# Login to testuser using wrong password for 5 times to lock 
for i in {1..5}
do
   login_wrong_password
done

# Login to testuser to check if the user is locked
retry 5 check_user_locked

# Try to login, should be a successful login after trying for 30s(lockout duration)
retry 8 login_testuser

# Test 2: Verify that the user does not get locked, when user lockout feature is
# disabled

# Disable lockout using auth tune and verify if user lockout is disabled
retry 5 $binpath  auth tune -user-lockout-disable=true userpass/

# Login to testuser using wrong password for 5 times to lock 
for i in {1..5}
do
   login_wrong_password
done

# Try to login, should be a successful as the user is unlocked
retry 5 login_testuser

# Test 3: Verify that the user can be unlocked using unlock api if locked 

# Enable user lockout feature
retry 5 $binpath  auth tune -user-lockout-disable=false userpass/

# Setting lockout duration to 15 mins to ensure that the user is not unlocked due
# to 30s lockout duration that we set before
retry 5 $binpath  auth tune -user-lockout-duration=15m userpass/

# Get the userpass accessor for the mount, we need this for unlock api
userpass_accessor="$($binpath auth list | awk '/^userpass/ {print $3}')"

# Login to testuser using wrong password for 5 times to lock 
for i in {1..5}
do
   login_wrong_password
done

# Login to testuser to check if the user is locked
retry 5 check_user_locked

# Unlock the user 
retry 5 $binpath  write -force sys/locked-users/$userpass_accessor/unlock/testuser

# Try to login, should be a successful as the user is unlocked
retry 8 login_testuser
