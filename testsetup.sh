#!/bin/sh

vault mount -description="RDS DEV" -path=rds.dev mysql
vault write rds.dev/config/connection connection_url="root:lco9Cwuoh64b97FW4nUL@tcp(rds.dev.crosschx.com:3306)/"
vault write rds.dev/config/lease lease=10s lease_max=24h
#vault write rds.dev/roles/identity-api-dev revoke_sql="REVOKE ALL PRIVILEGES, GRANT OPTION FROM '{{name}}'@'10.0.0.1'; DROP USER '{{name}}'@'10.0.0.1';" sql="CREATE USER '{{name}}'@'10.0.0.1' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'10.0.0.1';"
vault write rds.dev/roles/identity-api-dev sql="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}';GRANT SELECT ON *.* TO '{{name}}'@'%';"
vault read rds.dev/roles/identity-api-dev
date ; vault read rds.dev/creds/identity-api-dev
