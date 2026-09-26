module github.com/hashicorp/vault/sdk

go 1.26.0

require (
	cloud.google.com/go/cloudsqlconn v1.21.0
	github.com/armon/go-radix v1.0.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/containerd/errdefs v1.0.0
	github.com/evanphx/json-patch/v5 v5.9.11
	github.com/fatih/structs v1.1.0
	github.com/go-ldap/ldap/v3 v3.4.13
	github.com/go-test/deep v1.1.1
	github.com/golang/protobuf v1.5.4
	github.com/golang/snappy v1.0.0
	github.com/hashicorp/cap/ldap v0.0.0-20250911140431-44d01434c285
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v1.6.3
	github.com/hashicorp/go-immutable-radix v1.3.1
	github.com/hashicorp/go-kms-wrapping/entropy/v2 v2.0.1
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.26
	github.com/hashicorp/go-metrics v0.6.1
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.8.0
	github.com/hashicorp/go-retryablehttp v0.7.8
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2
	github.com/hashicorp/go-secure-stdlib/cryptoutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.3
	github.com/hashicorp/go-secure-stdlib/parseutil v0.2.0
	github.com/hashicorp/go-secure-stdlib/password v0.1.5
	github.com/hashicorp/go-secure-stdlib/permitpool v1.0.0
	github.com/hashicorp/go-secure-stdlib/plugincontainer v0.5.0
	github.com/hashicorp/go-secure-stdlib/regexp v1.0.0
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.3
	github.com/hashicorp/go-sockaddr v1.0.7
	github.com/hashicorp/go-uuid v1.0.4
	github.com/hashicorp/go-version v1.9.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/hashicorp/hcl v1.0.1-vault-7
	github.com/hashicorp/vault/api v1.16.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c
	github.com/moby/go-archive v0.3.0
	github.com/moby/moby/api v1.56.0
	github.com/moby/moby/client v0.6.0
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/robfig/cron/v3 v3.0.1
	github.com/ryanuber/go-glob v1.0.0
	github.com/stretchr/testify v1.12.1
	github.com/tink-crypto/tink-go/v2 v2.7.0
	go.uber.org/atomic v1.12.0
	golang.org/x/crypto v0.57.0
	golang.org/x/net v0.59.0
	golang.org/x/text v0.42.0
	google.golang.org/grpc v1.83.2
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/bufbuild/protocompile v0.14.2-0.20260716165721-bb5762d29672 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/docker/go-connections v0.7.0 // indirect
	github.com/go-sql-driver/mysql v1.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jimlambrt/gldap v0.1.14 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sync v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260819154853-08b0e4226688 // indirect
)

require (
	cloud.google.com/go/auth v0.23.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/Azure/go-ntlmssp v0.1.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.8 // indirect
	github.com/go-jose/go-jose/v4 v4.1.5 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.20 // indirect
	github.com/googleapis/gax-go/v2 v2.24.0 // indirect
	github.com/hashicorp/go-hmac-drbg v0.0.0-20210916214228-a6e5a68489f6 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/joshlf/go-acl v0.0.0-20200411065538-eae00ae38531 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/sequential v0.7.0 // indirect
	github.com/moby/sys/user v0.4.1 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oklog/run v1.2.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/petermattis/goid v0.0.0-20260716134002-a9b348f0a2b9 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	github.com/sasha-s/go-deadlock v0.3.5
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0 // indirect
	go.opentelemetry.io/otel v1.46.0 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/trace v1.46.0 // indirect
	golang.org/x/exp v0.0.0-20260709172345-9ea1abe57597 // indirect
	golang.org/x/oauth2 v0.37.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	google.golang.org/api v0.297.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260921155816-b14227669459 // indirect
)
