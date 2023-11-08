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
go 1.21.3

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/api/auth/approle => ./api/auth/approle

replace github.com/hashicorp/vault/api/auth/kubernetes => ./api/auth/kubernetes

replace github.com/hashicorp/vault/api/auth/userpass => ./api/auth/userpass

replace github.com/hashicorp/vault/sdk => ./sdk

require (
	cloud.google.com/go/cloudsqlconn v1.5.0
	cloud.google.com/go/monitoring v1.16.3
	cloud.google.com/go/spanner v1.51.0
	cloud.google.com/go/storage v1.34.1
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0
	github.com/Azure/azure-storage-blob-go v0.15.0
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/Azure/go-autorest/autorest/adal v0.9.23
	github.com/ProtonMail/go-crypto v0.0.0-20230923063757-afb1ddc0824c
	github.com/SAP/go-hdb v1.5.9
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/aerospike/aerospike-client-go/v5 v5.11.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.605
	github.com/aliyun/aliyun-oss-go-sdk v3.0.1+incompatible
	github.com/apple/foundationdb/bindings/go v0.0.0-20231107221918-967a546e15e8
	github.com/armon/go-metrics v0.4.1
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/aws/aws-sdk-go v1.47.7
	github.com/aws/aws-sdk-go-v2/config v1.22.3
	github.com/axiomhq/hyperloglog v0.0.0-20230201085229-3ddf4bad03dc
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/client9/misspell v0.3.4
	github.com/cockroachdb/cockroach-go v2.0.1+incompatible
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/duosecurity/duo_api_golang v0.0.0-20230418202038-096d3306c029
	github.com/dustin/go-humanize v1.0.1
	github.com/fatih/structs v1.1.0
	github.com/favadi/protoc-go-inject-tag v1.4.0
	github.com/gammazero/workerpool v1.1.3
	github.com/go-errors/errors v1.5.1
	github.com/go-git/go-git/v5 v5.7.0
	github.com/go-jose/go-jose/v3 v3.0.1
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/go-sql-driver/mysql v1.7.1
	github.com/go-test/deep v1.1.0
	github.com/go-zookeeper/zk v1.0.3
	github.com/gocql/gocql v1.6.0
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/protobuf v1.5.3
	github.com/golangci/revgrep v0.0.0-20220804021717-745bb2f7c2e6
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.6.0
	github.com/google/tink/go v1.7.0
	github.com/hashicorp/cap v0.4.0
	github.com/hashicorp/cap/ldap v0.0.0-20231012003312-273118a6e3b8
	github.com/hashicorp/consul-template v0.33.0
	github.com/hashicorp/consul/api v1.26.1
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/eventlogger v0.2.6
	github.com/hashicorp/go-bexpr v0.1.13
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20230724184603-e89ebd1b2f65
	github.com/hashicorp/go-gcp-common v0.8.0
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-kms-wrapping/entropy/v2 v2.0.0
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.14
	github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2 v2.0.8
	github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms/v2 v2.0.2
	github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2 v2.0.8
	github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2 v2.0.9
	github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2 v2.0.9
	github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2 v2.0.8
	github.com/hashicorp/go-kms-wrapping/wrappers/transit/v2 v2.0.9
	github.com/hashicorp/go-memdb v1.3.4
	github.com/hashicorp/go-msgpack v1.1.5
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-plugin v1.5.2
	github.com/hashicorp/go-raftchunking v0.7.0
	github.com/hashicorp/go-retryablehttp v0.7.5
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-secure-stdlib/awsutil v0.2.3
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2
	github.com/hashicorp/go-secure-stdlib/gatedwriter v0.1.1
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.3
	github.com/hashicorp/go-secure-stdlib/nonceutil v0.1.0
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8
	github.com/hashicorp/go-secure-stdlib/password v0.1.3
	github.com/hashicorp/go-secure-stdlib/reloadutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.3
	github.com/hashicorp/go-sockaddr v1.0.5
	github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/hashicorp/hcl v1.0.1-vault-5
	github.com/hashicorp/hcl/v2 v2.16.2
	github.com/hashicorp/hcp-link v0.1.0
	github.com/hashicorp/hcp-scada-provider v0.2.1
	github.com/hashicorp/hcp-sdk-go v0.72.0
	github.com/hashicorp/nomad/api v0.0.0-20231109160603-b61a31c38f1a
	github.com/hashicorp/raft v1.5.0
	github.com/hashicorp/raft-autopilot v0.2.0
	github.com/hashicorp/raft-boltdb/v2 v2.2.2
	github.com/hashicorp/raft-snapshot v1.0.4
	github.com/hashicorp/vault-plugin-auth-alicloud v0.16.0
	github.com/hashicorp/vault-plugin-auth-azure v0.16.2
	github.com/hashicorp/vault-plugin-auth-centrify v0.15.1
	github.com/hashicorp/vault-plugin-auth-cf v0.15.1
	github.com/hashicorp/vault-plugin-auth-gcp v0.16.1
	github.com/hashicorp/vault-plugin-auth-jwt v0.17.2
	github.com/hashicorp/vault-plugin-auth-kerberos v0.10.1
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.17.1
	github.com/hashicorp/vault-plugin-auth-oci v0.14.2
	github.com/hashicorp/vault-plugin-database-couchbase v0.9.4
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.13.3
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.10.1
	github.com/hashicorp/vault-plugin-database-redis v0.2.2
	github.com/hashicorp/vault-plugin-database-redis-elasticache v0.2.3
	github.com/hashicorp/vault-plugin-database-snowflake v0.9.0
	github.com/hashicorp/vault-plugin-mock v0.16.1
	github.com/hashicorp/vault-plugin-secrets-ad v0.16.1
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.15.1
	github.com/hashicorp/vault-plugin-secrets-azure v0.16.3
	github.com/hashicorp/vault-plugin-secrets-gcp v0.17.0
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.15.2
	github.com/hashicorp/vault-plugin-secrets-kubernetes v0.6.0
	github.com/hashicorp/vault-plugin-secrets-kv v0.16.2
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.10.2
	github.com/hashicorp/vault-plugin-secrets-openldap v0.11.2
	github.com/hashicorp/vault-plugin-secrets-terraform v0.7.3
	github.com/hashicorp/vault-testing-stepwise v0.1.4
	github.com/hashicorp/vault/api v1.10.0
	github.com/hashicorp/vault/api/auth/approle v0.5.0
	github.com/hashicorp/vault/api/auth/userpass v0.1.0
	github.com/hashicorp/vault/command v0.0.0-20231108163531-6c44b0d888db
	github.com/hashicorp/vault/sdk v0.10.2
	github.com/hashicorp/vault/vault/hcp_link/proto v0.0.0-20230201201504-b741fa893d77
	github.com/influxdata/influxdb1-client v0.0.0-20220302092344-a9ab5670611c
	github.com/jackc/pgx/v4 v4.18.1
	github.com/jcmturner/gokrb5/v8 v8.4.4
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.1
	github.com/joyent/triton-go v1.8.5
	github.com/klauspost/compress v1.17.2
	github.com/kr/pretty v0.3.1
	github.com/michaelklishin/rabbit-hole/v2 v2.15.0
	github.com/mikesmitty/edkey v0.0.0-20170222072505-3356ea4e686a
	github.com/mitchellh/cli v1.1.5
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/ncw/swift v1.0.53
	github.com/oklog/run v1.1.0
	github.com/okta/okta-sdk-golang/v2 v2.20.0
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.10.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pires/go-proxyproto v0.7.0
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/common v0.45.0
	github.com/rboyer/safeio v0.2.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/sasha-s/go-deadlock v0.3.1
	github.com/sethvargo/go-limiter v0.7.2
	github.com/shirou/gopsutil/v3 v3.23.10
	github.com/stretchr/testify v1.8.4
	go.etcd.io/bbolt v1.3.8
	go.etcd.io/etcd/client/pkg/v3 v3.5.10
	go.etcd.io/etcd/client/v2 v2.305.10
	go.etcd.io/etcd/client/v3 v3.5.10
	go.mongodb.org/atlas v0.35.0
	go.mongodb.org/mongo-driver v1.13.0
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/sdk v1.19.0
	go.opentelemetry.io/otel/trace v1.19.0
	go.uber.org/atomic v1.11.0
	go.uber.org/goleak v1.2.1
	golang.org/x/crypto v0.15.0
	golang.org/x/exp v0.0.0-20231108232855-2478ac86f678
	golang.org/x/net v0.18.0
	golang.org/x/oauth2 v0.14.0
	golang.org/x/sync v0.5.0
	golang.org/x/sys v0.14.0
	golang.org/x/term v0.14.0
	golang.org/x/text v0.14.0
	golang.org/x/tools v0.15.0
	google.golang.org/api v0.150.0
	google.golang.org/grpc v1.59.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gotest.tools/gotestsum v1.10.0
	honnef.co/go/tools v0.4.3
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b
	layeh.com/radius v0.0.0-20230922032716-6579be8edf5d
	mvdan.cc/gofumpt v0.3.1
	nhooyr.io/websocket v1.8.10
)

