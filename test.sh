set -exo pipefail


# branchName="vault-17777-fix2"
branchName2="release/vault-17777-fix3"
# if [[ -n "$branchName" && "$branchName" == release/* ]] || [[ "$branchName2" == release/* ]]

# if [[ "${{ github.base_ref }}" = "release/*" ]] || [[ "${{ github.ref_name }}" = "release/*" ]]
if [[ "$branchName" = release/* ]] || [[ "$branchName2" = release/* ]]
then
    RERUN_FAILS="--rerun-fails"
fi
go run gotest.tools/gotestsum --format=short-verbose \
              --junitfile test-results/go-test/results-1.xml \
              --jsonfile test-results/go-test/results-1.json \
              $RERUN_FAILS \
              --packages "github.com/hashicorp/vault/vault/external_tests/sealmigration github.com/hashicorp/vault/physical/azure github.com/hashicorp/vault/builtin/plugin github.com/hashicorp/vault/physical/cassandra github.com/hashicorp/vault/physical/etcd github.com/hashicorp/vault/helper/random github.com/hashicorp/vault/tools/codechecker/pkg/gonilnilfunctions github.com/hashicorp/vault/command/token github.com/hashicorp/vault/helper/metricsutil github.com/hashicorp/vault/command/agentproxyshared/auth/jwt github.com/hashicorp/vault/helper/monitor github.com/hashicorp/vault/helper/policies github.com/hashicorp/vault/shamir github.com/hashicorp/vault/helper/timeutil" \
              -- \
              -timeout=60m 