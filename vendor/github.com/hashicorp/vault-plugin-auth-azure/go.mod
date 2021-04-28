module github.com/hashicorp/vault-plugin-auth-azure

go 1.15

require (
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v0.15.0
	github.com/hashicorp/vault/api v1.0.5-0.20200215224050-f6547fa8e820
	github.com/hashicorp/vault/sdk v0.1.14-0.20200215224050-f6547fa8e820
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)
