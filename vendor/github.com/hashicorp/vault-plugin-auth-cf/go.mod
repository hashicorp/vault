module github.com/hashicorp/vault-plugin-auth-cf

go 1.12

require (
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190201205600-f136f9222381
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/vault/api v1.0.5-0.20190814205728-e9c5cd8aca98
	github.com/hashicorp/vault/sdk v0.1.14-0.20190814205504-1cad00d1133b
	github.com/pkg/errors v0.8.1
)
