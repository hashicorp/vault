set -exo pipefail
test_packages=()
test_packages[1]+=' github.com/hashicorp/vault/api'
test_packages[1]+=' github.com/hashicorp/vault/command'
test_packages[1]+=' github.com/hashicorp/vault/sdk/helper/keysutil'

# combine the list to a string separated by space and quote it 
bar=$(IFS=" " ; echo "${test_packages[1]}")
bar1='"'"$bar"'"'
rerunFails="--rerun-fails --packages "${test_packages[1]}""
go run gotest.tools/gotestsum --format=short-verbose \
              --junitfile results-1.xml \
              --jsonfile results-1.json \
              --jsonfile-timing-events failure-summary-1.json \
              $rerunFails \
              -- \
              -tags ,deadlock

# go run gotest.tools/gotestsum --format=short-verbose --rerun-fails --packages "github.com/hashicorp/vault/api github.com/hashicorp/vault/command github.com/hashicorp/vault/sdk/helper/keysutil" --junitfile results-1.xml --jsonfile results-1.json --jsonfile-timing-events failure-summary-1.json
# --rerun-fails --packages "$bar"