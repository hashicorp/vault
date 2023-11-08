module github.com/hashicorp/vault/command

go 1.21.3

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230923063757-afb1ddc0824c
	github.com/fatih/color v1.16.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-test/deep v1.1.0
	github.com/google/go-cmp v0.6.0
	github.com/google/tink/go v1.7.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-secure-stdlib/gatedwriter v0.1.1
	github.com/hashicorp/go-secure-stdlib/kv-builder v0.1.2
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/hcl v1.0.1-vault-5
	github.com/hashicorp/vault v1.15.1
	github.com/hashicorp/vault-plugin-auth-alicloud v0.16.0
	github.com/hashicorp/vault-plugin-auth-centrify v0.15.1
	github.com/hashicorp/vault-plugin-auth-cf v0.15.1
	github.com/hashicorp/vault-plugin-auth-gcp v0.16.1
	github.com/hashicorp/vault-plugin-auth-jwt v0.17.2
	github.com/hashicorp/vault-plugin-auth-kerberos v0.10.1
	github.com/hashicorp/vault-plugin-auth-oci v0.14.2
	github.com/hashicorp/vault-plugin-secrets-kv v0.16.2
	github.com/hashicorp/vault/api v1.10.0
	github.com/hashicorp/vault/api/auth/approle v0.5.0
	github.com/hashicorp/vault/sdk v0.10.2
	github.com/kr/text v0.2.0
	github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty v0.0.20
	github.com/mholt/archiver/v3 v3.5.1
	github.com/mitchellh/cli v1.1.5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/natefinch/atomic v1.0.1
	github.com/oklog/run v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.14.0
	nhooyr.io/websocket v1.8.10
)

