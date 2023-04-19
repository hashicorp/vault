module github.com/hashicorp/vault/sdk

go 1.19

// TODO revert this, convenient though it is
replace github.com/hashicorp/vault/api => ../api

require (
	github.com/armon/go-metrics v0.3.9
	github.com/armon/go-radix v1.0.0
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/docker/docker v23.0.1+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/evanphx/json-patch/v5 v5.5.0
	github.com/fatih/structs v1.1.0
	github.com/go-ldap/ldap/v3 v3.1.10
	github.com/go-test/deep v1.0.2
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v0.16.2
	github.com/hashicorp/go-immutable-radix v1.3.1
	github.com/hashicorp/go-kms-wrapping/entropy/v2 v2.0.0
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.8
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.4.5
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.1
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.1
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6
	github.com/hashicorp/go-secure-stdlib/password v0.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.2
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/vault/api v1.9.0
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-testing-interface v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pierrec/lz4 v2.5.2+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/stretchr/testify v1.8.2
	go.uber.org/atomic v1.9.0
	golang.org/x/crypto v0.6.0
	golang.org/x/net v0.8.0
	golang.org/x/text v0.8.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/containerd/containerd v1.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/go-asn1-ber/asn1-ber v1.3.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2.0.20221005185240-3a7f492d3f1b // indirect
	github.com/opencontainers/runc v1.1.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.4.0 // indirect
)
