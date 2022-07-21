## Setup PostgresSQL in HA (made easy by Bitnami https://bitnami.com/stack/postgresql-ha/containers)


```bash
docker network create postgres-lab --driver bridge
```


```bash
# PG-0 is primary
docker run --detach --name pg-0 \
  --network postgres-lab \
  -p 5431:5432 \
  --env REPMGR_PARTNER_NODES=pg-0,pg-1 \
  --env REPMGR_NODE_NAME=pg-0 \
  --env REPMGR_NODE_NETWORK_NAME=pg-0 \
  --env REPMGR_PRIMARY_HOST=pg-0 \
  --env REPMGR_PASSWORD=repmgrpass \
  --env POSTGRESQL_PASSWORD=secretpass \
  bitnami/postgresql-repmgr:latest
```


```bash
# PG-1 is standby
docker run --detach --name pg-1 \
  --network postgres-lab \
  -p 5432:5432 \
  --env REPMGR_PARTNER_NODES=pg-0,pg-1 \
  --env REPMGR_NODE_NAME=pg-1 \
  --env REPMGR_NODE_NETWORK_NAME=pg-1 \
  --env REPMGR_PRIMARY_HOST=pg-0 \
  --env REPMGR_PASSWORD=repmgrpass \
  --env POSTGRESQL_PASSWORD=secretpass \
  bitnami/postgresql-repmgr:latest
```


```bash
# Create role for Vault dynamic users

docker exec -i pg-0 \
    psql postgresql://postgres:secretpass@localhost:5432 -c "CREATE ROLE \"ro\" NOINHERIT;"

```


```bash
# Give permissions to ro role

docker exec -i pg-0 \
    psql postgresql://postgres:secretpass@localhost:5432 -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"ro\";"

```


```bash
# Create static role for Vault static user

docker exec -i pg-0 \
    psql postgresql://postgres:secretpass@localhost:5432 -c "CREATE ROLE staticuser WITH LOGIN PASSWORD 'staticuser' INHERIT;
GRANT ro TO staticuser;"

```


```bash
# Creation statement for dynamic users

tee readonly.sql <<EOF
CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}' INHERIT;
GRANT ro TO "{{name}}";
EOF

```


```bash
# Rotation statement for static users

tee rotation.sql <<EOF
ALTER USER "{{name}}" WITH PASSWORD '{{password}}';
EOF

```

## Vault Test


```bash
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=root
```


```bash
vault secrets enable database
```


```bash
# this database config will use the modified postgressql secret engine with pgx

vault write database/config/postgresql-new \
     plugin_name=postgres-new \
     connection_url="postgresql://{{username}}:{{password}}@localhost:5431,localhost:5432/postgres?target_session_attrs=read-write&sslmode=disable" \
     allowed_roles=readonly-new,staticrole-new \
     username="postgres" \
     password="secretpass"

```


```bash
# this database config will use the standard postgressql secret engine with pq

vault write database/config/postgresql \
     plugin_name=postgresql-database-plugin \
     connection_url="postgresql://{{username}}:{{password}}@localhost:5431/postgres?sslmode=disable" \
     allowed_roles=readonly,staticrole \
     username="postgres" \
     password="secretpass"


```


```bash
# this role for creation of dynamic users using the modified secret engine

vault write database/roles/readonly-new \
      db_name=postgresql-new \
      creation_statements=@readonly.sql \
      default_ttl=1h \
      max_ttl=24h
```


```bash
# this role for creation of static user using the modified secret engine

vault write database/static-roles/staticrole-new \
    db_name=postgresql-new \
    rotation_statements=@rotation.sql \
    username="staticuser" \
    rotation_period=86400
```


```bash
# this role for creation of dynamic users using the standard secret engine

vault write database/roles/readonly \
      db_name=postgresql \
      creation_statements=@readonly.sql \
      default_ttl=1h \
      max_ttl=24h
```


```bash
# this role for creation of static user using the standard secret engine

vault write database/static-roles/staticrole \
    db_name=postgresql \
    rotation_statements=@rotation.sql \
    username="staticuser" \
    rotation_period=86400
```

## Testing

### Both PG-0 (Active) and PG-1 (Standby) are up


```bash
# output the ip address of the active postgres instance (supports read-write)

psql -Atx postgresql://postgres:secretpass@localhost:5431/postgres?target_session_attrs=read-write -c 'select inet_server_addr();'
```


```bash
# output the ip address of the standby postgres instance (supports read-only)

psql -Atx postgresql://postgres:secretpass@localhost:5432/postgres?target_session_attrs=read-only -c 'select inet_server_addr();'
```


```bash
vault read database/creds/readonly-new
```


```bash
vault read database/creds/readonly
```


```bash
vault read database/static-creds/staticrole
```


```bash
vault write -f database/rotate-role/staticrole
```


```bash
vault read database/static-creds/staticrole-new
```


```bash
vault write -f database/rotate-role/staticrole-new
```

### PG-0 (Active) Down and PG-1 (Standby) is up


```bash
# output the ip address of the active postgres instance (supports read-write)

psql -Atx postgresql://postgres:secretpass@localhost:5431/postgres?target_session_attrs=read-write -c 'select inet_server_addr();'
```


```bash
# output the ip address of the standby postgres instance (supports read-only)

psql -Atx postgresql://postgres:secretpass@localhost:5432/postgres?target_session_attrs=read-only -c 'select inet_server_addr();'
```


```bash
docker stop pg-0
```


```bash
vault read database/creds/readonly-new
```


```bash
vault read database/creds/readonly
```


```bash
vault read database/static-creds/staticrole
```


```bash
vault write -f database/rotate-role/staticrole
```


```bash
vault read database/static-creds/staticrole-new
```


```bash
vault write -f database/rotate-role/staticrole-new
```

### PG-0 (Standby) Restored and PG-1 (Active) is up


```bash
docker start pg-0
```


```bash
# output the ip address of the active postgres instance (supports read-write)

psql -Atx postgresql://postgres:secretpass@localhost:5431/postgres?target_session_attrs=read-write -c 'select inet_server_addr();'
```


```bash
# output the ip address of the standby postgres instance (supports read-only)

psql -Atx postgresql://postgres:secretpass@localhost:5432/postgres?target_session_attrs=read-only -c 'select inet_server_addr();'
```


```bash
vault read database/creds/readonly-new
```


```bash
vault read database/creds/readonly
```


```bash
vault read database/static-creds/staticrole
```


```bash
vault write -f database/rotate-role/staticrole
```


```bash
vault read database/static-creds/staticrole-new
```


```bash
vault write -f database/rotate-role/staticrole-new
```


```bash

```

### PG-0 (Standby) promoted to Active and PG-1 (Active) is shutdown


```bash
docker stop pg-1
```


```bash
# output the ip address of the active postgres instance (supports read-write)

psql -Atx postgresql://postgres:secretpass@localhost:5431/postgres?target_session_attrs=read-write -c 'select inet_server_addr();'
```


```bash
# output the ip address of the standby postgres instance (supports read-only)

psql -Atx postgresql://postgres:secretpass@localhost:5432/postgres?target_session_attrs=read-only -c 'select inet_server_addr();'
```


```bash
vault read database/creds/readonly-new
```


```bash
vault read database/creds/readonly
```


```bash
vault read database/static-creds/staticrole
```


```bash
vault write -f database/rotate-role/staticrole
```


```bash
vault read database/static-creds/staticrole-new
```


```bash
vault write -f database/rotate-role/staticrole-new
```


```bash

```
