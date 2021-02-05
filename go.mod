module github.com/hashicorp/vault

go 1.13

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/sdk => ./sdk

require (
	cloud.google.com/go v0.60.0
	cloud.google.com/go/spanner v1.7.0
	cloud.google.com/go/storage v1.10.0
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/adal v0.9.11
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/DataDog/datadog-go v4.3.1+incompatible // indirect
	github.com/Jeffail/gabs/v2 v2.6.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/NYTimes/gziphandler v1.1.1
	github.com/SAP/go-hdb v0.102.7
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/aerospike/aerospike-client-go v4.1.0+incompatible
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.913
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20200923215132-ac86123a3f01 // indirect
	github.com/apple/foundationdb/bindings/go v0.0.0-20210204204729-e4a55908ff04
	github.com/armon/go-metrics v0.3.6
	github.com/armon/go-proxyproto v0.0.0-20200108142055-f0b8253b1507
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/aws/aws-sdk-go v1.37.5
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/centrify/cloud-golang-sdk v0.0.0-20200612223121-348d1cfa8842 // indirect
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/circonus-labs/circonusllhist v0.1.4 // indirect
	github.com/client9/misspell v0.3.4
	github.com/cockroachdb/cockroach-go v2.0.1+incompatible
	github.com/containerd/cgroups v0.0.0-20210114181951-8a68de567b68 // indirect
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	// github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/couchbase/gocb/v2 v2.2.0 // indirect
	github.com/couchbase/gocbcore/v9 v9.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/denverdino/aliyungo v0.0.0-20210113054000-11eaa932b667 // indirect
	github.com/digitalocean/godo v1.57.0 // indirect
	github.com/dnaeon/go-vcr v1.1.0 // indirect
	github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/duosecurity/duo_api_golang v0.0.0-20201112143038-0e07e9f869e3
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/fatih/color v1.10.0
	github.com/fatih/structs v1.1.0
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/fullsailor/pkcs7 v0.0.0-20190404230743-d7302db945fa
	github.com/gammazero/deque v0.0.0-20201010052221-3932da5530cc // indirect
	github.com/gammazero/workerpool v1.1.1 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-asn1-ber/asn1-ber v1.5.3 // indirect
	github.com/go-errors/errors v1.1.1
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-resty/resty/v2 v2.4.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.7
	github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/flatbuffers v1.12.0 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.2.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/gophercloud/gophercloud v0.15.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20210202160940-bed99a852dfe // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/hashicorp/consul-template v0.25.2-0.20210123001810-166043f8559d
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-bindata v3.0.8-0.20180209072458-bf7910af8997+incompatible
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20200812215701-c4b85f6ed31f
	github.com/hashicorp/go-gcp-common v0.6.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/hashicorp/go-kms-wrapping v0.6.0
	github.com/hashicorp/go-memdb v1.3.1
	github.com/hashicorp/go-msgpack v1.1.5
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-plugin v1.4.0 // indirect
	github.com/hashicorp/go-raftchunking v0.6.3-0.20191002164813-7e9e8525653a
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hashicorp/hcl v1.0.1-vault
	github.com/hashicorp/mdns v1.0.3 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20210205152330-744c4470b9c7
	github.com/hashicorp/raft v1.2.0
	github.com/hashicorp/raft-snapshot v1.0.3
	github.com/hashicorp/vault-plugin-auth-alicloud v0.7.0
	github.com/hashicorp/vault-plugin-auth-azure v0.6.0
	github.com/hashicorp/vault-plugin-auth-centrify v0.7.0
	github.com/hashicorp/vault-plugin-auth-cf v0.7.0
	github.com/hashicorp/vault-plugin-auth-gcp v0.8.0
	github.com/hashicorp/vault-plugin-auth-jwt v0.8.1
	github.com/hashicorp/vault-plugin-auth-kerberos v0.2.0
	// github.com/hashicorp/vault-plugin-auth-kerberos v0.2.0
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.8.0
	github.com/hashicorp/vault-plugin-auth-oci v0.6.0
	github.com/hashicorp/vault-plugin-database-couchbase v0.2.1
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.6.1
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.2.1
	github.com/hashicorp/vault-plugin-database-snowflake v0.1.1
	github.com/hashicorp/vault-plugin-mock v0.16.1
	github.com/hashicorp/vault-plugin-secrets-ad v0.8.0
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.7.0
	github.com/hashicorp/vault-plugin-secrets-azure v0.8.0
	github.com/hashicorp/vault-plugin-secrets-gcp v0.8.2
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.7.0
	github.com/hashicorp/vault-plugin-secrets-kv v0.7.0
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.2.0
	github.com/hashicorp/vault-plugin-secrets-openldap v0.1.6-0.20210201204049-4f0f91977798
	github.com/hashicorp/vault/api v1.0.5-0.20201001211907-38d91b749c77
	github.com/hashicorp/vault/sdk v0.1.14-0.20210127185906-6b455835fa8c
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/influxdata/influxdb v1.8.4
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.1
	github.com/jhump/protoreflect v1.8.1 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/joyent/triton-go v1.8.5
	github.com/keybase/go-crypto v0.0.0-20200123153347-de78d2cb44f4
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/kr/pretty v0.2.1
	github.com/kr/text v0.2.0
	github.com/lib/pq v1.9.0
	github.com/linode/linodego v1.0.0 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-shellwords v1.0.11 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/michaelklishin/rabbit-hole v0.0.0-20191008194146-93d9988f0cd5
	github.com/miekg/dns v1.1.38 // indirect
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/copystructure v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mitchellh/pointerstructure v1.1.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.1
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/natefinch/atomic v0.0.0-20200526193002-18c0533a5b09
	github.com/ncw/swift v1.0.53
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/nxadm/tail v1.4.6 // indirect
	github.com/oklog/run v1.1.0
	github.com/okta/okta-sdk-golang/v2 v2.3.0
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.6.3
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/packethost/packngo v0.6.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/pquerna/otp v1.3.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/procfs v0.4.0 // indirect
	github.com/rboyer/safeio v0.2.1
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/sasha-s/go-deadlock v0.2.0
	github.com/sethvargo/go-limiter v0.6.0
	github.com/shirou/gopsutil v3.21.1+incompatible
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/snowflakedb/gosnowflake v1.4.0 // indirect
	github.com/softlayer/softlayer-go v1.0.2 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.233+incompatible // indirect
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vmware/govmomi v0.24.0 // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.etcd.io/bbolt v1.3.5
	go.etcd.io/etcd v3.3.25+incompatible
	go.mongodb.org/atlas v0.7.2
	go.mongodb.org/mongo-driver v1.4.6
	go.opencensus.io v0.22.5 // indirect
	go.uber.org/atomic v1.7.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	golang.org/x/tools v0.1.0
	google.golang.org/api v0.29.0
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.1-0.20200805231151-a709e31e5d12
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gopkg.in/square/go-jose.v2 v2.5.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.0.3 // indirect
	honnef.co/go/tools v0.1.1 // indirect
	layeh.com/radius v0.0.0-20201203135236-838e26d0c9be
)
