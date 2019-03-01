#!/usr/bin/env bash

vault secrets enable database
vault write database/config/my-postgresql-database \
    plugin_name=postgresql-database-plugin \
    allowed_roles="my-role" \
    connection_url="postgresql://postgres:fred@localhost:5432?sslmode=disable"
vault write database/roles/my-role \
    db_name=my-postgresql-database \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}'; \
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
    default_ttl="24h" \
    max_ttl="24h" \
    username="george" \
    rotation_period="10s" \
    rotation_statements="ALTER USER {{name}} WITH PASSWORD '{{password}}';"

