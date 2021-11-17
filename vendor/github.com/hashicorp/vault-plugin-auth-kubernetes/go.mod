module github.com/hashicorp/vault-plugin-auth-kubernetes

go 1.12

require (
	github.com/briankassouf/jose v0.9.2-0.20180619214549-d2569464773f
	github.com/go-test/deep v1.0.8
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.1
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/vault/api v1.2.0
	github.com/hashicorp/vault/sdk v0.2.1
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/mitchellh/mapstructure v1.4.2
	k8s.io/api v0.0.0-20190409092523-d687e77c8ae9
	k8s.io/apimachinery v0.22.2
)
