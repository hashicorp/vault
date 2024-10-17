module github.com/hashicorp/vault

// The go version directive value isn't consulted when building our production binaries,
// and the vault module isn't intended to be imported into other projects.  As such the
// impact of this setting is usually rather limited.  Note however that in some cases the
// Go project introduces new semantics for handling of go.mod depending on the value.
//
// The general policy for updating it is: when the Go major version used on the branch is
// updated. If we choose not to do so at some point (e.g. because we don't want some new
// semantic related to Go module handling), this comment should be updated to explain that.
//
// Whenever this value gets updated, sdk/go.mod should be updated to the same value.
go 1.22.5

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/api/auth/approle => ./api/auth/approle

replace github.com/hashicorp/vault/api/auth/kubernetes => ./api/auth/kubernetes

replace github.com/hashicorp/vault/api/auth/userpass => ./api/auth/userpass

replace github.com/hashicorp/vault/sdk => ./sdk

require (
	cloud.google.com/go/cloudsqlconn v1.4.3
	cloud.google.com/go/monitoring v1.21.0
	cloud.google.com/go/spanner v1.67.0
	cloud.google.com/go/storage v1.43.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.14.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0
	github.com/Azure/azure-storage-blob-go v0.15.0
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/Azure/go-autorest/autorest/adal v0.9.23
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371
	github.com/SAP/go-hdb v1.10.1
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/aerospike/aerospike-client-go/v5 v5.6.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.12
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/apple/foundationdb/bindings/go v0.0.0-20190411004307-cd5c9d91fad2
	github.com/armon/go-metrics v0.4.1
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2/config v1.27.11
	github.com/axiomhq/hyperloglog v0.0.0-20220105174342-98591331716a
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/cockroachdb/cockroach-go/v2 v2.3.8
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/docker/docker v26.1.5+incompatible
	github.com/duosecurity/duo_api_golang v0.0.0-20190308151101-6c680f768e74
	github.com/dustin/go-humanize v1.0.1
	github.com/fatih/color v1.17.0
	github.com/fatih/structs v1.1.0
	github.com/gammazero/workerpool v1.1.3
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.5.1
	github.com/go-git/go-git/v5 v5.11.0
	github.com/go-jose/go-jose/v3 v3.0.3
	github.com/go-ldap/ldap/v3 v3.4.8
	github.com/go-sql-driver/mysql v1.7.1
	github.com/go-test/deep v1.1.1
	github.com/go-zookeeper/zk v1.0.3
	github.com/gocql/gocql v1.0.0
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.2.0
	github.com/hashicorp-forge/bbolt v1.3.8-hc3
	github.com/hashicorp/cap v0.7.0
	github.com/hashicorp/cap/ldap v0.0.0-20240403125925-c0418810d10e
	github.com/hashicorp/cli v1.1.6
	github.com/hashicorp/consul-template v0.39.1
	github.com/hashicorp/consul/api v1.29.1
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/eventlogger v0.2.10
	github.com/hashicorp/go-bexpr v0.1.12
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20210818145131-c573d69da192
	github.com/hashicorp/go-gcp-common v0.9.0
	github.com/hashicorp/go-hclog v1.6.3
	github.com/hashicorp/go-kms-wrapping/entropy/v2 v2.0.1
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.16
	github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2 v2.0.9
	github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms/v2 v2.0.3
	github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2 v2.0.9
	github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2 v2.0.11
	github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2 v2.0.12
	github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2 v2.0.7
	github.com/hashicorp/go-kms-wrapping/wrappers/transit/v2 v2.0.12
	github.com/hashicorp/go-memdb v1.3.4
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.6.1
	github.com/hashicorp/go-raftchunking v0.6.3-0.20191002164813-7e9e8525653a
	github.com/hashicorp/go-retryablehttp v0.7.7
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-secure-stdlib/awsutil v0.3.0
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2
	github.com/hashicorp/go-secure-stdlib/gatedwriter v0.1.1
	github.com/hashicorp/go-secure-stdlib/kv-builder v0.1.2
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.3
	github.com/hashicorp/go-secure-stdlib/nonceutil v0.1.0
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8
	github.com/hashicorp/go-secure-stdlib/password v0.1.1
	github.com/hashicorp/go-secure-stdlib/reloadutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.3
	github.com/hashicorp/go-sockaddr v1.0.6
	github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version v1.7.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/hashicorp/hcl v1.0.1-vault-5
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/hashicorp/hcp-link v0.2.1
	github.com/hashicorp/hcp-scada-provider v0.2.2
	github.com/hashicorp/hcp-sdk-go v0.101.0
	github.com/hashicorp/nomad/api v0.0.0-20240213164230-c364cb57298d
	github.com/hashicorp/raft v1.7.1
	github.com/hashicorp/raft-autopilot v0.2.0
	github.com/hashicorp/raft-boltdb/v2 v2.3.0
	github.com/hashicorp/raft-snapshot v1.0.4
	github.com/hashicorp/raft-wal v0.4.0
	github.com/hashicorp/vault-hcp-lib v0.0.0-20240704151836-a5c058ac604c
	github.com/hashicorp/vault-plugin-auth-alicloud v0.19.0
	github.com/hashicorp/vault-plugin-auth-azure v0.19.1
	github.com/hashicorp/vault-plugin-auth-cf v0.19.0
	github.com/hashicorp/vault-plugin-auth-gcp v0.19.0
	github.com/hashicorp/vault-plugin-auth-jwt v0.22.0
	github.com/hashicorp/vault-plugin-auth-kerberos v0.13.0
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.20.0
	github.com/hashicorp/vault-plugin-auth-oci v0.17.0
	github.com/hashicorp/vault-plugin-database-couchbase v0.12.0
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.16.0
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.13.0
	github.com/hashicorp/vault-plugin-database-redis v0.4.0
	github.com/hashicorp/vault-plugin-database-redis-elasticache v0.5.0
	github.com/hashicorp/vault-plugin-database-snowflake v0.12.0
	github.com/hashicorp/vault-plugin-mock v0.16.1
	github.com/hashicorp/vault-plugin-secrets-ad v0.19.0
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.18.0
	github.com/hashicorp/vault-plugin-secrets-azure v0.20.1
	github.com/hashicorp/vault-plugin-secrets-gcp v0.20.0
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.19.0
	github.com/hashicorp/vault-plugin-secrets-kubernetes v0.9.0
	github.com/hashicorp/vault-plugin-secrets-kv v0.20.0
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.13.0
	github.com/hashicorp/vault-plugin-secrets-openldap v0.14.2
	github.com/hashicorp/vault-plugin-secrets-terraform v0.10.0
	github.com/hashicorp/vault-testing-stepwise v0.3.1
	github.com/hashicorp/vault/api v1.15.0
	github.com/hashicorp/vault/api/auth/approle v0.1.0
	github.com/hashicorp/vault/api/auth/userpass v0.1.0
	github.com/hashicorp/vault/sdk v0.14.0
	github.com/hashicorp/vault/vault/hcp_link/proto v0.0.0-20230201201504-b741fa893d77
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jcmturner/gokrb5/v8 v8.4.4
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.1
	github.com/joyent/triton-go v1.7.1-0.20200416154420-6801d15b779f
	github.com/klauspost/compress v1.17.8
	github.com/kr/pretty v0.3.1
	github.com/kr/text v0.2.0
	github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty v0.0.20
	github.com/michaelklishin/rabbit-hole/v2 v2.12.0
	github.com/mikesmitty/edkey v0.0.0-20170222072505-3356ea4e686a
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/ncw/swift v1.0.47
	github.com/oklog/run v1.1.0
	github.com/okta/okta-sdk-golang/v4 v4.1.2
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.10.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pires/go-proxyproto v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/pquerna/otp v1.2.1-0.20191009055518-468c2dd2b58d
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/common v0.55.0
	github.com/rboyer/safeio v0.2.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/sasha-s/go-deadlock v0.2.0
	github.com/sethvargo/go-limiter v0.7.1
	github.com/shirou/gopsutil/v3 v3.22.6
	github.com/stretchr/testify v1.9.0
	github.com/tink-crypto/tink-go/v2 v2.2.0
	go.etcd.io/bbolt v1.3.10
	go.etcd.io/etcd/client/pkg/v3 v3.5.13
	go.etcd.io/etcd/client/v2 v2.305.5
	go.etcd.io/etcd/client/v3 v3.5.13
	go.mongodb.org/atlas v0.37.0
	go.mongodb.org/mongo-driver v1.16.1
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
	go.opentelemetry.io/otel/trace v1.30.0
	go.uber.org/atomic v1.11.0
	go.uber.org/goleak v1.3.0
	golang.org/x/crypto v0.27.0
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
	golang.org/x/net v0.29.0
	golang.org/x/oauth2 v0.23.0
	golang.org/x/sync v0.8.0
	golang.org/x/sys v0.25.0
	golang.org/x/term v0.24.0
	golang.org/x/text v0.18.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
	google.golang.org/api v0.197.0
	google.golang.org/grpc v1.66.1
	google.golang.org/protobuf v1.34.2
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	k8s.io/apimachinery v0.31.0
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8
	layeh.com/radius v0.0.0-20231213012653-1006025d24f8
	nhooyr.io/websocket v1.8.11
)

