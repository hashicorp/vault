module github.com/hashicorp/vault/api

go 1.12

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/hashicorp/vault/sdk v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)
