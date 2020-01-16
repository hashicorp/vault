module github.com/hashicorp/vault-plugin-auth-azure

go 1.12

require (
	github.com/Azure/azure-sdk-for-go v36.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/vault/api v1.0.5-0.20191119041037-cccda49b3962
	github.com/hashicorp/vault/sdk v0.1.14-0.20191108161836-82f2b5571044
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
)
