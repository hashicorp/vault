set -exo pipefail
# test_packages=()
# test_packages[1]+=' github.com/hashicorp/vault/api'
# test_packages[1]+=' github.com/hashicorp/vault/command'
# test_packages[1]+=' github.com/hashicorp/vault/sdk/helper/keysutil'
test_packages="github.com/hashicorp/vault/plugins/database/mysql github.com/hashicorp/vault/physical/cassandra github.com/hashicorp/vault/serviceregistration/consul github.com/hashicorp/vault/tools/codechecker/pkg/godoctests github.com/hashicorp/vault/vault/external_tests/expiration github.com/hashicorp/vault/helper/logging github.com/hashicorp/vault/command/token github.com/hashicorp/vault/helper/hostutil github.com/hashicorp/vault/helper/metricsutil github.com/hashicorp/vault/internalshared/configutil github.com/hashicorp/vault/command/agentproxyshared/auth/token-file github.com/hashicorp/vault/command/agentproxyshared/cache/cachememdb github.com/hashicorp/vault/physical/swift"

# combine the list to a string separated by space and quote it 
bar=$(IFS=" " ; echo "${test_packages[1]}")
bar1='"'"$bar"'"'
# rerunFails="--rerun-fails --packages "${test_packages[1]}""
rerunFails="--rerun-fails"
go run gotest.tools/gotestsum --format=short-verbose \
              --junitfile results-1.xml \
              --jsonfile results-1.json \
              --jsonfile-timing-events failure-summary-1.json \
              --packages "${test_packages}" \
              -- \
              -tags ,deadlock \



# go run gotest.tools/gotestsum --format=short-verbose --rerun-fails --packages "github.com/hashicorp/vault/api github.com/hashicorp/vault/command github.com/hashicorp/vault/sdk/helper/keysutil" --junitfile results-1.xml --jsonfile results-1.json --jsonfile-timing-events failure-summary-1.json
# --rerun-fails --packages "$bar"