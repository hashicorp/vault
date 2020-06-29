module github.com/hashicorp/vault/api

go 1.13

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/go-test/deep v1.0.2
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-retryablehttp v0.6.2
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/vault/sdk v0.1.14-0.20200519221838-e0cfd64bc267
	github.com/mitchellh/mapstructure v1.2.2
	golang.org/x/net v0.0.0-20200519113804-d87ec0cfa476
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gopkg.in/square/go-jose.v2 v2.4.1
)
