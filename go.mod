module github.com/hashicorp/vault

go 1.17

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/api/auth/approle => ./api/auth/approle

replace github.com/hashicorp/vault/api/auth/userpass => ./api/auth/userpass

replace github.com/hashicorp/vault/sdk => ./sdk

replace go.etcd.io/etcd/client/pkg/v3 v3.5.0 => go.etcd.io/etcd/client/pkg/v3 v3.0.0-20210928084031-3df272774672

require (
	cloud.google.com/go v0.65.0
	cloud.google.com/go/spanner v1.5.1
	cloud.google.com/go/storage v1.10.0
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/Azure/go-autorest/autorest v0.11.21
	github.com/Azure/go-autorest/autorest/adal v0.9.14
	github.com/NYTimes/gziphandler v1.1.1
	github.com/SAP/go-hdb v0.14.1
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/aerospike/aerospike-client-go v3.1.1+incompatible
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190620160927-9418d7b0cd0f
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/apple/foundationdb/bindings/go v0.0.0-20190411004307-cd5c9d91fad2
	github.com/armon/go-metrics v0.3.10
	github.com/armon/go-proxyproto v0.0.0-20210323213023-7e956b284f0a
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/aws/aws-sdk-go v1.37.19
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/client9/misspell v0.3.4
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/denisenkom/go-mssqldb v0.11.0
	github.com/docker/docker v20.10.10+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/duosecurity/duo_api_golang v0.0.0-20190308151101-6c680f768e74
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/favadi/protoc-go-inject-tag v1.3.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.4.1
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.8
	github.com/gocql/gocql v0.0.0-20210401103645-80ab1e13e309
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.2.0
	github.com/hashicorp/cap v0.1.1
	github.com/hashicorp/consul-template v0.27.2-0.20211014231529-4ff55381f1c4
	github.com/hashicorp/consul/api v1.11.0
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20210818145131-c573d69da192
	github.com/hashicorp/go-gcp-common v0.7.0
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/go-kms-wrapping v0.6.8
	github.com/hashicorp/go-memdb v1.3.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-raftchunking v0.6.3-0.20191002164813-7e9e8525653a
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-secure-stdlib/awsutil v0.1.5
	github.com/hashicorp/go-secure-stdlib/base62 v0.1.2
	github.com/hashicorp/go-secure-stdlib/gatedwriter v0.1.1
	github.com/hashicorp/go-secure-stdlib/kv-builder v0.1.1
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.1
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.2
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
	github.com/hashicorp/vault-plugin-auth-alicloud v0.10.0
	github.com/hashicorp/vault-plugin-auth-azure v0.9.2
	github.com/hashicorp/vault-plugin-auth-centrify v0.10.0
	github.com/hashicorp/vault-plugin-auth-cf v0.10.0
	github.com/hashicorp/vault-plugin-auth-gcp v0.11.2
	github.com/hashicorp/vault-plugin-auth-jwt v0.11.2
	github.com/hashicorp/vault-plugin-auth-kerberos v0.5.0
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.11.3
	github.com/hashicorp/vault-plugin-auth-oci v0.9.0
	github.com/hashicorp/vault-plugin-database-couchbase v0.5.1
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.9.1
	github.com/hashicorp/vault-plugin-database-mongodbatlas v0.5.1
	github.com/hashicorp/vault-plugin-database-snowflake v0.3.1
	github.com/hashicorp/vault-plugin-mock v0.16.1
	github.com/hashicorp/vault-plugin-secrets-ad v0.11.1
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.10.2
	github.com/hashicorp/vault-plugin-secrets-azure v0.11.1
	github.com/hashicorp/vault-plugin-secrets-gcp v0.11.0
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.10.0
	github.com/hashicorp/vault-plugin-secrets-kv v0.10.1
	github.com/hashicorp/vault-plugin-secrets-mongodbatlas v0.5.1
	github.com/hashicorp/vault-plugin-secrets-openldap v0.6.0
	github.com/hashicorp/vault-plugin-secrets-terraform v0.3.0
	github.com/hashicorp/vault-testing-stepwise v0.1.2
	github.com/hashicorp/vault/api v1.3.0
	github.com/hashicorp/vault/api/auth/approle v0.0.0-00010101000000-000000000000
	github.com/hashicorp/vault/api/auth/userpass v0.0.0-00010101000000-000000000000
	github.com/hashicorp/vault/sdk v0.3.0
	github.com/influxdata/influxdb v0.0.0-20190411212539-d24b7ba8c4c4
	github.com/jcmturner/gokrb5/v8 v8.4.2
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.0
	github.com/joyent/triton-go v1.7.1-0.20200416154420-6801d15b779f
	github.com/keybase/go-crypto v0.0.0-20190403132359-d65b6b94177f
	github.com/kr/pretty v0.3.0
	github.com/kr/text v0.2.0
	github.com/lib/pq v1.10.3
	github.com/mattn/go-colorable v0.1.11
	github.com/mholt/archiver/v3 v3.5.1
	github.com/michaelklishin/rabbit-hole/v2 v2.11.0
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
	github.com/ory/dockertest/v3 v3.8.0
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
	go.etcd.io/bbolt v1.3.6
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
	golang.org/x/net v0.0.0-20211020060615-d418f374d309
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	golang.org/x/tools v0.1.5
	google.golang.org/api v0.30.0
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gopkg.in/square/go-jose.v2 v2.6.0
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b
	layeh.com/radius v0.0.0-20190322222518-890bc1058917
	mvdan.cc/gofumpt v0.1.1
)

require (
	github.com/Microsoft/hcsshim v0.9.0 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/miekg/dns v1.1.40 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	go.uber.org/zap v1.19.1 // indirect
	k8s.io/client-go v0.22.2 // indirect
)
