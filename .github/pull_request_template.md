### Description
What does this PR do?

### TODO
- [ ] Make sure there's a changelog. Brand new features have a differently
  formatted changelog than other changelogs, so pay attention to that. If this
  PR is the CE portion of an ENT PR, put the changelog here, _not_ in your ENT
  PR.
- [ ] Add a milestone. If this change is going to be backported, each backport
  PR should receive the milestone for the exact version it's being backported
  to. For this PR, just pick the biggest version.
- [ ] If this PR is the CE portion of an ENT change, and that ENT change is
  getting backported to N-2, use the new style `backport/ent/x.x.x+ent` labels
  instead of the old style `backport/x.x.x` labels.
- [ ] If this PR is a CE only change, it can only be backported to N, so use
  the normal `backport/x.x.x` label (there should be only 1).
- [ ] If this PR either 1) removes a public function OR 2) changes the signature
  of a public function, even if that change is in a CE file, _double check_ that
  applying the patch for this PR to the ENT repo and running tests doesn't
  break any tests. Sometimes ENT only tests rely on public functions in CE
  files.
