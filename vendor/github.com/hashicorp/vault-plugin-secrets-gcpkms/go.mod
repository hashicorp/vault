module github.com/hashicorp/vault-plugin-secrets-gcpkms

go 1.12

require (
	cloud.google.com/go v0.37.4
	github.com/gammazero/deque v0.0.0-20190130191400-2afb3858e9c7 // indirect
	github.com/gammazero/workerpool v0.0.0-20190406235159-88d534f22b56
	github.com/golang/protobuf v1.3.1
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/vault/api v1.0.0
	github.com/hashicorp/vault/sdk v0.1.5
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2
	github.com/satori/go.uuid v1.2.0
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	google.golang.org/api v0.3.2
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.20.0
)
