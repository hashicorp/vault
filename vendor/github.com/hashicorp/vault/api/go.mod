module github.com/hashicorp/vault/api

go 1.12

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-retryablehttp v0.5.3
	github.com/hashicorp/go-rootcerts v1.0.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/vault/sdk v0.1.8
	github.com/mitchellh/mapstructure v1.1.2
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	gopkg.in/square/go-jose.v2 v2.3.1
)
