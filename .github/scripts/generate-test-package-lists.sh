# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This script is meant to be sourced into the shell running in a Github
# workflow.

# This script is a temporary measure until we implement a dynamic test-splitting
# solution. It distributes the entire set of test packages into 16 sublists,
# which should roughly take an equal amount of time to complete.

test_packages=()

base="github.com/hashicorp/vault"

# Total time: 526
test_packages[1]+=" $base/api"
test_packages[1]+=" $base/command"
test_packages[1]+=" $base/sdk/helper/keysutil"

# Total time: 1160
test_packages[2]+=" $base/sdk/helper/ocsp"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[2]+=" $base/vault/external_tests/replication-perf"
fi

# Total time: 1009
test_packages[3]+=" $base/builtin/credential/approle"
test_packages[3]+=" $base/command/agent/sink/file"
test_packages[3]+=" $base/command/agent/template"
test_packages[3]+=" $base/helper/random"
test_packages[3]+=" $base/helper/storagepacker"
test_packages[3]+=" $base/sdk/helper/certutil"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[3]+=" $base/vault/external_tests/entropy"
fi
test_packages[3]+=" $base/vault/external_tests/raft"

# Total time: 830
test_packages[4]+=" $base/builtin/plugin"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[4]+=" $base/enthelpers/fsm"
fi
test_packages[4]+=" $base/http"
test_packages[4]+=" $base/sdk/helper/pluginutil"
test_packages[4]+=" $base/serviceregistration/kubernetes"
test_packages[4]+=" $base/tools/godoctests/pkg/analyzer"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[4]+=" $base/vault/external_tests/apilock"
    test_packages[4]+=" $base/vault/external_tests/filteredpaths"
    test_packages[4]+=" $base/vault/external_tests/perfstandby"
    test_packages[4]+=" $base/vault/external_tests/replication-dr"
fi


# Total time: 258
test_packages[5]+=" $base/builtin/credential/aws"
test_packages[5]+=" $base/builtin/credential/cert"
test_packages[5]+=" $base/builtin/logical/aws"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[5]+=" $base/enthelpers/logshipper"
    test_packages[5]+=" $base/enthelpers/merkle"
fi
test_packages[5]+=" $base/helper/hostutil"
test_packages[5]+=" $base/helper/pgpkeys"
test_packages[5]+=" $base/sdk/physical/inmem"
test_packages[5]+=" $base/vault/activity"
test_packages[5]+=" $base/vault/diagnose"
test_packages[5]+=" $base/vault/external_tests/pprof"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[5]+=" $base/vault/external_tests/resolver"
fi
test_packages[5]+=" $base/vault/external_tests/response"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[5]+=" $base/vault/external_tests/seal"
fi
test_packages[5]+=" $base/vault/external_tests/sealmigration"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[5]+=" $base/vault/external_tests/transform"
fi

# Total time: 588
test_packages[6]+=" $base"
test_packages[6]+=" $base/audit"
test_packages[6]+=" $base/builtin/audit/file"
test_packages[6]+=" $base/builtin/credential/github"
test_packages[6]+=" $base/builtin/credential/okta"
test_packages[6]+=" $base/builtin/logical/database/dbplugin"
test_packages[6]+=" $base/command/agent/auth/cert"
test_packages[6]+=" $base/command/agent/auth/jwt"
test_packages[6]+=" $base/command/agent/auth/kerberos"
test_packages[6]+=" $base/command/agent/auth/kubernetes"
test_packages[6]+=" $base/command/agent/auth/token-file"
test_packages[6]+=" $base/command/agent/cache"
test_packages[6]+=" $base/command/agent/cache/cacheboltdb"
test_packages[6]+=" $base/command/agent/cache/cachememdb"
test_packages[6]+=" $base/command/agent/cache/keymanager"
test_packages[6]+=" $base/command/agent/config"
test_packages[6]+=" $base/command/config"
test_packages[6]+=" $base/command/token"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[6]+=" $base/enthelpers/namespace"
    test_packages[6]+=" $base/enthelpers/replicatedpaths"
    test_packages[6]+=" $base/enthelpers/sealrewrap"
