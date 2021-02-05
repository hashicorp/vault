module github.com/hashicorp/vault

go 1.13

replace github.com/hashicorp/vault/api => ./api

replace github.com/hashicorp/vault/sdk => ./sdk

require (
	bazil.org/fuse v0.0.0-20200524192727-fb710f7dfd05 // indirect
	cloud.google.com/go v0.76.0
	cloud.google.com/go/bigquery v1.15.0 // indirect
	cloud.google.com/go/bigtable v1.7.1 // indirect
	cloud.google.com/go/datastore v1.4.0 // indirect
	cloud.google.com/go/pubsub v1.9.1 // indirect
	cloud.google.com/go/spanner v1.13.0
	cloud.google.com/go/storage v1.13.0
	collectd.org v0.5.0 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20201218220906-28db891af037 // indirect
	git.apache.org/thrift.git v0.13.0 // indirect
	github.com/Azure/azure-sdk-for-go v51.1.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/adal v0.9.11
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // indirect
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/DataDog/datadog-go v4.3.1+incompatible // indirect
	github.com/Jeffail/gabs v1.4.0 // indirect
	github.com/Jeffail/gabs/v2 v2.6.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/Microsoft/hcsshim v0.8.14 // indirect
	github.com/NYTimes/gziphandler v1.1.1
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/SAP/go-hdb v0.102.7
	github.com/Sectorbob/mlab-ns2 v0.0.0-20171030222938-d3aa0c295a8a
	github.com/Shopify/sarama v1.27.2 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/aerospike/aerospike-client-go v4.1.0+incompatible
	github.com/ajstarks/svgo v0.0.0-20200725142600-7a3c8b57fecb // indirect
	github.com/alecthomas/units v0.0.0-20201120081800-1786d5ef83d4 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.913
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20210205162634-132918040f1c // indirect
	github.com/apple/foundationdb/bindings/go v0.0.0-20210204204729-e4a55908ff04
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2 // indirect
	github.com/armon/go-metrics v0.3.6
	github.com/armon/go-proxyproto v0.0.0-20200108142055-f0b8253b1507
	github.com/armon/go-radix v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/aws/aws-lambda-go v1.22.0 // indirect
	github.com/aws/aws-sdk-go v1.37.5
	github.com/aws/aws-sdk-go-v2 v1.1.0 // indirect
	github.com/benbjohnson/immutable v0.3.0 // indirect
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/c-bata/go-prompt v0.2.5 // indirect
	github.com/c2h5oh/datasize v0.0.0-20200825124411-48ed595a09d2 // indirect
	github.com/casbin/casbin/v2 v2.23.0 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/centrify/cloud-golang-sdk v0.0.0-20200612223121-348d1cfa8842 // indirect
	github.com/checkpoint-restore/go-criu v4.0.0+incompatible // indirect
	github.com/chrismalek/oktasdk-go v0.0.0-20181212195951-3430665dfaa0
	github.com/cilium/ebpf v0.3.0 // indirect
	github.com/circonus-labs/circonusllhist v0.1.4 // indirect
	github.com/client9/misspell v0.3.4
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20201123235753-4f46d6348a05 // indirect
	github.com/cncf/udpa/go v0.0.0-20201211205326-cc1b757b3edd // indirect
	github.com/cockroachdb/cockroach-go v2.0.1+incompatible
	github.com/cockroachdb/errors v1.8.2 // indirect
	github.com/cockroachdb/redact v1.0.9 // indirect
	github.com/containerd/cgroups v0.0.0-20210114181951-8a68de567b68 // indirect
	github.com/containerd/console v1.0.1 // indirect
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/containerd/continuity v0.0.0-20201208142359-180525291bb7 // indirect
	github.com/containerd/fifo v0.0.0-20210129194248-f8e8fdba47ef // indirect
	github.com/containerd/go-runc v0.0.0-20201020171139-16b287bc67d0 // indirect
	github.com/containerd/ttrpc v1.0.2 // indirect
	github.com/containerd/typeurl v1.0.1 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/couchbase/gocb/v2 v2.2.0 // indirect
	github.com/couchbase/gocbcore/v9 v9.1.1 // indirect
	github.com/creack/pty v1.1.11 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/dave/jennifer v1.4.1 // indirect
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/denverdino/aliyungo v0.0.0-20210113054000-11eaa932b667 // indirect
	github.com/dgryski/go-sip13 v0.0.0-20200911182023-62edffca9245 // indirect
	github.com/digitalocean/godo v1.57.0 // indirect
	github.com/dnaeon/go-vcr v1.1.0 // indirect
	github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dsnet/golib v1.0.2 // indirect
	github.com/duosecurity/duo_api_golang v0.0.0-20201112143038-0e07e9f869e3
	github.com/eclipse/paho.mqtt.golang v1.3.1 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/envoyproxy/protoc-gen-validate v0.4.1 // indirect
	github.com/fatih/color v1.10.0
	github.com/fatih/structs v1.1.0
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/franela/goblin v0.0.0-20210113153425-413781f5e6c8 // indirect
	github.com/frankban/quicktest v1.11.3 // indirect
	github.com/fullsailor/pkcs7 v0.0.0-20190404230743-d7302db945fa
	github.com/gammazero/deque v0.0.0-20201010052221-3932da5530cc // indirect
	github.com/gammazero/workerpool v1.1.1 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/glycerine/go-unsnap-stream v0.0.0-20210130063903-47dfef350d96 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.3 // indirect
	github.com/go-errors/errors v1.1.1
	github.com/go-gl/glfw v0.0.0-20201108214237-06ea97f0c265 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20201108214237-06ea97f0c265 // indirect
	github.com/go-latex/latex v0.0.0-20210118124228-b3d85cf34e07 // indirect
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-openapi/swag v0.19.13 // indirect
	github.com/go-resty/resty/v2 v2.4.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-test/deep v1.0.7
	github.com/gobuffalo/attrs v1.0.0 // indirect
	github.com/gobuffalo/depgen v0.2.0 // indirect
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/flect v0.2.2 // indirect
	github.com/gobuffalo/genny v0.6.0 // indirect
	github.com/gobuffalo/gogen v0.2.0 // indirect
	github.com/gobuffalo/mapi v1.2.1 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.8.1 // indirect
	github.com/gobuffalo/syncx v0.1.0 // indirect
	github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/geo v0.0.0-20210108004804-a63082ebfb66 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/flatbuffers v1.12.0 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-metrics-stackdriver v0.2.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210125172800-10e9aeb4a998 // indirect
	github.com/google/renameio v1.0.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/gophercloud/gophercloud v0.15.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20210202160940-bed99a852dfe // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/consul-template v0.25.2-0.20210123001810-166043f8559d
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-bindata v3.0.8-0.20180209072458-bf7910af8997+incompatible
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-discover v0.0.0-20201029210230-738cb3105cd0
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
	github.com/hashicorp/raft-boltdb v0.0.0-20191021154308-4207f1bf0617 // indirect
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
	github.com/huaweicloud/golangsdk v0.0.0-20210205071117-066cac2eec52 // indirect
	github.com/iancoleman/strcase v0.1.3 // indirect
	github.com/influxdata/flux v0.105.0 // indirect
	github.com/influxdata/influxdb v1.8.4
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab // indirect
	github.com/influxdata/line-protocol v0.0.0-20201012155213-5f565037cbc9 // indirect
	github.com/influxdata/tdigest v0.0.1 // indirect
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/jcmturner/aescts v2.0.0+incompatible // indirect
	github.com/jcmturner/dnsutils v2.0.0+incompatible // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2
	github.com/jefferai/isbadcipher v0.0.0-20190226160619-51d2077c035f
	github.com/jefferai/jsonx v1.0.1
	github.com/jhump/protoreflect v1.8.1 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/joyent/triton-go v1.8.5
	github.com/jsternberg/zap-logfmt v1.2.0 // indirect
	github.com/jung-kurt/gofpdf v1.16.2 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/keybase/go-crypto v0.0.0-20200123153347-de78d2cb44f4
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/crc32 v1.2.0 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kr/logfmt v0.0.0-20210122060352-19f9bcb100e6 // indirect
	github.com/kr/pretty v0.2.1
	github.com/kr/pty v1.1.8 // indirect
	github.com/kr/text v0.2.0
	github.com/lestrrat-go/jwx v1.1.1 // indirect
	github.com/lib/pq v1.9.0
	github.com/lightstep/lightstep-tracer-common/golang/gogo v0.0.0-20200310182322-adf4263e074b // indirect
	github.com/lightstep/lightstep-tracer-go v0.24.0 // indirect
	github.com/linode/linodego v1.0.0 // indirect
	github.com/lyft/protoc-gen-star v0.5.2 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/markbates/oncer v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/mattn/go-shellwords v1.0.11 // indirect
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/michaelklishin/rabbit-hole v0.0.0-20191008194146-93d9988f0cd5
	github.com/miekg/dns v1.1.38 // indirect
	github.com/mitchellh/cli v1.1.2
	github.com/mitchellh/copystructure v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/gox v1.0.1
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mitchellh/pointerstructure v1.1.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.1
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/montanaflynn/stats v0.6.4 // indirect
	github.com/mrunalp/fileutils v0.5.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/atomic v0.0.0-20200526193002-18c0533a5b09
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	github.com/ncw/swift v1.0.53
	github.com/neelance/astrewrite v0.0.0-20160511093645-99348263ae86 // indirect
	github.com/neelance/sourcemap v0.0.0-20200213170602-2833bce08e4c // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/nxadm/tail v1.4.6 // indirect
	github.com/oklog/run v1.1.0
	github.com/okta/okta-sdk-golang v1.1.0 // indirect
	github.com/okta/okta-sdk-golang/v2 v2.3.0
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/opencontainers/selinux v1.8.0 // indirect
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.5 // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.6.3
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/packethost/packngo v0.6.0 // indirect
	github.com/pact-foundation/pact-go v1.5.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/peterh/liner v1.2.1 // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pierrec/lz4 v2.6.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0 // indirect
	github.com/posener/complete v1.2.3
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/pquerna/otp v1.3.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/prometheus/procfs v0.4.0 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rboyer/safeio v0.2.1
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/renier/xmlrpc v0.0.0-20191022213033-ce560eccbd00 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.7.0 // indirect
	github.com/rs/zerolog v1.20.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/ryanuber/go-glob v1.0.0
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/sasha-s/go-deadlock v0.2.0
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/segmentio/kafka-go v0.4.9 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sethvargo/go-limiter v0.6.0
	github.com/shirou/gopsutil v3.21.1+incompatible
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/snowflakedb/gosnowflake v1.4.0 // indirect
	github.com/softlayer/softlayer-go v1.0.2 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	github.com/square/go-jose v2.5.1+incompatible // indirect
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/streadway/handy v0.0.0-20200128134331-0f66f006fb2e // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.233+incompatible // indirect
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/tinylib/msgp v1.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c // indirect
	github.com/ugorji/go v1.2.4 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/vmware/govmomi v0.24.0 // indirect
	github.com/vmware/vmw-guestinfo v0.0.0-20200218095840-687661b8bd8e // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xlab/treeprint v1.0.0 // indirect
	github.com/yandex-cloud/go-genproto v0.0.0-20210204150426-0e4d55ea0f0c // indirect
	github.com/yandex-cloud/go-sdk v0.0.0-20210204095654-4382e50bda0f // indirect
	github.com/yuin/goldmark v1.3.1 // indirect
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da // indirect
	go.etcd.io/bbolt v1.3.5
	go.etcd.io/etcd v3.3.25+incompatible
	go.mongodb.org/atlas v0.7.2
	go.mongodb.org/mongo-driver v1.4.6
	go.opencensus.io v0.22.6 // indirect
	go.uber.org/atomic v1.7.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/exp v0.0.0-20210201131500-d352d2db2ceb // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	golang.org/x/tools v0.1.0
	gonum.org/v1/gonum v0.8.2 // indirect
	gonum.org/v1/netlib v0.0.0-20201012070519-2390d26c3658 // indirect
	gonum.org/v1/plot v0.8.1 // indirect
	google.golang.org/api v0.39.0
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460 // indirect
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/grpc/examples v0.0.0-20210205041354-b753f4903c1b // indirect
	google.golang.org/protobuf v1.25.1-0.20200805231151-a709e31e5d12
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/ory-am/dockertest.v3 v3.3.4
	gopkg.in/square/go-jose.v2 v2.5.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools/v3 v3.0.3 // indirect
	honnef.co/go/tools v0.1.1 // indirect
	layeh.com/radius v0.0.0-20201203135236-838e26d0c9be
	rsc.io/sampler v1.99.99 // indirect
)
