module github.com/hashicorp/vault-plugin-auth-jwt

go 1.12

require (
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/go-test/deep v1.0.2-0.20181118220953-042da051cf31
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-uuid v1.0.2-0.20191001231223-f32f5fe8d6a8
	github.com/hashicorp/vault/api v1.0.5-0.20191216174727-9d51b36f3ae4
	github.com/hashicorp/vault/sdk v0.1.14-0.20191216174727-9d51b36f3ae4
	github.com/mitchellh/pointerstructure v0.0.0-20190430161007-f252a8fd71c8
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	gopkg.in/square/go-jose.v2 v2.3.1
)
