printf "Non-database plugins that live in Vault except kv, cubbyhole, identity, pki, and transit\n"
grep -ro "\"github.com/hashicorp/vault/builtin/[a-z]*/[a-z]*.*\"" \
--exclude="*_test.go" \
--exclude="*main.go" \
--exclude="*registry*.go" \
--exclude="*test*helper*" \
--exclude="*find-plugin-refs*" \
| grep -v \
-e "github.com/hashicorp/vault/builtin/logical/transit.*" \
-e "github.com/hashicorp/vault/builtin/logical/pki.*" \
-e "github.com/hashicorp/vault/builtin/logical/database.*" \
| sort -t: -k2 | column -t -s ':'

printf "\n\nDatabase plugins that live in Vault\n"
grep -ro "\"github.com/hashicorp/vault/plugins/database/[a-z]*.*\"" \
--exclude="*_test.go" \
--exclude="*main.go" \
--exclude="*registry*.go" \
--exclude="*test*helper*" \
--exclude="*find-plugin-refs*" \
| sort -t: -k2 | column -t -s ':'

printf "\n\nIndependent plugins that imported to Vault\n"
grep -ro "\"github.com/hashicorp/vault-plugin-[a-z]*-[a-z]*.*\"" \
--exclude="*_test.go" \
--exclude="*main.go" \
--exclude="*registry*.go" \
--exclude="*test*helper*" \
--exclude="*find-plugin-refs*" \
| grep -v \
-e "github.com/hashicorp/vault-plugin-secrets-kv.*" \
| sort -t: -k2 | column -t -s ':'

printf "\n\nStorage plugins - all currently live in Vault\n"
grep -ro "\"github.com/hashicorp/vault/physical/[a-z]*.*\"" \
--exclude="*_test.go" \
--exclude="*main.go" \
--exclude="*test*helper*" \
--exclude="*find-plugin-refs*" \
| grep -v \
-e "\"github.com/hashicorp/vault/physical/raft.*\"" \
-e "\"github.com/hashicorp/vault/physical/consul.*\"" \
| sort -t: -k2 | column -t -s ':'