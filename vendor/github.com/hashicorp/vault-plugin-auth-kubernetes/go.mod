module github.com/hashicorp/vault-plugin-auth-kubernetes

go 1.12

require (
	github.com/briankassouf/jose v0.9.2-0.20180619214549-d2569464773f
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/vault/api v1.0.1
	github.com/hashicorp/vault/sdk v0.1.12-0.20190619234858-76b551f81856
	github.com/mitchellh/mapstructure v1.1.2
	k8s.io/api v0.0.0-20190409092523-d687e77c8ae9
	k8s.io/apimachinery v0.0.0-20190409092423-760d1845f48b
)
