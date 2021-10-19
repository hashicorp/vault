module github.com/hashicorp/vault

go 1.17

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/sdk => ./sdk

replace go.etcd.io/etcd/client/pkg/v3 v3.5.0 => go.etcd.io/etcd/client/pkg/v3 v3.0.0-20210928084031-3df272774672

require (
	cloud.google.com/go v0.56.0
	cloud.google.com/go/spanner v1.5.1
	cloud.google.com/go/storage v1.6.0
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/adal v0.9.11
	github.com/NYTimes/gziphandler v1.1.1
	github.com/SAP/go-hdb v0.14.1
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/aerospike/aerospike-client-go v3.1.1+incompatible
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190620160927-9418d7b0cd0f
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/apple/foundationdb/bindings/go v0.0.0-20190411004307-cd5c9d91fad2
	github.com/armon/go-metrics v0.3.9
	github.com/armon/go-proxyproto v0.0.0-20210323213023-7e956b284f0a
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/aws/aws-sdk-go v1.37.19
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/client9/misspell v0.3.4
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/denisenkom/go-mssqldb v0.0.0-20200428022330-06a60b6afbbc
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/duosecurity/duo_api_golang v0.0.0-20190308151101-6c680f768e74
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/favadi/protoc-go-inject-tag v1.3.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.0.1
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.7
	github.com/gocql/gocql v0.0.0-20210401103645-80ab1e13e309
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.2.0
	github.com/hashicorp/cap v0.1.0
	github.com/hashicorp/consul-template v0.27.2-0.20211014231529-4ff55381f1c4
	github.com/hashicorp/consul/api v1.11.0
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20210818145131-c573d69da192
	github.com/hashicorp/go-gcp-common v0.7.0
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/go-kms-wrapping v0.6.7
	github.com/hashicorp/go-memdb v1.3.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-raftchunking v0.6.3-0.20191002164813-7e9e8525653a
	github.com/hashicorp/go-retryablehttp v0.6.7
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-secure-stdlib/awsutil v0.1.5
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.1
	github.com/hashicorp/go-secure-stdlib/gatedwriter v0.1.1
	github.com/hashicorp/go-secure-stdlib/kv-builder v0.1.1
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.1
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/password v0.1.1
	github.com/hashicorp/go-secure-stdlib/reloadutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.1
	github.com/hashicorp/go-secure-stdlib/tlsutil v0.1.1
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hashicorp/hcl v1.0.1-vault-3
	github.com/hashicorp/nomad/api v0.0.0-20211006193434-215bf04bc650
	github.com/hashicorp/raft v1.3.1
	github.com/hashicorp/raft-autopilot v0.1.3
	github.com/hashicorp/raft-boltdb/v2 v2.0.0-20210421194847-a7e34179d62c
	github.com/hashicorp/raft-snapshot v1.0.3
	github.com/hashicorp/vault-plugin-auth-alicloud v0.9.0
	github.com/hashicorp/vault-plugin-auth-azure v0.8.0
	github.com/hashicorp/vault-plugin-auth-centrify v0.9.0
	github.com/hashicorp/vault-plugin-auth-cf v0.9.0
	github.com/hashicorp/vault-plugin-auth-gcp v0.10.0
	github.com/hashicorp/vault-plugin-auth-jwt v0.10.1
	github.com/hashicorp/vault-plugin-auth-kerberos v0.4.0
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.11.1-0.20210929181055-821e911b1751
	github.com/hashicorp/vault-plugin-auth-oci v0.8.0
	github.com/hashicorp/vault-plugin-database-couchbase v0.3.1-0.20210902192635-c3ee7c5bc378
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.8.0
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.4.0
	github.com/hashicorp/vault-plugin-database-snowflake v0.2.1
	github.com/hashicorp/vault-plugin-mock v0.16.1
	github.com/hashicorp/vault-plugin-secrets-ad v0.10.0
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.9.0
	github.com/hashicorp/vault-plugin-secrets-azure v0.6.3-0.20210924190759-58a034528e35
	github.com/hashicorp/vault-plugin-secrets-gcp v0.10.2
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.9.0
	github.com/hashicorp/vault-plugin-secrets-kv v0.5.7-0.20211013154503-eec8a1c892fb
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.4.0
	github.com/hashicorp/vault-plugin-secrets-openldap v0.4.1-0.20210921171411-e86105e4986d
	github.com/hashicorp/vault-plugin-secrets-terraform v0.1.1-0.20210715043003-e02ca8f6408e
	github.com/hashicorp/vault-testing-stepwise v0.1.1
	github.com/hashicorp/vault/api v1.1.1
	github.com/hashicorp/vault/sdk v0.2.2-0.20211004171540-a8c7e135dd6a
	github.com/influxdata/influxdb v0.0.0-20190411212539-d24b7ba8c4c4
	github.com/jcmturner/gokrb5/v8 v8.0.0
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.0
	github.com/joyent/triton-go v1.7.1-0.20200416154420-6801d15b779f
	github.com/keybase/go-crypto v0.0.0-20190403132359-d65b6b94177f
	github.com/kr/pretty v0.3.0
	github.com/kr/text v0.2.0
	github.com/lib/pq v1.10.3
	github.com/mattn/go-colorable v0.1.11
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/michaelklishin/rabbit-hole v0.0.0-20191008194146-93d9988f0cd5
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.0
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mitchellh/reflectwalk v1.0.2
	github.com/mongodb/go-client-mongodb-atlas v0.1.2
	github.com/natefinch/atomic v0.0.0-20150920032501-a62ce929ffcc
	github.com/ncw/swift v1.0.47
	github.com/oklog/run v1.0.0
	github.com/okta/okta-sdk-golang/v2 v2.0.0
	github.com/oracle/oci-go-sdk v13.1.0+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.6.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/pquerna/otp v1.2.1-0.20191009055518-468c2dd2b58d
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/rboyer/safeio v0.2.1
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
	github.com/sasha-s/go-deadlock v0.2.0
	github.com/sethvargo/go-limiter v0.7.1
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/stretchr/testify v1.7.0
	go.etcd.io/bbolt v1.3.5
	go.etcd.io/etcd/client/pkg/v3 v3.5.0
	go.etcd.io/etcd/client/v2 v2.305.0
	go.etcd.io/etcd/client/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.7.3
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	go.uber.org/atomic v1.9.0
	go.uber.org/goleak v1.1.11-0.20210813005559-691160354723
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/net v0.0.0-20210928044308-7d9f5e0b762b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20211004093028-2c5d950f24ef
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
	golang.org/x/tools v0.1.5
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20210928142010-c7af6a1a74c9 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gopkg.in/square/go-jose.v2 v2.5.1
	layeh.com/radius v0.0.0-20190322222518-890bc1058917
	mvdan.cc/gofumpt v0.1.1
)

