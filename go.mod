module github.com/hashicorp/vault

go 1.12

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/sdk => ./sdk

require (
	cloud.google.com/go v0.39.0
	github.com/Azure/azure-sdk-for-go v29.0.0+incompatible
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Azure/go-autorest v11.7.1+incompatible
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/NYTimes/gziphandler v1.1.1
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/SAP/go-hdb v0.14.1
	github.com/abdullin/seq v0.0.0-20160510034733-d5467c17e7af // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190620160927-9418d7b0cd0f
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/apple/foundationdb/bindings/go v0.0.0-20190411004307-cd5c9d91fad2
	github.com/armon/go-metrics v0.0.0-20190430140413-ec5e00d3c878
	github.com/armon/go-proxyproto v0.0.0-20190211145416-68259f75880e
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/aws/aws-sdk-go v1.19.39
	github.com/bitly/go-hostpool v0.0.0-20171023180738-a3a6125de932 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/boombuler/barcode v1.0.0 // indirect
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c
	github.com/containerd/continuity v0.0.0-20181203112020-004b46473808 // indirect
	github.com/coreos/go-semver v0.2.0
	github.com/denisenkom/go-mssqldb v0.0.0-20190412130859-3b1d194e553a
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/duosecurity/duo_api_golang v0.0.0-20190308151101-6c680f768e74
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/fatih/color v1.7.0
	github.com/fatih/structs v1.1.0
	github.com/fullsailor/pkcs7 v0.0.0-20190404230743-d7302db945fa
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.0.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-test/deep v1.0.2
	github.com/gocql/gocql v0.0.0-20190402132108-0e1d5de854df
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/consul/api v1.0.1
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-gcp-common v0.5.0
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-memdb v1.0.2
	github.com/hashicorp/go-msgpack v0.5.5
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-raftchunking v0.6.1
	github.com/hashicorp/go-rootcerts v1.0.1
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/golang-lru v0.5.1
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/nomad/api v0.0.0-20190412184103-1c38ced33adf
	github.com/hashicorp/raft v1.1.1
	github.com/hashicorp/raft-snapshot v1.0.1
	github.com/hashicorp/vault-plugin-auth-alicloud v0.5.2-0.20190725165955-db428c1b0976
	github.com/hashicorp/vault-plugin-auth-azure v0.5.2-0.20190725170003-0fc4bd518aa2
	github.com/hashicorp/vault-plugin-auth-centrify v0.5.2-0.20190725170023-9d8d5e0b03e6
	github.com/hashicorp/vault-plugin-auth-gcp v0.5.2-0.20190725170032-0aa7d1c92039
	github.com/hashicorp/vault-plugin-auth-jwt v0.5.2-0.20190725170041-1cfee03e8d3a
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.5.2-0.20190725170047-354505be0ecf
	github.com/hashicorp/vault-plugin-auth-pcf v0.0.0-20190725170053-826a135618c1
	github.com/hashicorp/vault-plugin-database-elasticsearch v0.0.0-20190725170059-5c0d558eb59d
	github.com/hashicorp/vault-plugin-secrets-ad v0.5.3-0.20190725170108-e1b17ad0c772
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.5.2-0.20190725170114-7d66a3fa0600
	github.com/hashicorp/vault-plugin-secrets-azure v0.5.2-0.20190725170121-541440395211
	github.com/hashicorp/vault-plugin-secrets-gcp v0.5.3-0.20190725170127-aa49df112140
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.5.2-0.20190725170135-aaf270943731
	github.com/hashicorp/vault-plugin-secrets-kv v0.5.2-0.20190725170141-1c4dac87f383
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/influxdata/influxdb v0.0.0-20190411212539-d24b7ba8c4c4
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.3.0+incompatible // indirect
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.0
	github.com/joyent/triton-go v0.0.0-20190112182421-51ffac552869
	github.com/keybase/go-crypto v0.0.0-20190403132359-d65b6b94177f
	github.com/kr/pretty v0.1.0
	github.com/kr/text v0.1.0
	github.com/lib/pq v1.1.1
	github.com/mattn/go-colorable v0.0.9
	github.com/michaelklishin/rabbit-hole v1.5.0
	github.com/mitchellh/cli v1.0.0
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/reflectwalk v1.0.0
	github.com/ncw/swift v1.0.47
	github.com/oklog/run v1.0.0
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.4+incompatible
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/posener/complete v1.2.1
	github.com/pquerna/otp v1.1.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/common v0.2.0
	github.com/ryanuber/columnize v2.1.0+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94 // indirect
	github.com/stretchr/testify v1.3.0
	go.etcd.io/bbolt v1.3.2
	go.etcd.io/etcd v0.0.0-20190412021913-f29b1ada1971
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	google.golang.org/api v0.5.0
	google.golang.org/genproto v0.0.0-20190513181449-d00d292a067c
	google.golang.org/grpc v1.22.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gopkg.in/square/go-jose.v2 v2.3.1
	gotest.tools v2.2.0+incompatible // indirect
	layeh.com/radius v0.0.0-20190322222518-890bc1058917
)