require (
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/cloudsqlconn v1.4.3 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/kms v1.15.1 // indirect
	cloud.google.com/go/monitoring v1.15.1 // indirect
	cloud.google.com/go/spanner v1.47.0 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.7.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.3.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 // indirect
	github.com/Azure/azure-storage-blob-go v0.15.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.29 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.12 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.1.1 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/DataDog/datadog-go v3.2.0+incompatible // indirect
	github.com/Jeffail/gabs v1.1.1 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/SAP/go-hdb v0.14.1 // indirect
	github.com/aerospike/aerospike-client-go/v5 v5.6.0 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.521 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/apache/arrow/go/v12 v12.0.1 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/apple/foundationdb/bindings/go v0.0.0-20190411004307-cd5c9d91fad2 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.44.331 // indirect
	github.com/aws/aws-sdk-go-v2 v1.17.7 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.59 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.31.0 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/axiomhq/hyperloglog v0.0.0-20220105174342-98591331716a // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/centrify/cloud-golang-sdk v0.0.0-20210923165758-a8c48d049166 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0 // indirect
	github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible // indirect
	github.com/circonus-labs/circonusllhist v0.1.3 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20220930021109-9c4e6c59ccf1 // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-oidc/v3 v3.5.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/couchbase/gocb/v2 v2.6.3 // indirect
	github.com/couchbase/gocbcore/v10 v10.2.3 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/denverdino/aliyungo v0.0.0-20190125010748-a747050bb1ba // indirect
	github.com/dgryski/go-metro v0.0.0-20180109044635-280f6062b5bc // indirect
	github.com/digitalocean/godo v1.7.5 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.5+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/duosecurity/duo_api_golang v0.0.0-20190308151101-6c680f768e74 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.10.1 // indirect
	github.com/envoyproxy/go-control-plane v0.11.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/gammazero/workerpool v1.1.3 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/go-errors/errors v1.5.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-ldap/ldap/v3 v3.4.4 // indirect
	github.com/go-ldap/ldif v0.0.0-20200320164324-fd88d9b715b3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/analysis v0.20.0 // indirect
	github.com/go-openapi/errors v0.20.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.20.2 // indirect
	github.com/go-openapi/runtime v0.19.24 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/strfmt v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.20.2 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gocql/gocql v1.0.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v23.1.21+incompatible // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-metrics-stackdriver v0.2.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.5 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.5 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/cap v0.3.4 // indirect
	github.com/hashicorp/consul/api v1.23.0 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/eventlogger v0.2.3 // indirect
	github.com/hashicorp/go-bexpr v0.1.12 // indirect
	github.com/hashicorp/go-discover v0.0.0-20210818145131-c573d69da192 // indirect
	github.com/hashicorp/go-gcp-common v0.8.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-kms-wrapping/entropy/v2 v2.0.0 // indirect
	github.com/hashicorp/go-kms-wrapping/v2 v2.0.13 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2 v2.0.7-1 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms/v2 v2.0.1 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2 v2.0.7 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2 v2.0.8 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2 v2.0.8 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2 v2.0.7 // indirect
	github.com/hashicorp/go-kms-wrapping/wrappers/transit/v2 v2.0.8 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/go-metrics v0.5.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.5.2 // indirect
	github.com/hashicorp/go-raftchunking v0.6.3-0.20191002164813-7e9e8525653a // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/awsutil v0.2.3 // indirect
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2 // indirect
	github.com/hashicorp/go-secure-stdlib/fileutil v0.1.0 // indirect
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.3 // indirect
	github.com/hashicorp/go-secure-stdlib/nonceutil v0.1.0 // indirect
	github.com/hashicorp/go-secure-stdlib/password v0.1.1 // indirect
	github.com/hashicorp/go-secure-stdlib/plugincontainer v0.2.2 // indirect
	github.com/hashicorp/go-secure-stdlib/reloadutil v0.1.1 // indirect
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.2 // indirect
	github.com/hashicorp/go-slug v0.12.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.5 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/go-tfe v1.33.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcp-sdk-go v0.23.0 // indirect
	github.com/hashicorp/jsonapi v0.0.0-20210826224640-ee7dae0fb22d // indirect
	github.com/hashicorp/mdns v1.0.4 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20230519153805-2275a83cbfdf // indirect
	github.com/hashicorp/raft v1.3.10 // indirect
	github.com/hashicorp/raft-autopilot v0.2.0 // indirect
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210421194847-a7e34179d62c // indirect
	github.com/hashicorp/raft-snapshot v1.0.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/vault-plugin-auth-azure v0.16.2 // indirect
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.17.1 // indirect
	github.com/hashicorp/vault-plugin-database-couchbase v0.9.4 // indirect
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.13.3 // indirect
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.10.1 // indirect
	github.com/hashicorp/vault-plugin-database-redis v0.2.2 // indirect
	github.com/hashicorp/vault-plugin-database-redis-elasticache v0.2.3 // indirect
	github.com/hashicorp/vault-plugin-database-snowflake v0.9.0 // indirect
	github.com/hashicorp/vault-plugin-secrets-ad v0.16.1 // indirect
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.15.1 // indirect
	github.com/hashicorp/vault-plugin-secrets-azure v0.16.3 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcp v0.17.0 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.15.1 // indirect
	github.com/hashicorp/vault-plugin-secrets-kubernetes v0.6.0 // indirect
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.10.1 // indirect
	github.com/hashicorp/vault-plugin-secrets-openldap v0.11.2 // indirect
	github.com/hashicorp/vault-plugin-secrets-terraform v0.7.3 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx v3.3.0+incompatible // indirect
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2 // indirect
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f // indirect
	github.com/jefferai/jsonx v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/joyent/triton-go v1.7.1-0.20200416154420-6801d15b779f // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/linode/linodego v0.7.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mediocregopher/radix/v4 v4.1.4 // indirect
	github.com/michaelklishin/rabbit-hole/v2 v2.12.0 // indirect
	github.com/microsoft/kiota-abstractions-go v1.2.0 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.0.0 // indirect
	github.com/microsoft/kiota-http-go v1.1.0 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.0.0 // indirect
	github.com/microsoft/kiota-serialization-json-go v1.0.4 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.0.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go v1.13.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v1.0.0 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/mikesmitty/edkey v0.0.0-20170222072505-3356ea4e686a // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/pointerstructure v1.2.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.0.5 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncw/swift v1.0.47 // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/okta/okta-sdk-golang/v2 v2.12.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2.0.20221005185240-3a7f492d3f1b // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible // indirect
	github.com/oracle/oci-go-sdk/v60 v60.0.0 // indirect
	github.com/packethost/packngo v0.1.1-0.20180711074735-b9cb5096f54c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pires/go-proxyproto v0.6.1 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/pquerna/otp v1.2.1-0.20191009055518-468c2dd2b58d // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rboyer/safeio v0.2.1 // indirect
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sasha-s/go-deadlock v0.2.0 // indirect
	github.com/sethvargo/go-limiter v0.7.1 // indirect
	github.com/shirou/gopsutil/v3 v3.22.6 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/snowflakedb/gosnowflake v1.6.24 // indirect
	github.com/softlayer/softlayer-go v0.0.0-20180806151055-260589d94c7d // indirect
	github.com/sony/gobreaker v0.4.2-0.20210216022020-dd874f9dd33b // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.162 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vmware/govmomi v0.18.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.7 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.7 // indirect
	go.etcd.io/etcd/client/v2 v2.305.5 // indirect
	go.etcd.io/etcd/client/v3 v3.5.7 // indirect
	go.mongodb.org/atlas v0.33.0 // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.138.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230807174057-1744710a1577 // indirect
	google.golang.org/grpc v1.58.3 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.28.1 // indirect
	k8s.io/apimachinery v0.28.1 // indirect
	k8s.io/client-go v0.28.1 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2 // indirect
	layeh.com/radius v0.0.0-20190322222518-890bc1058917 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