fi
test_packages[6]+=" $base/helper/builtinplugins"
test_packages[6]+=" $base/helper/dhutil"
test_packages[6]+=" $base/helper/fairshare"
test_packages[6]+=" $base/helper/flag-kv"
test_packages[6]+=" $base/helper/flag-slice"
test_packages[6]+=" $base/helper/forwarding"
test_packages[6]+=" $base/helper/logging"
test_packages[6]+=" $base/helper/metricsutil"
test_packages[6]+=" $base/helper/namespace"
test_packages[6]+=" $base/helper/osutil"
test_packages[6]+=" $base/helper/parseip"
test_packages[6]+=" $base/helper/policies"
test_packages[6]+=" $base/helper/testhelpers/logical"
test_packages[6]+=" $base/helper/timeutil"
test_packages[6]+=" $base/helper/useragent"
test_packages[6]+=" $base/helper/versions"
test_packages[6]+=" $base/internalshared/configutil"
test_packages[6]+=" $base/internalshared/listenerutil"
test_packages[6]+=" $base/physical/alicloudoss"
test_packages[6]+=" $base/physical/gcs"
test_packages[6]+=" $base/physical/manta"
test_packages[6]+=" $base/physical/mssql"
test_packages[6]+=" $base/physical/oci"
test_packages[6]+=" $base/physical/s3"
test_packages[6]+=" $base/physical/spanner"
test_packages[6]+=" $base/physical/swift"
test_packages[6]+=" $base/physical/zookeeper"
test_packages[6]+=" $base/plugins/database/hana"
test_packages[6]+=" $base/plugins/database/redshift"
test_packages[6]+=" $base/sdk/database/dbplugin/v5"
test_packages[6]+=" $base/sdk/database/helper/credsutil"
test_packages[6]+=" $base/sdk/helper/authmetadata"
test_packages[6]+=" $base/sdk/helper/compressutil"
test_packages[6]+=" $base/sdk/helper/cryptoutil"
test_packages[6]+=" $base/sdk/helper/identitytpl"
test_packages[6]+=" $base/sdk/helper/kdf"
test_packages[6]+=" $base/sdk/helper/locksutil"
test_packages[6]+=" $base/sdk/helper/pathmanager"
test_packages[6]+=" $base/sdk/helper/roottoken"
test_packages[6]+=" $base/sdk/helper/testhelpers/schema"
test_packages[6]+=" $base/sdk/helper/xor"
test_packages[6]+=" $base/sdk/physical/file"
test_packages[6]+=" $base/sdk/plugin/pb"
test_packages[6]+=" $base/serviceregistration/kubernetes/client"
test_packages[6]+=" $base/shamir"
test_packages[6]+=" $base/vault/cluster"
test_packages[6]+=" $base/vault/eventbus"
test_packages[6]+=" $base/vault/external_tests/api"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[6]+=" $base/vault/external_tests/consistencyheaders"
fi
test_packages[6]+=" $base/vault/external_tests/expiration"
test_packages[6]+=" $base/vault/external_tests/hcp_link"
test_packages[6]+=" $base/vault/external_tests/kv"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[6]+=" $base/vault/external_tests/plugins"
fi
test_packages[6]+=" $base/vault/external_tests/quotas"
test_packages[6]+=" $base/vault/seal"

# Total time: 389
test_packages[7]+=" $base/builtin/credential/userpass"
test_packages[7]+=" $base/builtin/logical/pki"
test_packages[7]+=" $base/builtin/logical/transit"
test_packages[7]+=" $base/command/agent"
test_packages[7]+=" $base/helper/monitor"
test_packages[7]+=" $base/sdk/database/helper/connutil"
test_packages[7]+=" $base/sdk/database/helper/dbutil"
test_packages[7]+=" $base/sdk/helper/cidrutil"
test_packages[7]+=" $base/sdk/helper/custommetadata"
test_packages[7]+=" $base/sdk/helper/jsonutil"
test_packages[7]+=" $base/sdk/helper/ldaputil"
test_packages[7]+=" $base/sdk/helper/logging"
test_packages[7]+=" $base/sdk/helper/policyutil"
test_packages[7]+=" $base/sdk/helper/salt"
test_packages[7]+=" $base/sdk/helper/template"
test_packages[7]+=" $base/sdk/helper/useragent"
test_packages[7]+=" $base/sdk/logical"
test_packages[7]+=" $base/sdk/plugin/mock"
test_packages[7]+=" $base/sdk/queue"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[7]+=" $base/vault/autosnapshots"
    test_packages[7]+=" $base/vault/external_tests/activity"
