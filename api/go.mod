module github.com/hashicorp/vault/api

go 1.13

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/go-test/deep v1.0.2
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/vault/sdk v0.1.14-0.20200519221838-e0cfd64bc267
	github.com/mitchellh/mapstructure v1.3.2
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1
	gopkg.in/square/go-jose.v2 v2.5.1
)
