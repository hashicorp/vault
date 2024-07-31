### Description
What does this PR do?

### TODO only if you're a HashiCorp employee
- [ ] **Enterprise Labels:** If this PR is the CE portion of an ENT change, you should backport
  to N, N-1, and N-2, using the `backport/ent/x.x.x+ent` labels
  instead of the `backport/x.x.x` labels.
- [ ] **CE Labels:** If this PR is a CE-only change, you should only backport to N, so use
  the CE-only `backport/x.x.x` label (there should be only 1).
- [ ] **LTS Labels**: If this PR contains a fix for a critical security vulnerability or [severity 1](https://www.hashicorp.com/customer-success/enterprise-support) bug, it will also need to be backported to the current LTS branches. If an LTS version of Vault is further back than N-2, be sure to add the appropriate enterprise label (`backport/ent/x.x.x+ent`) for that branch.
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