require (
	code.cloudfoundry.org/gofileutils v0.0.0-20170111115228-4d0c80011a0f // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.0 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/DataDog/datadog-go v3.2.0+incompatible // indirect
	github.com/Jeffail/gabs v1.1.1 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.4.16-0.20201130162521-d1ffc52c7331 // indirect
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apache/arrow/go/arrow v0.0.0-20200601151325-b2287a20f230 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v12 v12.0.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.1.5 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.1.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.0.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.5.0 // indirect
	github.com/aws/smithy-go v1.3.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/briankassouf/jose v0.9.2-0.20180619214549-d2569464773f // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/centrify/cloud-golang-sdk v0.0.0-20190214225812-119110094d0f // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/circonus-labs/circonus-gometrics v2.3.1+incompatible // indirect
	github.com/circonus-labs/circonusllhist v0.1.3 // indirect
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20190201205600-f136f9222381 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-oidc/v3 v3.0.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/couchbase/gocb/v2 v2.3.1 // indirect
	github.com/couchbase/gocbcore/v10 v10.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denverdino/aliyungo v0.0.0-20170926055100-d3308649c661 // indirect
	github.com/digitalocean/godo v1.7.5 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/evanphx/json-patch/v5 v5.5.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/gammazero/deque v0.0.0-20190130191400-2afb3858e9c7 // indirect
	github.com/gammazero/workerpool v0.0.0-20190406235159-88d534f22b56 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/go-ldap/ldif v0.0.0-20200320164324-fd88d9b715b3 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v1.11.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/tink/go v1.4.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/cronexpr v1.1.0 // indirect
	github.com/hashicorp/go-fpe v0.0.0-20210209202841-f9bb8a2d8eab // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-kms-wrapping/entropy v0.1.0 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/go-plugin v1.4.3 // indirect
	github.com/hashicorp/go-slug v0.4.1 // indirect
	github.com/hashicorp/go-tfe v0.12.0 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/mdns v1.0.1 // indirect
	github.com/hashicorp/serf v0.9.5 // indirect
	github.com/hashicorp/vic v1.5.1-0.20190403131502-bbfe86ec9443 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jackc/pgx v3.3.0+incompatible // indirect
	github.com/jcmturner/aescts v1.0.1 // indirect
	github.com/jcmturner/dnsutils v1.0.1 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/rpc/v2 v2.0.2 // indirect
	github.com/jeffchao/backoff v0.0.0-20140404060208-9d7fd7aa17f2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmoiron/sqlx v1.3.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/linode/linodego v0.7.1 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/miekg/dns v1.1.40 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/colorstring v0.0.0-20150917214807-8631ce90f286 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/mitchellh/iochan v1.0.0 // indirect
	github.com/mitchellh/pointerstructure v1.0.0 // indirect
	github.com/moby/term v0.0.0-20200915141129-7f0af18e79f2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/mwielbut/pointy v1.1.0 // indirect
	github.com/nicolai86/scaleway-sdk v1.10.2-0.20180628010248-798f60e20bb2 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/okta/okta-sdk-golang v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc9 // indirect
	github.com/openlyinc/pointy v1.1.2 // indirect
	github.com/packethost/packngo v0.1.1-0.20180711074735-b9cb5096f54c // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/snowflakedb/gosnowflake v1.6.1 // indirect
	github.com/softlayer/softlayer-go v0.0.0-20180806151055-260589d94c7d // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/svanharmelen/jsonapi v0.0.0-20180618144545-0c0828c3f16d // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.162 // indirect
	github.com/tidwall/pretty v1.0.1 // indirect
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/tv42/httpunix v0.0.0-20150427012821-b75d8614f926 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vmware/govmomi v0.18.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	github.com/zclconf/go-cty v1.2.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.0 // indirect
	go.mongodb.org/atlas v0.7.1 // indirect
	go.opencensus.io v0.22.3 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
	k8s.io/api v0.18.2 // indirect
	k8s.io/apimachinery v0.18.2 // indirect
	k8s.io/client-go v0.18.2 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
