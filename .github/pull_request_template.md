### Description
What does this PR do?

### TODO only if you're a HashiCorp employee
- [ ] **Backport Labels:** If this PR is in the ENT repo and needs to be backported, backport  
  to N, N-1, and N-2, using the `backport/ent/x.x.x+ent` labels. If this PR is in the CE repo, you should only backport to N, using the `backport/x.x.x` label, not the enterprise labels.
    - [ ] If this fixes a critical security vulnerability or [severity 1](https://www.hashicorp.com/customer-success/enterprise-support) bug, it will also need to be backported to the current [LTS versions](https://developer.hashicorp.com/vault/docs/enterprise/lts#why-is-there-a-risk-to-updating-to-a-non-lts-vault-enterprise-version) of Vault. To ensure this, use **all** available enterprise labels.
- [ ] **ENT Breakage:** If this PR either 1) removes a public function OR 2) changes the signature
  of a public function, even if that change is in a CE file, _double check_ that
  applying the patch for this PR to the ENT repo and running tests doesn't
  break any tests. Sometimes ENT only tests rely on public functions in CE
  files.
- [ ] **Jira:** If this change has an associated Jira, it's referenced either
  in the PR description, commit message, or branch name.
- [ ] **RFC:** If this change has an associated RFC, please link it in the description.
- [ ] **ENT PR:** If this change has an associated ENT PR, please link it in the
  description. Also, make sure the changelog is in this PR, _not_ in your ENT PR.
