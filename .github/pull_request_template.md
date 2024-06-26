### Description
What does this PR do?

### TODO only if you're a HashiCorp employee
- [ ] **Labels:** If this PR is the CE portion of an ENT change, and that ENT change is
  getting backported to N-2, use the new style `backport/ent/x.x.x+ent` labels
  instead of the old style `backport/x.x.x` labels.
- [ ] **Labels:** If this PR is a CE only change, it can only be backported to N, so use
  the normal `backport/x.x.x` label (there should be only 1).
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
