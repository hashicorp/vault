#! /bin/bash
sed -i "s/RENEW_TOKEN/$RENEW_TOKEN/g" payload.json
sed -i "s/INCREMENT_VALUE/$INCREMENT_VALUE/g" payload.json
exec curl -L --header "X-Vault-Token: $ROOT_TOKEN" --request POST --data @payload.json http://$URL/v1/auth/token/renew
