module github.com/hashicorp/vault/api/auth/azure

go 1.16

replace github.com/hashicorp/vault/api => ../../../api

require (
	github.com/Azure/go-autorest/autorest v0.11.21 // indirect
	github.com/hashicorp/vault/api v1.2.0
)