require (
	cel.dev/expr v0.15.0 // indirect
	cloud.google.com/go/longrunning v0.6.0 // indirect
	github.com/containerd/containerd v1.7.20 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/hashicorp/go-secure-stdlib/httputil v0.1.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.29 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/xerrors v0.0.0-20240716161551-93cc26a95ae9 // indirect
)

require (
	cloud.google.com/go v0.115.1 // indirect
	cloud.google.com/go/auth v0.9.3 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.4 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	cloud.google.com/go/iam v1.2.0 // indirect
	cloud.google.com/go/kms v1.19.0 // indirect; indirect\
	dario.cat/mergo v1.0.1 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys v0.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.7.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2 v2.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.12 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/DataDog/datadog-go v3.2.0+incompatible // indirect
	github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.0 // indirect
	github.com/Jeffail/gabs/v2 v2.1.0 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apache/arrow/go/v15 v15.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.11 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.53.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.6 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/benbjohnson/immutable v0.4.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible // indirect
	github.com/circonus-labs/circonusllhist v0.1.3 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220930021109-9c4e6c59ccf1 // indirect
	github.com/cncf/xds/go v0.0.0-20240423153145-555b57ec207b // indirect
	github.com/containerd/continuity v0.4.3 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-oidc/v3 v3.11.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/coreos/pkg v0.0.0-20220810130054-c7d1c02cb6cf // indirect
	github.com/couchbase/gocb/v2 v2.9.1 // indirect
	github.com/couchbase/gocbcore/v10 v10.5.1 // indirect
	github.com/couchbase/gocbcoreps v0.1.3 // indirect
	github.com/couchbase/goprotostellar v1.0.2 // indirect
	github.com/couchbaselabs/gocbconnstr/v2 v2.0.0-20240607131231-fb385523de28 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denverdino/aliyungo v0.0.0-20190125010748-a747050bb1ba // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/digitalocean/godo v1.7.5 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v26.1.5+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/envoyproxy/go-control-plane v0.12.1-0.20240621013728-1eb8caab5155 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-ldap/ldif v0.0.0-20200320164324-fd88d9b715b3 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.2 // indirect
	github.com/hashicorp/go-secure-stdlib/fileutil v0.1.0 // indirect
	github.com/hashicorp/go-secure-stdlib/plugincontainer v0.4.0 // indirect
	github.com/hashicorp/go-slug v0.15.2 // indirect
	github.com/hashicorp/go-tfe v1.64.2 // indirect
	github.com/hashicorp/jsonapi v1.3.1 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/mdns v1.0.4 // indirect
	github.com/hashicorp/net-rpc-msgpackrpc/v2 v2.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/vault/api/auth/kubernetes v0.6.0 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.3 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/joshlf/go-acl v0.0.0-20200411065538-eae00ae38531 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linode/linodego v0.7.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/mediocregopher/radix/v4 v4.1.4 // indirect
	github.com/microsoft/kiota-abstractions-go v1.6.1 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.1.0 // indirect
	github.com/microsoft/kiota-http-go v1.4.4 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-json-go v1.0.8 // indirect
	github.com/microsoft/kiota-serialization-multipart-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go v1.47.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v1.2.1 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/pointerstructure v1.2.1 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/sys/user v0.2.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.1.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runc v1.2.0-rc.1 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/oracle/oci-go-sdk/v59 v59.0.0 // indirect
	github.com/oracle/oci-go-sdk/v60 v60.0.0 // indirect
	github.com/packethost/packngo v0.1.1-0.20180711074735-b9cb5096f54c // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.2.1 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 // indirect
	github.com/snowflakedb/gosnowflake v1.11.0 // indirect
	github.com/softlayer/softlayer-go v0.0.0-20180806151055-260589d94c7d // indirect
	github.com/sony/gobreaker v0.5.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/std-uritemplate/std-uritemplate/go v0.0.57 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.162 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c // indirect
	github.com/vmware/govmomi v0.18.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.13 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect; indirect\
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.31.0 // indirect
	k8s.io/client-go v0.31.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/ma314smith/signedxml v1.1.1 => github.com/moov-io/signedxml v1.1.1

// Support using the forked repository until https://github.com/pires/go-proxyproto/pull/110 merges
// and is released.
replace github.com/pires/go-proxyproto v1.0.0 => github.com/peteski22/go-proxyproto v1.0.0
