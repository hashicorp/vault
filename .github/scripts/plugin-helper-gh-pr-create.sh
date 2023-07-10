#!/usr/bin/env bash

set -e

# this script expects the following env vars to be set
# error if these are not set
[ ${GITHUB_TOKEN:?} ]
[ ${PLUGIN_REPO:?} ]
[ ${VAULT_BRANCH:?} ]
[ ${PLUGIN_BRANCH:?} ]
[ ${RUN_ID:?} ]

# we are using the GH API directly so that we can get the resulting
# PR URL from the JSON response

resp=$(curl -SL \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}"\
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/hashicorp/vault/pulls \
  --data @<(cat <<EOF
{
  "title":"[DO NOT MERGE]: $PLUGIN_REPO Automated plugin update check",
  "body":"Updates $PLUGIN_REPO to verify vault CI. Full log: https://github.com/hashicorp/vault/actions/runs/$RUN_ID",
  "head":"$VAULT_BRANCH",
  "base":"main"
}
EOF
)
)

echo "captured response:"
echo "$resp" | jq .

# get Vault PR number
vault_pr_num=$(echo "$resp" | jq -r '.number')
vault_pr_url=$(echo "$resp" | jq -r '.html_url')
echo "Vault PR number: $vault_pr_url"

# add labels to Vault PR - this requires a wider permission set than we currently have available as a repo token
#reviewers="austingebauer,fairclothjm,vinay-gopalan,maxcoulombe,robmonte,Zlaticanin,kpcraig,raymonstah"
#gh pr edit "$vault_pr_num" --add-label "dependencies,pr/no-changelong,pr/no-milestone" --repo hashicorp/vault
#gh pr edit "$vault_pr_num" --add-reviewer "$reviewers"

# get Plugin PR number
plugin_pr_num=$(gh pr list --head "$PLUGIN_BRANCH" --json number --repo "$PLUGIN_REPO" -q '.[0].number')

# make a comment on the plugin repo's PR
gh pr comment "$plugin_pr_num" --body "Vault CI check PR: $vault_pr_url" --repo "$PLUGIN_REPO"
