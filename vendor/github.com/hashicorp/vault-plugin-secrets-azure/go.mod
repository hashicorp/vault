module github.com/hashicorp/vault-plugin-secrets-azure

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/Azure/azure-sdk-for-go v29.0.0+incompatible
	github.com/Azure/go-autorest v11.7.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/go-test/deep v1.0.2-0.20181118220953-042da051cf31
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
)
