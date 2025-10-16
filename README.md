# Vault (ML-KEM fork) — Dev quickstart

## Prerequisites
- Go ≥ 1.25.2 on PATH
- `make`


## One-time bootstrap
```bash
make bootstrap
```

## Build
```bash
go generate ./...     # safe to run always
go mod vendor         # updates vendor from local replaces
make dev              # builds ./bin/vault
```

## Run dev server
```bash
VAULT_DEV_ROOT_TOKEN_ID=root ./bin/vault server -dev -dev-listen-address=127.0.0.1:8200
```

## Enable transit
```bash
./bin/vault secrets enable transit
```

## HTTP requests
Use these with VS Code REST Client or JetBrains HTTP Client.

- Env file: [`dev/http/http-client.env.json`](dev/http/http-client.env.json)
- Requests: [`dev/http/vault_transit.http`](dev/http/vault_transit.http)
