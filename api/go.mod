module github.com/hashicorp/vault/api

go 1.12

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/go-test/deep v1.0.2-0.20181118220953-042da051cf31
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-retryablehttp v0.6.2
	github.com/hashicorp/go-rootcerts v1.0.1
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/vault/sdk v0.1.14-0.20191108161836-82f2b5571044
	github.com/mitchellh/mapstructure v1.1.2
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	gopkg.in/square/go-jose.v2 v2.3.1
)