require (
	cloud.google.com/go v0.110.10 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/kms v1.15.5 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys v0.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.7.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.0 // indirect
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
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/Jeffail/gabs/v2 v2.7.0 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/andybalholm/brotli v1.0.6 // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/apache/thrift v0.19.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.22.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.42.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.17.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.1 // indirect
	github.com/aws/smithy-go v1.16.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/centrify/cloud-golang-sdk v0.0.0-20220926200933-ed5f25b01f45 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible // indirect
	github.com/circonus-labs/circonusllhist v0.1.5 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/cloudflare/circl v1.3.6 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220930021109-9c4e6c59ccf1 // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20231109132714-523115ebc101 // indirect
	github.com/containerd/containerd v1.7.0 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-oidc/v3 v3.7.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/couchbase/gocb/v2 v2.6.5 // indirect
	github.com/couchbase/gocbcore/v10 v10.2.9 // indirect
	github.com/danieljoos/wincred v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denverdino/aliyungo v0.0.0-20230411124812-ab98a9173ace // indirect
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/digitalocean/godo v1.105.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/cli v20.10.20+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v24.0.7+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/envoyproxy/go-control-plane v0.11.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/evanphx/json-patch/v5 v5.7.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.4.1 // indirect
	github.com/go-ldap/ldif v0.0.0-20200320164324-fd88d9b715b3 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.4 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.26.0 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/strfmt v0.21.7 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/go-resty/resty/v2 v2.10.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.1.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v1.7.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.1 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.0.0 // indirect
	github.com/hashicorp/go-secure-stdlib/fileutil v0.1.0 // indirect
	github.com/hashicorp/go-secure-stdlib/kv-builder v0.1.2 // indirect
	github.com/hashicorp/go-secure-stdlib/plugincontainer v0.2.2 // indirect
	github.com/hashicorp/go-slug v0.13.1 // indirect
	github.com/hashicorp/go-tfe v1.39.1 // indirect
	github.com/hashicorp/jsonapi v0.0.0-20231023233540-b6a3d216e521 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/mdns v1.0.5 // indirect
	github.com/hashicorp/net-rpc-msgpackrpc/v2 v2.0.0 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/vault/api/auth/kubernetes v0.4.1 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linode/linodego v1.25.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20231016141302-07b5767bb0ed // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-ieproxy v0.0.11 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/mediocregopher/radix/v4 v4.1.4 // indirect
	github.com/mholt/archiver/v3 v3.5.1 // indirect
	github.com/microsoft/kiota-abstractions-go v1.5.0 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.0.1 // indirect
	github.com/microsoft/kiota-http-go v1.1.0 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-json-go v1.0.4 // indirect
	github.com/microsoft/kiota-serialization-multipart-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go v1.24.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v1.0.0 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/pointerstructure v1.2.1 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.5 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc5 // indirect
	github.com/opencontainers/runc v1.1.6 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/oracle/oci-go-sdk/v60 v60.0.0 // indirect
	github.com/packethost/packngo v0.30.0 // indirect
	github.com/petermattis/goid v0.0.0-20230904192822-1876fd5063bc // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.1.1 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 // indirect
	github.com/snowflakedb/gosnowflake v1.6.25 // indirect
	github.com/softlayer/softlayer-go v1.1.2 // indirect
	github.com/softlayer/xmlrpc v0.0.0-20200409220501-5f089df7cb7e // indirect
	github.com/sony/gobreaker v0.5.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/std-uritemplate/std-uritemplate/go v0.0.46 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.233+incompatible // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vmware/govmomi v0.33.1 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	github.com/yuin/gopher-lua v1.1.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.10 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20221208152030-732eee02a75a // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.3 // indirect
	k8s.io/apimachinery v0.28.3 // indirect
	k8s.io/client-go v0.28.3 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/kube-openapi v0.0.0-20231010175941-2dd684a91f00 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/ma314smith/signedxml v1.1.1 => github.com/moov-io/signedxml v1.1.1
