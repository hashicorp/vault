set -exo pipefail
test_packages=()
test_packages[1]+=' github.com/hashicorp/vault/api'
test_packages[1]+=' github.com/hashicorp/vault/command'
test_packages[1]+=' github.com/hashicorp/vault/sdk/helper/keysutil'

# combine the list to a string separated by space and quote it 
bar=$(IFS=" " ; echo "${test_packages[1]}")
bar1='"'"$bar"'"'
rerunFails="--rerun-fails --packages $bar"
go run gotest.tools/gotestsum --format=short-verbose  $rerunFails