fi
test_packages[7]+=" $base/vault/external_tests/approle"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[7]+=" $base/vault/external_tests/kmip"
fi
test_packages[7]+=" $base/vault/external_tests/mfa"
test_packages[7]+=" $base/vault/external_tests/misc"
test_packages[7]+=" $base/vault/quotas"

# Total time: 779
test_packages[8]+=" $base/builtin/credential/aws/pkcs7"
test_packages[8]+=" $base/builtin/logical/totp"
test_packages[8]+=" $base/command/agent/auth"
test_packages[8]+=" $base/physical/raft"
test_packages[8]+=" $base/sdk/framework"
test_packages[8]+=" $base/sdk/plugin"
test_packages[8]+=" $base/vault"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[8]+=" $base/vault/external_tests/barrier"
    test_packages[8]+=" $base/vault/external_tests/cubbyholes"
fi
test_packages[8]+=" $base/vault/external_tests/metrics"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[8]+=" $base/vault/external_tests/replication"
fi
test_packages[8]+=" $base/vault/external_tests/router"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[8]+=" $base/vault/external_tests/system"
    test_packages[8]+=" $base/vault/managed_key"
fi

# Total time: 310
test_packages[9]+=" $base/vault/hcp_link/capabilities/api_capability"
test_packages[9]+=" $base/vault/external_tests/plugin"

# Total time: 925
test_packages[10]+=" $base/builtin/credential/ldap"
test_packages[10]+=" $base/builtin/logical/database"
test_packages[10]+=" $base/physical/etcd"
test_packages[10]+=" $base/physical/postgresql"

# Total time: 851
test_packages[11]+=" $base/builtin/logical/rabbitmq"
test_packages[11]+=" $base/physical/dynamodb"
test_packages[11]+=" $base/plugins/database/influxdb"
test_packages[11]+=" $base/vault/external_tests/identity"
test_packages[11]+=" $base/vault/external_tests/token"

# Total time: 340
test_packages[12]+=" $base/builtin/logical/consul"
test_packages[12]+=" $base/physical/couchdb"
test_packages[12]+=" $base/plugins/database/mongodb"
test_packages[12]+=" $base/plugins/database/mssql"
test_packages[12]+=" $base/plugins/database/mysql"

# Total time: 704
test_packages[13]+=" $base/builtin/logical/pkiext"
test_packages[13]+=" $base/command/server"
test_packages[13]+=" $base/physical/aerospike"
test_packages[13]+=" $base/physical/cockroachdb"
test_packages[13]+=" $base/plugins/database/postgresql"
test_packages[13]+=" $base/plugins/database/postgresql/scram"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[13]+=" $base/vault/external_tests/filteredpathsext"
fi
test_packages[13]+=" $base/vault/external_tests/policy"

# Total time: 374
test_packages[14]+=" $base/builtin/credential/radius"
test_packages[14]+=" $base/builtin/logical/ssh"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[14]+=" $base/enthelpers/wal"
fi
test_packages[14]+=" $base/physical/azure"
test_packages[14]+=" $base/serviceregistration/consul"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[14]+=" $base/vault/external_tests/quotas-docker"
fi
test_packages[14]+=" $base/vault/external_tests/raftha"

# Total time: 362
test_packages[15]+=" $base/builtin/logical/nomad"
test_packages[15]+=" $base/physical/mysql"
test_packages[15]+=" $base/plugins/database/cassandra"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[15]+=" $base/vault/external_tests/namespaces"
fi
test_packages[15]+=" $base/vault/external_tests/sealmigrationext"

# Total time: 635
test_packages[16]+=" $base/physical/cassandra"
test_packages[16]+=" $base/physical/consul"
if [ "${ENTERPRISE:+x}" == "x" ] ; then
    test_packages[16]+=" $base/vault/external_tests/autosnapshots"
    test_packages[16]+=" $base/vault/external_tests/replicationext"
    test_packages[16]+=" $base/vault/external_tests/sealext"
fi
